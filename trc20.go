package tronWallet

import (
	"fmt"
	"github.com/Amirilidan78/tron-wallet/enums"
	"github.com/Amirilidan78/tron-wallet/grpcClient"
	"github.com/Amirilidan78/tron-wallet/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

const (
	trc20NameSignature           = "0x06fdde03"
	trc20SymbolSignature         = "0x95d89b41"
	trc20DecimalsSignature       = "0x313ce567"
	trc20BalanceOf               = "0x70a08231"
	trc20TransferMethodSignature = "0xa9059cbb"
)

type Token struct {
	Wallet          *TronWallet
	ContractAddress enums.ContractAddress
}

func (t *Token) GetSymbol() (string, error) {

	c, err := grpcClient.GetGrpcClient(t.Wallet.Node)
	if err != nil {
		return "", err
	}

	result, err := c.TRC20Call(t.Wallet.Address, t.ContractAddress.Base58(), trc20SymbolSignature, true, 0)
	if err != nil {
		return "", err
	}

	data := util.ToHex(result.GetConstantResult()[0])

	return c.ParseTRC20StringProperty(data)
}

func (t *Token) GetName() (string, error) {

	c, err := grpcClient.GetGrpcClient(t.Wallet.Node)
	if err != nil {
		return "", err
	}

	result, err := c.TRC20Call(t.Wallet.Address, t.ContractAddress.Base58(), trc20NameSignature, true, 0)
	if err != nil {
		return "", err
	}

	data := util.ToHex(result.GetConstantResult()[0])

	return c.ParseTRC20StringProperty(data)
}

func (t *Token) GetDecimal() (*big.Int, error) {

	c, err := grpcClient.GetGrpcClient(t.Wallet.Node)
	if err != nil {
		return nil, err
	}

	result, err := c.TRC20Call(t.Wallet.Address, t.ContractAddress.Base58(), trc20DecimalsSignature, true, 0)
	if err != nil {
		return nil, err
	}

	data := util.ToHex(result.GetConstantResult()[0])

	return c.ParseTRC20NumericProperty(data)
}

func (t *Token) GetBalance() (*big.Int, error) {

	c, err := grpcClient.GetGrpcClient(t.Wallet.Node)
	if err != nil {
		return nil, err
	}

	req := trc20BalanceOf + "0000000000000000000000000000000000000000000000000000000000000000"[len(t.Wallet.Address)-2:] + t.Wallet.Address[2:]

	result, err := c.TRC20Call(t.Wallet.Address, t.ContractAddress.Base58(), req, true, 0)
	if err != nil {
		return nil, err
	}

	data := util.ToHex(result.GetConstantResult()[0])

	r, err := c.ParseTRC20NumericProperty(data)
	if err != nil {
		return nil, fmt.Errorf("contract address %s: %v", t.ContractAddress.Base58(), err)
	}
	if r == nil {
		return nil, fmt.Errorf("contract address %s: invalid balance of %s", t.ContractAddress.Base58(), t.Wallet.Address)
	}

	return r, nil
}

func (t *Token) Transfer(toAddress util.Address, amount *big.Int, feeLimit int64) (string, error) {

	c, err := grpcClient.GetGrpcClient(t.Wallet.Node)
	if err != nil {
		return "", err
	}

	ab := common.LeftPadBytes(amount.Bytes(), 32)

	req := trc20TransferMethodSignature + "0000000000000000000000000000000000000000000000000000000000000000"[len(toAddress.Hex())-4:] + toAddress.Hex()[4:]

	req += common.Bytes2Hex(ab)

	tx, err := c.TRC20Call(t.Wallet.Address, t.ContractAddress.Base58(), req, false, feeLimit)
	if err != nil {
		return "", err
	}

	privateKey, err := t.Wallet.PrivateKeyRCDSA()
	if err != nil {
		return "", err
	}

	signedTx, err := signTransaction(tx, privateKey)
	if err != nil {
		return "", err
	}

	err = broadcastTransaction(t.Wallet.Node, signedTx)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(tx.GetTxid())[2:], nil
}
