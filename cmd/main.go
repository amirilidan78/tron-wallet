package main

import (
	"fmt"
	tronWallet "github.com/Amirilidan78/tron-wallet"
	"github.com/Amirilidan78/tron-wallet/enums"
)

func main() {

	w, _ := tronWallet.CreateTronWallet(enums.SHASTA_NODE, "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27")

	c := &tronWallet.Crawler{
		Node: enums.SHASTA_NODE,
		Addresses: []string{
			w.AddressBase58,
		},
	}

	fmt.Println(c.ScanBlocksFromTo(28905305, 28905307))
}
