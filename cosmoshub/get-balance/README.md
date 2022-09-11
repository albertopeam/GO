# Cosmos get balance

## How to install?

* Clone repository
* run `go mod tidy` to fetch deps
* if fails try `go clean -modcache` and rerun `go mod tidy`

## How to run?

* `go run main.go`
* if you want to query you wallet address change `addr` property and run again

## How to find an URL for connect cosmos?

* review chain cosmos chain registry on [github](https://github.com/cosmos/chain-registry/blob/master/cosmoshub/chain.json)
* if needed an url for other chain then review registry index on [github](https://github.com/cosmos/chain-registry)
* Inside `peers` -> `seeds` or `persistent_seeds` use address for selecting the URL to connect
