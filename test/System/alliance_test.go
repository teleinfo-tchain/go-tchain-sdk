package System

import (
	"errors"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"strconv"
	"testing"
)

func connectWithSig() (*bif.Bif, *system.SysTxParams, error) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		return nil, nil, err
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressAlliance, block.LATEST)
	if err != nil {
		return nil, nil, err
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressAllianceFile)
	if err != nil {
		return nil, nil, err
	}
	if len(keyFileData) == 0 {
		return nil, nil, errors.New("keyFileData can't be empty")
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = resources.NotSm2
	sysTxParams.Password = resources.SystemPassword
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(45)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	return connection, sysTxParams, nil
}

func connectBif() (*bif.Bif, error) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func TestRegisterDirector(t *testing.T) {
	con, sigPara, err := connectWithSig()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	registerDirector := new(dto.AllianceInfo)
	registerDirector.Id = resources.TestAddressAlliance
	registerDirector.PublicKey = resources.AlliancePublicKey
	registerDirector.CompanyName = "teleInfo"
	registerDirector.CompanyCode = "110112"

	registerDirectorHash, err := ali.RegisterDirector(sigPara, registerDirector)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(registerDirectorHash, err)
}

func TestUpgradeDirector(t *testing.T) {
	con, sigPara, err := connectWithSig()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	transactionHash, err := ali.UpgradeDirector(sigPara, resources.TestAddressAlliance)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestRevoke(t *testing.T) {
	con, sigPara, err := connectWithSig()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	revokeReason := "不合法"

	transactionHash, err := ali.Revoke(sigPara, resources.TestAddressAlliance, revokeReason)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestSetWeights(t *testing.T) {
	con, sigPara, err := connectWithSig()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	transactionHash, err := ali.SetWeights(sigPara, 1, 2, 3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestGetAllDirectors(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	directors, err := ali.GetAllDirectors()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("directors is %v \n", directors)
}

func TestGetAllVices(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	vices, err := ali.GetAllVices()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("vices is %v \n", vices)
}

func TestGetAlliance(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	aliance, err := ali.GetAlliance(resources.TestAddressAlliance)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("aliance is %v \n", aliance)
}

func TestGetWeights(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	weights, err := ali.GetWeights()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("weights is %v \n", weights)
}
