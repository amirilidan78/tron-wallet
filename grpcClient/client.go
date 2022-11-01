package grpcClient

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tronWallet/enums"
)

func GetGrpcClient(node enums.Node) (*client.GrpcClient, error) {

	c := client.NewGrpcClient(string(node))

	err := c.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))

	return c, err
}
