package tronWallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"github.com/Amirilidan78/tron-wallet/enums"
	"github.com/Amirilidan78/tron-wallet/grpcClient"
	"github.com/Amirilidan78/tron-wallet/grpcClient/proto/api"
	"github.com/ethereum/go-ethereum/crypto"
)

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

func createTransactionInput(node enums.Node, fromAddressBase58 string, toAddressBase58 string, amountInSun int64) (*api.TransactionExtention, error) {

	c, err := grpcClient.GetGrpcClient(node)
	if err != nil {
		return nil, err
	}

	return c.Transfer(fromAddressBase58, toAddressBase58, amountInSun)
}

func signTransaction(transaction *api.TransactionExtention, privateKey *ecdsa.PrivateKey) (*api.TransactionExtention, error) {

	rawData, err := proto.Marshal(transaction.Transaction.GetRawData())
	if err != nil {
		return transaction, fmt.Errorf("proto marshal tx raw data error: %v", err)
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return transaction, fmt.Errorf("sign error: %v", err)
	}

	transaction.Transaction.Signature = append(transaction.Transaction.Signature, signature)
	return transaction, nil
}

func broadcastTransaction(node enums.Node, transaction *api.TransactionExtention) error {

	c, err := grpcClient.GetGrpcClient(node)
	if err != nil {
		return err
	}

	res, err := c.Broadcast(transaction.Transaction)
	if err != nil {
		return err
	}

	if res.Result != true {
		return errors.New(res.Code.String())
	}

	return nil
}
