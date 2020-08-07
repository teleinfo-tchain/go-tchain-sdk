package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

func TestCoreGetPendingTransactions(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))



	coinBase, err := connection.Core.GetCoinBase()
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinBase
	transaction.To = resources.Address51
	transaction.Value = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e15))
	transaction.Gas = big.NewInt(40000)
	transaction.Data = "Transfer test"

	_, err = connection.Core.SendTransaction(transaction)
	if err != nil{
		t.FailNow()
	}
	pendingTransactions, _ := connection.Core.GetPendingTransactions()

	if len(pendingTransactions)!=0{
		t.Logf("%#v \n", pendingTransactions[0])
	}
}
