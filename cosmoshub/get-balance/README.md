# Cosmoshub(Gaia) blockchain info

Gaia info on [github](https://github.com/cosmos/gaia)

Golang cosmos [sdk](https://pkg.go.dev/github.com/cosmos/cosmos-sdk)

Chain registry info on [github](https://github.com/cosmos/chain-registry)

Tutorial on [cosmos](https://docs.cosmos.network/master/run-node/txs.html)

## How to make a connection

Info and examples on [cosmos doc](https://docs.cosmos.network/v0.46/run-node/interact-node.html#using-grpc)

It is possible to create connections via: GRPCurl, GO, JS and Rest endpoints

### GRPCurl

Checking the list of gRPC services: `grpcurl -plaintext <host:port> list`,

* ie: `grpcurl -plaintext cosmoshub.strange.love:9090 list`

For a full description of all the services: `grpcurl -plaintext <host:port> describe`

* ie: `grpcurl -plaintext cosmoshub.strange.love:9090 describe`

For a full description of one service: `grpcurl -plaintext <host:port> describe <method>`

* ie: `grpcurl -plaintext cosmoshub.strange.love:9090 describe cosmos.bank.v1beta1.Query`

Querying a service

* ie: `grpcurl -plaintext -d '{"address":"cosmos196ax4vc0lwpxndu9dyhvca7jhxp70rmcfhxsrt"}' cosmoshub.strange.love:9090 cosmos.bank.v1beta1.Query/AllBalances`

More info on [doc](https://docs.cosmos.network/v0.46/run-node/interact-node.html#grpcurl)

### GO

Example in the `main.go` file. Download dependencies: `go mod tidy`. Run: `go run main.go`

More info on [doc](https://docs.cosmos.network/v0.46/run-node/interact-node.html#programmatically-via-go)

### CosmJS

Not included.

More info on [doc](https://docs.cosmos.network/v0.46/run-node/interact-node.html#cosmjs)

### Rest endpoints(gRPC-gateway)

Querying a service via curl.

* ie: `curl -X GET -H "Content-Type: application/json" http://cosmoshub.strange.love:1317/cosmos/bank/v1beta1/balances/cosmos196ax4vc0lwpxndu9dyhvca7jhxp70rmcfhxsrt`

The list of Swagger endpoints are available is in `host:1317/swagger`. It can be disabled.

More info on [doc](https://docs.cosmos.network/v0.46/run-node/interact-node.html#using-the-rest-endpoints)

## How to sign a transaction?

Info and examples on [cosmos query blockchain](https://docs.cosmos.network/master/run-node/txs.html), [cosmos create transactions](https://github.com/cosmos/cosmos-sdk/blob/main/docs/run-node/txs.md)
Testnets & Faucets [cosmos doc](https://github.com/cosmos/testnets)

### Load an account

IMPORTANT: Never load this seeds into your wallets as this data is exposed publicly on github.com

Accounts used in from/to files were generated using using this [website](https://iancoleman.io/bip39)

* Mnemonic from `write sense wage direct salute north now dog divorce inflict pole provide spike welcome bring sister fetch upset chimney direct siren trash cruise mother` must generate `cosmos19kzdcmysekqu926fwdcjg5pdqlx3saujcldys5` for path `m/44'/118'/0'/0/0`(ATOM)
* Mnemonic to `sugar cereal decorate hip jelly choose milk cave rally liquid angry hat blood movie rare shadow skate drop giant insane argue shock mimic plate` must generate `cosmos1kc4zwgea50n6404untq05qsnlx9wayceknujcu` for path `m/44'/118'/0'/0/0`(ATOM)
  
### Create an account

Creating an acount:

* Generate a mnemonic(using bip39 package)
* Transform to a seed this mnemonic(using bip39)
* Create to a masterkey from the seed(using hd package)
* Derive the masterkey to an Atom private key with Account - 0 / External - 0(using hd package). Path: `"m/44'/118'/0'/0/0"`, more info on [BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#change)
* Generate Public Key and Address from Atom private key

### Make a Transaction

Transaction lyfecicle [cosmos doc](https://docs.cosmos.network/master/basics/tx-lifecycle.html)
Info on [cosmos sdk](https://docs.cosmos.network/master/core/transactions.html#transaction-generation)

Making a transaction:

* Create txConfig using: codec & SIGN_MODE_DIRECT
  * [NewProtoCodec](https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.46.0/codec#ProtoCodec)
  * [SIGN_MODE_DIRECT / Protobuf](https://docs.cosmos.network/master/core/transactions.html#sign-mode-direct-preferred)
* Create txBuilder from txConfig
* Set the message, Gas limit and other transaction parameters
* Sign the transaction(The API requires us to first perform a round of SetSignatures() with empty signatures, only to populate SignerInfos, and a second round of SetSignatures() to actually sign the correct payload):
  * Populate the SignerInfo
  * Sign the SignDoc (the payload to be signed)
* Generate transaction bytes
* Create the grpc connection
* Send the transaction bytes via grpc connection
* Wait for the result and check the destination address in a [bigdipper explorer](https://testnet.cosmos.bigdipper.live/) or [cosmostation explorer](https://cosmoshub-testnet.mintscan.io/cosmoshub-testnet)

// TODO: Change/Inject from command line parameters
// TODO: investigate how to get current network avg gas price
// TODO: investigate fee
// TODO: not works printing to json