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
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/compiler"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
	"path"
	"runtime"
	"strconv"
	"testing"
	"time"
)

//  部署合约及合约的call方法调用
func TestCoreContract(t *testing.T) {
	solcDir  := path.Join(bif.GetCurrentAbPath(), "compiler", "tmp")
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows"{
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	}else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solFile  := path.Join(bif.GetCurrentAbPath(), "compiler", "contract", "simple-token.sol")

	solc := compiler.NewSolidityCompiler(solcDir)

	output, err := solc.Compile(solFile)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	type Contract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}
	var contract Contract

	for _, v := range output.Contracts{
		contract.Abi = v.Abi
		contract.ByteCode = v.Bin
	}


	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(resources.Addr1)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  200000,
		Sender:    &sender,
		Recipient: nil,
		Amount:    nil,
	}

	hash, err := contractObj.Deploy(tx, false, resources.Addr1Pri, contract.ByteCode)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txHash is ", hash)

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


	transaction := new(dto.TransactionParameters)
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// txHash is  0x0e57cc103a20ebdbd3ac36530e440e95af68b4e6b3ce2f55d845218c7a7b1add
	// Contract Address:  did:bid:qwer:sfzxtY15dmcYSKN9r5sWs3GDBFvCU2Sq

	// chainId, _ := connection.Core.GetChainId()
	// nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)
	transaction.ChainId = chainId
	transaction.Sender = generator
	transaction.Recipient = receipt.ContractAddress
	transaction.AccountNonce = nonce.Uint64()

	result, err := contractObj.Call(transaction, "getName")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		val, _ := result.ToComplexString()
		fmt.Println(val)
	}
}

func ballotDeploy(t *testing.T) string {
	solcDir  := path.Join(bif.GetCurrentAbPath(), "compiler", "tmp")
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows"{
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	}else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solFile  := path.Join(bif.GetCurrentAbPath(), "compiler", "contract", "ballot.sol")

	solc := compiler.NewSolidityCompiler(solcDir)

	output, err := solc.Compile(solFile)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	type Contract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}
	var contract Contract

	for _, v := range output.Contracts{
		contract.Abi = v.Abi
		contract.ByteCode = v.Bin
	}


	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(resources.Addr1)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  200000,
		Sender:    &sender,
		Recipient: nil,
		Amount:    nil,
	}

	proposalNames := "0100000000000000000000000000000000000000000000000000000000000000"

	hash, err := contractObj.Deploy(tx, false, resources.Addr1Pri, contract.ByteCode, proposalNames)

	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Log("txHash is ", hash)

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

	return hash
}
// 测试合约的交互(只是为了测试合约交互，实际使用contract中的Call或者Send)
func TestCoreSendTransactionInteractContract(t *testing.T) {
	txHash := ballotDeploy(t)
	fmt.Println(txHash)
}
