package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"strconv"
	"testing"
)

func TestGetRawTransactionByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	// generator, err := connection.Core.GetGenerator()
	//
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	//
	// transaction := new(dto.TransactionParameters)
	// transaction.Sender = generator
	// transaction.Recipient = generator
	// transaction.Amount = big.NewInt(10)
	// transaction.GasPrice = big.NewInt(40000)
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
