package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"strconv"
	"testing"
)

func TestAddWords(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressSen, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressSenFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = resources.NotSm2
	sysTxParams.Password = resources.SystemPassword
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(45)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressSen, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressSenFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = resources.NotSm2
	sysTxParams.Password = resources.SystemPassword
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()

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
