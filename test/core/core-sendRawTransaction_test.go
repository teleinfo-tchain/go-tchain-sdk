package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

//test/core/core-sendtransaction_test.go
func TestCoreSendRawTransaction(t *testing.T){
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	transaction := new(dto.TransactionParameters)
	transaction.Nonce = big.NewInt(2)
	transaction.From = resources.CoinBase
	transaction.To = resources.AddressTwo
	transaction.Value = big.NewInt(0).Mul(big.NewInt(5), big.NewInt(1e17))
	transaction.Gas = big.NewInt(50000)
	transaction.GasPrice = big.NewInt(1)
	transaction.Data = "Sing Transfer Bifer test"

	txID, err := connection.Core.SignTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(txID.Raw)
	//
	//
	//util := utils.NewUtils()
	//hexStr, _ := util.ToHex(resources.AddressTwo[:8])
	//addressTwoHex := hexStr +resources.AddressTwo[8:]
	//
	//if txID.Transaction.To != addressTwoHex {
	//	t.Errorf(fmt.Sprintf("Expected %s | Got: %s", addressTwoHex, txID.Transaction.To))
	//	t.FailNow()
	//}
	//
	//if txID.Transaction.Value.Cmp(transaction.Value) != 0 {
	//	t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Value.Uint64(), txID.Transaction.Value.Uint64()))
	//	t.FailNow()
	//}
	//
	//if txID.Transaction.Gas.Cmp(transaction.Gas) != 0 {
	//	t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Gas.Uint64(), txID.Transaction.Gas.Uint64()))
	//	t.FailNow()
	//}
	//if txID.Transaction.GasPrice.Cmp(transaction.GasPrice) != 0 {
	//	t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.GasPrice.Uint64(), txID.Transaction.GasPrice.Uint64()))
	//	t.FailNow()
	//}
	//
	//fmt.Println(txID.Raw)
	//fmt.Println(txID.Transaction.Gas)
	//fmt.Println(txID.Transaction.GasPrice)
	//fmt.Println(txID.Transaction.Hash)
	//fmt.Println(txID.Transaction.Input)
	//fmt.Println(txID.Transaction.Nonce)
	//fmt.Println(txID.Transaction.To)
	//fmt.Println(txID.Transaction.Value)
	rawTx := "0xf87b010a831e8480946469643a6269643a73f6a70d05af2141dd4ad995946469643a6269643a73890cf407f6c883e9a427358307a120808082039ea05572544e70662d7bd0b8be692db8c455ebd3ca42aef28ed5e09e43061c2d4a82a01a9a8e218a3ddfc1ed081018fa9cb78c84cffb8101265141987f4e5a0a215f48"
	txIDRaw, err := connection.Core.SendRawTransaction(rawTx)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txIDRaw)
}
