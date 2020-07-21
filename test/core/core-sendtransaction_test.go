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
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"io/ioutil"
	"math/big"
	"testing"
	"time"
)

// test transfer bifer
func TestCoreSendTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	toAddress := resources.AddressTwo
	balance, err := connection.Core.GetBalance(toAddress, block.LATEST)
	if err == nil {
		util := utils.NewUtils()
		balBif, _ := util.FromWei(balance.String())
		fmt.Printf("toAddress bal is %s \n", balBif)
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinBase
	transaction.To = toAddress
	transaction.Value = big.NewInt(0).Mul(big.NewInt(5), big.NewInt(1e17))
	transaction.Gas = big.NewInt(40000)
	transaction.Data = "Transfer Bifer test"

	txID, err := connection.Core.SendTransaction(transaction)

	// Wait for a block
	time.Sleep(time.Second)

	if err != nil {
		t.Errorf("Failed SendTransaction")
		t.Error(err)
		t.FailNow()
	}

	// if success, get transaction hash
	t.Log(txID)
}

// test deploy contract
func TestCoreSendTransactionDeployContract(t *testing.T) {
	content, err := ioutil.ReadFile("../resources/simple-contract.json")

	type Contract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}

	var unmarshalResponse Contract

	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	byteCode, err := utils.NewUtils().ByteCodeDeploy(unmarshalResponse.Abi, unmarshalResponse.ByteCode, big.NewInt(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	fromAddress, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = fromAddress

	transaction.Data = types.ComplexString(byteCode)

	// estimate the gas required to deploy the contract
	gas, err := connection.Core.EstimateGas(transaction)
	if err != nil {
		t.Errorf("Failed EstimateGas")
		t.Error(err)
		t.FailNow()
	}
	t.Log(gas)

	//transaction.Gas = big.NewInt(1000000)
	transaction.Gas = gas
	txID, err := connection.Core.SendTransaction(transaction)

	// Wait for a block
	time.Sleep(time.Second)

	if err != nil {
		t.Errorf("Failed Deploy Contract")
		t.Error(err)
		t.FailNow()
	}

	t.Log(txID)

	time.Sleep(time.Second * 5)

	receipt, err := connection.Core.GetTransactionReceipt(txID)
	if err != nil {
		t.Errorf("Failed GetTransactionReceipt")
		t.Error(err)
		t.FailNow()
	}

	//did:bid:29e92743850b4c7473f2aafa
	t.Log("contract Address is ", receipt.ContractAddress)

}

// test interaction with contract
// its content is decided by specific contract
func TestCoreSendTransactionInteractContract(t *testing.T) {
	content, err := ioutil.ReadFile("../resources/simple-contract.json")

	type Contract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}

	var unmarshalResponse Contract

	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// for more information, please check https://solidity.readthedocs.io/en/v0.6.10/abi-spec.html
	byteCode, err := utils.NewUtils().ByteCodeInteract(unmarshalResponse.Abi, "multiply", big.NewInt(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	fromAddress, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = fromAddress
	// To is contract address
	transaction.To = "did:bid:23bfc0ad7cfa209523193dfd"
	transaction.Data = types.ComplexString(byteCode)

	txID, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Errorf("Failed SendTransaction")
		t.Error(err)
		t.FailNow()
	}

	t.Log("transaction hash is ", txID)

	// wait too long
	time.Sleep(time.Second * 8)
	receipt, err := connection.Core.GetTransactionReceipt(txID)
	if err != nil {
		t.Errorf("Failed GetTransactionReceipt")
		t.Error(err)
		t.FailNow()
	}

	t.Log("receipt logs data is ", receipt.Logs[0].Data)
}

func TestCoreSendTransactionResult(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	tx, err := connection.Core.GetTransactionByHash("0x6c04bc524b5376d5b6d6ef580ae1bb9573bf3e316806cfca72e6fc71566705c7")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx.BlockNumber)
}
