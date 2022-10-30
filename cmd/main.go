package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"tronWallet"
	"tronWallet/enums"
	"tronWallet/trongridClient"
	"tronWallet/util"
)

func main() {

	network := enums.NetworkSHASTA
	privateKey := "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
	toAddressBase58 := "TVwt3HTg6PjP5bbb5x1GtSvTe1J5FYM2BT"
	amount := int64(1000000) // 1000000sun = 1trx
	toAddress := util.Base58ToHex(toAddressBase58)

	w := tronWallet.CreateTronWallet(network, privateKey)

	tx, err := trongridClient.CreateTransaction(network, w.Address, toAddress, amount)

	rawData, _ := json.Marshal(tx)
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	privateKeyECDSA, _ := w.PrivateKeyRCDSA()
	signature, err := crypto.Sign(hash, privateKeyECDSA)

	signatureHash := hex.EncodeToString(signature)

	tt := txSSS{
		Signature: []string{signatureHash},
		TxID:      tx.TxID,
		RawData:   tx.RawData,
	}

	dd, _ := json.Marshal(tt)

	txHash := hex.EncodeToString(dd)

	res, err := trongridClient.BroadcastTransaction(network, txHash)

	fmt.Println(res)
	fmt.Println(err)

}

type txSSS struct {
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
	RawData   struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Amount       int    `json:"amount"`
					OwnerAddress string `json:"owner_address"`
					ToAddress    string `json:"to_address"`
				} `json:"value"`
				TypeUrl string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    int64  `json:"expiration"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"raw_data"`
}
