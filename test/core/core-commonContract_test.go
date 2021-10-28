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
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	solcDir := path.Join(bif.GetCurrentAbPath(), "compiler", "tmp")
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows" {
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	} else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solFile := path.Join(bif.GetCurrentAbPath(), "compiler", "contract", "simple-token.sol")

	solc := compiler.NewSolidityCompiler(solcDir)

	output, err := solc.Compile(solFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	type Contract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}
	var contract Contract

	for _, v := range output.Contracts {
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

	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(singAddr)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  200000,
		Sender:    &sender,
		Recipient: nil,
		Amount:    nil,
	}

	hash, err := contractObj.Deploy(tx, false, singAddrPriKey, contract.ByteCode)

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
	// nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)
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

type Contract struct {
	Abi      string `json:"abi"`
	ByteCode string `json:"byteCode"`
}

func ballotContract() Contract {
	solcDir := path.Join(bif.GetCurrentAbPath(), "compiler", "tmp")
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows" {
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	} else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solFile := path.Join(bif.GetCurrentAbPath(), "compiler", "contract", "ballot.sol")

	solc := compiler.NewSolidityCompiler(solcDir)

	output, err := solc.Compile(solFile)
	if err != nil {
		panic(err)
	}

	var contract Contract

	for _, v := range output.Contracts {
		contract.Abi = v.Abi
		contract.ByteCode = v.Bin
	}

	return contract
}

func ballotDeploy(t *testing.T) string {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	contract := ballotContract()

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(singAddr)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  2000000,
		Sender:    &sender,
		Recipient: nil,
		Amount:    nil,
	}

	var key1 [32]byte
	var key2 [32]byte
	copy(key1[:], []byte("test1"))
	copy(key1[:], []byte("test2"))
	proposalNames := []utils.Hash{key1, key2}

	hash, err := contractObj.Deploy(tx, false, singAddrPriKey, contract.ByteCode, proposalNames)

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

	return hash
}

// 合约交互 查询 call 方法
func ballotCallWinnerName(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	transaction := new(dto.TransactionParameters)
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contractAddress := "did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)
	transaction.ChainId = chainId
	transaction.Sender = generator
	transaction.Recipient = contractAddress
	transaction.AccountNonce = nonce.Uint64()

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	result, err := contractObj.Call(transaction, "winnerName")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		val, _ := result.ToComplexString()
		fmt.Println(val)
	}
}

func ballotCallVoters(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	transaction := new(dto.TransactionParameters)
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contractAddress := "did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)
	transaction.ChainId = chainId
	transaction.Sender = generator
	transaction.Recipient = contractAddress
	transaction.AccountNonce = nonce.Uint64()

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// voter := utils.StringToAddress("did:bid:qwer:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n")

	result, err := contractObj.Call(transaction, "chairperson")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		val, _ := result.ToComplexString()
		fmt.Println(val)
	}
}

// 合约交互 交易 send 方法
func ballotSend(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractAddress := "did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)

	sender := utils.StringToAddress(singAddr)
	recipient := utils.StringToAddress(contractAddress)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  2000000,
		Sender:    &sender,
		Recipient: &recipient,
		Amount:    nil,
	}

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// voter := utils.StringToAddress("did:bid:qwer:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n")
	// voter := utils.StringToAddress("did:bid:qwer:sfCXQusR8SEWgp8fQ9BQu61riWdDLCMN")
	voter := utils.StringToAddress("did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY")

	txHash, err := contractObj.Send(tx, false, singAddrPriKey, "giveRightToVote", voter)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("txHash is ", txHash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txHash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试合约的交互(只是为了测试合约交互，实际使用contract中的Call或者Send)
func TestCoreSendTransactionInteractContract(t *testing.T) {
	// 合约部署
	// txHash := ballotDeploy(t)
	// txHash is  0x8e17880962519fa1421f1eea19a3503290758db10ce37524239030b9b7aa17c3
	// Contract Address:  did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8

	// 合约交互 查询 call 方法
	// ballotCallWinnerName(t)

	// // 合约交互 交易 send 方法
	// ballotSend(t)
	//
	// // call voters
	// // ballotCallVoters(t)
	// var singAddr = "did:bid:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	// // 0x000000000000000000000000b39229c015c595850c16948d72f940303c8f3e11
	// // 0x66323558474251553845387747466f3977474b6f39356a55677459504d323459
	// // 0x66323558474251553845387747466f3977474b6f39356a55677459504d323459
	// 0xadFfd5A56Fc679aa9E8BfD89a7438b614f317ad1
	// fmt.Println(utils.StringToHash(singAddr).String())

	// // res, err := connection.Core.GetTransactionReceipt("0x86401d8ffac5ce804cef47ff89f5297a33ad33f3e9dad5bfb270352b330e4169")
	// res, err := connection.Core.GetTransactionReceipt("0x8451fa14eb1e301a40e5eb01ae075b7186cc6e0747a13600c44ae550394c72b0")
	// fmt.Println(res)
	// fmt.Println(err)
	// res, _ := hexutil.Decode("0xb39229c015c595850c16948d72f940303c8f3e11")
	// fmt.Println(res)
	// var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	// //
	// sender := utils.StringToAddress(singAddr)
	// fmt.Println(sender)
}
