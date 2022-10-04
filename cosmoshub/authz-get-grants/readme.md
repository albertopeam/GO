# Stake auth

## ADR-30 authz

[ADR-30 authz](https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-030-authz-module.md#adr-030-authorization-module)

## Github Authz doc

[main doc](https://github.com/cosmos/cosmos-sdk/blob/v0.46.1/x/authz/spec/README.md)

[core concepts](https://github.com/cosmos/cosmos-sdk/blob/v0.46.1/x/authz/spec/01_concepts.md)

## Cosmos Authz doc

[Authz package](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz)

[Stake authorization type](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/staking/types#StakeAuthorization)

## Steps to grant to the granter a stake redelegation authorization

1. Create [grpc connection](https://pkg.go.dev/google.golang.org/grpc#ClientConn)

2. Create the [query client](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#NewQueryClient) for the module. [Query client interface](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#QueryClient)

3. Create the [QueryGranterGrantsRequest](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#QueryGranterGrantsRequest) to obtain a list of GrantAuthorization

4. Invoke `GranterGrants` for [QueryClient](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#QueryClient)

5. Print [Grants](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#Grant): Auth+Expiration if success. Cosmos SDK Built-in Authorizations: [`GenericAuthorization`](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#GenericAuthorization), `SendAuthorization` and [`StakeAuthorization`](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/staking/types#StakeAuthorization)