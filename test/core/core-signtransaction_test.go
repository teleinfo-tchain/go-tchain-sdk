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
	"github.com/bif/bif-sdk-go/common"
	"testing"

	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"math/big"
)

func TestCoreSignTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:44002", 10, false))

	add := common.StringToAddress("did:bid:73890cf407f6c883e9a42735")
	transaction := new(dto.TransactionParameters)
	transaction.Nonce = big.NewInt(2)
	transaction.From = "0x6469643a6269643a73890cf407f6c883e9a42735"
	transaction.To = "0x6469643a6269643a73890cf407f6c883e9a42735"
	transaction.Value = big.NewInt(100000)
	transaction.Gas = big.NewInt(50000)
	transaction.GasPrice = big.NewInt(1)
	transaction.Data = ""

	txID, err := connection.Core.SignTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("to:", txID.Transaction.To)
	fmt.Println("add:", common.Bytes2Hex(add.Bytes()))

	if txID.Transaction.Nonce.Cmp(transaction.Nonce) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", 5, txID.Transaction.Nonce.Uint64()))
		t.FailNow()
	}

	if txID.Transaction.To != "0x6469643a6269643a73890cf407f6c883e9a42735" {
		t.Errorf(fmt.Sprintf("Expected %s | Got: %s", add.Hex(), txID.Transaction.To))
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
