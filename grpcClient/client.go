package grpcClient

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"tronWallet/enums"
)

func GetGrpcClient(node enums.Node) {

	c := client.NewGrpcClient(string(node))

	c.Start()

}
