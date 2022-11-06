package main

import (
	"fmt"
	tronWallet "github.com/Amirilidan78/tron-wallet"
	"github.com/Amirilidan78/tron-wallet/enums"
	"github.com/Amirilidan78/tron-wallet/util"
	"math/big"
)

func main() {

	toAddress, _ := util.Base58ToAddress("TCTfPcF1wTRDK1qFAb7w2zGjF11WE9v8DA")

	w := tronWallet.CreateTronWallet(enums.SHASTA_NODE, "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27")
	t := tronWallet.Token{
		Wallet:          w,
		ContractAddress: enums.SHASTA_Tether_USDT,
	}

	txId, _ := t.Transfer(toAddress, big.NewInt(1000000), 10000000)

	fmt.Println(txId)
}
