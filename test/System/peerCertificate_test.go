package System

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"strconv"
	"testing"
)

const publicKeyTest = "0xc70357f11062a2a41a8a0907e9724659829729615826a294bb7cd45a11c4f7ed7a1deb36b3715b60c7665f5966526029103fb048ee7e38a7dda2326a622f8eba"

// 经过超级节点投票，已经将基础信任锚did:bid:c935bd29a90fbeea87badf3e 激活
// 节点证书的颁布，
// failed to parse peer ID: cid too short??????
//failed to parse peer ID: selected encoding not supported
//failed to parse peer ID: error while conversion: strconv.ParseInt: parsing "0000000x": invalid syntax
// 公钥的问题
func TestPeerRegisterCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	peerCer, err := connection.System.NewPeerCertificate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	registerCertificateInfo := new(system.RegisterCertificateInfo)
	registerCertificateInfo.PublicKey = publicKeyTest
	registerCertificateInfo.Period = 2
	registerCertificateInfo.IP = resources.IP
	port, _ := strconv.ParseUint(resources.Port, 10,64)
	registerCertificateInfo.Port = port
	registerCertificateInfo.CompanyName = "tele info"
	registerCertificateInfo.CompanyCode = "11250"
	fmt.Println(hexutil.Encode(common.FromHex(coinBase)))
	//registerCertificateInfo
	peerRegisterCertificateHash, err := peerCer.RegisterCertificate(common.StringToAddress(coinBase), registerCertificateInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(peerRegisterCertificateHash, err)
}

func TestPeerRevokedCertificate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	peerCer, err := connection.System.NewPeerCertificate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	transactionHash, err := peerCer.RevokedCertificate(common.StringToAddress(coinBase), publicKeyTest)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)
}

// 如果是0的话，是否提示错误，默认是0
func TestPeerGetPeriod(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	peerCer, err := connection.System.NewPeerCertificate()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	period, err := peerCer.GetPeriod(common.StringToAddress(coinBase), publicKeyTest)
	if err != nil && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}
	if err == system.ErrCertificateNotExist {
		t.Log(err)
	}
	t.Log(period)
}

func TestPeerGetActive(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	peerCer, err := connection.System.NewPeerCertificate()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	isEnable, err := peerCer.GetActive(common.StringToAddress(coinBase), publicKeyTest)
	if err != nil && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}
	if err == system.ErrCertificateNotExist {
		t.Log(err)
	}
	t.Log(isEnable)
}