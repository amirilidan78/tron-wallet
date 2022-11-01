# tron-wallet
tron wallet package

### Main methods 
- generating tron wallet 
```
w := GenerateTronWallet(enums.SHASTA_NODE)
w.Address // strnig 
w.AddressBase58 // strnig 
w.PrivateKey // strnig 
w.PublicKey // strnig 
```
- creating tron wallet from private key 
```
w := CreateTronWallet(enums.SHASTA_NODE,privateKeyHex)
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
		Node: enums.SHASTA_NODE, // network -> maninet, shasta, nile
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
```
txId, err := w.Transfer(toAddressBase58, amount)
txId // string 
```

### Util methods 
- convert base58 address to hex
```
hex := util.Base58ToHex("TNvQe93ay9MACT26oC92sP9NkvVqqXm2Cw") // <- 41718de6b323652d1257437ace160c4f4198aae4e1

```
- convert hex address to base58
```
hex := util.HexToBase58("41718de6b323652d1257437ace160c4f4198aae4e1") // <- TNvQe93ay9MACT26oC92sP9NkvVqqXm2Cw
```

### Supported networks
- Main net
- Shasta
- Nile