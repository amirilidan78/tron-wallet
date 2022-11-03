package grpcClient

import (
	"bytes"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"tronWallet/grpcClient/proto/core"
)

func (g *GrpcClient) GetAccount(addr string) (*core.Account, error) {
	account := new(core.Account)
	var err error

	account.Address, err = common.DecodeCheck(addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := g.getContext()
	defer cancel()

	acc, err := g.Client.GetAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(acc.Address, account.Address) {
		return nil, fmt.Errorf("account not found")
	}
	return acc, nil
}
