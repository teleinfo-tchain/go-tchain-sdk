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
	bid = "did:bid:c117c1794fc7a27bd301ae52"
	// 53位的字符
	publicKeyTest              = "16Uiu2HAmT5Fb81HeCsZGU4rRzNQR8YDtuGMuhacdRcYsxqRrEhFt"
	isSM2PeerCertificate       = false
	passwordPeerCertificate    = "teleinfo"
	testAddressPeerCertificate = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
)

func TestPeerRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
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
	sysTxParams.IsSM2 = false
	sysTxParams.Password = passwordPeerCertificate
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	peerCer := connection.System.NewPeerCertificate()

	registerCertificateInfo := new(dto.RegisterCertificateInfo)
	registerCertificateInfo.Id = bid
	registerCertificateInfo.Apply = ""
	registerCertificateInfo.PublicKey = publicKeyTest
	registerCertificateInfo.NodeName = "testTELE"
	registerCertificateInfo.NodeType = 0
	registerCertificateInfo.Period = 1
	registerCertificateInfo.IP = resources.IP55
	registerCertificateInfo.Port = uint64(44012)
	registerCertificateInfo.CompanyName = "tele info"
	registerCertificateInfo.CompanyCode = "10000000000"

	peerRegisterCertificateHash, err := peerCer.RegisterCertificate(sysTxParams, registerCertificateInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(peerRegisterCertificateHash)
}

func TestPeerRevokedCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressPeerCertificate, block.LATEST)
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
	sysTxParams.IsSM2 = isSM2PeerCertificate
	sysTxParams.Password = passwordPeerCertificate
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	peerCer := connection.System.NewPeerCertificate()
	transactionHash, err := peerCer.RevokedCertificate(sysTxParams, bid)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestPeerGetPeriod(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	period, err := peerCer.GetPeriod(bid)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(period)
}

func TestPeerGetActive(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	isEnable, err := peerCer.GetActive(bid)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}

func TestGetPeerCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	peerCertificate, err := peerCer.GetPeerCertificate(bid)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("detail info is %#v \n", peerCertificate)
}

func TestGetPeerCertificateIdList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	peerCer := connection.System.NewPeerCertificate()
	peerCertificateIdList, err := peerCer.GetPeerCertificateIdList(bid)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(peerCertificateIdList)
}
