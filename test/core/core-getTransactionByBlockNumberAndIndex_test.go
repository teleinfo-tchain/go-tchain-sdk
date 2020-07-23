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
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
	"time"
)

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	coinBase, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txVal := big.NewInt(2000000)

	transaction := new(dto.TransactionParameters)
	transaction.From = coinBase
	transaction.To = coinBase
	transaction.Value = big.NewInt(0).Mul(big.NewInt(500), big.NewInt(1E18))
	transaction.Value = txVal
	transaction.Gas = big.NewInt(40000)

	txID, err := connection.Core.SendTransaction(transaction)

	t.Log("Tx Submitted: ", txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  wait for a block
	time.Sleep(time.Second)

	//  如果交易没有执行，则用已有的交易hash测试
	// txID := "0x99e87282b319c80682ad9cdd621b726bd1ee6c583b0cef0e9a05ba6d2f0fbc21"
	txFromHash, err := connection.Core.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// if it fails, it may be that the time is too short and the transaction has not been executed
	tx, err := connection.Core.GetTransactionByBlockNumberAndIndex(hexutil.EncodeBig(txFromHash.BlockNumber), txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if tx.From != coinBase || tx.To != coinBase || tx.Value.Cmp(txVal) != 0 || tx.Hash != txID {
		t.Errorf("Incorrect transaction from hash and index")
		t.FailNow()
	}
}
