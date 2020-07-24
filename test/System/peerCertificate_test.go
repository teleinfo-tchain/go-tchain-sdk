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
	"strconv"
	"testing"
)

// 53位的字符
const publicKeyTest = ""

// 经过超级节点投票，已经将基础信任锚did:bid:c935bd29a90fbeea87badf3e 激活
//	MessageSha3   什么消息的sha3？？？   Signature 怎么拿到？？？  Id的53个字符，怎么生成的？？？
func TestPeerRegisterCertificate(t *testing.T) {
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

	peerCer := connection.System.NewPeerCertificate()

	registerCertificateInfo := new(dto.RegisterCertificateInfo)
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
