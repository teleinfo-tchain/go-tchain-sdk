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
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common/types"
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

// 测试转账
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
		balBif, _ := util.FromWei(balance)
		t.Log("toAddress balance is ", balBif)
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinBase
	transaction.To = toAddress
	transaction.Value = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e17))
	transaction.Gas = big.NewInt(40000)
	transaction.Data = "Transfer test"

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

// 测试部署合约
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
	t.Log("Estimate gas is ", gas)

	// transaction.Gas = big.NewInt(1000000)
	transaction.Gas = gas
	txID, err := connection.Core.SendTransaction(transaction)

	// Wait for a block
	// 等待时间较短，交易可能还未执行，导致测试失败
	time.Sleep(time.Second)

	if err != nil {
		t.Errorf("Failed Deploy Contract")
		t.Error(err)
		t.FailNow()
	}

	t.Log("transaction hash is ", txID)

	time.Sleep(time.Second * 5)

	// txID := "0xa7be6f1c34f2fa4e6d36f33be9e8f461763f81a7bbd8eb502e9e4507c3704197"
	receipt, err := connection.Core.GetTransactionReceipt(txID)
	if err != nil {
		t.Errorf("Failed GetTransactionReceipt")
		t.Error(err)
		t.FailNow()
	}

	// did:bid:6f7a7de13fb193f10a76255e
	t.Log("contract Address is ", receipt.ContractAddress)

}

// 测试合约的交互
func TestCoreSendTransactionInteractContract(t *testing.T) {
	content, err := ioutil.ReadFile("../resources/simple-contract.json")
	const contractAddress = "did:bid:6f7a7de13fb193f10a76255e"

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
	transaction.To = contractAddress
	transaction.Data = types.ComplexString(byteCode)

	txID, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Errorf("Failed SendTransaction")
		t.Error(err)
		t.FailNow()
	}

	t.Log("transaction hash is ", txID)

	// wait a block
	// 可能等待时间不足，导致交易执行失败
	time.Sleep(time.Second * 8)
	// txID := "0x93e2b5eba500c7c6bd43b8c6f09eb86f1e872ddcadd874255b622e304811ca61"
	receipt, err := connection.Core.GetTransactionReceipt(txID)
	if err != nil {
		t.Errorf("Failed GetTransactionReceipt")
		t.Error(err)
		t.FailNow()
	}

	t.Log("receipt logs data is ", receipt.Logs[0].Data)
}

