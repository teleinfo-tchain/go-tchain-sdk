package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"reflect"
	"strings"
)

type System struct {
	provider providers.ProviderInterface
}

func NewSystem(provider providers.ProviderInterface) *System {
	system := new(System)
	system.provider = provider
	return system
}

func (sys *System) PrePareTransaction(from common.Address, systemContract string, data types.ComplexString) *dto.TransactionParameters{
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = systemContract
	transaction.Data = data
	return transaction
}

func (sys *System) Call(transaction *dto.TransactionParameters) (*dto.RequestResult, error) {

	params := make([]interface{}, 2)
	params[0] = transaction.Transform()
	params[1] = block.LATEST

	pointer := &dto.RequestResult{}

	err := sys.provider.SendRequest(&pointer, "core_call", params)

	if err != nil {
		return nil, err
	}

	return pointer, err

}

// convert a struct to interface
func (sys *System) StructToInterface(convert interface{}, values []interface{}) []interface{}{
	elem := reflect.ValueOf(convert)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Kind() == reflect.Struct{
			values = append(sys.StructToInterface(field.Interface(), values))
		} else {
			values = append(values, field.Interface())
		}
	}
	return values
}

func (sys *System) SendTransaction(transaction *dto.TransactionParameters) (string, error) {

	params := make([]*dto.RequestTransactionParameters, 1)
	params[0] = transaction.Transform()

	pointer := &dto.RequestResult{}

	err := sys.provider.SendRequest(&pointer, "core_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

type LogData struct {
	Method string
	Status bool
	Result string
}

// 解析交易hash的Log data，这个后续应该简化，实现快速解析，而不是调用abi
func DecodeTxHash(transactionHash string) (*LogData, error){
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

func (sys *System) SystemLogDecode(transactionHash string) (*LogData, error){
	params := make([]string, 1)
	params[0] = transactionHash

	pointer := &dto.RequestResult{}

	err := sys.provider.SendRequest(pointer, "core_getTransactionReceipt", params)

	if err != nil {
		return nil, err
	}

	receipt,err :=  pointer.ToTransactionReceipt()
	if err != nil {
		return nil, err
	}

	return DecodeTxHash(receipt.Logs[0].Data)
}
