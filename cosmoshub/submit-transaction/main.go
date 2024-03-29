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
	"time"

	"cosmossdk.io/math"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typestx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	accounts "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/go-bip39"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tenderminttypes "github.com/tendermint/tendermint/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// info on types https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/crypto/types
type account struct {
	privKey types.PrivKey
	pubKey  types.PubKey
	address sdk.AccAddress
}

func (a account) String() string {
	return fmt.Sprintf("Private Key %s\n Public Key %s\n Address %s\n", a.privKey, a.privKey, a.address.String())
}

func main() {
	// load or create accounts
	newAccounts := false //TODO: Change/Inject from command line parameters
	var from, to account
	if newAccounts {
		from, to = createAccounts()
	} else {
		from, to = loadAccounts("from.txt", "to.txt")
	}

	// create grpc connection
	grpcConn := createGrpcConn()
	defer grpcConn.Close()

	// print used accounts
	printAccounts(from, to)

	// verify balance before any transaction
	verifyBalance(grpcConn, from.address, "from before")
	verifyBalance(grpcConn, to.address, "to before")

	// wait to have funds on from address
	waitForUserToTransferCoinsTo(from)

	// send transaction
	tx := sendTransaction(grpcConn, from, to)

	// wait for transaction
	subscribeToTransactionConfirmation(tx)

	// verify balance after the transaction
	verifyBalance(grpcConn, from.address, "from after")
	verifyBalance(grpcConn, to.address, "to after")
}

func printAccounts(from account, to account) {
	fmt.Println("from", from)
	fmt.Println("to", to)
}

// Load from disk if exists the files or creates new accounts
// DOC
// https://en.wikipedia.org/wiki/Digital_signature
// https://pkg.go.dev/github.com/cosmos/go-bip39#section-readme
// https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/crypto/keys/secp256k1
// https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/crypto/hd
// BIP39
// https://github.com/cosmos/go-bip39 (COSMOS FORK)
// https://iancoleman.io/bip39/#english
// BIP44
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
	// path: check BIP44 standard https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose or https://iancoleman.io/bip39/#english to get the path
	atomPath := "m/44'/118'/0'/0/0" // we could use "sdk.FullFundraiserPath" instead. https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types#pkg-constants
	atomPriv, err := hd.DerivePrivateKeyForPath(master, ch, atomPath)
	if err != nil {
		log.Fatalf("hd.DerivePrivateKeyForPath error %s", err)
	}

	// Keys to create an account
	secp256k1Algo := hd.Secp256k1                   // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/crypto/hd#pkg-variables
	secp256k1GenerateFn := secp256k1Algo.Generate() // https://github.com/cosmos/cosmos-sdk/blob/main/crypto/hd/algo.go
	privKey := secp256k1GenerateFn(atomPriv)
	var pubKey types.PubKey = privKey.PubKey()
	var address sdk.AccAddress = sdk.AccAddress(pubKey.Address().Bytes())
	return account{privKey: privKey, pubKey: pubKey, address: address}
}

func waitForUserToTransferCoinsTo(from account) {
	input := bufio.NewScanner(os.Stdin)
	fmt.Printf("Waiting until %s has enough balance to make a transaction. Then tap enter to continue!\n", from.address.String())
	input.Scan()
}

func sendTransaction(grpcConn *grpc.ClientConn, from account, to account) *sdk.TxResponse {
	// retrieve account number and sequence number.
	var account = getAccount(grpcConn, from.address)

	// create the transaction
	txBytes := createTransaction(from, to, account.GetAccountNumber(), account.GetSequence())

	// broadcast transaction
	return broadcastTransaction(grpcConn, txBytes)
}

func createGrpcConn() *grpc.ClientConn {
	// Connect to testnet https://hub.cosmos.network/main/hub-tutorials/join-testnet.html
	grpcUrl := "rpc.sentry-01.theta-testnet.polypore.xyz:9090" // https://github.com/cosmos/testnets/tree/master/v7-theta/public-testnet
	grpcConn, err := grpc.Dial(
		grpcUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		log.Fatalf("grpc.Dial Error %s", err)
	}
	fmt.Println("Testnet URL", grpcUrl)
	return grpcConn
}

