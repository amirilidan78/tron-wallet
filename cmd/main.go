package main

import (
	"tronWallet"
	"tronWallet/enums"
)

func main() {

	network := enums.NetworkSHASTA
	privateKey := "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
	toAddressBase58 := "TVwt3HTg6PjP5bbb5x1GtSvTe1J5FYM2BT"
	amount := int64(1000000) // 1000000sun = 1trx
	//toAddress := util.Base58ToHex(toAddressBase58)

	w := tronWallet.CreateTronWallet(network, privateKey)

	_, _ = w.Transfer(toAddressBase58, amount)

}
