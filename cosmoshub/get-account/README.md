# Get Cosmos Account

## Dependencies

go.mod contains replacement for protobuf, more info [here](https://docs.cosmos.network/master/run-node/interact-node.html#programmatically-via-go)

## Getting cosmos account

* The account to get the info is hardcoded
* Create grpc connection
* Use grpc to query accounts
* Unpack protobuf response into an `cosmos.auth.v1beta1.BaseAccount`, more info [here](https://docs.cosmos.network/v0.46/core/encoding.html#interface-encoding-and-usage-of-any). Usage of Protobuf in cosmos, [ADR 019](https://docs.cosmos.network/master/architecture/adr-019-protobuf-state-encoding.html)
* Print account
