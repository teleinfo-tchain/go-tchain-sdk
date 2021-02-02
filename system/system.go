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
	"regexp"
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
	From        string   // 交易的发起方，与私钥对应的地址可以相同也可不同
	IsSM2       bool     // 私钥生成是否使用国密，true为国密；false为非国密
	Password    string   // 解密私钥的密码
	KeyFileData []byte   // keystore文件内容
	GasPrice    *big.Int // 交易的gas价格，默认是网络gas价格的平均值
	Gas         uint64   // 交易可使用的gas，未使用的gas会退回
	Nonce       uint64   // 从该账户发起交易的Nonce值
	ChainId     uint64   // 链的ChainId
	Version     uint64
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

	sender := utils.StringToAddress(signTxParams.From)
	recipient := utils.StringToAddress(contractAddr)
	signTx := &account.SignTxParams{
		Sender:       &sender,
		Recipient:    &recipient,
		AccountNonce: signTxParams.Nonce,
		GasPrice:     signTxParams.GasPrice,
		GasLimit:     signTxParams.Gas,
		Amount:       nil,
		Payload:      payLoad,
		ChainId:      signTxParams.ChainId,
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
 	CN - 将已签名的交易添加到交易池中。交易发送方负责签署交易并使用正确的随机数（AccountNonce）
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

func isValidHexAddress(address string) bool {
	return utils.StringToAddress(address).EqualString(address)
}

// email verify
func verifyEmailFormat(email string) bool {
	// pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 检测是否空白字符
func isBlankCharacter(judge string) bool {
	reg, _ := regexp.Compile(`^\s+$`)
	return reg.MatchString(judge)
}

// 检测IP是否合法
// (1~255).(0~255).(0~255).(0~255)
func isLegalIP(ip string) bool {
	// ip地址范围：(1~255).(0~255).(0~255).(0~255)
	// ipRegEx := "^([1-9]|([1-9][0-9])|(1[0-9][0-9])|(2[0-4][0-9])|(25[0-5]))(\\.([0-9]|([1-9][0-9])|(1[0-9][0-9])|(2[0-4][0-9])|(25[0-5]))){3}$"
	ipRegEx := "^([1-9]|([1-9]\\d)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))(\\.(\\d|([1-9]\\d)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))){3}$"
	// ipRegEx := "^(([1-9]\\d?)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))(\\.(0|([1-9]\\d?)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))){3}$"
	// Pattern
	reg, _ := regexp.Compile(ipRegEx)
	// Matcher
	return reg.MatchString(ip)
}