func getAccount(grpcConn *grpc.ClientConn, address sdk.Address) accounts.BaseAccount {
	// query account https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/auth/types#QueryClient
	accountRequest := accounts.QueryAccountRequest{Address: address.String()}
	accountClient := accounts.NewQueryClient(grpcConn)
	fmt.Println(accountRequest, accountClient)
	accRes, err := accountClient.Account(context.Background(), &accountRequest)
	if err != nil {
		log.Fatalf("accountClient.Account %s", err)
	}
	fmt.Println("Response.AccountType", accRes.Account.TypeUrl) // Type of the returned proto inside Account Any: /cosmos.auth.v1beta1.BaseAccount
	// Unpack any from protobuff explanation
	// https://docs.cosmos.network/v0.46/core/encoding.html#interface-encoding-and-usage-of-any
	// https://docs.cosmos.network/master/architecture/adr-019-protobuf-state-encoding.html
	// https://pkg.golang.ir/github.com/cosmos/cosmos-sdk/x/auth/types#QueryAccountResponse.UnpackInterfaces
	// https://docs.cosmos.network/v0.46/core/encoding.html
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry()) // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/codec#NewProtoCodec
	var acc accounts.BaseAccount                                  // returned concrete type inside Response.Account
	cdc.Unmarshal(accRes.Account.Value, &acc)                     // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/codec#ProtoCodec.Unmarshal
	fmt.Println("AccountNumber", acc.GetAccountNumber())
	fmt.Println("AccountSequence", acc.GetSequence())
	return acc
}

func createTransaction(from account, to account, accountNumber uint64, sequence uint64) []byte {
	//create the transaction builder
	signinModes := []signing.SignMode{signing.SignMode_SIGN_MODE_DIRECT} // https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/types/tx/signing#SignMode
	registry := codectypes.NewInterfaceRegistry()                        // https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/codec/types#NewInterfaceRegistry
	codec := codec.NewProtoCodec(registry)                               // https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/codec#ProtoCodecMarshaler https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/codec#ProtoCodec
	txConfig := tx.NewTxConfig(codec, signinModes)                       // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/auth/tx#NewTxConfig
	txBuilder := txConfig.NewTxBuilder()                                 // https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/client#TxConfig

	// sets amount to send in the txBuilder
	coin := sdk.NewCoin("uatom", math.NewInt(10000))             // 0.01 ATOM https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/types#NewCoin
	coins := sdk.NewCoins(coin)                                  // https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/types#NewCoins
	msg := banktypes.NewMsgSend(from.address, to.address, coins) // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/bank/types#NewMsgSend
	err := txBuilder.SetMsgs(msg)                                // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types#Msg
	if err != nil {
		log.Fatalf("txBuilder.SetMsgs error %s", err)
	}

	// TODO: investigate gasMeter
	// https://docs.cosmos.network/master/basics/gas-fees.html
	// https://docs.cosmos.network/master/basics/tx-lifecycle.html

	// setup limits https://docs.cosmos.network/v0.44/core/transactions.html#transaction-generation
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100)))) // the maximum amount the user is willing to pay in fees.
	txBuilder.SetGasLimit(200_000)                                               // max units of gas

	// Sign in transaction
	// main info https://docs.cosmos.network/master/run-node/txs.html
	// accounts https://docs.cosmos.network/master/basics/accounts.html
	// https://docs.cosmos.network/v0.46/modules/auth/02_state.html

	// First round: we gather all the signer infos. We use the "set empty signature" hack to do that.
	sigV2 := signing.SignatureV2{
		PubKey: from.pubKey,
		Data: &signing.SingleSignatureData{
			SignMode:  txConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: sequence,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		log.Fatalf("txBuilder.SetSignatures Populate SignerInfo error %s", err)
	}
	// Second round: all signer infos are set, so each signer can sign.
	chainID := "theta-testnet-001"         // https://github.com/cosmos/testnets/tree/master/v7-theta/public-testnet
	signerData := xauthsigning.SignerData{ // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/auth/signing
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
	sigV2, err = clienttx.SignWithPrivKey( // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/client/tx#SignWithPrivKey
		txConfig.SignModeHandler().DefaultMode(),
		signerData,
		txBuilder,
		from.privKey,
		txConfig,
		sequence)
	if err != nil {
		log.Fatalf("tx.SignWithPrivKey Sign error %s", err)
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		log.Fatalf("txBuilder.SetSignatures error %s", err)
	}

	// generate transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		log.Fatalf("txConfig.TxEncoder error %s", err)
	}
	fmt.Println("RAW transaction", string(txBytes))

	//TODO: not works printing to json
	// jsonBytes, err := txConfig.TxJSONEncoder()(txBuilder.GetTx())
	// if err != nil {
	// 	log.Fatalf("txConfig.TxJSONEncoder error %s", err)
	// }
	// txJSON := string(jsonBytes)
	// fmt.Println("json transaction", txJSON)

	return txBytes
}

