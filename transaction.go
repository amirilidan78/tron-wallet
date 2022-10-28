package tronWallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"time"
	"tronWallet/enums"
	"tronWallet/trongridClient"
	protoTron "tronWallet/trongridClient/proto"
	"tronWallet/util"
)

func createTransactionInput(network enums.Network, fromAddressHex string, fromPrivateKey []byte, toAddressHex string, amountInSun int64, feeInSun int64) (*protoTron.TronSigningInput, error) {

	blockHeader, err := makeTransactionBlockHeader(network)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	timestamp := now.Unix() * 1000
	expirationTimeStamp := blockHeader.Timestamp + 60*60*1000

	transferContract := &protoTron.TronTransferContract{
		OwnerAddress: fromAddressHex,
		ToAddress:    toAddressHex,
		Amount:       amountInSun,
	}

	txContract := &protoTron.TronTransaction_Transfer{
		Transfer: transferContract,
	}

	tx := &protoTron.TronTransaction{
		Timestamp:     timestamp,
		Expiration:    expirationTimeStamp,
		BlockHeader:   blockHeader,
		FeeLimit:      feeInSun,
		ContractOneof: txContract,
	}

	return &protoTron.TronSigningInput{
		Transaction: tx,
		PrivateKey:  fromPrivateKey,
	}, nil

}

func signTransaction(pb proto.Message, privateKey *ecdsa.PrivateKey) ([]byte, string, error) {

	rawData, err := proto.Marshal(pb)
	if err != nil {
		return nil, "", fmt.Errorf("proto marshal tx raw data error: %v", err)
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, "", fmt.Errorf("sign error: %v", err)
	}

	return signature, hex.EncodeToString(hash), nil
}

func getRawTransaction(signed []byte) (string, error) {

	so := &protoTron.TronSigningOutput{}

	err := proto.Unmarshal(signed, so)
	if err != nil {
		return "", err
	}

	return so.GetJson(), nil
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
