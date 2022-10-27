package util

import (
	"crypto/sha256"
	"encoding/hex"
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
