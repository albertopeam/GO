package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/go-bip39"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type account struct {
	privKey secp256k1.PrivKey
	pubKey  types.PubKey
	address sdk.AccAddress
}

func (a account) String() string {
	return fmt.Sprintf("Private Key %s\n Public Key %s\n Address %s\n", a.privKey, a.privKey, a.address.String())
}

func main() {
	// Read State in mainnet
	queryMainnetState()

	// Write state in testnet
	newAccounts := true //TODO: Change/Inject from command line parameters
	var from, to account
	if newAccounts {
		from, to = createAccounts()
	} else {
		from, to = loadAccounts("from.txt", "to.txt")
	}
	printAccounts(from, to)
	waitForUserToTransferCoinsTo(from)
	verifyBalance(to, "before")
	sendTransaction(from, to)
	verifyBalance(to, "after")
}

// full tutorial https://docs.cosmos.network/v0.46/run-node/interact-node.html
func queryMainnetState() {
	// create an addr. doc https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types
	sg1Addr, err := sdk.AccAddressFromBech32("cosmos196ax4vc0lwpxndu9dyhvca7jhxp70rmcfhxsrt")
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	fmt.Println("using", sg1Addr.String())

	// Create a connection to the gRPC server. doc https://pkg.go.dev/google.golang.org/grpc
	// more info on how to find a host and port on https://github.com/cosmos/chain-registry/blob/master/cosmoshub/chain.json
	grpcConn, err := grpc.Dial(
		"cosmoshub.strange.love:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	defer grpcConn.Close()
	fmt.Println("open gRPC connection", grpcConn.Target())

	// This creates a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)

	// query uatom balance for an account
	bankRes, err := bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: sg1Addr.String(), Denom: "uatom"})
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	fmt.Println("atom balance", bankRes.String())

	// query all balances for an account
	bankAllRes, err := bankClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: sg1Addr.String()})
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	fmt.Println("all balances", bankAllRes.String())
	fmt.Println("-----------------------------------")
}

func printAccounts(from account, to account) {
	fmt.Println("from", from)
	fmt.Println("to", to)
}

// Load from disk if exists the files or creates new accounts
//DOC
// https://en.wikipedia.org/wiki/Digital_signature
// https://pkg.go.dev/github.com/cosmos/go-bip39#section-readme
// https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/crypto/keys/secp256k1
// https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/crypto/hd
//BIP39
// https://github.com/cosmos/go-bip39 (COSMOS FORK)
// https://iancoleman.io/bip39/#english
//BIP44
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose
func loadAccounts(fromFile string, toFile string) (account, account) {
	// Read from disk
	fromBody, err := ioutil.ReadFile(fromFile)
	if err != nil && err != io.EOF {
		log.Fatalf("ioutil.ReadFile() %s %s", fromFile, err)
	}
	toBody, err := ioutil.ReadFile(toFile)
	if err != nil && err != io.EOF {
		log.Fatalf("ioutil.ReadFile() %s %s", toFile, err)
	}
	// Generate seed from mnemonic
	toSeed := generateSeed(toBody)
	fmt.Println("BIP39 seed 'to'", string(fromBody), hex.EncodeToString(toSeed))
	fromSeed := generateSeed(fromBody)
	fmt.Println("BIP39 seed 'from'", string(toBody), hex.EncodeToString(fromSeed))
	// Derive Atom account from seed
	fromAcc := deriveAtomAccountFromSeed(fromSeed)
	toAcc := deriveAtomAccountFromSeed(toSeed)

	return fromAcc, toAcc
}

func createAccounts() (account, account) {
	return createAccount(), createAccount()
}

func createAccount() account {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		log.Fatalf("bip39.NewEntropy error %s", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatalf("bip39.NewMnemonic error %s", err)
	}
	fmt.Println("mnemonic", mnemonic)
	// Generate seed from mnemonic
	seed := generateSeed([]byte(mnemonic))
	fmt.Println("BIP39 seed", hex.EncodeToString(seed))
	// Derive Atom account from seed
	newAccount := deriveAtomAccountFromSeed(seed)

	return newAccount
}

func generateSeed(mnemonicBytes []byte) []byte {
	mnemonic := string(mnemonicBytes)
	password := ""
	seedBytes := bip39.NewSeed(mnemonic, password)
	return seedBytes
}

func deriveAtomAccountFromSeed(seed []byte) account {
	// Derivation Path
	master, ch := hd.ComputeMastersFromSeed([]byte(seed))
	atomPath := "m/44'/118'/0'/0/0" // check BIP44 standard https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose or https://iancoleman.io/bip39/#english to get the path
	atomPriv, err := hd.DerivePrivateKeyForPath(master, ch, atomPath)
	if err != nil {
		log.Fatalf("hd.DerivePrivateKeyForPath error %s", err)
	}
	// Keys
	var privKey secp256k1.PrivKey = secp256k1.PrivKey{Key: atomPriv}
	var pubKey types.PubKey = privKey.PubKey()
	var address sdk.AccAddress = sdk.AccAddress(pubKey.Address().Bytes())
	return account{privKey: privKey, pubKey: pubKey, address: address}
}

func waitForUserToTransferCoinsTo(from account) {
	input := bufio.NewScanner(os.Stdin)
	fmt.Printf("Waiting until %s has enough balance to make a transaction. Then tap enter to continue!", from.address.String())
	input.Scan()
}

func sendTransaction(from account, to account) {
	// connect to testnet

	// send a transaction
}

func verifyBalance(to account, tag string) {
	// verify balance changed
}
