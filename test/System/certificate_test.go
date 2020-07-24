package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

// 经过超级节点投票，已经将基础信任锚did:bid:c935bd29a90fbeea87badf3e 激活
// 证书的颁布，
func TestRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Log(err)
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
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	cer := connection.System.NewCertificate()

	registerCertificate := new(dto.RegisterCertificate)
	registerCertificate.Id = common.StringToAddress(coinBase).String()
	registerCertificate.Context = "context_test"
	registerCertificate.Subject = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
	registerCertificate.Period = 3
	registerCertificate.IssuerAlgorithm = ""
	registerCertificate.IssuerSignature = ""
	registerCertificate.SubjectPublicKey = ""
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Log(err)
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
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	cer := connection.System.NewCertificate()

	transactionHash, err := cer.RevokedCertificate(sysTxParams, coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestRevokedCertificates(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Log(err)
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
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	period, err := cer.GetPeriod(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(period)
}

func TestGetActive(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	isEnable, err := cer.GetActive(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}

func TestGetCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	certificate, err := cer.GetCertificate(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(certificate)
}

func TestGetIssuer(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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
