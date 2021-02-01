/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
	"time"
)

func TestCoreGetTransactionReceipt(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	toAddress := resources.NewAddrE

	transaction := new(dto.TransactionParameters)
	transaction.Sender = generator
	transaction.Recipient = toAddress
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e17))
	transaction.GasLimit = big.NewInt(40000)
	transaction.Payload = "Transfer test"

	txID, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt
	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txID)
	}
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	if len(receipt.ContractAddress) == 0 {
		t.Log("No contract address")
	}

	if len(receipt.TransactionHash) == 0 {
		t.Error("No transaction hash")
		t.FailNow()
	}

	if receipt.TransactionIndex == nil {
		t.Error("No transaction index")
		t.FailNow()
	}

	if len(receipt.BlockHash) == 0 {
		t.Error("No block hash")
		t.FailNow()
	}

	if receipt.BlockNumber == nil || receipt.BlockNumber.Cmp(big.NewInt(0)) == 0 {
		t.Error("No block number")
		t.FailNow()
	}

	if receipt.Logs == nil || len(receipt.Logs) == 0 {
		t.Log("No logs")
	}

	if !receipt.Status {
		t.Error("False status")
		t.FailNow()
	}
}
