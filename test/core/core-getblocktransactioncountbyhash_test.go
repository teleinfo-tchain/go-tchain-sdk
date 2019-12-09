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
 * @file core-getblocktransactioncountbyhash_test.go
 * @authors:
 *   Sigma Prime <sigmaprime.io>
 * @date 2018
 */

package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"math/big"
	"testing"
	"time"
)

func TestGetBlockTransactionCountByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	blockNumber, err := connection.Core.GetBlockNumber()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	block, err := connection.Core.GetBlockByNumber(blockNumber, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txCount, err := connection.Core.GetBlockTransactionCountByHash(block.Hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

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

	txCount, err = connection.Core.GetBlockTransactionCountByHash(tx.BlockHash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount.Int64() != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}

}
