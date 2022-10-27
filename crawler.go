package tronWallet

import (
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"strconv"
	"strings"
	"tronWallet/api"
	"tronWallet/enums"
)

type Crawler struct {
	Network   enums.Network
	Addresses []string
}

type CrawlResult struct {
	Address      string
	Transactions []CrawlTransaction
}

type CrawlTransaction struct {
	TxId        string
	FromAddress string
	ToAddress   string
	Amount      uint64
}

func (c *Crawler) ScanBlocks(count int) ([]CrawlResult, error) {

	var allTransactions [][]CrawlTransaction

	block, err := api.CurrentBlock(c.Network)
	if err != nil {
		return nil, err
	}

	// check block for transaction
	txs, err := c.getBlockTransactions(block)
	allTransactions = append(allTransactions, txs)
	if err != nil {
		return nil, err
	}

	for i := count; i > 0; i-- {

		block, err = api.GetBlock(c.Network, block.BlockHeader.RawData.ParentHash)
		if err != nil {
			return nil, err
		}

		fmt.Println("Scanning block number " + strconv.FormatInt(block.BlockHeader.RawData.Number, 10))

		// check block for transaction
		txs, err = c.getBlockTransactions(block)
		allTransactions = append(allTransactions, txs)
		if err != nil {
			return nil, err
		}
	}

	return c.prepareCrawlResultFromTransactions(allTransactions), nil
}

func (c *Crawler) getBlockTransactions(block api.BlockResponseBody) ([]CrawlTransaction, error) {

	var txs []CrawlTransaction

	for _, transaction := range block.Transactions {

		// if transaction is not success
		if transaction.Ret[0].ContractRet != enums.SuccessTransactionRetStatus {
			fmt.Println(transaction.TxID + " transaction is not success")
			continue
		}

		// if transaction is not tron transfer
		if transaction.RawData.Contract[0].Type != enums.TransferContract {
			fmt.Println(transaction.TxID + " transaction is token transfer")
			continue
		}

		// if address is hex convert to base58
		toAddress := transaction.RawData.Contract[0].Parameter.Value.ToAddress
		if strings.HasPrefix(toAddress, "41") == true {
			toAddress = address.HexToAddress(toAddress).String()
		}

		// if address is hex convert to base58
		fromAddress := transaction.RawData.Contract[0].Parameter.Value.OwnerAddress
		if strings.HasPrefix(fromAddress, "41") == true {
			fromAddress = address.HexToAddress(fromAddress).String()
		}

		for _, ourAddress := range c.Addresses {
			if ourAddress == toAddress {
				fmt.Println(toAddress + " address is our address")
				txs = append(txs, CrawlTransaction{
					TxId:        transaction.TxID,
					FromAddress: fromAddress,
					ToAddress:   toAddress,
					Amount:      transaction.RawData.Contract[0].Parameter.Value.Amount,
				})
			} else {
				fmt.Println(toAddress + " address is not our address")
			}
		}
	}

	return txs, nil
}

func (c *Crawler) prepareCrawlResultFromTransactions(transactions [][]CrawlTransaction) []CrawlResult {

	var result []CrawlResult

	for _, transaction := range transactions {
		for _, tx := range transaction {

			if c.addressExistInResult(result, tx.ToAddress) {
				id, res := c.getAddressCrawlInResultList(result, tx.ToAddress)
				res.Transactions = append(res.Transactions, tx)
				result[id] = res

			} else {
				result = append(result, CrawlResult{
					Address:      tx.ToAddress,
					Transactions: []CrawlTransaction{tx},
				})
			}
		}
	}

	return result
}

func (c *Crawler) addressExistInResult(result []CrawlResult, address string) bool {
	for _, res := range result {
		if res.Address == address {
			return true
		}
	}
	return false
}

func (c *Crawler) getAddressCrawlInResultList(result []CrawlResult, address string) (int, CrawlResult) {
	for id, res := range result {
		if res.Address == address {
			return id, res
		}
	}
	panic("crawl result not found")
}
