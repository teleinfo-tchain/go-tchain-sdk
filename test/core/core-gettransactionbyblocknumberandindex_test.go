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
 * @file core-sendtransaction_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */
package test

import (
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"math/big"
)

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	coinbase, err := connection.Core.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinbase
	transaction.To = coinbase
	transaction.Value = big.NewInt(0).Mul(big.NewInt(500), big.NewInt(1e18))
	transaction.Gas = big.NewInt(40000)
	transaction.Data = types.ComplexString("p2p transaction")

	//txID, err := connection.Core.SendTransaction(transaction)

	//t.Log(txID)

	blockNumber, err := connection.Core.GetBlockNumber()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	tx, err := connection.Core.GetTransactionByBlockNumberAndIndex(blockNumber, big.NewInt(0))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx.Hash)
	t.Log(tx.BlockHash)
	t.Log(tx.BlockNumber)
	t.Log(tx.TransactionIndex)
}
