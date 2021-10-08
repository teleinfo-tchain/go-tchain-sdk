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
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"testing"
)

func TestCoreSignTransaction(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	transaction := new(dto.TransactionParameters)
	transaction.AccountNonce = uint64(2)
	transaction.Sender = resources.TestAddr
	transaction.Recipient = resources.NewAddrE
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(5), big.NewInt(1e17))
	transaction.GasLimit = uint64(50000)
	transaction.GasPrice = big.NewInt(1)
	transaction.Payload = "Sign Transfer bif test"

	txID, err := connection.Core.SignTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txID.Raw)
	// t.Logf("%#v \n", txID.Transaction)

	toAddress := utils.StringToAddress(transaction.Recipient)
	if txID.Transaction.To != hexutil.Encode(toAddress[:]) {
		t.Errorf(fmt.Sprintf("Expected %s | Got: %s", transaction.Recipient, txID.Transaction.To))
		t.FailNow()
	}

	if txID.Transaction.Value.Cmp(transaction.Amount) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Amount.Uint64(), txID.Transaction.Value.Uint64()))
		t.FailNow()
	}

	//if txID.Transaction.GasLimit.Cmp(transaction.GasLimit) != 0 {
	//	t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.GasLimit.Uint64(), txID.Transaction.GasLimit.Uint64()))
	//	t.FailNow()
	//}
	if txID.Transaction.GasPrice.Cmp(transaction.GasPrice) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.GasPrice.Uint64(), txID.Transaction.GasPrice.Uint64()))
		t.FailNow()
	}
}
