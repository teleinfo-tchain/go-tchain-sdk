package privacytest

import (
	"io/ioutil"
	"math/big"
	"strconv"
	"testing"

	bif "github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/test/privacycontract"
	"github.com/tchain/go-tchain-sdk/test/resources"
)

// 部署合約
func TestFinacialDeploy(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP52+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount("did:bid:sfb9jjhpMreqtwjTQh1ETKjH65RFPJLP", block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressManagerFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	txParams := new(privacycontract.BaseTxParams)
	txParams.IsSM2 = resources.NotSm2
	txParams.Password = resources.SystemPassword
	txParams.KeyFileData = keyFileData
	txParams.GasPrice = big.NewInt(0)
	txParams.Gas = 2000000
	txParams.Nonce = nonce
	txParams.ChainId = chainId

	finacial := connection.PrivacyContract.NewFinacial()

	txHash, err := finacial.FinacialDeploy(txParams)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

// 發行接口 did:bid:qwer:zf21JocUsuvd4gMNFuK3qkCRS1AigY46Y
// 日志：Issue
func TestFinacialIssue(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP52+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount("did:bid:qwer:zf21JocUsuvd4gMNFuK3qkCRS1AigY46Y", block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressManagerFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	txParams := new(privacycontract.BaseTxParams)
	txParams.IsSM2 = resources.IsSM2
	txParams.Password = resources.SystemPassword
	txParams.KeyFileData = keyFileData
	txParams.GasPrice = big.NewInt(0)
	txParams.Gas = 2000000
	txParams.Nonce = nonce
	txParams.ChainId = chainId

	manager := connection.PrivacyContract.NewFinacial()

	contractAddress := resources.TestContractAddress

	txHash, err := manager.FinacialIssue(txParams, contractAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

// 轉移接口
func TestFinacialTransfer(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP52+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount("did:bid:qwer:zf21JocUsuvd4gMNFuK3qkCRS1AigY46Y", block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressManagerFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	txParams := new(privacycontract.BaseTxParams)
	txParams.IsSM2 = resources.IsSM2
	txParams.Password = resources.SystemPassword
	txParams.KeyFileData = keyFileData
	txParams.GasPrice = big.NewInt(0)
	txParams.Gas = 2000000
	txParams.Nonce = nonce
	txParams.ChainId = chainId

	manager := connection.PrivacyContract.NewFinacial()

	contractAddress := resources.TestContractAddress

	txHash, err := manager.FinacialTransfer(txParams, contractAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

// 查詢代筆名接口
func TestFinacialQueryname(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP52+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount("did:bid:qwer:zf21JocUsuvd4gMNFuK3qkCRS1AigY46Y", block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressManagerFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	txParams := new(privacycontract.BaseTxParams)
	txParams.IsSM2 = resources.IsSM2
	txParams.Password = resources.SystemPassword
	txParams.KeyFileData = keyFileData
	txParams.GasPrice = big.NewInt(0)
	txParams.Gas = 2000000
	txParams.Nonce = nonce
	txParams.ChainId = chainId

	manager := connection.PrivacyContract.NewFinacial()

	contractAddress := resources.TestContractAddress
	txHash, err := manager.FinacialQueryName(txParams, contractAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

// 查詢代筆標志接口
func TestFinacialQuerysymbol(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP52+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount("did:bid:qwer:zf21JocUsuvd4gMNFuK3qkCRS1AigY46Y", block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressManagerFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	txParams := new(privacycontract.BaseTxParams)
	txParams.IsSM2 = resources.IsSM2
	txParams.Password = resources.SystemPassword
	txParams.KeyFileData = keyFileData
	txParams.GasPrice = big.NewInt(0)
	txParams.Gas = 2000000
	txParams.Nonce = nonce
	txParams.ChainId = chainId

	manager := connection.PrivacyContract.NewFinacial()

	contractAddress := resources.TestContractAddress

	txHash, err := manager.FinacialQuerySymbol(txParams, contractAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}
