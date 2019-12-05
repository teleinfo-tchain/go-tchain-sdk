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
 * @file core-signtransaction_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */
package test

import (
	"testing"

	"fmt"
	"github.com/bif/bifGo"
	"github.com/bif/bifGo/complex/types"
	"github.com/bif/bifGo/dto"
	"github.com/bif/bifGo/providers"
	"math/big"
)

func TestCoreSignTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	coinbase, err := connection.Core.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.Nonce = big.NewInt(5)
	transaction.From = coinbase
	transaction.To = coinbase
	transaction.Value = big.NewInt(0).Mul(big.NewInt(500), big.NewInt(1e18))
	transaction.Gas = big.NewInt(40000)
	transaction.GasPrice = big.NewInt(1e9)
	transaction.Data = types.ComplexString("p2p transaction")

	txID, err := connection.Core.SignTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txID.Transaction.Nonce.Cmp(transaction.Nonce) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", 5, txID.Transaction.Nonce.Uint64()))
		t.FailNow()
	}

	if txID.Transaction.To != coinbase {
		t.Errorf(fmt.Sprintf("Expected %s | Got: %s", coinbase, txID.Transaction.To))
		t.FailNow()
	}

	if txID.Transaction.Value.Cmp(transaction.Value) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Value.Uint64(), txID.Transaction.Value.Uint64()))
		t.FailNow()
	}

	if txID.Transaction.Gas.Cmp(transaction.Gas) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Gas.Uint64(), txID.Transaction.Gas.Uint64()))
		t.FailNow()
	}
	if txID.Transaction.GasPrice.Cmp(transaction.GasPrice) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.GasPrice.Uint64(), txID.Transaction.GasPrice.Uint64()))
		t.FailNow()
	}
}
