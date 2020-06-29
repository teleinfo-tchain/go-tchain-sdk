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
 * @file core-sendtransaction_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */
package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
	"testing"
	"time"
)

// test transfer bifer
func TestCoreSendTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinbase, err := connection.Core.GetCoinbase()
	// coinbase address is did:bid:6cc796b8d6e2fbebc9b3cf9e
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	toAddress := "did:bid:c117c1794fc7a27bd301ae52"
	balance, err := connection.Core.GetBalance(toAddress, block.LATEST)
	if err == nil {
		util := utils.NewUtils()
		balBif, _ := util.FromWei(balance.String())
		fmt.Printf("toAddress bal is %s \n", balBif)
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = coinbase
	transaction.To = toAddress
	transaction.Value = big.NewInt(0).Mul(big.NewInt(5), big.NewInt(1e17))
	transaction.Gas = big.NewInt(40000)
	transaction.Data = types.ComplexString("Transfer Bifer test")

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

	time.Sleep(time.Second)

	tx, err := connection.Core.GetTransactionByHash("0xc8e0815167040fc13421e3b3d2ec38188ee0ac4a8afad6c5fe71a4f5451b7691")

	if err != nil {
		t.Errorf("Failed GetTransactionByHash")
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx.BlockNumber)
}

// test deploy contract
func TestCoreSendTransactionDeployContract(t *testing.T){
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	fromAddress, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.From = fromAddress
	// data is compiled by contract
	// "contract Multiply7 { event Print(uint); function multiply(uint input) returns (uint) { Print(input * 7); return input * 7; } }
	transaction.Data = "0x6060604052605f8060106000396000f3606060405260e060020a6000350463c6888fa18114601a575b005b60586004356007810260609081526000907f24abdb5865df5079dcc5ac590ff6f01d5c16edbc5fab4e195d9febd1114503da90602090a15060070290565b5060206060f3"

	// estimate the gas required to deploy the contract
	gas, err := connection.Core.EstimateGas(transaction)
	if err != nil {
		t.Errorf("Failed EstimateGas")
		t.Error(err)
		t.FailNow()
	}
	t.Log(gas)

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

	time.Sleep(time.Second*5)

	receipt, err := connection.Core.GetTransactionReceipt(txID)
	if err != nil {
		t.Errorf("Failed GetTransactionReceipt")
		t.Error(err)
		t.FailNow()
	}

	//did:bid:0ad33772c600c61b184919d7
	t.Log("contract Address is ", receipt.ContractAddress)

}

// test interaction with contract
// its content is decided by specific contract
func TestCoreSendTransactionInteractContract(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	fromAddress, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// for more information, please check https://solidity.readthedocs.io/en/v0.6.10/abi-spec.html
	util := utils.NewUtils()
	funcEncode, _ := util.Sha3("multiply(uint256)")
	// 32 bytes
	paramsEncode, _ := util.ToTwosComplement(6)

	transaction := new(dto.TransactionParameters)
	transaction.From = fromAddress
	// To is contract address
	transaction.To = "did:bid:0ad33772c600c61b184919d7"
	transaction.Data = types.ComplexString("0x" + funcEncode[2:10] + paramsEncode[2:])

	txID, err := connection.Core.SendTransaction(transaction)

	// Wait for a block
	time.Sleep(time.Second)

	if err != nil {
		t.Errorf("Failed SendTransaction")
		t.Error(err)
		t.FailNow()
	}

	t.Log(txID)

	// wait too long
	time.Sleep(time.Second*8)
	receipt, err := connection.Core.GetTransactionReceipt(txID)
	if err != nil {
		t.Errorf("Failed GetTransactionReceipt")
		t.Error(err)
		t.FailNow()
	}

	t.Log("receipt logs data is ", receipt.Logs[0].Data)
}