package api

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
//func createTransactionInput(network enums.Network, fromAddressBase58 string, fromPrivateKey []byte, toAddressBase58 string, amountInSun int64, feeInSun int64) (*proto.TronSigningInput, error) {
//
//	blockHeader, err := makeTransactionBlockHeader(network)
//	if err != nil {
//		return nil, err
//	}
//
//	now := time.Now()
//	timestamp := now.Unix() * 1000
//	expirationTimeStamp := blockHeader.Timestamp + 60*60*1000
//
//	transferContract := &proto.TronTransferContract{
//		OwnerAddress: fromAddressBase58,
//		ToAddress:    toAddressBase58,
//		Amount:       amountInSun,
//	}
//
//	txContract := &proto.TronTransaction_Transfer{
//		Transfer: transferContract,
//	}
//
//	tx := &proto.TronTransaction{
//		Timestamp:     timestamp,
//		Expiration:    expirationTimeStamp,
//		BlockHeader:   blockHeader,
//		FeeLimit:      feeInSun,
//		ContractOneof: txContract,
//	}
//
//	return &proto.TronSigningInput{
//		Transaction: tx,
//		PrivateKey:  fromPrivateKey,
//	}, nil
//
//}
//
//func signTransaction(input *proto.TronSigningInput) {
//
//	goProto.Marshal(input)
//}
//
//func makeTransactionBlockHeader(network enums.Network) (*proto.TronBlockHeader, error) {
//
//	nowBlockResponseBody, err := CurrentBlock(network)
//	if err != nil {
//		return nil, err
//	}
//
//	blockHeaderRaw := nowBlockResponseBody.BlockHeader.RawData
//
//	txTrieRootHex, errTxTrieRootHex := util.StringToHex(blockHeaderRaw.TxTrieRoot)
//
//	if errTxTrieRootHex != nil {
//		return nil, errTxTrieRootHex
//	}
//
//	parentHash, errParentHash := util.StringToHex(blockHeaderRaw.ParentHash)
//
//	if errParentHash != nil {
//		return nil, errParentHash
//	}
//
//	witnessAddress, errWitnessAddress := util.StringToHex(blockHeaderRaw.WitnessAddress)
//
//	if errWitnessAddress != nil {
//		return nil, errWitnessAddress
//	}
//
//	blockHeader := &proto.TronBlockHeader{
//		Timestamp:      blockHeaderRaw.Timestamp,
//		TxTrieRoot:     txTrieRootHex,
//		ParentHash:     parentHash,
//		Number:         blockHeaderRaw.Number,
//		WitnessAddress: witnessAddress,
//		Version:        blockHeaderRaw.Version,
//	}
//
//	return blockHeader, nil
//}
