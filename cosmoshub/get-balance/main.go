package main

import (
	"context"
	"encoding/hex"
	"fmt"
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

func main() {
	//queryState()
	createOrLoadAccounts(true)
	sendTransaction()
	verifyBalbance()
}

// full tutorial https://docs.cosmos.network/v0.46/run-node/interact-node.html
func queryState() {
	// create an addr. doc https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types
	sg1Addr, err := sdk.AccAddressFromBech32("cosmos196ax4vc0lwpxndu9dyhvca7jhxp70rmcfhxsrt")
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
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
		fmt.Println("Error", err)
		os.Exit(1)
	}
	defer grpcConn.Close()
	fmt.Println("open gRPC connection", grpcConn.Target())

	// This creates a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)

	// query uatom balance for an account
	bankRes, err := bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: sg1Addr.String(), Denom: "uatom"})
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	fmt.Println("atom balance", bankRes.String())

	// query all balances for an account
	bankAllRes, err := bankClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: sg1Addr.String()})
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	fmt.Println("all balances", bankAllRes.String())
}

func createOrLoadAccounts(generated bool) {
	//INFO 	 https://en.wikipedia.org/wiki/Digital_signature
	//OPT1 	 https://pkg.go.dev/github.com/tyler-smith/go-bip32#section-readme
	//OPT2 	 https://pkg.go.dev/github.com/cosmos/go-bip39#section-readme
	//COSMOS crypto/hd doc https://pkg.go.dev/github.com/cosmos/cosmos-sdk/crypto/hd

	//TODO: read from disc if exists. For this example we don't need to encrypt it

	//COSMOS
	// seed := bip39.NewSeed("blast about old claw current first paste risk involve victory edit current", "")
	// master, ch := hd.ComputeMastersFromSeed(seed)
	// path := "m/44'/1022'/0'/0/0'"
	// priv, err := hd.DerivePrivateKeyForPath(master, ch, path)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// var privKey secp256r1.PrivKey
	// var pubKey types.PubKey = privKey.PubKey()
	// var addr sdk.AccAddress = sdk.AccAddress(pubKey.Address().Bytes())

	// fmt.Println(privKey)
	// fmt.Println(pubKey)
	// fmt.Println(addr)

	//OPT1
	// Generate a seed to determine all keys from.
	// This should be persisted, backed up, and secured
	// seed, err := bip32.NewSeed()
	// if err != nil {
	// 	fmt.Println("bip32.NewSeed returned", err)
	// 	os.Exit(1)
	// }

	// Create master private key from seed. Type *bip32.Key
	// masterKey, err := bip32.NewMasterKey(seed)
	// if err != nil {
	// 	fmt.Println("bip32.NewMasterKey returned", err)
	// 	os.Exit(1)
	// }

	// There is a very small chance a given child index is invalid
	// If so your real program should handle this by skipping the index
	// Type *bip32.Key
	// privKey, err := masterKey.NewChildKey(0)
	// if err != nil {
	// 	fmt.Println("masterKey.NewChildKey returned", err)
	// 	os.Exit(1)
	// }

	// pubKey := privKey.PublicKey()

	// fmt.Println("seed", seed)
	// fmt.Println("privKey", privKey)
	// fmt.Println("pubKey", pubKey)
	// fmt.Println("account", pubKey.account) account is a cosmos concept. DOC?

	//DOC
	// https://pkg.go.dev/github.com/cosmos/go-bip39#section-readme
	// https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/crypto/keys/secp256k1
	// https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/crypto/hd
	//BIP39
	// https://github.com/cosmos/go-bip39 (COSMOS FORK)
	// https://iancoleman.io/bip39/#english
	//BIP44
	// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose

	var seed string
	if generated {
		// Generate a mnemonic for memorization or user-friendly seeds
		entropy, err := bip39.NewEntropy(256)
		if err != nil {
			fmt.Println("bip39.NewEntropy error", err)
			os.Exit(1)
		}
		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			fmt.Println("bip39.NewMnemonic error", err)
			os.Exit(1)
		}
		fmt.Println("mnemonic", mnemonic)
		// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
		password := ""
		seedBytes := bip39.NewSeed(mnemonic, password)
		seed = string(seedBytes)
		fmt.Println("BIP39 seed(validate in https://iancoleman.io/bip39)", seed) // validate the seed in https://iancoleman.io/bip39
	} else {
		// it MUST generate addr cosmos19kzdcmysekqu926fwdcjg5pdqlx3saujcldys5 for path "m/44'/118'/0'/0/0"
		mnemonic := "write sense wage direct salute north now dog divorce inflict pole provide spike welcome bring sister fetch upset chimney direct siren trash cruise mother" // generated using https://iancoleman.io/bip39
		password := ""
		seedBytes := bip39.NewSeed(mnemonic, password)
		seed = string(seedBytes)
		fmt.Println("BIP39 seed", hex.EncodeToString(seedBytes))
	}

	// Derivation Path
	master, ch := hd.ComputeMastersFromSeed([]byte(seed))
	atomPath := "m/44'/118'/0'/0/0" // check BIP44 standard https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose or https://iancoleman.io/bip39/#english to get the path
	atomPriv, err := hd.DerivePrivateKeyForPath(master, ch, atomPath)
	if err != nil {
		fmt.Println("hd.DerivePrivateKeyForPath error", err)
		os.Exit(1)
	}
	// Keys
	var privKey secp256k1.PrivKey = secp256k1.PrivKey{Key: atomPriv}
	var pubKey types.PubKey = privKey.PubKey()
	var address sdk.AccAddress = sdk.AccAddress(pubKey.Address().Bytes())
	fmt.Println("Private Key", privKey)
	fmt.Println("Public Key", pubKey)
	fmt.Println("Address", address.String())
}

func sendTransaction() {
	// create/load accounts from disk if exists https://pkg.go.dev/github.com/cosmos/go-bip39#section-readme

	// faucet to get coins

	// connect to testnet

	// send a transaction

	// verify balance changed
}

func verifyBalbance() {

}
