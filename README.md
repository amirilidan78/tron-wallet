# tron-wallet
tron wallet package for creating and generating wallet, transferring TRX, getting wallet balance and crawling blocks to find wallet transactions

### Installation 
```
go get github.com/Amirilidan78/tron-wallet@v0.1.0
```

### Wallet methods 
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
- getting wallet trx balance 
```
balanceInSun,err := w.Balance()
balanceInSun // int64 
```
- getting wallet trc20 balance
```
balanceInToken,err := w.BalanceTRC20(token)
balanceInToken // int64 
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
- transfer trx from wallet 
```
txId, err := w.Transfer(toAddressBase58, amount)
txId // string 
```
- transfer trc20 from wallet
```
txId, err := w.TransferTRC20(token, toAddressBase58, amount)
txId // string 
```

### Token methods 
- declaring token 
```
token := &tronWallet.Token{
    ContractAddress: enums.SHASTA_Tether_USDT,
}
```
- Getting token name 
```
token.GetName(w.Node, w.AddressBase58) // return string,error
``` 
- Getting token symbol
```
token.GetSymbol(w.Node, w.AddressBase58) // return string,error
```
- Getting token decimals 
```
token.GetDecimal(w.Node, w.AddressBase58) // return int64,error
```

### Supported networks
check `enums/nodes` file
alternatively you can create your own node
```
node := enums.CreateNode("grpc.test.com:50051")
```

### Supported contracts
check `enums/contracts` file
alternatively you can create your own contract
```
contractAddress := enums.CreateContractAddress("TPYmHEhy5n8TCEfYGqW2rPxsghSfzghPDn")
```

### TRX Faucet
Follow TronTestnet Twitter account
@TronTest2
.
Write your address in your tweet and
@TronTest2
.
They will transfer 10,000 test TRX (usually within five minutes).
Each address can only be obtained once a day.
If you need TRX for the nile testnet, please add "NILE" in your tweet.

### TRC20 Faucet
Go to https://developers.tron.network/ and connect to the discord community.
You can than ask for usdt in #faucet channel.
Just type !shasta_usdt YOUR_WALLET_ADDRESS and send. TronFAQ bot will send you 5000  USDT (SASHTA) soon.


### Important
I simplified this repository https://github.com/fbsobreira repository to create this package
You can check go tron sdk for better examples and functionalities
and do not use this package in production, I created this package for education purposes.


