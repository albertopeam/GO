# Stake auth

## ADR-30 authz

[ADR-30 authz](https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-030-authz-module.md#adr-030-authorization-module)

## Github Authz doc

[main doc](https://github.com/cosmos/cosmos-sdk/blob/v0.46.1/x/authz/spec/README.md)

[core concepts](https://github.com/cosmos/cosmos-sdk/blob/v0.46.1/x/authz/spec/01_concepts.md)

## Cosmos Authz doc

[Authz package](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz)

[Stake authorization type](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/staking/types#StakeAuthorization)

## Cosmos testnet info

[Cosmos public testnet](https://github.com/cosmos/testnets/tree/master/public)

## Steps to check granter grants

1. Create [grpc connection](https://pkg.go.dev/google.golang.org/grpc#ClientConn)

2. Create the [query client](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#NewQueryClient) for the module. [QueryClient interface](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#QueryClient)

3. Create the [QueryGranterGrantsRequest](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#QueryGranterGrantsRequest) to obtain a list of GrantAuthorization

4. Invoke `GranterGrants` for [QueryClient](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#QueryClient)

5. Print [Grants](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#Grant): Auth+Expiration if success. Cosmos SDK Built-in Authorizations: [`GenericAuthorization`](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#GenericAuthorization), `SendAuthorization` and [`StakeAuthorization`](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/staking/types#StakeAuthorization)

## Steps to grant to the grantee a stake redelegation authorization from the granter

Before start we need to find a validator to be the grantee. Check the testnet explorers to do that:

* [Testnet github doc](https://github.com/cosmos/testnets/tree/master/public) 

* [Testnet explorer big dipper](https://explorer.theta-testnet.polypore.xyz/) | [Testnet explorer mintscan](https://cosmoshub-testnet.mintscan.io/cosmoshub-testnet) | [Testnet explorer big dipper live](https://testnet.cosmos.bigdipper.live/)

* [Testnet validators big dipper](https://explorer.theta-testnet.polypore.xyz/validators). 

    We pick [stakely.io](https://explorer.theta-testnet.polypore.xyz/validator/cosmosvaloper1c28cfmvvne62n5347h3nptar7ka0dffxam8ct3) validator, so `grantee` will be `cosmos1c28cfmvvne62n5347h3nptar7ka0dffxc0nd8z`(USE AUTODELEGATION ADDRESS)

Grant steps:

1. Create [grpc connection](https://pkg.go.dev/google.golang.org/grpc#ClientConn)

2. Create the [MsgClient](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#NewMsgClient) to grant/revoke authorizations. [MsgClient interface](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#MsgClient)

3. Create the [Grant](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#NewGrant) and [MsgGrant](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#MsgGrant) data

4. Run the [Grant via MsgClient](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#MsgClient)

5. Print [MsgGrantResponse](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#MsgGrantResponse)