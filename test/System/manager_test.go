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
	isSM2Manager       = false
	passwordManager    = "teleinfo"
	testAddressManager = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
)

func TestContractEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressManager, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Manager
	sysTxParams.Password = passwordManager
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	manager := connection.System.NewManager()

	contractAddress := resources.Address55

	txHash, err := manager.Enable(sysTxParams, contractAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

func TestContractDisable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressManager, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Manager
	sysTxParams.Password = passwordManager
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	manager := connection.System.NewManager()

	contractAddress := resources.Address55

	txHash, err := manager.Disable(sysTxParams, contractAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

func TestSetPower(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressManager, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Manager
	sysTxParams.Password = passwordManager
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	manager := connection.System.NewManager()

	userAddress := resources.Address55
	var power uint64 = 1

	txHash, err := manager.SetPower(sysTxParams, userAddress, power)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txHash)
}

func TestAllContracts(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	manager := connection.System.NewManager()
	contracts, err := manager.GetAllContracts()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", contracts)

	if len(contracts) != 0{
		t.Logf("%#v \n", contracts[0])
	}
}

func TestContractIsEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	manager := connection.System.NewManager()
	contractAddress := resources.Address55
	isEnable, err := manager.IsEnable(contractAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}

func TestGetPower(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	manager := connection.System.NewManager()
	userAddress := resources.Address55
	power, err := manager.GetPower(userAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(power)
}
