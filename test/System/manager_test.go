package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

func TestContractEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddress, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2
	sysTxParams.Password = password
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddress, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2
	sysTxParams.Password = password
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddress, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2
	sysTxParams.Password = password
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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

	t.Log(contracts)
}

func TestContractIsEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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