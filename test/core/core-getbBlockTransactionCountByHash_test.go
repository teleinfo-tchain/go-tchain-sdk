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

package test

import (
	"github.com/bif/bif-sdk-go"
	block2 "github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
	"time"
)

func TestGetBlockTransactionCountByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	block, err := connection.Core.GetBlockByNumber(block2.LATEST, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	txCount, err := connection.Core.GetBlockTransactionCountByHash(block.(*dto.BlockNoDetails).Hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txCount:", txCount)

	// submit a transaction, wait for the block and there should be 1 tx.
	coinBase, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinBase
	transaction.To = coinBase
	transaction.Value = big.NewInt(200000)
	transaction.Gas = big.NewInt(40000)

	txID, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("txID:", txID)

	tx, err := connection.Core.GetTransactionByHash(txID)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for tx.BlockNumber.Uint64() == 0 {
		time.Sleep(time.Second*2)
		tx, _ = connection.Core.GetTransactionByHash(txID)
	}

	txCount, err = connection.Core.GetBlockTransactionCountByHash(tx.BlockHash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  注意如果有其他用户发送交易，可能会大于1，该测试可能会失败
	if txCount != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}
	t.Log("txCount:", txCount)
}
