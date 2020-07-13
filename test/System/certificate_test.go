package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

// 如果是0的话，是否提示错误，默认是0
func TestGetPeriod(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer, err := connection.System.NewCertificate()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	period, err := cer.GetPeriod(common.StringToAddress(coinBase), coinBase)
	if err != nil && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}
	if err == system.ErrCertificateNotExist {
		t.Log(err)
	}
	t.Log(period)
}

// 如何激活证书？？？
func TestGetActive(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer, err := connection.System.NewCertificate()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	isEnable, err := cer.GetActive(common.StringToAddress(coinBase), coinBase)
	if err != nil && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}
	if err == system.ErrCertificateNotExist {
		t.Log(err)
	}
	t.Log(isEnable)
}

// 如果是0000000000000000000000000000000000000000，的话是否提示错误还是？？？
func TestGetIssuer(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer, err := connection.System.NewCertificate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	issuer, err := cer.GetIssuer(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(issuer)
}

// 如果是0000000000000000000000000000000000000000，的话是否提示错误还是？？？
func TestGetIssuerSignature(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer, err := connection.System.NewCertificate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	issuerSignature, err := cer.GetIssuerSignature(common.StringToAddress(coinBase), coinBase)
	if err != nil  && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}

	if err == system.ErrCertificateNotExist {
		t.Log(err)
	}

	if issuerSignature != nil{
		t.Log(issuerSignature.Id)
		t.Log(issuerSignature.PublicKey)
		t.Log(issuerSignature.Algorithm)
		t.Log(issuerSignature.Signature)
	}
}

func TestGetSubjectSignature(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer, err := connection.System.NewCertificate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	subjectSignature, err := cer.GetSubjectSignature(common.StringToAddress(coinBase), coinBase)
	if err != nil  && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}

	if err == system.ErrCertificateNotExist {
		t.Log(err)
	}

	if subjectSignature != nil{
		t.Log(subjectSignature.Id)
		t.Log(subjectSignature.PublicKey)
		t.Log(subjectSignature.Algorithm)
		t.Log(subjectSignature.Signature)
	}
}

// 经过超级节点投票，已经将基础信任锚did:bid:c935bd29a90fbeea87badf3e 激活
// 证书的颁布，
func TestRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer, err := connection.System.NewCertificate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	registerCertificate := new(system.RegisterCertificate)
	registerCertificate.Id = common.StringToAddress(coinBase).String()
	registerCertificate.Context = "context_test"
	registerCertificate.Subject = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
	registerCertificate.Period = 3
	registerCertificate.IssuerAlgorithm = ""
	registerCertificate.IssuerSignature = ""
	registerCertificate.SubjectPublicKey = ""
	registerCertificate.SubjectAlgorithm = ""
	registerCertificate.SubjectSignature = ""
	//registerCertificate
	registerCertificateHash, err := cer.RegisterCertificate(common.StringToAddress(coinBase), registerCertificate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(registerCertificateHash, err)
}

func TestRevokedCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	cer, err := connection.System.NewCertificate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	transactionHash, err := cer.RevokedCertificate(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}



