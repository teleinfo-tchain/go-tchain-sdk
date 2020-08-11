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
	isSM2Doc       = false
	passwordDoc    = "teleinfo"
	testAddressDoc = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
	testAddressDocPublicKey = "0x04862d66d1255c7ac9c6d86920081069a4d7c295c838140217644ac6f4fcddd0d2edef63951148affb107223fbd1331363e275125bc0cfeb34b1daad26328f8b20"
)

func TestInit(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(55)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.SetBidName(sysTxParams, testAddressDoc, "       ")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestGetDocument(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	document, err := doc.GetDocument(testAddressDoc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(document)
}

func TestAddPublic(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
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

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()
	txHash, err := doc.AddPublic(sysTxParams, testAddressDoc, "secp256k1", "all", testAddressDocPublicKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelPublic(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelPublic(sysTxParams, testAddressDoc, " ")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddAuth(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddAuth(sysTxParams, testAddressDoc, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelAuth(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelAuth(sysTxParams, testAddressDoc, "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddService(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddService(sysTxParams, testAddressDoc, "123", "0", "serviceEndpoint")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelService(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelService(sysTxParams, testAddressDoc, "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddProof(sysTxParams, testAddressDoc, "0", "testProof", "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.DelProof(sysTxParams, testAddressDoc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestAddExtra(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.AddExtra(sysTxParams, testAddressDoc, "attr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDelExtra(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.Enable(sysTxParams, testAddressDoc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestDisable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressDoc, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Doc
	sysTxParams.Password = passwordDoc
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	doc := connection.System.NewDoc()

	txHash, err := doc.Disable(sysTxParams, testAddressDoc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txHash)
}

func TestIsEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	isEnable, err := doc.IsEnable(testAddressDoc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}
