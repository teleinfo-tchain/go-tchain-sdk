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
	Abi "github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/types"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
)

// 测试转账
func TestCoreSendTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	// generator, err := connection.Core.GetGenerator()
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	generator := resources.TestAddressAlliance
	toAddress := resources.RegisterAllianceOne
	balance, err := connection.Core.GetBalance(toAddress, block.LATEST)
	if err == nil {
		balBif, _ := utils.FromWei(balance)
		t.Log("toAddress balance is ", balBif)
	}

	transaction := new(dto.TransactionParameters)
	transaction.Sender = generator
	transaction.Recipient = toAddress
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e17))
	transaction.GasLimit = 40000
	transaction.Payload = "Transfer test"

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

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	byteCode := unmarshalResponse.ByteCode
	parsedAbi, err := Abi.JSON(strings.NewReader(unmarshalResponse.Abi))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inputEncode, err := parsedAbi.Pack("", big.NewInt(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)

	// generator, err := connection.Core.GetGenerator()
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }

	transaction.Sender = resources.Addr1
	transaction.Payload = types.ComplexString(byteCode) + types.ComplexString(utils.Bytes2Hex(inputEncode))
	// estimate the gas required to deploy the contract
	gas, err := connection.Core.EstimateGas(transaction)
	if err != nil {
		t.Errorf("Failed EstimateGas")
		t.Error(err)
		t.FailNow()
	}
	t.Log("Estimate gas is ", gas)

	// transaction.GasPrice = big.NewInt(1000000)
	transaction.GasLimit = gas.Uint64()
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

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	parsedAbi, err := Abi.JSON(strings.NewReader(unmarshalResponse.Abi))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inputEncode, err := parsedAbi.Pack("multiply", big.NewInt(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// fromAddress, err := connection.Core.GetGenerator()
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	fromAddress := resources.Addr1

	transaction := new(dto.TransactionParameters)
	transaction.Sender = fromAddress
	// Recipient is contract address
	transaction.Recipient = contractAddress
	transaction.Payload = types.ComplexString("0x" + utils.Bytes2Hex(inputEncode))
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
