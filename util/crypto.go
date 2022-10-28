package util

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"github.com/mr-tron/base58"
)

func S256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func StringToHex(str string) ([]byte, error) {

	if len(str)%2 == 1 {
		str = "0" + str
	}

	return hex.DecodeString(str)
}

func ReverseByte(input []byte) []byte {

	temp := len(input)

	var output []byte

	for i := temp - 1; i >= 0; i-- {
		output = append(output, input[i])
	}

	return output
}

func HexToBase58(str string) string {

	addb, _ := hex.DecodeString(str)
	hash1 := S256(S256(addb))
	secret := hash1[:4]
	for _, v := range secret {
		addb = append(addb, v)
	}
	return base58.Encode(addb)
}

func Base58ToHex(str string) string {

	temp, err := base58.Decode(str)
	if err != nil {
		panic(err)
	}

	temp = temp[:len(temp)-4]

	return hex.EncodeToString(temp)
}

func ZeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}
