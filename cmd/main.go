package main

import (
	"fmt"
	"tronWallet"
	"tronWallet/enums"
)

func main() {

	network := enums.NetworkSHASTA
	privateKey := "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
	toAddress := "TVwt3HTg6PjP5bbb5x1GtSvTe1J5FYM2BT"
	amount := int64(1000000) // 1000000sun = 1trx

	w := tronWallet.CreateTronWallet(network, privateKey)

	tx, err := w.Transfer(toAddress, amount)

	fmt.Println(tx)
	fmt.Println(err)

}
