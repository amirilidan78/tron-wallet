package main

import (
	"fmt"
	"github.com/JFJun/trx-sign-go/grpcs"
	"github.com/JFJun/trx-sign-go/sign"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
	"tronWallet"
	"tronWallet/enums"
)

func main() {

	test()

}

func grpc_broadcast() {
	node := "grpc.shasta.trongrid.io:50051"
	privateKey := "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
	fromAddressBase58 := "TSw5FSuWhTAcaJmBUVFY9fUY4ihwx588b6"
	toAddressBase58 := "TVwt3HTg6PjP5bbb5x1GtSvTe1J5FYM2BT"
	amount := int64(1000000) // 1000000sun = 1trx

	c := new(grpcs.Client)
	c.GRPC = client.NewGrpcClient(node)
	_ = c.GRPC.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))

	tx, err := c.Transfer(fromAddressBase58, toAddressBase58, amount)
	if err != nil {
		fmt.Println(111)
		fmt.Println(err)
	}

	signTx, err := sign.SignTransaction(tx.Transaction, privateKey)
	if err != nil {
		fmt.Println(err)
	}

	err = c.BroadcastTransaction(signTx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(tx)
	fmt.Println(signTx)
}

func test() {
	node := enums.SHASTA_NODE
	privateKey := "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
	toAddressBase58 := "TVwt3HTg6PjP5bbb5x1GtSvTe1J5FYM2BT"
	amount := int64(1000000) // 1000000sun = 1trx
	//toAddress := util.Base58ToHex(toAddressBase58)

	w := tronWallet.CreateTronWallet(node, privateKey)

	_, err := w.Transfer(toAddressBase58, amount)
	fmt.Println(err)

	time.Sleep(time.Second * 3)

	c := &tronWallet.Crawler{
		Node:      node,
		Addresses: []string{w.AddressBase58},
	}

	res, err := c.ScanBlocks(3)

	fmt.Println(res)
	fmt.Println(err)

}
