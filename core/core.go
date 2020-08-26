/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package core

import (
	"errors"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"strings"
)

// Core - The Core Module
// Core - Core 模块
type Core struct {
	provider providers.ProviderInterface
}

// NewCore - Core Module constructor to set the default provider
// NewCore - Core Module 构造函数来初始化
func NewCore(provider providers.ProviderInterface) *Core {
	core := new(Core)
	core.provider = provider
	return core
}

func (core *Core) Contract(jsonInterface string) (*Contract, error) {
	return core.NewContract(jsonInterface)
}

/*
  GetProtocolVersion:
   	EN - Returns the current bif protocol version.
 	CN - 返回当前的Bif协议版本。
  Params:
  	- None

  Returns:
  	- string, 当前的Bif协议版本
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetProtocolVersion() (uint64, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_protocolVersion", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()

}

/*
  IsSyncing:
   	EN - Returns an object with data about the sync status or false.
 	CN - 返回同步状态的数据对象。
  Params:
  	- None

  Returns:
  	- *dto.SyncingResponse, 如果没有同步，则返回&{<nil> <nil> <nil>};如果正在同步则返回：
 		StartingBlock *big.Int - 同步时，导入的起始区块
		CurrentBlock  *big.Int - 当前区块，与GetBlockNumber效果相同
		HighestBlock  *big.Int - 当前估计的最高区块
 	- error

  Call permissions: Anyone
*/
func (core *Core) IsSyncing() (*dto.SyncingResponse, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_syncing", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToSyncingResponse()

}

