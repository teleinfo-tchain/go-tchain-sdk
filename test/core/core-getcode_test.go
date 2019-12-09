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
 * @file core-getcode_test.go
 * @authors:
 *   Junjie Chen <chuckjunjchen@gmail.com>
 * @date 2018
 */

package test

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/bif/bif-sdk-go/core/block"

	bif "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
)

func TestCoreGetcode(t *testing.T) {

	content, err := ioutil.ReadFile("../resources/simple-token.json")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	type TruffleContract struct {
		Abi              string `json:"abi"`
		Bytecode         string `json:"bytecode"`
		DeployedBytecode string `json:"deployedBytecode"`
	}

	var unmarshalResponse TruffleContract

	json.Unmarshal(content, &unmarshalResponse)

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))
	bytecode := unmarshalResponse.Bytecode
	deployedBytecode := unmarshalResponse.DeployedBytecode

	contract, err := connection.Core.NewContract(unmarshalResponse.Abi)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	coinbase, err := connection.Core.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

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

	address := receipt.ContractAddress
	code, err := connection.Core.GetCode(address, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if deployedBytecode != code {
		t.Error("Contract code not expected")
		t.FailNow()
	}
}
