package main

import (
	"fmt"
	tronWallet "github.com/Amirilidan78/tron-wallet"
	"github.com/Amirilidan78/tron-wallet/enums"
)

func main() {

	w := tronWallet.CreateTronWallet(enums.SHASTA_NODE, "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27")

	token := &tronWallet.Token{
		ContractAddress: enums.SHASTA_Tether_USDT,
	}

	fmt.Println(w.BalanceTRC20(token))
	fmt.Println(token.GetName(w.Node, w.AddressBase58))
	fmt.Println(token.GetSymbol(w.Node, w.AddressBase58))
	fmt.Println(token.GetDecimal(w.Node, w.AddressBase58))
}
