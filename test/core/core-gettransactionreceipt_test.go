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
	"encoding/json"
	web3 "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestCoreGetTransactionReceipt(t *testing.T) {

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

	var connection = web3.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	byteCode := unmarshalResponse.ByteCode
	contract, err := connection.Core.NewContract(unmarshalResponse.Abi)

	transaction := new(dto.TransactionParameters)
	coinBase, err := connection.Core.GetCoinBase()
	transaction.From = coinBase
	transaction.Gas = big.NewInt(4000000)

	hash, err := contract.Deploy(transaction, byteCode, nil)

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
