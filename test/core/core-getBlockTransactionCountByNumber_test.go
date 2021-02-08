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
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"testing"
	"time"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
)

func TestGetBlockTransactionCountByNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	// submit a transaction, wait for the block and there should be 1 tx.
	generator, err := connection.Core.GetGenerator()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.Sender = generator
	transaction.Recipient = generator
	transaction.Amount = big.NewInt(200000)
	transaction.GasLimit = uint64(40000)

	txID, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tx, err := connection.Core.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for tx.BlockNumber.Uint64() == 0 {
		time.Sleep(time.Second*2)
		tx, _ = connection.Core.GetTransactionByHash(txID)
	}

	blockNumber := hexutil.EncodeBig(tx.BlockNumber)

	txCount, err := connection.Core.GetBlockTransactionCountByNumber(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}

	txCount, err = connection.Core.GetBlockTransactionCountByNumber(block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}
}
