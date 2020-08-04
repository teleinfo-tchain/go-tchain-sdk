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
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestCoreContract(t *testing.T) {

	content, err := ioutil.ReadFile("../resources/simple-token.json")

	type TruffleContract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}

	var unmarshalResponse TruffleContract

	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	byteCode := unmarshalResponse.ByteCode
	contract, err := connection.Core.NewContract(unmarshalResponse.Abi)

	transaction := new(dto.TransactionParameters)
	coinBase, err := connection.Core.GetCoinBase()
	transaction.From = coinBase
	transaction.Gas = big.NewInt(4000000)

	hash, err := contract.Deploy(transaction, byteCode)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("hash is ", hash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Contract Address: ", receipt.ContractAddress)

	transaction.To = receipt.ContractAddress

	result, err := contract.Call(transaction, "name")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		name, _ := result.ToComplexString()
		if name.ToString() != "Simple ERC20 Token" {
			t.Errorf(fmt.Sprintf("Name not expected; [Expected %s | Got %s]", "SimpleToken", name.ToString()))
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

	result, err = contract.Call(transaction, "balanceOf", coinBase)
	if result != nil && err == nil {
		balance, _ := result.ToBigInt()
		if balance.Cmp(bigInt) != 0 {
			t.Errorf("Balance not expected")
			t.FailNow()
		}
	}

	hash, err = contract.Send(transaction, "approve", utils.StringToAddress(coinBase), big.NewInt(10))
	if err != nil {
		t.Log(err)
		t.Errorf("Can't send approve transaction")
		t.FailNow()
	}

	t.Log(hash)

	reallyBigInt, _ := big.NewInt(0).SetString("20", 16)
	_, err = contract.Send(transaction, "approve", utils.StringToAddress(coinBase), reallyBigInt)
	if err != nil {
		t.Errorf("Can't send approve transaction")
		t.FailNow()
	}

}
