package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tronWallet/enums"
	"tronWallet/httpClient"
)

type BlockResponseBody struct {
	BlockID      string        `json:"blockID"`
	BlockHeader  BlockHeader   `json:"block_header"`
	Transactions []Transaction `json:"transactions"`
}

type BlockHeader struct {
	WitnessSignature string             `json:"witness_signature"`
	RawData          BlockHeaderRawData `json:"raw_data"`
}

type BlockHeaderRawData struct {
	Number         int64  `json:"number"`
	TxTrieRoot     string `json:"txTrieRoot"`
	WitnessAddress string `json:"witness_address"`
	ParentHash     string `json:"parentHash"`
	Version        int32  `json:"version"`
	Timestamp      int64  `json:"timestamp"`
}

type GetBlockById struct {
	Value   string `json:"value"`
	Visible bool   `json:"visible"`
}

func CurrentBlock(network enums.Network) (BlockResponseBody, error) {

	url := string(network) + "/wallet/getnowblock"

	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	responseBody := BlockResponseBody{}

	httpResponse, _, statusCode, err := httpClient.NewHttpClient().HttpPost(url, nil, header)

	if statusCode != http.StatusOK {
		return responseBody, errors.New(fmt.Sprintf("http status code is not 200 it is %d", statusCode))
	}

	if err != nil {
		return responseBody, err
	}

	err = json.Unmarshal(httpResponse, &responseBody)
	if err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

func GetBlock(network enums.Network, hex string) (BlockResponseBody, error) {

	url := string(network) + "/wallet/getblockbyid"

	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	requestBody := GetBlockById{
		Value:   hex,
		Visible: true,
	}

	responseBody := BlockResponseBody{}

	httpResponse, _, statusCode, err := httpClient.NewHttpClient().HttpPost(url, requestBody, header)

	if statusCode != http.StatusOK {
		return responseBody, errors.New(fmt.Sprintf("http status code is not 200 it is %d", statusCode))
	}

	if err != nil {
		return responseBody, err
	}

	err = json.Unmarshal(httpResponse, &responseBody)
	if err != nil {
		return responseBody, err
	}

	return responseBody, nil
}
