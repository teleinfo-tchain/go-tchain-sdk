package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"strconv"
	"testing"
)

// 53位的字符
const publicKeyTest = ""

//	MessageSha3   什么消息的sha3？？？   Signature 怎么拿到？？？  Id的53个字符，怎么生成的？？？
func TestPeerRegisterCertificate(t *testing.T) {
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

	peerCer := connection.System.NewPeerCertificate()

	registerCertificateInfo := new(dto.RegisterCertificateInfo)
	registerCertificateInfo.Id = ""
	registerCertificateInfo.Apply = ""
	registerCertificateInfo.PublicKey = publicKeyTest
	registerCertificateInfo.NodeName = "testTELE"
	registerCertificateInfo.MessageSha3 = "testTELE"
	registerCertificateInfo.Signature = "testTELE"
	registerCertificateInfo.NodeType = 0
	registerCertificateInfo.Period = system.Year
	registerCertificateInfo.IP = resources.IP
	port, _ := strconv.ParseUint(resources.Port, 10, 64)
	registerCertificateInfo.Port = port
	registerCertificateInfo.CompanyName = "tele info"
	registerCertificateInfo.CompanyCode = "341004600017214"

	// registerCertificateInfo
	peerRegisterCertificateHash, err := peerCer.RegisterCertificate(sysTxParams, registerCertificateInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(peerRegisterCertificateHash, err)
}

func TestPeerRevokedCertificate(t *testing.T) {
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

	peerCer := connection.System.NewPeerCertificate()
	transactionHash, err := peerCer.RevokedCertificate(sysTxParams, publicKeyTest)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestPeerGetPeriod(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	period, err := peerCer.GetPeriod(publicKeyTest)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(period)
}

func TestPeerGetActive(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	isEnable, err := peerCer.GetActive(publicKeyTest)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}

func TestGetPeerCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	peerCertificate, err := peerCer.GetPeerCertificate(publicKeyTest)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(peerCertificate)
}

func TestGetPeerCertificateIdList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	peerCertificateIdList, err := peerCer.GetPeerCertificateIdList(publicKeyTest)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(peerCertificateIdList)
}
