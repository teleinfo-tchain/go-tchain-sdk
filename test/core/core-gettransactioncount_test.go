/********************************************************************************
   This file is part of go-bif.
   go-bif is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-bif is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-bif.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

/**
 * @file core-gettransactioncount_test.go
 * @authors:
 * 		Sigma Prime <sigmaprime.io>
 * @date 2017
 */

package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"math/big"
	"testing"
	"time"
)

func TestCoreGetTransactionCount(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("172.20.3.21:44032", 10, false))

	coinbase, _ := connection.Core.GetCoinbase()

	count, err := connection.Core.GetTransactionCount(coinbase, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	countTwo, err := connection.Core.GetTransactionCount(coinbase, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// count should not change
	if count.Cmp(countTwo) != 0 {
		t.Errorf("Count incorrect, changed between calls")
		t.FailNow()
	}
	// send a transaction and the count should increase

	t.Log("Starting Count:", count)
	transaction := new(dto.TransactionParameters)
	transaction.From = coinbase
	transaction.To = coinbase
	transaction.Value = big.NewInt(0).Mul(big.NewInt(500), big.NewInt(1e18))
	transaction.Gas = big.NewInt(40000)
	transaction.Data = types.ComplexString("p2p transaction")

	_, err = connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Errorf("Failed to send tx")
		t.FailNow()
	}

	time.Sleep(time.Second)

	newCount, err := connection.Core.GetTransactionCount(coinbase, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if newCount.Int64() != (countTwo.Int64() + 1) {
		t.Errorf(fmt.Sprintf("Incorrect count retrieved; [Expected %d | Got %d]", countTwo.Int64()+1, newCount))
		t.FailNow()
	}

	t.Log("Final Count: ", newCount)
}
