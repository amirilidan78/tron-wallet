package tronWallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"tronWallet/api"
	"tronWallet/enums"
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
	addressBase58 := convertAddressToBase58(address)

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
	addressBase58 := convertAddressToBase58(address)

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

func convertAddressToBase58(address string) string {

	addb, _ := hex.DecodeString(address)
	hash1 := util.S256(util.S256(addb))
	secret := hash1[:4]
	for _, v := range secret {
		addb = append(addb, v)
	}

	return base58.Encode(addb)
}

// balance

func (t *TronWallet) Balance() (api.GetAccountResponseBody, error) {
	return api.GetAddressBalance(t.Network, t.AddressBase58)
}
