package tronWallet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
	"strings"
	"sync"
	"time"
	"tronWallet/enums"
	"tronWallet/grpcClient"
)

type Crawler struct {
	Node      enums.Node
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

	client, err := grpcClient.GetGrpcClient(c.Node)
	if err != nil {
		return nil, err
	}

	block, err := client.GetNowBlock()
	if err != nil {
		return nil, err
	}

	// check block for transaction
	allTransactions = append(allTransactions, c.extractOurTransactionsFromBlock(block))
	if err != nil {
		return nil, err
	}

	blockNumber := block.BlockHeader.RawData.Number

	for i := count; i > 0; i-- {
		wg.Add(1)
		blockNumber = blockNumber - 1
		// sleep to avoid 503 error
		time.Sleep(100 * time.Millisecond)
		go c.getBlockData(&wg, client, &allTransactions, blockNumber)
	}

	wg.Wait()

	return c.prepareCrawlResultFromTransactions(allTransactions), nil
}

func (c *Crawler) getBlockData(wg *sync.WaitGroup, client *client.GrpcClient, allTransactions *[][]CrawlTransaction, num int64) {

	defer wg.Done()

	block, err := client.GetBlockByNum(num)
	if err != nil {
		fmt.Println(err)
		return
	}

	// check block for transaction
	*allTransactions = append(*allTransactions, c.extractOurTransactionsFromBlock(block))
}

func (c *Crawler) extractOurTransactionsFromBlock(block *api.BlockExtention) []CrawlTransaction {

	var txs []CrawlTransaction

	for _, t := range block.Transactions {

		transaction := t.Transaction

		// if transaction is not success
		if transaction.Ret[0].ContractRet != core.Transaction_Result_SUCCESS {
			fmt.Println("transaction is not success")
			continue
		}

		// if transaction is not tron transfer
		if transaction.RawData.Contract[0].Type != core.Transaction_Contract_TransferContract {
			fmt.Println("transaction is not trx transfer")
			continue
		}

		contract := &core.TransferContract{}
		err := proto.Unmarshal(transaction.RawData.Contract[0].Parameter.Value, contract)
		if err != nil {
			fmt.Println(err)
			continue
		}

		amount := contract.Amount

		// if address is hex convert to base58
		toAddress := hexutil.Encode(contract.ToAddress)[2:]
		if strings.HasPrefix(toAddress, "41") == true {
			toAddress = address.HexToAddress(toAddress).String()
		}

		// if address is hex convert to base58
		fromAddress := hexutil.Encode(contract.OwnerAddress)[2:]
		if strings.HasPrefix(fromAddress, "41") == true {
			fromAddress = address.HexToAddress(fromAddress).String()
		}

		for _, ourAddress := range c.Addresses {
			if ourAddress == toAddress || ourAddress == fromAddress {
				txs = append(txs, CrawlTransaction{
					TxId:        hexutil.Encode(t.GetTxid())[2:],
					FromAddress: fromAddress,
					ToAddress:   toAddress,
					Amount:      uint64(amount),
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