/*
  GetGenerator:
   	EN - Returns the client generator address
 	CN - 返回客户端的出块奖励地址
  Params:
  	- None

  Returns:
  	- string, 20 bytes的字符串, 出块奖励地址
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetGenerator() (string, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_generator", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  Generating:
   	EN - Returns true if client is actively generating new blocks.
 	CN - 如果客户端正在积极挖掘新块，则返回true
  Params:
  	- None

  Returns:
  	- bool, true 正在出块；false 出块停止
 	- error

  Call permissions: Anyone
*/
func (core *Core) Generating() (bool, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_generating", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}

/*
  GetHashRate:
   	EN - Returns the number of hashes per second that the node is mining with.
 	CN - 返回节点每秒挖掘的哈希数
  Params:
  	- None

  Returns:
  	- *big.Int, 每秒哈希数
 	- error

  Call permissions: Anyone

*/
func (core *Core) GetHashRate() (*big.Int, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_hashrate", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

/*
  GetGasPrice:
   	EN - Returns the current price per gas in bif
 	CN - 返回以bif为单位的当前GasPrice
  Params:
  	- None

  Returns:
  	- *big.Int, 以bif为单位的当前gasPrice
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetGasPrice() (*big.Int, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_gasPrice", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

/*
  GetAccounts:
   	EN - Returns a list of addresses owned by client.
 	CN - 返回当前客户端所有的账户地址列表。
  Params:
  	- None

  Returns:
  	- []string, 当前客户端拥有的账户地址列表
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetAccounts() ([]string, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_accounts", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()

}

/*
  GetBlockNumber:
   	EN - Returns the number of most recent block.
 	CN - 返回当前区块链的区块数
  Params:
  	- None

  Returns:
  	- *big.Int, 当前区块链的区块数
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetBlockNumber() (*big.Int, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_blockNumber", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

/*
  GetBalance:
   	EN - Returns the balance of the account of given address.
 	CN - 返回给定地址的帐户余额。
  Params:
  	- address, string, 20 bytes
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- *big.Int， 给定地址的余额
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetBalance(address string, blockNumber string) (*big.Int, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = blockNumber

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBalance", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

/*
  GetTransactionCount:
   	EN - Returns the number of transactions the given address has sent for the given block number
 	CN - 返回在指定区块号下，给定地址已发送的交易数量
  Params:
	- address, string, 20 bytes
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- *big.Int, 指定区块号下，给定地址已发送的交易数量
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetTransactionCount(address string, blockNumber string) (*big.Int, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = blockNumber

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionCount", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

/*
  EstimateGas:
   	EN - Returns an estimate of the amount of gas needed to execute the given transaction, which won't be added to the blockchain and returns the used gas
 	CN - 返回执行该交易所需消耗gas的估算值，该交易不会被添加到区块链中
  Params:
  	- transaction, *dto.TransactionParameters, 详细参阅Call中参数，二者一致。
		如果未指定gas，则将使用pending区块中的gas值；如果执行交易所需的gas超过限制，则返回的gas评估值可能不足以执行交易

  Returns:
  	- *big.Int, gas消耗的量
 	- error

  Call permissions: Anyone
*/
func (core *Core) EstimateGas(transaction *dto.TransactionParameters) (*big.Int, error) {

	params := make([]*dto.RequestTransactionParameters, 1)

	params[0] = transaction.Transform()

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(&pointer, "core_estimateGas", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

/*
  SendTransaction:
	EN - Creates a transaction for the given argument, sign it and submit it to the transaction pool,return transaction hash
 	CN - 对给定参数创建交易，对其进行签名并将其提交到交易池，返回交易哈希
  Params:
  	- transaction: 要发送的交易对象(*dto.TransactionParameters)
		from: string，20 Bytes - 指定的发送者的地址。
		to: string，20 Bytes - （可选）交易消息的目标地址，如果是合约创建，则不填.
		gas: *big.Int - （可选）默认是自动，交易可使用的gas，未使用的gas会退回。
		gasPrice: *big.Int - （可选）默认是自动确定，交易的gas价格，默认是网络gas价格的平均值 。
		data: string - （可选）或者包含相关数据的字节字符串，如果是合约创建，则是初始化要用到的代码。
		value: *big.Int - （可选）交易携带的货币量，以bifer为单位。如果合约创建交易，则为初始的基金
		nonce: *big.Int - （可选）整数，使用此值，可以允许你覆盖你自己的相同nonce的，待pending中的交易

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
 	- error

  Call permissions: Anyone
*/
func (core *Core) SendTransaction(transaction *dto.TransactionParameters) (string, error) {

	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(&pointer, "core_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  SendRawTransaction:
   	EN - Add the signed transaction to the transaction pool.The sender is responsible for signing the transaction and using the correct nonce
 	CN - 将已签名的交易添加到交易池中。交易发送方负责签署交易并使用正确的随机数（Nonce）
  Params:
  	- encodedTx: string, 已签名的交易数据

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
 	- error

  Call permissions: Anyone
*/
func (core *Core) SendRawTransaction(encodedTx string) (string, error) {

	params := make([]string, 1)
	params[0] = encodedTx

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(&pointer, "core_sendRawTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  SignTransaction:
   	EN - sign the given transaction with the from account.
 	CN - 使用交易发起方的帐户（地址）签署给定的交易
  Params:
  	- transaction，*dto.TransactionParameters，交易构造的对象
	  - From     string                    交易的发起方
	  - To       string                    交易的接收方
	  - Nonce    *big.Int                  （可选）整数，使用此值，可以允许你覆盖你自己的相同nonce的，待pending中的交易
	  - Gas      *big.Int                  （可选）默认是自动，交易可使用的gas，未使用的gas会退回。
	  - GasPrice *big.Int                  （可选）默认是自动确定，交易的gas价格，默认是网络gas价格的平均值 。
	  - Value    *big.Int                  （可选）交易携带的货币量，以bifer为单位。如果合约创建交易，则为初始的基金
	  - Data     types.ComplexString       （可选）或者包含相关数据的字节字符串，如果是合约创建，则是初始化要用到的代码。

  Returns:
  	- *dto.SignTransactionResponse，
		Raw         string                   已签名的RLP编码的交易
		Transaction SignedTransactionParams  transaction object
		  - Gas      *big.Int                交易发起方约定的gas
		  - GasPrice *big.Int                交易发起方约定的gasPrice
		  - Hash     string  			     交易哈希
		  - Input    string   				 随交易发送的数据
		  - Nonce    *big.Int                交易发起者之前发起交易的次数
		  - S        string                  ？？？
		  - R        string                  ？？？
		  - V        *big.Int                ？？？
		  - To       string                  交易的接收方，如果是合约创建则为空
		  - Value    *big.Int                转移的bif数量
 	- error

  Call permissions: 交易的发起方的账户处于解锁状态

*/
func (core *Core) SignTransaction(transaction *dto.TransactionParameters) (*dto.SignTransactionResponse, error) {
	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(&pointer, "core_signTransaction", params)

	if err != nil {
		return &dto.SignTransactionResponse{}, err
	}

	return pointer.ToSignTransactionResponse()
}

/*
  Call:
   	EN - Executes a new message call immediately without creating a transaction on the block chain.
 	CN - 执行新的消息调用，而无需在区块链上创建交易，它不会改变区块链的状态，一般用于检索。
  Params:
	- transaction，*dto.TransactionParameters，交易Call的对象
 	  - From     string                    交易的发起方
 	  - To       string                    交易的接收方
 	  - Nonce    *big.Int                  （可选）整数，使用此值，可以允许你覆盖你自己的相同nonce的，待pending中的交易
 	  - Gas      *big.Int                  （可选）默认是自动，交易可使用的gas，未使用的gas会退回。
 	  - GasPrice *big.Int                  （可选）默认是自动确定，交易的gas价格，默认是网络gas价格的平均值 。
 	  - Value    *big.Int                  （可选）交易携带的货币量，以bifer为单位。如果合约创建交易，则为初始的基金
 	  - Data     types.ComplexString       （可选）或者包含相关数据的字节字符串，如果是合约创建，则是初始化要用到的代码。

  Returns:
  	- 已执行合约的返回值
 	- error

  Call permissions: Anyone
  Bug 待测试，需要比对rpc中的callArgs和sendTxArgs！！！！！，数据结构
*/
func (core *Core) Call(transaction *dto.TransactionParameters) (*dto.RequestResult, error) {

	params := make([]interface{}, 2)
	params[0] = transaction.Transform()
	params[1] = block.LATEST

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(&pointer, "core_call", params)

	if err != nil {
		return nil, err
	}

	return pointer, err

}

/*
  GetTransactionReceipt:
   	EN - Returns the transaction receipt for the given transaction hash
 	CN - 返回给定交易哈希的交易收据
  Params:
  	- hash,string 32 Bytes 交易哈希

  Returns:
  	- *dto.TransactionReceipt 交易收据对象
	  - TransactionHash   string            交易哈希
	  - TransactionIndex  *big.Int          交易在区块中的索引（位置）
	  - BlockHash         string            区块哈希
	  - BlockNumber       *big.Int          区块
	  - From              string            交易发起方
	  - To                string            交易接收方
	  - CumulativeGasUsed *big.Int          在区块中执行此交易时使用的gas总量。
	  - GasUsed           *big.Int          仅此特定交易使用的gas量。
	  - ContractAddress   string            如果交易是合约创建，则为创建的合约地址，否则为空
	  - Logs              []TransactionLogs 交易生成的日志对象数组
	  - LogsBloom         string            ？？？
	  - Status            bool              ？？？
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetTransactionReceipt(hash string) (*dto.TransactionReceipt, error) {

	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionReceipt", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionReceipt()

}

/*
  GetBlockTransactionCountByHash:
   	EN - Returns the number of transactions in the block with the given hash
 	CN - 根据区块哈希获取该区块内包含的交易数
  Params:
  	- hash，32 bytes - block hash

  Returns:
  	- uint64, 区块包含的交易数
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetBlockTransactionCountByHash(hash string) (uint64, error) {
	// ensure that the hash is correctly formatted
	if strings.HasPrefix(hash, "0x") {
		if len(hash) != 66 {
			return 0, errors.New("malformed block hash")
		}
	} else {
		if len(hash) != 64 {
			return 0, errors.New("malformed block hash")
		}
		hash = "0x" + hash
	}

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockTransactionCountByHash", []string{hash})

	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
  GetBlockTransactionCountByNumber:
   	EN - Returns the number of transactions in a block matching the given block number
 	CN - 返回与给定区块编号匹配的区块中的交易数量
  Params:
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- uint64, 区块包含的交易数
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetBlockTransactionCountByNumber(blockNumber string) (uint64, error) {

	params := make([]string, 1)
	params[0] = blockNumber

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockTransactionCountByNumber", params)

	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
  GetBlockByNumber:
   	EN - Returns the information about a block requested by blockNumber
 	CN - 根据区块号返回区块的信息
  Params:
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions
  	- transactionDetails,bool, 如果为True，返回区块内详细的交易信息和其他信息；如果为false则仅返回区块内交易hash和其他信息

  Returns:
  	- interface{}, 如果transactionDetails为true，则是*dto.BlockDetails；如果为false，则是*dto.BlockNoDetails
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetBlockByNumber(blockNumber string, transactionDetails bool) (interface{}, error) {

	params := make([]interface{}, 2)
	params[0] = blockNumber
	params[1] = transactionDetails

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockByNumber", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBlock(transactionDetails)
}

/*
  GetBlockByHash:
   	EN - Returns the information about a block requested by block hash
 	CN - 根据区块哈希返回区块的信息
  Params:
	- blockHash, string32 bytes - Hash of a block
  	- transactionDetails,bool, 如果为True，返回区块内详细的交易信息和其他信息；如果为false则仅返回区块内交易hash和其他信息

  Returns:
  	- interface{}, 如果transactionDetails为true，则是*dto.BlockDetails；如果为false，则是*dto.BlockNoDetails
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetBlockByHash(blockHash string, transactionDetails bool) (interface{}, error) {
	// ensure that the hash is correctly formatted
	if strings.HasPrefix(blockHash, "0x") {
		if len(blockHash) != 66 {
			return nil, errors.New("malformed block hash")
		}
	} else {
		blockHash = "0x" + blockHash
		if len(blockHash) != 62 {
			return nil, errors.New("malformed block hash")
		}
	}

	params := make([]interface{}, 2)
	params[0] = blockHash
	params[1] = transactionDetails

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockByHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBlock(transactionDetails)
}

/*
  GetCode:
   	EN - Returns the code stored at the given address in the state for the given block number.
 	CN - 返回以给定块号存储在给定地址的代码。？？
  Params:
  	- address,string 20 Bytes 账户地址
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- string，the code
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetCode(address string, blockNumber string) (string, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = blockNumber

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getCode", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
  GetTrustNumber:
   	EN - Check the numbers of trusted certificates at given address
 	CN - 检查账户有多少可信证书
  Params:
  	- address, string 账户地址

  Returns:
  	- uint64, 可信证书数量
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetTrustNumber(address string) (uint64, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = block.LATEST

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getCanTrust", params)

	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
  GetChainId:
   	EN - Returns the chain ID of the current connected node
 	CN - 返回当前连接节点的链ID
  Params:
  	- None

  Returns:
  	- uint64, chain ID
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetChainId() (uint64, error) {
	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_chainId", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
  描述:
   	EN - Returns the Merkle-proof for a given account and optionally some storage keys.
 	CN - 返回给定帐户的Merkle证明  ????
  Params:
  	- address, string, 20 bytes
	- storageKeys,  []string, 一组storageKeys，应对其进行校验并包括在内
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- string
	- error

  Call permissions: Anyone
*/
func (core *Core) GetProof(address string, storageKeys []string, blockNumber string) (*dto.AccountResult, error) {

	params := make([]interface{}, 3)
	params[0] = address
	params[1] = storageKeys
	params[2] = blockNumber

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getProof", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToProof()
}

/*
  GetStorageAt:
   	EN - Returns the value from a storage at the given address, key and block number
 	CN - 从给定地址、位置和区块号的状态返回存储值???(是什么存储值)
  Params:
	- address, string, 20 bytes
	- key,  *big.Int, 指定的位置
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- string
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetStorageAt(address string, key *big.Int, blockNumber string) (string, error) {

	params := make([]string, 3)
	params[0] = address
	params[1] = hexutil.EncodeBig(key)
	params[2] = blockNumber

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getStorageAt", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
  GetPendingTransactions:
   	EN - Returns a list of pending transactions
 	CN - 返回未执行交易的列表(交易还未被打包)
  Params:
	- None

  Returns:
  	- []*dto.TransactionResponse
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetPendingTransactions() ([]*dto.TransactionResponse, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_pendingTransactions", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToPendingTransactions()
}


/*
  GetTransactionByHash:
   	EN - Returns the information about a transaction requested by transaction hash.
 	CN - 根据交易hash返回交易的详细信息
  Params:
  	- hash,string,32 Bytes,hash of a transaction

  Returns:
  	- *dto.TransactionResponse
		Hash             string              交易哈希
		Nonce            *big.Int            交易发送方从该节点发送的交易数
		BlockHash        string              交易所在的区块的哈希，如果交易待处理，则为0x0000000000000000000000000000000000000000000000000000000000000000
		BlockNumber      *big.Int            交易所在的区块，如果交易待处理，则为0
		TransactionIndex *big.Int            交易在区块中的位置（索引），如果交易待处理，则为0
		From             string              交易发送方
		To               string              交易接收方，如果为合约创建则为空
		Input            string              随交易发送的数据
		Value            *big.Int            交易转移的bif数量，单位为bif
		GasPrice         *big.Int            发送方提供的GasPrice，单位为bif
		Gas              *big.Int            发送方提供的Gas
		Data             types.ComplexString ？？？？
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetTransactionByHash(hash string) (*dto.TransactionResponse, error) {

	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionByHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()

}

/*
  GetRawTransactionByHash:
   	EN - Returns the bytes of the transaction for the given hash
 	CN - 根据交易hash返回交易信息
  Params:
  	- hash,string,32 Bytes,hash of a transaction

  Returns:
  	- string
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetRawTransactionByHash(hash string) (string, error) {

	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getRawTransactionByHash", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  GetTransactionByBlockHashAndIndex:
   	EN - Returns the transaction for the given block hash and index
 	CN - 根据区块哈希和交易索引，返回交易信息
  Params:
  	- hash,string,32 Bytes,区块哈希
  	- index, *big.Int, 交易在区块中的索引

  Returns:
  	- *dto.TransactionResponse, 参照GetTransactionByHash的返回值，二者一致
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetTransactionByBlockHashAndIndex(hash string, index *big.Int) (*dto.TransactionResponse, error) {

	// ensure that the hash is correctly formatted
	if strings.HasPrefix(hash, "0x") {
		if len(hash) != 66 {
			return nil, errors.New("malformed block hash")
		}
	} else {
		if len(hash) != 64 {
			return nil, errors.New("malformed block hash")
		}

		hash = "0x" + hash
	}

	params := make([]string, 2)
	params[0] = hash
	params[1] = hexutil.EncodeBig(index)

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionByBlockHashAndIndex", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

/*
  GetRawTransactionByBlockHashAndIndex:
   	EN - Returns the bytes of the transaction for the given block hash and index
 	CN - 根据区块哈希和交易索引，返回交易信息
  Params:
  	- hash,string,32 Bytes,区块哈希
  	- index, *big.Int, 交易在区块中的索引

  Returns:
  	- string
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetRawTransactionByBlockHashAndIndex(hash string, index *big.Int) (string, error) {

	// ensure that the hash is correctly formatted
	if strings.HasPrefix(hash, "0x") {
		if len(hash) != 66 {
			return "", errors.New("malformed block hash")
		}
	} else {
		if len(hash) != 64 {
			return "", errors.New("malformed block hash")
		}

		hash = "0x" + hash
	}

	params := make([]string, 2)
	params[0] = hash
	params[1] = hexutil.EncodeBig(index)

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getRawTransactionByBlockHashAndIndex", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

/*
  GetTransactionByBlockNumberAndIndex:
   	EN - Returns the transaction for the given block number and transaction index
 	CN - 返回给定区块号和索引的交易
  Params:
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions
  	- index, *big.Int, 交易在区块中的索引

  Returns:
  	- *dto.TransactionResponse, 参照GetTransactionByHash的返回值，二者一致
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetTransactionByBlockNumberAndIndex(blockNumber string, index *big.Int) (*dto.TransactionResponse, error) {

	params := make([]string, 2)
	params[0] = blockNumber
	params[1] = hexutil.EncodeBig(index)

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionByBlockNumberAndIndex", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()

}

/*
  GetRawTransactionByBlockNumberAndIndex:
   	EN - Returns the bytes of the transaction for the given block number and index
 	CN - 返回给定区块号和索引的交易
  Params:
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions
  	- index, *big.Int, 交易在区块中的索引

  Returns:
  	- string (如果返回0x，则表示交易未执行或不存在)
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetRawTransactionByBlockNumberAndIndex(blockNumber string, index *big.Int) (string, error) {

	params := make([]string, 2)
	params[0] = blockNumber
	params[1] = hexutil.EncodeBig(index)

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_getRawTransactionByBlockNumberAndIndex", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
