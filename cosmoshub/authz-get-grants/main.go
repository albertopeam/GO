package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create GRPC connection
	mainNet := "54.180.225.240:9090" // cosmos mainnet
	grpcConn, err := grpc.Dial(
		mainNet,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		log.Fatalf("grpcConn Error %s", err)
	}
	defer grpcConn.Close()

	// Create the query client
	authzClient := authz.NewQueryClient(grpcConn)

	// Create the query
	address := "cosmos196ax4vc0lwpxndu9dyhvca7jhxp70rmcfhxsrt"
	request := authz.QueryGranterGrantsRequest{Granter: address}
	fmt.Println("QueryGrantsRequest", request)

	// Query grants
	response, err := authzClient.GranterGrants(context.Background(), &request)
	if err != nil {
		log.Fatalf("GranterGrants Error %s", err)
	}

	// Print grants
	if len(response.Grants) == 0 {
		fmt.Println("No grants for", address)
	} else {
		for _, grant := range response.Grants {
			fmt.Println("Grant", grant.Expiration, grant.Authorization.GetTypeUrl(), grant.Granter, grant.Grantee)
		}
	}
}
