package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestGetRawTransactionByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	// coinBase, err := connection.Core.GetCoinBase()
	//
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	//
	// transaction := new(dto.TransactionParameters)
	// transaction.From = coinBase
	// transaction.To = coinBase
	// transaction.Value = big.NewInt(10)
	// transaction.Gas = big.NewInt(40000)
	//
	// txID, err := connection.Core.SendTransaction(transaction)
	//
	// t.Log("txID:", txID)
	//
	// if err != nil {
	// 	t.Errorf("Failed SendTransaction")
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	//
	// // Wait for a block
	// time.Sleep(time.Second)
	txID := "0x4300bb5bb76f34c18e51eaa2d376fad9acc03f5b02ae64032ced3a4dd310852c"
	tx, err := connection.Core.GetRawTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx)
}
