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
// Core - Core ģ��
type Core struct {
	provider providers.ProviderInterface
}

// NewCore - Core Module constructor to set the default provider
// NewCore - Core Module ���캯������ʼ��
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
 	CN - ���ص�ǰ��BifЭ��汾��
  Params:
  	- None

  Returns:
  	- string, ��ǰ��BifЭ��汾
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
 	CN - ����ͬ��״̬�����ݶ���
  Params:
  	- None

  Returns:
  	- *dto.SyncingResponse, ���û��ͬ�����򷵻�&{<nil> <nil> <nil>};�������ͬ���򷵻أ�
 		StartingBlock *big.Int - ͬ��ʱ���������ʼ����
		CurrentBlock  *big.Int - ��ǰ���飬��GetBlockNumberЧ����ͬ
		HighestBlock  *big.Int - ��ǰ���Ƶ��������
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
  GetCoinBase:
   	EN - Returns the client coinbase address
 	CN - ���ؿͻ��˵ĳ��齱����ַ
  Params:
  	- None

  Returns:
  	- string, 20 bytes���ַ���, ���齱����ַ
 	- error

  Call permissions: Anyone
*/
func (core *Core) GetCoinBase() (string, error) {

	pointer := &dto.CoreRequestResult{}

	err := core.provider.SendRequest(pointer, "core_coinbase", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  Generating:
   	EN - Returns true if client is actively generating new blocks.
 	CN - ����ͻ������ڻ����ھ��¿飬�򷵻�true
  Params:
  	- None

  Returns:
  	- bool, true ���ڳ��飻false ����ֹͣ
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
 	CN - ���ؽڵ�ÿ���ھ�Ĺ�ϣ��
  Params:
  	- None

  Returns:
  	- *big.Int, ÿ���ϣ��
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
 	CN - ������bifΪ��λ�ĵ�ǰGasPrice
  Params:
  	- None

  Returns:
  	- *big.Int, ��bifΪ��λ�ĵ�ǰgasPrice
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
 	CN - ���ص�ǰ�ͻ������е��˻���ַ�б�
  Params:
  	- None

  Returns:
  	- []string, ��ǰ�ͻ���ӵ�е��˻���ַ�б�
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
 	CN - ���ص�ǰ��������������
  Params:
  	- None

  Returns:
  	- *big.Int, ��ǰ��������������
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
 	CN - ���ظ�����ַ���ʻ���
  Params:
  	- address, string, 20 bytes
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- *big.Int�� ������ַ�����
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
 	CN - ������ָ��������£�������ַ�ѷ��͵Ľ�������
  Params:
	- address, string, 20 bytes
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- *big.Int, ָ��������£�������ַ�ѷ��͵Ľ�������
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
 	CN - ����ִ�иý�����������gas�Ĺ���ֵ���ý��ײ��ᱻ��ӵ���������
  Params:
  	- transaction, *dto.TransactionParameters, ��ϸ����Call�в���������һ�¡�
		���δָ��gas����ʹ��pending�����е�gasֵ�����ִ�н��������gas�������ƣ��򷵻ص�gas����ֵ���ܲ�����ִ�н���

  Returns:
  	- *big.Int, gas���ĵ���
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
 	CN - �Ը��������������ף��������ǩ���������ύ�����׳أ����ؽ��׹�ϣ
  Params:
  	- transaction: Ҫ���͵Ľ��׶���(*dto.TransactionParameters)
		from: string��20 Bytes - ָ���ķ����ߵĵ�ַ��
		to: string��20 Bytes - ����ѡ��������Ϣ��Ŀ���ַ������Ǻ�Լ����������.
		gas: *big.Int - ����ѡ��Ĭ�����Զ������׿�ʹ�õ�gas��δʹ�õ�gas���˻ء�
		gasPrice: *big.Int - ����ѡ��Ĭ�����Զ�ȷ�������׵�gas�۸�Ĭ��������gas�۸��ƽ��ֵ ��
		data: string - ����ѡ�����߰���������ݵ��ֽ��ַ���������Ǻ�Լ���������ǳ�ʼ��Ҫ�õ��Ĵ��롣
		value: *big.Int - ����ѡ������Я���Ļ���������biferΪ��λ�������Լ�������ף���Ϊ��ʼ�Ļ���
		nonce: *big.Int - ����ѡ��������ʹ�ô�ֵ�����������㸲�����Լ�����ͬnonce�ģ���pending�еĽ���

  Returns:
  	- string, transactionHash��32 Bytes�����׹�ϣ����������в����ã���Ϊ���ϣ
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
 	CN - ����ǩ���Ľ�����ӵ����׳��С����׷��ͷ�����ǩ���ײ�ʹ����ȷ���������Nonce��
  Params:
  	- encodedTx: string, ��ǩ���Ľ�������

  Returns:
  	- string, transactionHash��32 Bytes�����׹�ϣ����������в����ã���Ϊ���ϣ
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
 	CN - ʹ�ý��׷��𷽵��ʻ�����ַ��ǩ������Ľ���
  Params:
  	- transaction��*dto.TransactionParameters�����׹���Ķ���
	  - From     string                    ���׵ķ���
	  - To       string                    ���׵Ľ��շ�
	  - Nonce    *big.Int                  ����ѡ��������ʹ�ô�ֵ�����������㸲�����Լ�����ͬnonce�ģ���pending�еĽ���
	  - Gas      *big.Int                  ����ѡ��Ĭ�����Զ������׿�ʹ�õ�gas��δʹ�õ�gas���˻ء�
	  - GasPrice *big.Int                  ����ѡ��Ĭ�����Զ�ȷ�������׵�gas�۸�Ĭ��������gas�۸��ƽ��ֵ ��
	  - Value    *big.Int                  ����ѡ������Я���Ļ���������biferΪ��λ�������Լ�������ף���Ϊ��ʼ�Ļ���
	  - Data     types.ComplexString       ����ѡ�����߰���������ݵ��ֽ��ַ���������Ǻ�Լ���������ǳ�ʼ��Ҫ�õ��Ĵ��롣

  Returns:
  	- *dto.SignTransactionResponse��
		Raw         string                   ��ǩ����RLP����Ľ���
		Transaction SignedTransactionParams  transaction object
		  - Gas      *big.Int                ���׷���Լ����gas
		  - GasPrice *big.Int                ���׷���Լ����gasPrice
		  - Hash     string  			     ���׹�ϣ
		  - Input    string   				 �潻�׷��͵�����
		  - Nonce    *big.Int                ���׷�����֮ǰ�����׵Ĵ���
		  - S        string                  ������
		  - R        string                  ������
		  - V        *big.Int                ������
		  - To       string                  ���׵Ľ��շ�������Ǻ�Լ������Ϊ��
		  - Value    *big.Int                ת�Ƶ�bif����
 	- error

  Call permissions: ���׵ķ��𷽵��˻����ڽ���״̬

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
 	CN - ִ���µ���Ϣ���ã����������������ϴ������ף�������ı���������״̬��һ�����ڼ�����
  Params:
	- transaction��*dto.TransactionParameters������Call�Ķ���
 	  - From     string                    ���׵ķ���
 	  - To       string                    ���׵Ľ��շ�
 	  - Nonce    *big.Int                  ����ѡ��������ʹ�ô�ֵ�����������㸲�����Լ�����ͬnonce�ģ���pending�еĽ���
 	  - Gas      *big.Int                  ����ѡ��Ĭ�����Զ������׿�ʹ�õ�gas��δʹ�õ�gas���˻ء�
 	  - GasPrice *big.Int                  ����ѡ��Ĭ�����Զ�ȷ�������׵�gas�۸�Ĭ��������gas�۸��ƽ��ֵ ��
 	  - Value    *big.Int                  ����ѡ������Я���Ļ���������biferΪ��λ�������Լ�������ף���Ϊ��ʼ�Ļ���
 	  - Data     types.ComplexString       ����ѡ�����߰���������ݵ��ֽ��ַ���������Ǻ�Լ���������ǳ�ʼ��Ҫ�õ��Ĵ��롣

  Returns:
  	- ��ִ�к�Լ�ķ���ֵ
 	- error

  Call permissions: Anyone
  Bug �����ԣ���Ҫ�ȶ�rpc�е�callArgs��sendTxArgs���������������ݽṹ
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
 	CN - ���ظ������׹�ϣ�Ľ����վ�
  Params:
  	- hash,string 32 Bytes ���׹�ϣ

  Returns:
  	- *dto.TransactionReceipt �����վݶ���
	  - TransactionHash   string            ���׹�ϣ
	  - TransactionIndex  *big.Int          �����������е�������λ�ã�
	  - BlockHash         string            �����ϣ
	  - BlockNumber       *big.Int          ����
	  - From              string            ���׷���
	  - To                string            ���׽��շ�
	  - CumulativeGasUsed *big.Int          ��������ִ�д˽���ʱʹ�õ�gas������
	  - GasUsed           *big.Int          �����ض�����ʹ�õ�gas����
	  - ContractAddress   string            ��������Ǻ�Լ��������Ϊ�����ĺ�Լ��ַ������Ϊ��
	  - Logs              []TransactionLogs �������ɵ���־��������
	  - LogsBloom         string            ������
	  - Status            bool              ������
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
 	CN - ���������ϣ��ȡ�������ڰ����Ľ�����
  Params:
  	- hash��32 bytes - block hash

  Returns:
  	- uint64, ��������Ľ�����
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
 	CN - ���������������ƥ��������еĽ�������
  Params:
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- uint64, ��������Ľ�����
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
 	CN - ��������ŷ����������Ϣ
  Params:
	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions
  	- transactionDetails,bool, ���ΪTrue��������������ϸ�Ľ�����Ϣ��������Ϣ�����Ϊfalse������������ڽ���hash��������Ϣ

  Returns:
  	- interface{}, ���transactionDetailsΪtrue������*dto.BlockDetails�����Ϊfalse������*dto.BlockNoDetails
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
 	CN - ���������ϣ�����������Ϣ
  Params:
	- blockHash, string32 bytes - Hash of a block
  	- transactionDetails,bool, ���ΪTrue��������������ϸ�Ľ�����Ϣ��������Ϣ�����Ϊfalse������������ڽ���hash��������Ϣ

  Returns:
  	- interface{}, ���transactionDetailsΪtrue������*dto.BlockDetails�����Ϊfalse������*dto.BlockNoDetails
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
 	CN - �����Ը�����Ŵ洢�ڸ�����ַ�Ĵ��롣����
  Params:
  	- address,string 20 Bytes �˻���ַ
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- string��the code
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
 	CN - ����˻��ж��ٿ���֤��
  Params:
  	- address, string �˻���ַ

  Returns:
  	- uint64, ����֤������
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
 	CN - ���ص�ǰ���ӽڵ����ID
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
  ����:
   	EN - Returns the Merkle-proof for a given account and optionally some storage keys.
 	CN - ���ظ����ʻ���Merkle֤��  ????
  Params:
  	- address, string, 20 bytes
	- storageKeys,  []string, һ��storageKeys��Ӧ�������У�鲢��������
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
 	CN - �Ӹ�����ַ��λ�ú�����ŵ�״̬���ش洢ֵ???(��ʲô�洢ֵ)
  Params:
	- address, string, 20 bytes
	- key,  *big.Int, ָ����λ��
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
 	CN - ����δִ�н��׵��б�(���׻�δ�����)
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
 	CN - ���ݽ���hash���ؽ��׵���ϸ��Ϣ
  Params:
  	- hash,string,32 Bytes,hash of a transaction

  Returns:
  	- *dto.TransactionResponse
		Hash             string              ���׹�ϣ
		Nonce            *big.Int            ���׷��ͷ��Ӹýڵ㷢�͵Ľ�����
		BlockHash        string              �������ڵ�����Ĺ�ϣ��������״�������Ϊ0x0000000000000000000000000000000000000000000000000000000000000000
		BlockNumber      *big.Int            �������ڵ����飬������״�������Ϊ0
		TransactionIndex *big.Int            �����������е�λ�ã���������������״�������Ϊ0
		From             string              ���׷��ͷ�
		To               string              ���׽��շ������Ϊ��Լ������Ϊ��
		Input            string              �潻�׷��͵�����
		Value            *big.Int            ����ת�Ƶ�bif��������λΪbif
		GasPrice         *big.Int            ���ͷ��ṩ��GasPrice����λΪbif
		Gas              *big.Int            ���ͷ��ṩ��Gas
		Data             types.ComplexString ��������
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
 	CN - ���ݽ���hash���ؽ�����Ϣ
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
 	CN - ���������ϣ�ͽ������������ؽ�����Ϣ
  Params:
  	- hash,string,32 Bytes,�����ϣ
  	- index, *big.Int, �����������е�����

  Returns:
  	- *dto.TransactionResponse, ����GetTransactionByHash�ķ���ֵ������һ��
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
 	CN - ���������ϣ�ͽ������������ؽ�����Ϣ
  Params:
  	- hash,string,32 Bytes,�����ϣ
  	- index, *big.Int, �����������е�����

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
 	CN - ���ظ�������ź������Ľ���
  Params:
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions
  	- index, *big.Int, �����������е�����

  Returns:
  	- *dto.TransactionResponse, ����GetTransactionByHash�ķ���ֵ������һ��
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
 	CN - ���ظ�������ź������Ľ���
  Params:
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions
  	- index, *big.Int, �����������е�����

  Returns:
  	- string
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
