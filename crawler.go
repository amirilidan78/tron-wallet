package tronWallet

import (
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"strings"
	"sync"
	"time"
	"tronWallet/enums"
	"tronWallet/trongridClient"
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

	var wg sync.WaitGroup

	var allTransactions [][]CrawlTransaction

	block, err := trongridClient.CurrentBlock(c.Network)
	if err != nil {
		return nil, err
	}

	// check block for transaction
	allTransactions = append(allTransactions, c.extractOurTransactionsFromBlock(block))
	if err != nil {
		return nil, err
	}

	blockNumber := int32(block.BlockHeader.RawData.Number)

	for i := count; i > 0; i-- {

		wg.Add(1)

		blockNumber = blockNumber - 1

		// sleep to avoid 503 error
		time.Sleep(100 * time.Millisecond)
		go c.getBlockData(&wg, &allTransactions, blockNumber)
	}

	wg.Wait()

	return c.prepareCrawlResultFromTransactions(allTransactions), nil
}

func (c *Crawler) getBlockData(wg *sync.WaitGroup, allTransactions *[][]CrawlTransaction, num int32) {

	defer wg.Done()

	block, err := trongridClient.GetBlockByNumber(c.Network, num)
	if err != nil {
		fmt.Println(err)
		return
	}

	// check block for transaction
	*allTransactions = append(*allTransactions, c.extractOurTransactionsFromBlock(block))
}

func (c *Crawler) extractOurTransactionsFromBlock(block trongridClient.BlockResponseBody) []CrawlTransaction {

	var txs []CrawlTransaction

	for _, transaction := range block.Transactions {

		// if transaction is not success
		if transaction.Ret[0].ContractRet != enums.SuccessTransactionRetStatus {
			continue
		}

		// if transaction is not tron transfer
		if transaction.RawData.Contract[0].Type != enums.TransferContract {
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
				txs = append(txs, CrawlTransaction{
					TxId:        transaction.TxID,
					FromAddress: fromAddress,
					ToAddress:   toAddress,
					Amount:      transaction.RawData.Contract[0].Parameter.Value.Amount,
				})
			}
		}
	}

	return txs
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
