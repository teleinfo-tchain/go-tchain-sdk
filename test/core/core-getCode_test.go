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
	"encoding/json"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"testing"
	"time"

	"github.com/bif/bif-sdk-go/core/block"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
)

func TestCoreGetCode(t *testing.T) {

	content, err := ioutil.ReadFile("../resources/simple-token.json")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	type TruffleContract struct {
		Abi              string `json:"abi"`
		ByteCode         string `json:"byteCode"`
		DeployedByteCode string `json:"deployedByteCode"`
	}

	var unmarshalResponse TruffleContract

	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	byteCode := unmarshalResponse.ByteCode
	deployedByteCode := unmarshalResponse.DeployedByteCode

	contract, err := connection.Core.NewContract(unmarshalResponse.Abi)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	generator, err := connection.Core.GetGenerator()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction.Sender = generator
	transaction.GasLimit = uint64(4000000)
	hash, err := contract.Deploy(transaction, byteCode)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
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

	if deployedByteCode != code {
		t.Error("Contract code not expected")
		t.FailNow()
	}

	t.Log("code is ", code)
}
