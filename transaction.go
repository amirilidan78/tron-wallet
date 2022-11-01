package tronWallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"tronWallet/enums"
	"tronWallet/trongridClient"
	protoTron "tronWallet/trongridClient/proto"
	"tronWallet/util"
)

func createTransactionInput(network enums.Network, fromAddressHex string, toAddressHex string, amountInSun int64) (*api.TransactionExtention, error) {

	c := client.NewGrpcClient(string(enums.SHASTA_NODE))

	c.Start()

	return c.Transfer(fromAddressHex, toAddressHex, amountInSun)
}

func signTransaction(pb trongridClient.CreateTransactionResponse, privateKey *ecdsa.PrivateKey) (string, error) {

	rawData, err := json.Marshal(pb.RawData)
	if err != nil {
		return "", fmt.Errorf("proto marshal tx raw data error: %v", err)
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", fmt.Errorf("sign error: %v", err)
	}

	signatureString := hexutil.Encode(signature)

	return signatureString[2:], nil
}

func getRawTransaction(input trongridClient.CreateTransactionResponse, signed []byte) (string, error) {

	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	so := trongridClient.TransactionBody{
		RawData: trongridClient.TransactionBodyRawData{
			Contract: []trongridClient.TransactionBodyContract{
				{
					Parameter: trongridClient.TransactionBodyContractParameter{
						TypeUrl: "type.googleapis.com/protocol.TransferContract",
						Value: trongridClient.TransactionBodyContractParameterValue{
							Amount:       input.RawData.Contract[0].Parameter.Value.Amount,
							OwnerAddress: input.RawData.Contract[0].Parameter.Value.OwnerAddress,
							ToAddress:    input.RawData.Contract[0].Parameter.Value.ToAddress,
						},
					},
					Type: "TransferContract",
				},
			},
			Expiration:    input.RawData.Expiration,
			FeeLimit:      enums.TrxTransferFeeLimit,
			RefBlockBytes: input.RawData.RefBlockBytes,
			RefBlockHash:  input.RawData.RefBlockHash,
			Timestamp:     input.RawData.Timestamp,
		},
		Signature:  [][]byte{signed},
		TxID:       input.TxID,
		RawDataHex: hexutil.Encode(data)[2:],
	}

	jsonStr, err := json.Marshal(so)
	if err != nil {
		return "", err
	}

	return string(jsonStr), nil
}

func makeTransactionBlockHeader(network enums.Network) (*protoTron.TronBlockHeader, error) {

	nowBlockResponseBody, err := trongridClient.CurrentBlock(network)
	if err != nil {
		return nil, err
	}

	blockHeaderRaw := nowBlockResponseBody.BlockHeader.RawData

	txTrieRootHex, errTxTrieRootHex := util.StringToHex(blockHeaderRaw.TxTrieRoot)

	if errTxTrieRootHex != nil {
		return nil, errTxTrieRootHex
	}

	parentHash, errParentHash := util.StringToHex(blockHeaderRaw.ParentHash)

	if errParentHash != nil {
		return nil, errParentHash
	}

	witnessAddress, errWitnessAddress := util.StringToHex(blockHeaderRaw.WitnessAddress)

	if errWitnessAddress != nil {
		return nil, errWitnessAddress
	}

	blockHeader := &protoTron.TronBlockHeader{
		Timestamp:      blockHeaderRaw.Timestamp,
		TxTrieRoot:     txTrieRootHex,
		ParentHash:     parentHash,
		Number:         blockHeaderRaw.Number,
		WitnessAddress: witnessAddress,
		Version:        blockHeaderRaw.Version,
	}

	return blockHeader, nil
}
