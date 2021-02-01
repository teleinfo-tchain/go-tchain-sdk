package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"testing"
	"time"
)

func TestGetRawTransactionByBlockNumberAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	generator, err := connection.Core.GetGenerator()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txVal := big.NewInt(2000000)

	transaction := new(dto.TransactionParameters)
	transaction.Sender = generator
	transaction.Recipient = generator
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(500), big.NewInt(1E18))
	transaction.Amount = txVal
	transaction.GasLimit = big.NewInt(40000)

	txID, err := connection.Core.SendTransaction(transaction)

	t.Log("Tx Submitted: ", txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  wait for a block
	// time.Sleep(time.Second*10)
	time.Sleep(time.Second)

	txFromHash, err := connection.Core.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// if it fails, it may be that the time is too short and the transaction has not been executed
	tx, err := connection.Core.GetRawTransactionByBlockNumberAndIndex(hexutil.EncodeBig(txFromHash.BlockNumber), txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx)
}
