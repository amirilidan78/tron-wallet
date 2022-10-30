package trongridClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tronWallet/enums"
	"tronWallet/httpClient"
)

type Transaction struct {
	Ret        []TransactionRet   `json:"ret"`
	Signature  []string           `json:"signature"`
	TxID       string             `json:"txID"`
	RawData    TransactionRawData `json:"raw_data"`
	RawDataHex string             `json:"raw_data_hex"`
}

type TransactionRet struct {
	ContractRet string `json:"contractRet"`
}

type TransactionRawData struct {
	Contract      []TransactionContract `json:"contract"`
	RefBlockBytes string                `json:"ref_block_bytes"`
	RefBlockHash  string                `json:"ref_block_hash"`
	Expiration    uint64                `json:"expiration"`
	FeeLimit      uint64                `json:"fee_limit"`
	Timestamp     uint64                `json:"timestamp"`
}

type TransactionContract struct {
	Parameter TransactionContractParameter `json:"parameter"`
	Type      string                       `json:"type"`
}

type TransactionContractParameter struct {
	Value   TransactionContractParameterValue `json:"value"`
	TypeUrl string                            `json:"type_url"`
}

type TransactionContractParameterValue struct {
	Amount          uint64 `json:"amount"`
	OwnerAddress    string `json:"owner_address"`
	ToAddress       string `json:"to_address"`
	Data            string `json:"data"`
	ContractAddress string `json:"contract_address"`
}

type GetTransactionById struct {
	Value   string `json:"value"`
	Visible bool   `json:"visible"`
}

type CreateTransactionRequest struct {
	OwnerAddress string `json:"owner_address"`
	ToAddress    string `json:"to_address"`
	Amount       int64  `json:"amount"`
}

type CreateTransactionResponse struct {
	TxID    string `json:"txID"`
	RawData struct {
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

type BroadcastTransactionRequest struct {
	Transaction string `json:"transaction"`
}

func BroadcastTransaction(network enums.Network, hex string) (Transaction, error) {

	url := string(network) + "/wallet/broadcasttransaction"

	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	requestBody := BroadcastTransactionRequest{
		Transaction: hex,
	}

	responseBody := Transaction{}

	httpResponse, _, statusCode, err := httpClient.NewHttpClient().HttpPost(url, requestBody, header)

	fmt.Println(string(httpResponse))

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

	fmt.Println(string(httpResponse))

	return responseBody, nil
}

func GetTransaction(network enums.Network, txHash string) (Transaction, error) {

	url := string(network) + "/wallet/gettransactionbyid"

	requestBody := GetTransactionById{
		Value:   txHash,
		Visible: true,
	}

	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	responseBody := Transaction{}

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

func CreateTransaction(network enums.Network, fromAddressHex string, toAddressHex string, amount int64) (CreateTransactionResponse, error) {

	url := string(network) + "/wallet/createtransaction"

	requestBody := CreateTransactionRequest{
		OwnerAddress: fromAddressHex,
		ToAddress:    toAddressHex,
		Amount:       amount,
	}

	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	responseBody := CreateTransactionResponse{}

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

//func BroadcastTransaction(network enums.Network, fromAddressBase58 string, fromPrivateKey []byte, toAddressBase58 string, amountInSun int64, feeInSun int64) (map[string]interface{}, error) {
//
//	url := string(network) + "/wallet/broadcasttransaction"
//
//	input, err := createTransactionInput(network, fromAddressBase58, fromPrivateKey, toAddressBase58, amountInSun, feeInSun)
//	if err != nil {
//		return nil, err
//	}
//
//	var requestBody map[string]interface{}
//	err := json.Unmarshal(tx, &requestBody)
//	if err != nil {
//		return nil, err
//	}
//
//	header := map[string]string{
//		"Content-Type": "application/json",
//		"Accept":       "application/json",
//	}
//
//	var responseBody map[string]interface{}
//
//	httpResponse, _, statusCode, err := httpClient.NewHttpClient().HttpPost(url, requestBody, header)
//
//	if statusCode != http.StatusOK {
//		return responseBody, errors.New(fmt.Sprintf("http status code is not 200 it is %d", statusCode))
//	}
//
//	if err != nil {
//		return responseBody, err
//	}
//
//	err = json.Unmarshal(httpResponse, &responseBody)
//	if err != nil {
//		return responseBody, err
//	}
//
//	return responseBody, nil
//}
//
