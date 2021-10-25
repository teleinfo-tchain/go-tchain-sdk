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
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/testutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"path"
	"strconv"
	"testing"
	"time"
)

func TestCoreContract(t *testing.T) {

	cc := &testutil.Contract{}

	file := path.Join(bif.GetCurrentAbPath(), "test", "resources", "simple-token.sol")
	artifact, err := cc.Compile(file)
	assert.NoError(t, err)

	type TruffleContract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}

	var unmarshalResponse TruffleContract

	unmarshalResponse.Abi = artifact.Abi
	unmarshalResponse.ByteCode = artifact.Bin

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	byteCode := unmarshalResponse.ByteCode
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

	transaction.ChainId, err = connection.Core.GetChainId()
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
	t.Log("hash is ", hash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Contract Address: ", receipt.ContractAddress)

	transaction.Recipient = receipt.ContractAddress

	result, err := contract.Call(transaction, "getName")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		name, _ := result.ToComplexString()
		if name.ToString() != "" {
			t.Errorf(fmt.Sprintf("Name not expected; [Expected %s | Got %s]", "demo", name.ToString()))
			t.FailNow()
		}
	}

	result, err = contract.Call(transaction, "symbol")
	if result != nil && err == nil {
		symbol, _ := result.ToComplexString()
		if symbol.ToString() != "SET" {
			t.Errorf("Symbol not expected")
			t.FailNow()
		}
	}

	result, err = contract.Call(transaction, "decimals")
	if result != nil && err == nil {
		decimals, _ := result.ToBigInt()
		if decimals.Int64() != 18 {
			t.Errorf("Decimals not expected")
			t.FailNow()
		}
	}

	bigInt, _ := new(big.Int).SetString("000000000000000000000000000000000000000000000d3c21bcecceda1000000", 16)

	result, err = contract.Call(transaction, "totalSupply")
	if result != nil && err == nil {
		total, _ := result.ToBigInt()
		if total.Cmp(bigInt) != 0 {
			t.Errorf("Total not expected")
			t.FailNow()
		}
	}

	result, err = contract.Call(transaction, "balanceOf", generator)
	if result != nil && err == nil {
		balance, _ := result.ToBigInt()
		if balance.Cmp(bigInt) != 0 {
			t.Errorf("Balance not expected")
			t.FailNow()
		}
	}

	// hash, err = contract.Send(transaction, "approve", utils.StringToAddress(generator), big.NewInt(10))
	// if err != nil {
	//	t.Log(err)
	//	t.Errorf("Can't send approve transaction")
	//	t.FailNow()
	// }
	//
	// t.Log(hash)
	//
	// receipt = nil
	// for receipt == nil {
	//	time.Sleep(time.Second)
	//	receipt, err = connection.Core.GetTransactionReceipt(hash)
	// }
	// t.Log(receipt.Logs[0].Data)
	//
	// reallyBigInt, _ := big.NewInt(0).SetString("20", 16)
	// _, err = contract.Send(transaction, "approve", utils.StringToAddress(generator), reallyBigInt)
	// if err != nil {
	//	t.Errorf("Can't send approve transaction")
	//	t.FailNow()
	// }
}

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

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
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