package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestInit(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.Init(sysTxParams, 2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestSetBidName(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.SetBidName(sysTxParams, "1", "testTele2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestGetDocument(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	document, err := doc.GetDocument(true, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(document)
}

func TestAddPublic(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddPublic(sysTxParams, "test", "test", "1", "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelPublic(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelPublic(sysTxParams, "test", "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddAuth(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)

	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddAuth(sysTxParams, "test", "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelAuth(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelAuth(sysTxParams, "test", "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddService(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddService(sysTxParams, "test", "123", "0", "serviceEndpoint")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelService(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelService(sysTxParams, "test", "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddProof(sysTxParams, "123", "0", "testProof", "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)
	doc := connection.System.NewDoc()

	txHash, err := doc.DelProof(sysTxParams, "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddExtra(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddExtra(sysTxParams, "testAttr", "attr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelExtra(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelExtra(sysTxParams, "testAttr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.Enable(sysTxParams, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDisable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.Disable(sysTxParams, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestIsEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	isEnable, err := doc.IsEnable(true, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}
