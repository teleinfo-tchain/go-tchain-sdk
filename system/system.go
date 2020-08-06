package system

import (
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"reflect"
	"strings"
)

// System - The System Module
type System struct {
	provider providers.ProviderInterface
	acc      *account.Account
}

type LogData struct {
	Method string
	Status bool
	Result string
}

// 系统合约交易构建参数
// TODO： 在使用此本地签署交易时，注意签署的内容是否要增加，其参数用于prePareSignTransaction，涉及account.TxData中的交易构建！！！（后续可能会增加）
type SysTxParams struct {
	IsSM2       bool     // 私钥生成是否使用国密，true为国密；false为非国密
	Password    string   // 解密私钥的密码
	KeyFileData []byte   // keystore文件内容
	GasPrice    *big.Int // 交易的gas价格，默认是网络gas价格的平均值
	Gas         uint64   // 交易可使用的gas，未使用的gas会退回
	Nonce       uint64   // 从该账户发起交易的Nonce值
	ChainId     *big.Int // 链的ChainId
}

// NewSystem - System Module constructor to set the default provider
func NewSystem(provider providers.ProviderInterface) *System {
	system := new(System)
	system.provider = provider
	system.acc = account.NewAccount(provider)
	return system
}

/*
	prePareSignTransaction - Construct transaction

	prePareSignTransaction - 构造交易
*/
func (sys *System) prePareSignTransaction(signTxParams *SysTxParams, payLoad []byte, contractAddr string) (string, error) {
	_, privateKey, err := sys.acc.Decrypt(signTxParams.KeyFileData, signTxParams.IsSM2, signTxParams.Password)
	if err != nil {
		return "", err
	}

	signTx := &account.SignTxParams{
		To:       contractAddr,
		Nonce:    signTxParams.Nonce,
		Gas:      signTxParams.Gas,
		GasPrice: signTxParams.GasPrice,
		Value:    nil,
		Data:     payLoad,
		ChainId:  signTxParams.ChainId,
	}
	signResult, err := sys.acc.SignTransaction(signTx, privateKey, signTxParams.IsSM2)
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

func decodeTxHash(transactionHash string) (*LogData, error) {
	unpacked := new(LogData)
	def := `[{ "name" : "log", "type": "function", "outputs": [{"name":"method","type":"string"},{"name":"status","type":"bool"},{"name":"result","type":"string"}]}]`
	ABI, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		return nil, err
	}
	err = ABI.Methods["log"].Outputs.Unpack(unpacked, utils.FromHex(transactionHash))
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
