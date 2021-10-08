package account

import (
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestGenerateKeyStore(t *testing.T) {
	for _, test := range []struct {
		storeKeyDir       string
		isSm2             bool
		password          string
		chainCode         string
		useLightweightKDF bool
	}{
		{"../resources/keystore", true, resources.PassWord, resources.ChainCode, false},
		{"../resources/keystore", false, resources.PassWord, resources.ChainCode, false},
	} {
		addr, err := account.GenerateKeyStore(test.storeKeyDir, test.isSm2, test.password, test.chainCode, test.useLightweightKDF)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("是否选择国密: %t , 生成的地址是: %s \n", test.isSm2, addr)
	}

}

func TestGetPrivateKeyFromFile(t *testing.T) {
	for _, test := range []struct {
		password string
		address  string
		keyDir   string
	}{
		{resources.PassWord, resources.Addr1, "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"},
	} {
		pri, err := account.GetPrivateKeyFromFile(test.address, test.keyDir, test.password)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log("pri:", pri)
	}
}

func TestPrivateKeyToKeyStoreFile(t *testing.T) {
	for _, test := range []struct {
		keyDir     string
		isSM2      bool
		privateKey string
		password   string
		chainCode  string

	}{
		{"../resources/keystore",false,resources.Addr1Pri, resources.PassWord, resources.ChainCode},
	} {
		isSuccess, err := account.PrivateKeyToKeyStoreFile(test.keyDir, test.isSM2, test.privateKey, test.password, test.chainCode)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if isSuccess {
			t.Log("privateKey Recipient KeyStoreFile success")
		}
	}

}

func TestGetAddressFromPrivate(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		isSM2      bool
	}{
		{resources.Addr1Pri, true},
		{resources.Addr2Pri, false},
	} {
		accountAddress, err := account.GetAddressFromPrivate(test.privateKey, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log("accountAddress: ", accountAddress)
	}
}

func TestGetPublicKeyFromPrivate(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		isSM2      bool
		want       string
	}{
		{resources.Addr1Pri, true, resources.Addr1Hex},
		{resources.Addr2Pri, false, resources.Addr2Hex},
	} {
		publicKey, err := account.GetPublicKeyFromPrivate(test.privateKey, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// t.Log("publicKey len is ", len(publicKey[2:]))
		if publicKey != test.want {
			t.Logf("want publicKey is : %s , result is %s \n", test.want, publicKey)
		}
	}

}

func TestGetGetPublicKeyFromFile(t *testing.T) {
	for _, test := range []struct {
		privateKeyFilePath string
		password           string
		isSM2              bool
		want               string
	}{
		{"../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs", "teleInfo", false, "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"},
	} {
		publicKey, err := account.GetPublicKeyFromFile(test.privateKeyFilePath, test.password, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// t.Log("publicKey len is ", len(publicKey[2:]))
		if publicKey != test.want {
			t.Logf("want publicKey is : %s , result is %s \n", test.want, publicKey)
		}
	}

}

// keyStorePath, password, host string, isSM2 bool, port uint64
func TestGenerateNodeUrlFromKeyStore(t *testing.T) {
	for _, test := range []struct {
		nodePrivateKeyPath string
		password           string
		host               string
		port               uint64
	}{
		{"./keystore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb", "teleInfo", "127.0.0.1", 55555},
	} {
		nodeUrl, err := account.GenerateNodeUrlFromKeyStore(test.nodePrivateKeyPath, test.password, test.host, test.port)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("nodeUrl is : %s \n", nodeUrl)
	}

}

func TestGenerateNodeUrlFromPrivateKey(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		host       string
		port       uint64
	}{
		{resources.Addr1Pri, resources.IP00, resources.Port},
		{resources.Addr2Pri, resources.IP00, resources.Port},
	} {
		nodeUrl, err := account.GenerateNodeUrlFromPrivateKey(test.privateKey, test.host, test.port)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("nodeUrl is : %s \n", nodeUrl)
	}

}
