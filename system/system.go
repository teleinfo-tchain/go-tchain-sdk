package system

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/keystore"
	"math/big"
	"reflect"
	"strings"
)

// System - The System Module
type System struct {
	provider providers.ProviderInterface
}

type LogData struct {
	Method string
	Status bool
	Result string
}

// 系统合约交易构建参数
type SysTxParams struct {
	From        common.Address // 系统合约发起方账户地址
	Password    string         // 解密私钥的密码
	KeyFileData []byte         // keystore文件内容
	PrivateKey  string         // 发起方账户对应的私钥
	GasPrice    *big.Int       // 交易的gas价格，默认是网络gas价格的平均值
	Gas         uint64         // 交易可使用的gas，未使用的gas会退回
	Nonce       uint64         // 从该账户发起交易的Nonce值
	ChainId     *big.Int       // 链的ChainId
}

// from: string，20 Bytes - 指定的发送者的地址。
//		to: string，20 Bytes - （可选）交易消息的目标地址，如果是合约创建，则不填.
//		gas: *big.Int - （可选）默认是自动，交易可使用的gas，未使用的gas会退回。
//		gasPrice: *big.Int - （可选）默认是自动确定，交易的gas价格，默认是网络gas价格的平均值 。
//		data: string - （可选）或者包含相关数据的字节字符串，如果是合约创建，则是初始化要用到的代码。
//		value: *big.Int - （可选）交易携带的货币量，以bifer为单位。如果合约创建交易，则为初始的基金
//		nonce: *big.Int - （可选）整数，使用此值，可以允许你覆盖你自己的相同nonce的，待pending中的交易

// NewSystem - System Module constructor to set the default provider
func NewSystem(provider providers.ProviderInterface) *System {
	system := new(System)
	system.provider = provider
	return system
}

/*
	prePareTransaction - Construct transaction

	prePareTransaction - 构造交易
*/
func (sys *System) prePareTransaction(from common.Address, systemContract string, data types.ComplexString) *dto.TransactionParameters {
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = systemContract
	transaction.Data = data
	return transaction
}

func getPrivateKeyFromFile(addrParse common.Address, keyJson []byte, password string) (string, error) {
	var key *keystore.Key
	var err error

	if bytes.HasPrefix(addrParse.Bytes(), []byte("did:bid:")) && addrParse[8] == 115 {
		key, err = keystore.DecryptKey(keyJson, password, crypto.SM2)
	} else {
		key, err = keystore.DecryptKey(keyJson, password, crypto.SECP256K1)
	}
	if err != nil {
		return "", err
	}
	privateKey := hex.EncodeToString(key.PrivateKey.D.Bytes())
	addrRes := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)
	if addrParse != addrRes {
		return "", errors.New("addrParse Not Match keyStoreFile")
	}
	return privateKey, nil
}

/*
	prePareSignTransaction - Construct transaction

	prePareSignTransaction - 构造交易
*/
func (sys *System) prePareSignTransaction(signTxParams *SysTxParams, payLoad []byte, contractAddr common.Address) (string, error) {
	if signTxParams.PrivateKey == "" {
		 privateKey, err := getPrivateKeyFromFile(signTxParams.From, signTxParams.KeyFileData, signTxParams.Password)
		if err!= nil{
			return "", err
		}
		signTxParams.PrivateKey = privateKey
	}

	signTx := &utils.TxData{
		AccountNonce: signTxParams.Nonce,
		Price:        signTxParams.GasPrice,
		GasLimit:     signTxParams.Gas,
		Sender:       &signTxParams.From,
		Recipient:    &contractAddr,
		Amount:       new(big.Int),
		Payload:      payLoad,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
		T:            big.NewInt(0),
	}
	signResult, err := utils.SignTransaction(signTx, signTxParams.PrivateKey, signTxParams.ChainId)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signResult.Raw), nil
}

// structToInterface - 将结构体转换为[]interface{}
func (sys *System) structToInterface(convert interface{}, values []interface{}) []interface{} {
	elem := reflect.ValueOf(convert)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Kind() == reflect.Struct {
			values = append(sys.structToInterface(field.Interface(), values))
		} else {
			values = append(values, field.Interface())
		}
	}
	return values
}

/*
  sendRawTransaction:
   	EN - Add the signed transaction to the transaction pool.The sender is responsible for signing the transaction and using the correct nonce
 	CN - 将已签名的交易添加到交易池中。交易发送方负责签署交易并使用正确的随机数（Nonce）
  Params:
  	- encodedTx: string, 已签名的交易数据

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
 	- error

  Call permissions: Anyone
*/
func (sys *System) sendRawTransaction(encodedTx string) (string, error) {

	params := make([]string, 1)
	params[0] = encodedTx

	pointer := &dto.CoreRequestResult{}

	err := sys.provider.SendRequest(&pointer, "core_sendRawTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

func (sys *System) sendTransaction(transaction *dto.TransactionParameters) (string, error) {

	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := sys.provider.SendRequest(&pointer, "core_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

func decodeTxHash(transactionHash string) (*LogData, error) {
	unpacked := new(LogData)
	def := `[{ "name" : "log", "type": "function", "outputs": [{"name":"method","type":"string"},{"name":"status","type":"bool"},{"name":"result","type":"string"}]}]`
	ABI, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		return nil, err
	}
	err = ABI.Methods["log"].Outputs.Unpack(unpacked, common.FromHex(transactionHash))
	if err != nil {
		return nil, err
	}
	return unpacked, nil
}

// Deprecated:解析交易hash的Log data，这个后续应该简化，实现快速解析，而不是调用abi
func (sys *System) SystemLogDecode(transactionHash string) (*LogData, error) {
	params := make([]string, 1)
	params[0] = transactionHash

	pointer := &dto.RequestResult{}

	err := sys.provider.SendRequest(pointer, "core_getTransactionReceipt", params)

	if err != nil {
		return nil, err
	}

	receipt, err := pointer.ToTransactionReceipt()
	if err != nil {
		return nil, err
	}

	if len(receipt.Logs) == 0 {
		return nil, errors.New("method log error")
	}
	return decodeTxHash(receipt.Logs[0].Data)
}
