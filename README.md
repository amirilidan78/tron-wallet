# tron-wallet
tron wallet package

### Main methods 
- generating tron wallet 
```
w := GenerateTronWallet(enums.NetworkSHASTA)
w.Address // strnig 
w.AddressBase58 // strnig 
w.PrivateKey // strnig 
w.PublicKey // strnig 
```
- creating tron wallet from private key 
```
w := CreateTronWallet(enums.NetworkSHASTA,privateKeyHex)
w.Address // strnig 
w.AddressBase58 // strnig 
w.PrivateKey // strnig 
w.PublicKey // strnig 
```
- getting wallet balance 
```
balanceInSun,err := w.Balance()
balanceInSun // int64 
```
- crawl blocks for addresses transactions 
```

c := &Crawler{
		Network: enums.NetworkSHASTA, // network -> maninet, shasta, nile
		Addresses: []string{
			"TY3PJu3VY8xVUc5BjYwJtyRgP7TfivV666", // list of your addresses
		},
	}
	
res, err := c.ScanBlocks(40) // scan latest 40 block on block chain and extract addressess transactions 

Example 
// *
{
    {
        "address": "TY3PJu3VY8xVUc5BjYwJtyRgP7TfivV666",
        "tranasctions": {
            {
                "tx_id": "6afbc5758d49e8d8bedddd903edbfc01c5f11ebfbaa6237e887294a6fc9394a2",
                "from_address": "TMaxbcktAkSQiwvz9eQ7GqvVBa99n5U555",
                "to_address": "TY3PJu3VY8xVUc5BjYwJtyRgP7TfivV666",
                "amount": "195500",
            }
        }
    },
    ...
}
* // 
	
```
- transfer from wallet - /// TODO : should be implemented 


### Supported networks 
- Main net  
- Shasta 
- Nile

### Api methods 
- `GetAddressBalance(network enums.Network, address string) (GetAccountResponseBody, error)`
- `CurrentBlock(network enums.Network) (BlockResponseBody, error)`
- `GetTransaction(network enums.Network, txHash string) (Transaction, error)`
- ...