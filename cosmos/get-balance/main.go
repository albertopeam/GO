package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// full tutorial https://docs.cosmos.network/v0.46/run-node/interact-node.html
func main() {
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
