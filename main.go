package tronWallet

import (
	"fmt"
	"tronWallet/enums"
)

func main() {

	c := &Crawler{
		Network: enums.NetworkSHASTA,
		Addresses: []string{
			"TY3PJu3VY8xVUc5BjYwJtyRgP7TfivV666",
		},
	}

	res, err := c.ScanBlocks(40)
	fmt.Println(res, err)
}
