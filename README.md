# DEPRECATED

 [![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)

This project is no longer supported, please consider using [go-ethereum](https://github.com/ethereum/go-ethereum) instead.

[go-ethereum](https://github.com/ethereum/go-ethereum) has all the features of this project(and more) and it's development is much more robust.

# Ethereum Go Client

[![Build Status](https://travis-ci.org/bif/go-bif-sdk.svg?branch=master)](https://travis-ci.org/bif/go-bif-sdk)

This is a Ethereum compatible Go Client

## Status

## DEPRECATED

This package is not currently under active development. It is not already stable and the infrastructure is not complete and there are still several RPCs left to implement.

## Usage

#### Deploying a contract

```go

bytecode := ... #contract bytecode
abi := ... #contract abi

var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))
contract, err := connection.core.NewContract(abi)

transaction := new(dto.TransactionParameters)
coinbase, err := connection.core.GetCoinbase()
transaction.From = coinbase
transaction.Gas = big.NewInt(4000000)

hash, err := contract.Deploy(transaction, bytecode, nil)

fmt.Println(hash)
	
```

#### Using contract public functions

```go

result, err = contract.Call(transaction, "balanceOf", coinbase)
if result != nil && err == nil {
	balance, _ := result.ToComplexIntResponse()
	fmt.Println(balance.ToBigInt())
}
	
```

#### Using contract payable functions

```go

hash, err = contract.Send(transaction, "approve", coinbase, 10)
	
```

#### Using RPC commands

GetBalance

```go

balance, err := connection.core.GetBalance(coinbase, block.LATEST)

```

SendTransaction

```go

transaction := new(dto.TransactionParameters)
transaction.From = coinbase
transaction.To = coinbase
transaction.Value = big.NewInt(10)
transaction.Gas = big.NewInt(40000)
transaction.Data = types.ComplexString("p2p transaction")

txID, err := connection.core.SendTransaction(transaction)

```


## Contribute!

#### Before a Pull Request:
- Create at least one test for your implementation.
- Don't change the import path to your github username.
- run `go fmt` for all your changes.
- run `go test -v ./...`

#### After a Pull Request:
- Please use the travis log if an error occurs.

### In Progress = ![](https://placehold.it/15/FFFF00/000000?text=+)
### Partially implemented = ![](https://placehold.it/15/008080/000000?text=+)

TODO List

- [x] bif_clientVersion                      
- [x] bif_sha3                               
- [x] net_version                             
- [x] net_peerCount                           
- [x] net_listening                           
- [x] core_protocolVersion                     
- [x] core_syncing                             
- [x] core_coinbase                            
- [x] core_mining                              
- [x] core_hashrate                            
- [x] core_gasPrice                            
- [x] core_accounts                            
- [x] core_blockNumber                         
- [x] core_getBalance                          
- [x] core_getStorageAt (deprecated)
- [x] core_getTransactionCount                 
- [x] core_getBlockTransactionCountByHash      
- [x] core_getBlockTransactionCountByNumber    
- [x] core_getUncleCountByBlockHash            
- [x] core_getUncleCountByBlockNumber          
- [x] core_getCode                             
- [x] core_sign                                
- [x] core_sendTransaction                     
- [ ] core_sendRawTransaction                  
- [x] core_call                                
- [x] core_estimateGas                         
- [x] core_getBlockByHash                      
- [x] core_getBlockByNumber                    
- [x] core_getTransactionByHash                
- [x] core_getTransactionByBlockHashAndIndex   
- [x] core_getTransactionByBlockNumberAndIndex 
- [x] core_getTransactionReceipt               
- [ ] core_getUncleByBlockHashAndIndex         
- [ ] core_getUncleByBlockNumberAndIndex       
- [ ] core_getCompilers                        
- [ ] core_compileLLL                          
- [x] core_compileSolidity (deprecated)                    
- [ ] core_compileSerpent                      
- [ ] core_newFilter                           
- [ ] core_newBlockFilter                      
- [ ] core_newPendingTransactionFilter         
- [ ] core_uninstallFilter                     
- [ ] core_getFilterChanges                    
- [ ] core_getFilterLogs                       
- [ ] core_getLogs                             
- [ ] core_getWork                             
- [ ] core_submitWork                          
- [ ] core_submitHashrate                      
- [ ] db_putString                            
- [ ] db_getString                            
- [ ] db_putHex                               
- [ ] db_getHex                               
- [ ] shh_post                                
- [ ] shh_version                             
- [ ] shh_newIdentity                         
- [ ] shh_hasIdentity                         
- [ ] shh_newGroup                            
- [ ] shh_addToGroup                          
- [ ] shh_newFilter                           
- [ ] shh_uninstallFilter                     
- [ ] shh_getFilterChanges                    
- [ ] shh_getMessages                         
- [x] personal_listAccounts                   
- [x] personal_newAccount                     
- [x] personal_sendTransaction                
- [x] personal_unlockAccount                  

## Installation

### go get

```bash
go get -u github.com/bif/bif-sdk-go
```

### glide

```bash
glide get github.com/bif/bif-sdk-go
```

### Requirements

* go ^1.8.3
* golang.org/x/net

## Testing

Node running in dev mode:

```bash
geth --dev --shh --ws --wsorigins="*" --rpc --rpcapi admin,db,core,debug,miner,net,shh,txpool,personal,web3 --mine
```

Full test:

```bash
go test -v ./test/...
```

Individual test:
```bash
go test -v test/modulename/filename.go
```

## License

Package go-web3 is licensed under the [GPLv3](https://www.gnu.org/licenses/gpl-3.0.en.html) License.
