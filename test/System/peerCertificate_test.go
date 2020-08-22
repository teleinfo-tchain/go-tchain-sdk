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
	bid = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	// 53位的字符
	publicKeyTest              = "16Uiu2HAkwviNXPoPHBkZxpg8nURQPiNVeCB9HrocfhXTRCs8j34z"
	isSM2PeerCertificate       = false
	passwordPeerCertificate    = "teleinfo"
	testAddressPeerCertificate = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	testAddressPeerCertificateFile = "../resources/superNodeKeyStore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
)

func TestPeerRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressPeerCertificateFile)
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
	sysTxParams.Version = 1


	peerCer := connection.System.NewPeerCertificate()

	registerCertificateInfo := new(dto.RegisterCertificateInfo)
	registerCertificateInfo.Id = bid
	// todo: 这个链是有问题的，因为其判断条件的问题
	registerCertificateInfo.Apply = "did:bid:590ed37615bdfefa496224c7"
	registerCertificateInfo.PublicKey = publicKeyTest
	registerCertificateInfo.NodeName = "testTELE"
	registerCertificateInfo.NodeType = 0
	registerCertificateInfo.Period = 1
	registerCertificateInfo.IP = resources.IP
	registerCertificateInfo.Port = uint64(44012)
	registerCertificateInfo.CompanyName = "tele info"
	registerCertificateInfo.CompanyCode = "10000000000"

	// 注册的ID（地址） 对应的keystore文件的密码
	idPassword := "teleinfo"
	// 注册的ID（地址）对应的keystore文件
	idKeyFileData, err := ioutil.ReadFile(testAddressPeerCertificateFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(idKeyFileData) == 0 {
		t.Errorf("idKeyFileData can't be empty")
		t.FailNow()
	}
	//  注册的ID（地址）对应的keyStore文件生成方式
	idIsSM2 := false

	peerRegisterCertificateHash, err := peerCer.RegisterCertificate(sysTxParams, registerCertificateInfo, idPassword, idKeyFileData, idIsSM2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(peerRegisterCertificateHash)
}

func TestPeerRevokedCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
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
	keyFileData, err := ioutil.ReadFile(testAddressPeerCertificateFile)
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
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	peerCer := connection.System.NewPeerCertificate()
	transactionHash, err := peerCer.RevokedCertificate(sysTxParams, bid)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

func TestPeerGetPeriod(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetGenerator()
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetGenerator()

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetGenerator()

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetGenerator()

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
