package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create GRPC connection
	mainNet := "rpc.sentry-01.theta-testnet.polypore.xyz:9090" // testnet // "54.180.225.240:9090" cosmos mainnet
	grpcConn, err := grpc.Dial(
		mainNet,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		log.Fatalf("grpcConn Error %s", err)
	}
	defer grpcConn.Close()

	// addresses
	granter := "cosmos19kzdcmysekqu926fwdcjg5pdqlx3saujcldys5"
	grantee := "cosmos1c28cfmvvne62n5347h3nptar7ka0dffxc0nd8z"
	validator := "cosmosvaloper1c28cfmvvne62n5347h3nptar7ka0dffxam8ct3"

	// Check granted Authorizations
	checkGrants(grpcConn, granter)

	// Enable grant(if already granted it will be overwritten)
	grant(grpcConn, granter, grantee, validator)

	// Check granted Authorizations(again)
	checkGrants(grpcConn, granter)
}

func checkGrants(grpcConn *grpc.ClientConn, address string) {
	// Create the query client
	queryClient := authz.NewQueryClient(grpcConn)

	// Create the query
	request := authz.QueryGranterGrantsRequest{Granter: address}
	fmt.Println("QueryGrantsRequest", request)

	// Query grants
	response, err := queryClient.GranterGrants(context.Background(), &request)
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

func grant(grpcConn *grpc.ClientConn, granter string, grantee string, validator string) {
	// Create the MsgClient
	msgClient := authz.NewMsgClient(grpcConn)

	// Create the MsgGrant
	var tokens *types.Coin // empty means no limit https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types#Coin
	valAddress, err := types.ValAddressFromBech32(validator)
	if err != nil {
		log.Fatalf("ValAddressFromBech32 Error %s", err)
	}
	allowed := []types.ValAddress{valAddress}
	denied := []types.ValAddress{}
	stakeAuth, err := stakingtypes.NewStakeAuthorization(allowed, denied, stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_REDELEGATE, tokens) // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/staking/types#NewStakeAuthorization https://github.com/cosmos/cosmos-sdk/blob/v0.46.0-rc1/x/staking/types/authz.go#L15-L35
	if err != nil {
		log.Fatalf("stakingtypes.NewStakeAuthorization Error %s", err)
	}
	granterAddr, err := types.AccAddressFromBech32(granter)
	if err != nil {
		log.Fatalf("GranterAddr Error %s", err)
	}
	granteeAddr, err := types.AccAddressFromBech32(grantee)
	if err != nil {
		log.Fatalf("GranteeAddr Error %s", err)
	}
	expiration := time.Now().Add(time.Hour * 24 * 30)
	msgGrant, err := authz.NewMsgGrant(granterAddr, granteeAddr, stakeAuth, &expiration) // https://pkg.go.dev/github.com/cosmos/cosmos-sdk/x/authz#NewMsgGrant https://github.com/cosmos/cosmos-sdk/blob/v0.46.1/x/authz/spec/03_messages.md#msggrant
	if err != nil {
		log.Fatalf("NewMsgGrant Error %s", err)
	}

	// Print data
	fmt.Println("Tokens", tokens)
	fmt.Println("Expiration", expiration)
	fmt.Println("Allowed addresses", valAddress.String())
	fmt.Println("Denied addresses", denied)
	fmt.Println("Granter address", granterAddr.String())
	fmt.Println("Grantee address", granteeAddr.String())

	// Run Grant
	//TODO: review discord https://discord.com/channels/669268347736686612/1019978171367559208
	//TODO: is needed to have stake before(I don't think, only a grant).
	//TODO: val address is needed? or is needed the default address?
	//TODO: is not impl in testnet?
	response, err := msgClient.Grant(context.Background(), msgGrant, grpc.FailFastCallOption{FailFast: false})
	if err != nil {
		log.Fatalf("MsgClient.Grant Error %s", err)
	}

	// Print the response(empty struct if success)
	fmt.Println("Grant success", response)
}