func broadcastTransaction(grpcConn *grpc.ClientConn, txBytes []byte) *sdk.TxResponse {
	// Broadcast the tx via gRPC. We create a new client for the Protobuf Tx service.
	txClient := typestx.NewServiceClient(grpcConn) // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types/tx#NewServiceClient
	grpcRes, err := txClient.BroadcastTx(          // We then call the BroadcastTx method on this client.
		context.Background(),
		&typestx.BroadcastTxRequest{ // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types/tx#BroadcastTxRequest
			Mode:    typestx.BroadcastMode_BROADCAST_MODE_SYNC, // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types/tx#BroadcastMode
			TxBytes: txBytes,                                   // Proto-binary of the signed transaction, see previous step.
		},
	)
	if err != nil {
		log.Fatalf("txClient.BroadcastTx Error %s", err)
	}
	fmt.Println("GRPCResponse TXResponse", grpcRes.TxResponse) // Should be `0` if the tx is successful https://grpc.github.io/grpc/core/md_doc_statuscodes.html
	if grpcRes.TxResponse.Code == 0 {
		fmt.Println("Transaction Submited correctly", grpcRes.TxResponse.TxHash)
	}
	return grpcRes.TxResponse
}

func subscribeToTransactionConfirmation(txRes *sdk.TxResponse) {
	// https://docs.cosmos.network/master/core/events.html
	// https://tutorials.cosmos.network/academy/2-main-concepts/events.html#subscribing-to-events
	// https://docs.tendermint.com/v0.34/tendermint-core/subscription.html

	// TENDERMINT RPC OPENAPI DOC:
	// https://docs.tendermint.com/v0.34/rpc/
	// 	CURL 	curl -X GET "https://rpc.sentry-01.theta-testnet.polypore.xyz/net_info" -H "accept: application/json"
	// 	OPENAPI https://docs.tendermint.com/v0.34/rpc/#/Websocket/subscribe
	//	GO DOC CLIENT RPC HTTP	https://pkg.go.dev/github.com/tendermint/tendermint/rpc/client/http
	// Alternatives:
	//  GO DOC CLIENT RPCJSON 	https://pkg.go.dev/github.com/tendermint/tendermint/rpc/jsonrpc/client
	//  ¿GO DOC GRPC? 			https://pkg.go.dev/github.com/tendermint/tendermint/rpc/grpc

	// Using CLIENT RPC HTTP OPTION https://pkg.go.dev/github.com/tendermint/tendermint@v0.35.9/rpc/client/http#HTTP

	// Create Connection https://pkg.go.dev/github.com/tendermint/tendermint@v0.35.9/rpc/client/http#New
	client, err := rpchttp.New("https://rpc.sentry-01.theta-testnet.polypore.xyz:26657", "/websocket")
	if err != nil {
		log.Fatalf("rpchttp.Start Error %s", err)
	}
	// Start the client
	err = client.Start()
	if err != nil {
		log.Fatalf("rpchttp.Start Error %s", err)
	}
	// Stop the client when the function ends
	defer client.Stop()

	// Subscribe via query to events
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// doc for events https://pkg.go.dev/github.com/tendermint/tendermint/types#pkg-constants
	// doc for queries https://pkg.go.dev/github.com/tendermint/tendermint/libs/pubsub/query
	// doc for query syntax https://pkg.go.dev/github.com/tendermint/tendermint/libs/pubsub/query/syntax
	query := fmt.Sprintf("tm.event = 'Tx' AND tx.hash = '%s'", txRes.TxHash)
	fmt.Println("subscribe query", query)
	txs, err := client.Subscribe(ctx, "wait for tx", query) // Subscribe https://pkg.go.dev/github.com/tendermint/tendermint@v0.35.9/rpc/client/http#HTTP.Subscribe
	if err != nil {
		log.Fatalf("rpchttp.Subscribe Error %s", err)
	}

	// Receive events
	event := <-txs
	fmt.Println("got ", event.Data.(tenderminttypes.EventDataTx))
}

func verifyBalance(grpcConn *grpc.ClientConn, account sdk.Address, tag string) {
	// This creates a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)

	// query uatom balance for an account
	bankRes, err := bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: account.String(), Denom: "uatom"})
	if err != nil {
		fmt.Println("bankClient.Balance Error", err)
	}
	fmt.Println("atom balance", account.String(), bankRes.String())
}
