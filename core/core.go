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

/**
 * @file core.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package core

import (
	"errors"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
	"strings"
)

// Core - The Core Module
type Core struct {
	provider providers.ProviderInterface
}

// NewEth - Core Module constructor to set the default provider
func NewEth(provider providers.ProviderInterface) *Core {
	core := new(Core)
	core.provider = provider
	return core
}

func (core *Core) Contract(jsonInterface string) (*Contract, error) {
	return core.NewContract(jsonInterface)
}

// GetProtocolVersion - Returns the current ethereum protocol version.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_protocolversion
// Parameters:
//    - none
// Returns:
// 	  - String - The current ethereum protocol version
func (core *Core) GetProtocolVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_protocolVersion", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

// IsSyncing - Returns an object with data about the sync status or false.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_syncing
// Parameters:
//    - none
// Returns:
// 	  - Object|Boolean, An object with sync status data or FALSE, when not syncing:
//    	- startingBlock: 	QUANTITY - The block at which the import started (will only be reset, after the sync reached his head)
//    	- currentBlock: 	QUANTITY - The current block, same as core_blockNumber
//    	- highestBlock: 	QUANTITY - The estimated highest block
func (core *Core) IsSyncing() (*dto.SyncingResponse, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_syncing", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToSyncingResponse()

}

// GetCoinbase - Returns the client coinbase address.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_coinbase
// Parameters:
//    - none
// Returns:
// 	  - DATA, 20 bytes - the current coinbase address.
func (core *Core) GetCoinbase() (string, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_coinbase", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

// IsMining - Returns true if client is actively mining new blocks.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_mining
// Parameters:
//    - none
// Returns:
// 	  - Boolean - returns true of the client is mining, otherwise false.
func (core *Core) Generating() (bool, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_generating", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}

// GetHashRate - Returns the number of hashes per second that the node is mining with.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_hashrate
// Parameters:
//    - none
// Returns:
// 	  - QUANTITY - number of hashes per second.
func (core *Core) GetHashRate() (*big.Int, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_hashrate", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// GetGasPrice - Returns the current price per gas in wei.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_gasprice
// Parameters:
//    - none
// Returns:
// 	  - QUANTITY - integer of the current gas price in wei.
func (core *Core) GetGasPrice() (*big.Int, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_gasPrice", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// ListAccounts - Returns a list of addresses owned by client.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_accounts
// Parameters:
//    - none
// Returns:
//    - Array of DATA, 20 Bytes - addresses owned by the client.
func (core *Core) ListAccounts() ([]string, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_accounts", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()

}

// GetBlockNumber - Returns the number of most recent block.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_blocknumber
// Parameters:
//    - none
// Returns:
// 	  - QUANTITY - integer of the current block number the client is on.
func (core *Core) GetBlockNumber() (*big.Int, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_blockNumber", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// GetBalance - Returns the balance of the account of given address.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getbalance
// Parameters:
//    - DATA, 20 Bytes - address to check for balance.
//	  - QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending", see the default block parameter: https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
// Returns:
// 	  - QUANTITY - integer of the current balance in wei.
func (core *Core) GetBalance(address string, defaultBlockParameter string) (*big.Int, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBalance", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// GetTransactionCount -  Returns the number of transactions sent from an address.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_gettransactionaccount
// Parameters:
//    - DATA, 20 Bytes - address to check for balance.
//	  - QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending", see the default block parameter: https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
// Returns:
// 	  - QUANTITY - integer of the number of transactions sent from this address
func (core *Core) GetTransactionCount(address string, defaultBlockParameter string) (*big.Int, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionCount", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// GetStorageAt - Returns the value from a storage position at a given address.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getstorageat
// Parameters:
//    - DATA, 20 Bytes - address of the storage.
//	  - QUANTITY - integer of the position in the storage.
//	  - QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending", see the default block parameter: https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter.
// Returns:
// 	  - DATA - the value at this storage position.
func (core *Core) GetStorageAt(address string, position *big.Int, defaultBlockParameter string) (string, error) {

	params := make([]string, 3)
	params[0] = address
	params[1] = utils.IntToHex(position)
	params[2] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getstorageat", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// EstimateGas - Makes a call or transaction, which won't be added to the blockchain and returns the used gas, which can be used for estimating the used gas.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_estimategas
// Parameters:
//    - See core_call parameters, expect that all properties are optional. If no gas limit is specified geth uses the block gas limit from the pending block as an
// 		upper bound. As a result the returned estimate might not be enough to executed the call/transaction when the amount of gas is higher than the pending block gas limit.
// Returns:
//    - QUANTITY - the amount of gas used.
func (core *Core) EstimateGas(transaction *dto.TransactionParameters) (*big.Int, error) {

	params := make([]*dto.RequestTransactionParameters, 1)

	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(&pointer, "core_estimateGas", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// GetTransactionByHash - Returns the information about a transaction requested by transaction hash.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_gettransactionbyhash
// Parameters:
//    - DATA, 32 Bytes - hash of a transaction
// Returns:
//    1. Object - A transaction object, or null when no transaction was found
//    - hash: DATA, 32 Bytes - hash of the transaction.
//    - nonce: QUANTITY - the number of transactions made by the sender prior to this one.
//    - blockHash: DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
//    - blockNumber: QUANTITY - block number where this transaction was in. null when its pending.
//    - transactionIndex: QUANTITY - integer of the transactions index position in the block. null when its pending.
//    - from: DATA, 20 Bytes - address of the sender.
//    - to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
//    - value: QUANTITY - value transferred in Wei.
//    - gasPrice: QUANTITY - gas price provided by the sender in Wei.
//    - gas: QUANTITY - gas provided by the sender.
//    - input: DATA - the data send along with the transaction.
func (core *Core) GetTransactionByHash(hash string) (*dto.TransactionResponse, error) {

	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionByHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()

}

// GetTransactionByBlockHashAndIndex - Returns the information about a transaction requested by block hash.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getTransactionByBlockNumberAndIndex
// Parameters:
//    - DATA, 32 Bytes - hash of a block
//    - QUANTITY, number - index of the transaction position
// Returns:
//    1. Object - A transaction object, or null when no transaction was found
//    - hash: DATA, 32 Bytes - hash of the transaction.
//    - nonce: QUANTITY - the number of transactions made by the sender prior to this one.
//    - blockHash: DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
//    - blockNumber: QUANTITY - block number where this transaction was in. null when its pending.
//    - transactionIndex: QUANTITY - integer of the transactions index position in the block. null when its pending.
//    - from: DATA, 20 Bytes - address of the sender.
//    - to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
//    - value: QUANTITY - value transferred in Wei.
//    - gasPrice: QUANTITY - gas price provided by the sender in Wei.
//    - gas: QUANTITY - gas provided by the sender.
//    - input: DATA - the data send along with the transaction.
func (core *Core) GetTransactionByBlockHashAndIndex(hash string, index *big.Int) (*dto.TransactionResponse, error) {

	// ensure that the hash is correctlyformatted
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
	params[1] = utils.IntToHex(index)

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionByBlockHashAndIndex", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()
}

// GetTransactionByBlockNumberAndIndex - Returns the information about a transaction requested by block index.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getTransactionByBlockNumberAndIndex
// Parameters:
//    - QUANTITY, number - block number
//    - QUANTITY, number - transaction index in block
// Returns:
//    1. Object - A transaction object, or null when no transaction was found
//    - hash: DATA, 32 Bytes - hash of the transaction.
//    - nonce: QUANTITY - the number of transactions made by the sender prior to this one.
//    - blockHash: DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
//    - blockNumber: QUANTITY - block number where this transaction was in. null when its pending.
//    - transactionIndex: QUANTITY - integer of the transactions index position in the block. null when its pending.
//    - from: DATA, 20 Bytes - address of the sender.
//    - to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
//    - value: QUANTITY - value transferred in Wei.
//    - gasPrice: QUANTITY - gas price provided by the sender in Wei.
//    - gas: QUANTITY - gas provided by the sender.
//    - input: DATA - the data send along with the transaction.
func (core *Core) GetTransactionByBlockNumberAndIndex(blockIndex *big.Int, index *big.Int) (*dto.TransactionResponse, error) {

	params := make([]string, 2)
	params[0] = utils.IntToHex(blockIndex)
	params[1] = utils.IntToHex(index)

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionByBlockNumberAndIndex", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionResponse()

}

// SendTransaction - Creates new message call transaction or a contract creation, if the data field contains code.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_sendtransaction
// Parameters:
//    1. Object - The transaction object
//    - from: 		DATA, 20 Bytes - The address the transaction is send from.
//    - to: 		DATA, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
//    - gas: 		QUANTITY - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
//    - gasPrice: 	QUANTITY - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas
//    - value: 		QUANTITY - (optional) Integer of the value send with this transaction
//    - data: 		DATA - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters. For details see Ethereum Contract ABI (https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI)
//    - nonce: 		QUANTITY - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
// Returns:
//	  - DATA, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.
// Use core_getTransactionReceipt to get the contract address, after the transaction was mined, when you created a contract.
func (core *Core) SendTransaction(transaction *dto.TransactionParameters) (string, error) {

	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(&pointer, "core_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

// encodedTx, The signed transaction data
func (core *Core) SendRawTransaction(encodedTx string) (string, error) {

	params := make([]string, 1)
	params[0] = encodedTx

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(&pointer, "core_sendRawTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

// SignTransaction - Signs transactions without dispatching it to the network. It can be later submitted using core_sendRawTransaction.
// Reference: https://wiki.parity.io/JSONRPC-core-module.html#core_signtransaction
// Parameters:
//    1. Object - The transaction call object
//    - from: 		DATA, 20 Bytes - The address the transaction is send from.
//    - to: 		DATA, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
//    - gas: 		QUANTITY - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
//    - gasPrice: 	QUANTITY - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas
//    - value: 		QUANTITY - (optional) Integer of the value send with this transaction
//    - data: 		DATA - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters. For details see Ethereum Contract ABI (https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI)
//    - nonce: 		QUANTITY - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
// Returns:
//    1. Object - A transaction sign result object
//    - raw: DATA - The signed, RLP encoded transaction.
//    - tx: Object - A transaction object
//      - hash: DATA, 32 Bytes - hash of the transaction.
//      - nonce: QUANTITY - the number of transactions made by the sender prior to this one.
//      - blockHash: DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
//      - blockNumber: QUANTITY - block number where this transaction was in. null when its pending.
//      - transactionIndex: QUANTITY - integer of the transactions index position in the block. null when its pending.
//      - from: DATA, 20 Bytes - address of the sender.
//      - to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
//      - value: QUANTITY - value transferred in Wei.
//      - gasPrice: QUANTITY - gas price provided by the sender in Wei.
//      - gas: QUANTITY - gas provided by the sender.
//      - input: DATA - the data send along with the transaction.
// Use core_sendRawTransaction to submit the transaction after it was signed.
func (core *Core) SignTransaction(transaction *dto.TransactionParameters) (*dto.SignTransactionResponse, error) {
	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(&pointer, "core_signTransaction", params)

	if err != nil {
		return &dto.SignTransactionResponse{}, err
	}

	return pointer.ToSignTransactionResponse()
}

// Call - Executes a new message call immediately without creating a transaction on the block chain.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_call
// Parameters:
//    1. Object - The transaction call object
//    - from: 		DATA, 20 Bytes - The address the transaction is send from.
//    - to: 		DATA, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
//    - gas: 		QUANTITY - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
//    - gasPrice: 	QUANTITY - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas
//    - value: 		QUANTITY - (optional) Integer of the value send with this transaction
//    - data: 		DATA - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters. For details see Ethereum Contract ABI (https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI)
//	  2. QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending", see the default block parameter: https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
// Returns:
//	  - DATA - the return value of executed contract.
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

// GetTransactionReceipt - Returns compiled solidity code.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_gettransactionreceipt
// Parameters:
//    1. DATA, 32 Bytes - hash of a transaction.
// Returns:
//	  1. Object - A transaction receipt object, or null when no receipt was found:
//    - transactionHash: 		DATA, 32 Bytes - hash of the transaction.
//    - transactionIndex: 		QUANTITY - integer of the transactions index position in the block.
//    - blockHash: 				DATA, 32 Bytes - hash of the block where this transaction was in.
//    - blockNumber:			QUANTITY - block number where this transaction was in.
//    - cumulativeGasUsed: 		QUANTITY - The total amount of gas used when this transaction was executed in the block.
//    - gasUsed: 				QUANTITY - The amount of gas used by this specific transaction alone.
//    - contractAddress: 		DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
//    - logs: 					Array - Array of log objects, which this transaction generated.
func (core *Core) GetTransactionReceipt(hash string) (*dto.TransactionReceipt, error) {

	params := make([]string, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getTransactionReceipt", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToTransactionReceipt()

}

// GetBlockByNumber - Returns the information about a block requested by number.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getblockbynumber
// Parameters:
//    - number, QUANTITY - number of block
//    - transactionDetails, bool - indicate if we should have or not the details of the transactions of the block
// Returns:
//    1. Object - A block object, or null when no transaction was found
//    2. error
func (core *Core) GetBlockByNumber(number *big.Int, transactionDetails bool) (interface{}, error) {

	params := make([]interface{}, 2)
	params[0] = utils.IntToHex(number)
	params[1] = transactionDetails

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockByNumber", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBlock(transactionDetails)
}

// GetBlockTransactionCountByHash
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getblocktransactioncountbyhash
// Parameters:
//    - DATA, 32 bytes - block hash
// Returns:
//    1. QUANTITY, number - number of transactions in the block
//    2. error
func (core *Core) GetBlockTransactionCountByHash(hash string) (*big.Int, error) {
	// ensure that the hash is correctlyformatted
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

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockTransactionCountByHash", []string{hash})

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// GetBlockTransactionCountByNumber - Returns the number of transactions in a block matching the given block number
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getblocktransactioncountbynumber
// Parameters:
//    - QUANTITY|TAG - integer of a block number, or the string "earliest", "latest" or "pending", as in the default block parameter
// Returns:
//    - QUANTITY - integer of the number of transactions in this block
func (core *Core) GetBlockTransactionCountByNumber(defaultBlockParameter string) (*big.Int, error) {

	params := make([]string, 1)
	params[0] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockTransactionCountByNumber", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

// GetBlockByHash - Returns information about a block by hash.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getblockbyhash
// Parameters:
//    - DATA, 32 bytes - Hash of a block
//    - transactionDetails, bool - indicate if we should have or not the details of the transactions of the block
// Returns:
//    1. Object - A block object, or null when no transaction was found
//    2. error
func (core *Core) GetBlockByHash(hash string, transactionDetails bool) (interface{}, error) {
	// ensure that the hash is correctlyformatted
	if strings.HasPrefix(hash, "0x") {
		if len(hash) != 66 {
			return nil, errors.New("malformed block hash")
		}
	} else {
		hash = "0x" + hash
		if len(hash) != 62 {
			return nil, errors.New("malformed block hash")
		}
	}

	params := make([]interface{}, 2)
	params[0] = hash
	params[1] = transactionDetails

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getBlockByHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToBlock(transactionDetails)
}

// GetCode - Returns code at a given address
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#core_getcode
// Parameters:
//    - DATA, 20 Bytes - address
//	  - QUANTITY|TAG - integer block number, or the string "latest", "earliest" or "pending", see the default block parameter: https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
// Returns:
//    - DATA - the code from the given address.
func (core *Core) GetCode(address string, defaultBlockParameter string) (string, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getCode", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

func (core *Core) GetAllCandidates() ([]*dto.CandidateResponse, error) {

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getAllCandidates", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToCandidatesResponse()
}

func (core *Core) GetVoter(address string) (*dto.VoterResponse, error) {

	params := make([]string, 1)
	params[0] = address

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getVoter", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToVoterResponse()
}

func (core *Core) GetStake(address string) (*dto.StakeResponse, error) {

	params := make([]string, 1)
	params[0] = address

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getStake", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToStakeResponse()
}

func (core *Core) QueryPeerCertificate(public string) (*dto.PeerCertificate, error) {

	params := make([]string, 1)
	params[0] = public

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_queryPeerCertificate", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToPeerCertificateResponse()
}

func (core *Core) GetCanTrust(address, defaultBlockParameter string) (*big.Int, error) {

	params := make([]string, 2)
	params[0] = address
	params[1] = defaultBlockParameter

	pointer := &dto.RequestResult{}

	err := core.provider.SendRequest(pointer, "core_getCanTrust", params)

	if err != nil {
		return nil, err
	}

	if strings.EqualFold(pointer.Data, "") {
		return common.Big0, nil
	}

	return pointer.ToBigInt()
}
