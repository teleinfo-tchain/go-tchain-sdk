package System

import (
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/system"
	"github.com/tchain/go-tchain-sdk/test/resources"
	"io/ioutil"
	"math/big"
	"path"
	"strconv"
	"testing"
	"time"
)

func TestRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressCertificate, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	file := path.Join(bif.GetCurrentAbPath(), "test", "resources", "keystore", resources.TestAddressCertificateFile)
	keyFileData, err := ioutil.ReadFile(file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = resources.TestAddressCertificate
	sysTxParams.IsSM2 = resources.NotSm2
	sysTxParams.Password = resources.SystemPassword
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce
	sysTxParams.ChainId = chainId

	cer := connection.System.NewCertificate()

	registerCertificate := new(dto.RegisterCertificate)
	registerCertificate.Id = resources.PersonCertificateId
	registerCertificate.Context = "test context"
	registerCertificate.Subject = resources.SubjectCertificatedId
	registerCertificate.Period = 1
	registerCertificate.IssuerAlgorithm = "test"
	registerCertificate.IssuerSignature = "test"
	registerCertificate.SubjectPublicKey = resources.SubjectCertificatedIdPubKey
	registerCertificate.SubjectAlgorithm = "test"
	registerCertificate.SubjectSignature = "test"
	// registerCertificate
	registerCertificateHash, err := cer.RegisterCertificate(sysTxParams, registerCertificate)
	if err != nil {
		t.Error(err)

		t.FailNow()
	}
	// 0x03b536ee4e2764aa78d8455e730971be00de89b65677395dea4c0bfa1ec7f753
	t.Log(registerCertificateHash, err)
	time.Sleep(10 * time.Second)
	log, err := connection.System.SystemLogDecode(registerCertificateHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}

}

func TestRevokedCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressCertificate, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("nonce is %d , chainId is %d ", nonce, chainId)

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressCertificateFile)
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
	sysTxParams.GasPrice = nil
	sysTxParams.Gas = 2000000

	cer := connection.System.NewCertificate()

	transactionHash, err := cer.RevokedCertificate(sysTxParams, resources.PersonCertificateId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestRevokedCertificates(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressCertificate, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressCertificateFile)
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

	cer := connection.System.NewCertificate()

	transactionHash, err := cer.RevokedCertificates(sysTxParams)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestGetPeriod(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	period, err := cer.GetPeriod(resources.PersonCertificateId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(period)
}

func TestGetActive(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	isEnable, err := cer.GetActive(resources.PersonCertificateId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}

func TestGetCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	certificate, err := cer.GetCertificate(resources.PersonCertificateId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(certificate)
}

func TestGetIssuer(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	issuer, err := cer.GetIssuer(resources.PersonCertificateId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(issuer)
}

func TestGetSubject(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	subjectSignature, err := cer.GetSubject(resources.PersonCertificateId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(subjectSignature)
}
