package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Read State in mainnet
	queryMainnetState()
}

// full tutorial https://docs.cosmos.network/v0.46/run-node/interact-node.html
func queryMainnetState() {
	// create an addr. doc https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types
	addr, err := sdk.AccAddressFromBech32("cosmos196ax4vc0lwpxndu9dyhvca7jhxp70rmcfhxsrt")
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	fmt.Println("using", addr.String())

	// Create a connection to the gRPC server. doc https://pkg.go.dev/google.golang.org/grpc
	// more info on how to find a host and port on https://github.com/cosmos/chain-registry/blob/master/cosmoshub/chain.json
	// if grpc fails while dialing then run the query `curl -X GET "https://rpc.cosmos.network/net_info" -H "accept: application/json"` and use `remote_ip` instead
	grpcConn, err := grpc.Dial(
		"54.180.225.240:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	defer grpcConn.Close()
	fmt.Println("open gRPC connection", grpcConn.Target())

	// This creates a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)

	// query uatom balance for an account
	bankRes, err := bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: addr.String(), Denom: "uatom"})
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	fmt.Println("atom balance", bankRes.String())

	// query all balances for an account
	bankAllRes, err := bankClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: addr.String()})
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	fmt.Println("all balances", bankAllRes.String())
	fmt.Println("-----------------------------------")
}
