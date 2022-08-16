# Cosmoshub(Gaia) blockchain info

Gaia info on [github](https://github.com/cosmos/gaia)

Chain registry info on [github](https://github.com/cosmos/chain-registry)

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

Info and examples on [cosmos doc](https://docs.cosmos.network/master/run-node/txs.html)
Testnets & Faucets [cosmos doc](https://github.com/cosmos/testnets)

### Create an account
