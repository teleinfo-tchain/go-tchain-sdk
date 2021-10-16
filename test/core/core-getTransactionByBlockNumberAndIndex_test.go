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
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

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
	transaction.GasLimit = uint64(40000)

	txID, err := connection.Core.SendTransaction(transaction)

	t.Log("Tx Submitted: ", txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  如果交易没有执行，则用已有的交易hash测试
	// txID := "0x1abaf67025a43fd5d3ecc19d3b67003ee1caf863aa642fca5248c153ca7ea5fc"
	var txFromHash *dto.TransactionResponse
	for {
		time.Sleep(time.Second*2)
		txFromHash, err = connection.Core.GetTransactionByHash(txID)
		if txFromHash!=nil && fmt.Sprintf("%d", txFromHash.BlockNumber) != "0"{
			break
		}
	}
	tx, err := connection.Core.GetTransactionByBlockNumberAndIndex(hexutil.EncodeBig(txFromHash.BlockNumber), txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if tx.From != generator || tx.To != generator || tx.Value.Cmp(txVal) != 0 || tx.Hash != txID {
		t.Errorf("Incorrect transaction from hash and index")
		t.FailNow()
	}
}
