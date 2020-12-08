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
	personCertificate = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	// 与personCertificate对应的公钥
	personCertificatePublicKey = "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"
	isSM2Certificate           = false
	passwordCertificate        = "teleinfo"
	testAddressCertificate     = "did:bid:EFTTQWPMdtghuZByPsfQAUuPkWkWYb"
	testAddressCertificateFile     = "../resources/superNodeKeyStore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb"
)

func TestRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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
	keyFileData, err := ioutil.ReadFile(testAddressCertificateFile)
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
	sysTxParams.ChainId = chainId

	cer := connection.System.NewCertificate()

	registerCertificate := new(dto.RegisterCertificate)
	registerCertificate.Id = personCertificate
	registerCertificate.Context = "test context"
	registerCertificate.Subject = personCertificate
	registerCertificate.Period = 1
	registerCertificate.IssuerAlgorithm = "test"
	registerCertificate.IssuerSignature = "test"
	registerCertificate.SubjectPublicKey = personCertificatePublicKey
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
}

func TestRevokedCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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
	keyFileData, err := ioutil.ReadFile(testAddressCertificateFile)
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

	cer := connection.System.NewCertificate()

	transactionHash, err := cer.RevokedCertificate(sysTxParams, personCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestRevokedCertificates(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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
	keyFileData, err := ioutil.ReadFile(testAddressCertificateFile)
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	issuer, err := cer.GetIssuer(personCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(issuer)
}

func TestGetSubject(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	subjectSignature, err := cer.GetSubject(personCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(subjectSignature)
}
