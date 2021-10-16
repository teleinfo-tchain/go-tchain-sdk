package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"strconv"
	"testing"
)

func TestCoreGetPendingTransactions(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))



	generator, err := connection.Core.GetGenerator()
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.Sender = generator
	transaction.Recipient = resources.Addr1
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e15))
	transaction.GasLimit = uint64(40000)
	transaction.Payload = "Transfer test"

	_, err = connection.Core.SendTransaction(transaction)
	if err != nil{
		t.FailNow()
	}
	pendingTransactions, _ := connection.Core.GetPendingTransactions()

	if len(pendingTransactions)!=0{
		t.Logf("%#v \n", pendingTransactions[0])
	}
}
