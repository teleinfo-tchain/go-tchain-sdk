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
	Abi "github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/types"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
	"time"
)

// 测试转账
func TestCoreSendTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	// coinBase, err := connection.Core.GetCoinBase()
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	coinBase := resources.TestAddr
	toAddress := resources.NewAddrE
	balance, err := connection.Core.GetBalance(resources.NewAddrE, block.LATEST)
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

// 测试部署合约(只是为了测试部署合约，实际使用contract中的Deploy)
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

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	byteCode := unmarshalResponse.ByteCode
	parsedAbi, err := Abi.JSON(strings.NewReader(unmarshalResponse.Abi))
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	inputEncode, err := parsedAbi.Pack("", big.NewInt(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)

	// coinBase, err := connection.Core.GetCoinBase()
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }

	transaction.From = resources.TestAddr
	transaction.Data = types.ComplexString(byteCode) + types.ComplexString(utils.Bytes2Hex(inputEncode))
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
	txHash, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Errorf("Failed Deploy Contract")
		t.Error(err)
		t.FailNow()
	}
	t.Log("transaction hash is ", txHash)


	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txHash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// did:bid:ace45606ce7b19c7da1143cb
	t.Log("contract Address is ", receipt.ContractAddress)

}

func TestGetReceipt(t *testing.T){
	txHash := "0x24d339d9d55ddc86041ecc8a5ca52e50800d81b94956d2b57b97092a2b750de6"
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	receipt, err := connection.Core.GetTransactionReceipt(txHash)
	fmt.Println("err", err)
	fmt.Printf("txHash %#v", receipt)
}

// 测试合约的交互(只是为了测试合约交互，实际使用contract中的Call或者Send)
func TestCoreSendTransactionInteractContract(t *testing.T) {
	content, err := ioutil.ReadFile("../resources/simple-contract.json")
	const contractAddress = "did:bid:EFT3zLGVg6PeULDZzHkTwKCuWwsc75m"

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

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	parsedAbi, err := Abi.JSON(strings.NewReader(unmarshalResponse.Abi))
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	inputEncode, err := parsedAbi.Pack("multiply", big.NewInt(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// fromAddress, err := connection.Core.GetCoinBase()
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	fromAddress := resources.TestAddr

	transaction := new(dto.TransactionParameters)
	transaction.From = fromAddress
	// To is contract address
	transaction.To = contractAddress
	transaction.Data = types.ComplexString("0x"+utils.Bytes2Hex(inputEncode))
	txHash, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Errorf("Failed SendTransaction")
		t.Error(err)
		t.FailNow()
	}
	t.Log("transaction hash is ", txHash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txHash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("receipt logs data is ", receipt.Logs[0].Data)
}

