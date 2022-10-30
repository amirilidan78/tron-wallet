package tronWallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"tronWallet/enums"
	"tronWallet/trongridClient"
	"tronWallet/util"
)

type TronWallet struct {
	Network       enums.Network
	Address       string
	AddressBase58 string
	PrivateKey    string
	PublicKey     string
}

// generating

func GenerateTronWallet(network enums.Network) *TronWallet {

	privateKey, _ := generatePrivateKey()
	privateKeyHex := convertPrivateKeyToHex(privateKey)

	publicKey, _ := getPublicKeyFromPrivateKey(privateKey)
	publicKeyHex := convertPublicKeyToHex(publicKey)

	address := getAddressFromPublicKey(publicKey)
	addressBase58 := util.HexToBase58(address)

	return &TronWallet{
		Network:       network,
		Address:       address,
		AddressBase58: addressBase58,
		PrivateKey:    privateKeyHex,
		PublicKey:     publicKeyHex,
	}
}

func CreateTronWallet(network enums.Network, privateKeyHex string) *TronWallet {

	privateKey, err := privateKeyFromHex(privateKeyHex)
	if err != nil {
		panic(err)
	}

	publicKey, _ := getPublicKeyFromPrivateKey(privateKey)
	publicKeyHex := convertPublicKeyToHex(publicKey)

	address := getAddressFromPublicKey(publicKey)
	addressBase58 := util.HexToBase58(address)

	return &TronWallet{
		Network:       network,
		Address:       address,
		AddressBase58: addressBase58,
		PrivateKey:    privateKeyHex,
		PublicKey:     publicKeyHex,
	}
}

// struct functions

func (t *TronWallet) PrivateKeyRCDSA() (*ecdsa.PrivateKey, error) {
	return privateKeyFromHex(t.PrivateKey)
}

func (t *TronWallet) PrivateKeyBytes() ([]byte, error) {

	priv, err := t.PrivateKeyRCDSA()
	if err != nil {
		return []byte{}, err
	}

	return crypto.FromECDSA(priv), nil
}

// private key

func generatePrivateKey() (*ecdsa.PrivateKey, error) {

	return crypto.GenerateKey()
}

func convertPrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {

	privateKeyBytes := crypto.FromECDSA(privateKey)

	return hexutil.Encode(privateKeyBytes)[2:]
}

func privateKeyFromHex(hex string) (*ecdsa.PrivateKey, error) {

	return crypto.HexToECDSA(hex)
}

// public key

func getPublicKeyFromPrivateKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error in getting public key")
	}

	return publicKeyECDSA, nil
}

func convertPublicKeyToHex(publicKey *ecdsa.PublicKey) string {

	privateKeyBytes := crypto.FromECDSAPub(publicKey)

	return hexutil.Encode(privateKeyBytes)[2:]
}

// address

func getAddressFromPublicKey(publicKey *ecdsa.PublicKey) string {

	address := crypto.PubkeyToAddress(*publicKey).Hex()

	address = "41" + address[2:]

	return address
}

// balance

func (t *TronWallet) Balance() (trongridClient.GetAccountResponseBody, error) {
	return trongridClient.GetAddressBalance(t.Network, t.AddressBase58)
}

// transaction

func (t *TronWallet) Transfer(toAddressBase58 string, amountInSun int64) (string, error) {

	toAddress := util.Base58ToHex(toAddressBase58)

	privateRCDSA, err := t.PrivateKeyRCDSA()
	if err != nil {
		return "", fmt.Errorf("RCDSA private key error: %v", err)
	}

	pb, err := trongridClient.CreateTransaction(t.Network, t.Address, toAddress, amountInSun)
	if err != nil {
		return "", fmt.Errorf("creating tx pb error: %v", err)
	}

	signed, txId, err := signTransaction(pb, privateRCDSA)
	if err != nil {
		return "", fmt.Errorf("signing transaction error: %v", err)
	}

	data := make(map[string]interface{})
	data["raw_data"] = pb.RawData
	data["txID"] = txId
	data["signature"] = []string{
		hexutil.Encode(signed),
	}
	data["visible"] = true

	res, err := trongridClient.BroadcastTransaction(t.Network, data)
	if err != nil {
		return "", fmt.Errorf("broadcast transaction error: %v", err)
	}

	return res.TxID, err
}
