# tron-wallet
tron wallet package

### Main methods 
- generating tron wallet `w := GenerateTronWallet(enums.NetworkSHASTA)`
- creating tron wallet from private key `w := CreateTronWallet(enums.NetworkSHASTA,privateKeyHex)`
- getting wallet balance `balanceInSun,err := w.Balance()`
- transfer from wallet - /// TODO : should be implemented 
- check list of addresses new transaction - /// TODO : should be implemented 


### Supported networks 
- Main net  
- Shasta 
- Nile

### Api methods 
- `GetAddressBalance(network enums.Network, address string) (GetAccountResponseBody, error)`
- `CurrentBlock(network enums.Network) (BlockResponseBody, error)`
- `GetTransaction(network enums.Network, txHash string) (Transaction, error)`
- ...