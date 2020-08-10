package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"testing"
)

const (
	// 注册证书的bid
	personCertificate      = "did:bid:c935bd29a90fbeea87badf3e"
	isSM2Certificate       = false
	passwordCertificate    = "teleinfo"
	testAddressCertificate = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
)

func TestRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressCertificate, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Certificate
	sysTxParams.Password = passwordCertificate
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	cer := connection.System.NewCertificate()

	registerCertificate := new(dto.RegisterCertificate)
	registerCertificate.Id = personCertificate
	registerCertificate.Context = "context_test"
	registerCertificate.Subject = personCertificate
	registerCertificate.Period = 1
	registerCertificate.IssuerAlgorithm = ""
	registerCertificate.IssuerSignature = ""
	registerCertificate.SubjectPublicKey = "0x23"
	registerCertificate.SubjectAlgorithm = ""
	registerCertificate.SubjectSignature = ""
	// registerCertificate
	registerCertificateHash, err := cer.RegisterCertificate(sysTxParams, registerCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(registerCertificateHash, err)
}

func TestRevokedCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressCertificate, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("nonce is %d , chainId is %d ", nonce, chainId)

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
	sysTxParams.IsSM2 = isSM2Certificate
	sysTxParams.Password = passwordCertificate
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = nil
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = 0
	sysTxParams.ChainId = nil

	cer := connection.System.NewCertificate()

	transactionHash, err := cer.RevokedCertificate(sysTxParams, testAddressCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestRevokedCertificates(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressCertificate, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2Certificate
	sysTxParams.Password = passwordCertificate
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(45)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	cer := connection.System.NewCertificate()

	transactionHash, err := cer.RevokedCertificates(sysTxParams)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestGetPeriod(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	period, err := cer.GetPeriod(personCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(period)
}

func TestGetActive(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	isEnable, err := cer.GetActive(personCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}

func TestGetCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	certificate, err := cer.GetCertificate(personCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(certificate)
}

func TestGetIssuer(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	issuer, err := cer.GetIssuer(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(issuer)
}

func TestGetSubject(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	subjectSignature, err := cer.GetSubject(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if subjectSignature != nil {
		t.Log(subjectSignature.Id)
		t.Log(subjectSignature.PublicKey)
		t.Log(subjectSignature.Algorithm)
		t.Log(subjectSignature.Signature)
	}
}
