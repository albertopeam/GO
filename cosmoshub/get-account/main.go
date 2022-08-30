package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	accounts "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcUrl := "rpc.sentry-01.theta-testnet.polypore.xyz:9090" // https://github.com/cosmos/testnets/tree/master/v7-theta/public-testnet
	fmt.Println("Testnet URL", grpcUrl)
	// GRPC conn https://pkg.go.dev/google.golang.org/grpc#Dial
	grpcConn, err := grpc.Dial(
		grpcUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		log.Fatalf("grpc.Dial Error %s", err)
	}
	defer grpcConn.Close()
	// query account https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/auth/types#QueryClient
	accountRequest := accounts.QueryAccountRequest{Address: "cosmos19kzdcmysekqu926fwdcjg5pdqlx3saujcldys5"}
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
	cdc := codec.NewProtoCodec(types.NewInterfaceRegistry()) // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/codec#NewProtoCodec
	var acc accounts.BaseAccount                             // returned concrete type inside Response.Account
	cdc.Unmarshal(accRes.Account.Value, &acc)                // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/codec#ProtoCodec.Unmarshal
	fmt.Println("Account.GetAdrress", acc.GetAddress().String())
	fmt.Println("Account.GetSequence", acc.GetSequence())
	fmt.Println("Account.GetAccountNumber", acc.GetAccountNumber())
}
