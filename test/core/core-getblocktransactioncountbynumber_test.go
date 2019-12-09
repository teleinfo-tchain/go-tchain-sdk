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
 * @file core-getblocktransactioncountbynumber.go
 * @authors:
 *   Junjie CHen <chuckjunjchen@gmail.com>
 * @date 2018
 */

package test

import (
	"math/big"
	"testing"
	"time"

	bif "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
)

func TestGetBlockTransactionCountByNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	// submit a transaction, wait for the block and there should be 1 tx.
	coinbase, err := connection.Core.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinbase
	transaction.To = coinbase
	transaction.Value = big.NewInt(200000)
	transaction.Gas = big.NewInt(40000)

	txID, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(time.Second)

	tx, err := connection.Core.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	blockNumber := block.NUMBER(tx.BlockNumber)

	txCount, err := connection.Core.GetBlockTransactionCountByNumber(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount.Int64() != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}

	txCount, err = connection.Core.GetBlockTransactionCountByNumber(block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount.Int64() != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}
}
