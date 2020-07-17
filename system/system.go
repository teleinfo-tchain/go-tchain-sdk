package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
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
	sendTransaction - 如果数据字段包含代码，则创建新的消息调用交易或合约创建。

	return：
		transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。
		(将交易哈希出传入GetTransactionReceipt可获取交易信息。)
*/
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

	return decodeTxHash(receipt.Logs[0].Data)
}
