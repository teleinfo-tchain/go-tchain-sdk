package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"testing"
)

const (
	isSM2Sen       = false
	passwordSen    = "teleinfo"
	testAddressSen = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
)

func TestAddWord(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressSen, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--2020-07-07T10-47-32.962000000Z--did-bid-6cc796b8d6e2fbebc9b3cf9e")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Sen
	sysTxParams.Password = passwordSen
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(55)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	sen := connection.System.NewSensitiveWord()

	addWord := "北京"

	txHash, err := sen.AddWord(sysTxParams, addWord)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

func TestAddWords(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressSen, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--2020-07-07T10-47-32.962000000Z--did-bid-6cc796b8d6e2fbebc9b3cf9e")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Sen
	sysTxParams.Password = passwordSen
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(45)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	sen := connection.System.NewSensitiveWord()

	wordsLi := []string{"北京", "上海"}

	txHash, err := sen.AddWords(sysTxParams, wordsLi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

func TestDelWord(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressSen, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--2020-07-07T10-47-32.962000000Z--did-bid-6cc796b8d6e2fbebc9b3cf9e")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Sen
	sysTxParams.Password = passwordSen
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	sen := connection.System.NewSensitiveWord()

	delWord := "北京"

	txHash, err := sen.DelWord(sysTxParams, delWord)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

func TestGetAllWords(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	sen := connection.System.NewSensitiveWord()
	wordsLi, err := sen.GetAllWords()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(wordsLi)
}

func TestIsContainWord(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	sen := connection.System.NewSensitiveWord()
	word := "北京"
	isContain, err := sen.IsContainWord(word)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isContain)
}
