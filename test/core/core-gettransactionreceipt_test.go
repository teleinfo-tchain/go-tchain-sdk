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

/**
 * @file core-gettransactionreceipt.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2018
 */

package test

import (
	"encoding/json"
	web3 "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestCoreGetTransactionReceipt(t *testing.T) {

	content, err := ioutil.ReadFile("../resources/simple-token.json")

	type TruffleContract struct {
		Abi      string `json:"abi"`
		Bytecode string `json:"bytecode"`
	}

	var unmarshalResponse TruffleContract

	json.Unmarshal(content, &unmarshalResponse)

	var connection = web3.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))
	bytecode := unmarshalResponse.Bytecode
	contract, err := connection.Core.NewContract(unmarshalResponse.Abi)

	transaction := new(dto.TransactionParameters)
	coinbase, err := connection.Core.GetCoinbase()
	transaction.From = coinbase
	transaction.Gas = big.NewInt(4000000)

	hash, err := contract.Deploy(transaction, bytecode, nil)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(receipt.ContractAddress) == 0 {
		t.Error("No contract address")
		t.FailNow()
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
		t.Error("No logs")
		t.FailNow()
	}

	if !receipt.Status {
		t.Error("False status")
		t.FailNow()
	}

}