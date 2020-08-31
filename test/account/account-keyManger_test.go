package account

import (
	"github.com/bif/bif-sdk-go/account"
	"testing"
)

func TestGenerateKeyStore(t *testing.T) {
	for _, test := range []struct {
		storeKeyDir       string
		isSm2             bool
		password          string
		useLightweightKDF bool
	}{
		{"./keystore", true, "teleinfo", false},
		{"./keystore", false, "teleinfo", false},
	} {
		addr, err := account.GenerateKeyStore(test.storeKeyDir, test.isSm2, test.password, test.useLightweightKDF)
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
		// {"teleinfo", "did:bid:ZFT4CziA2ktCNgfQPqSm1GpQxSck5q4", "./keystore/UTC--2020-08-19T05-48-44.625362500Z--did-bid-ZFT4CziA2ktCNgfQPqSm1GpQxSck5q4"},
		// {"teleinfo", "did:bid:EFTTQWPMdtghuZByPsfQAUuPkWkWYb", "./keystore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb"},
		{"teleinfo", "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs", "./keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"},
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
		password   string
		privateKey string
		isSM2      bool
		keyDir     string
	}{
		// {"teleinfo", "89b9c1cfc8ab8937cfda96393d4cf2f9789b824c75ff8eaeeeebd572193bec38", true, "./keystore"},
		// {"teleinfo", "e4b4a35bee3d92a0b07f16e3253ae8459e817305514dcd0ed0c64342312b41d8", false, "./keystore"},
		{"teleinfo", "41e46e858ea707453d8fc553805772165a4f66e6e18ca38220daa157534e0c0e", false, "./keystore"},
	} {
		isSuccess, err := account.PrivateKeyToKeyStoreFile(test.keyDir, test.isSM2, test.privateKey, test.password)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if isSuccess {
			t.Log("privateKey To KeyStoreFile success")
		}
	}

}

func TestGetAddressFromPrivate(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		isSM2      bool
	}{
		{"89b9c1cfc8ab8937cfda96393d4cf2f9789b824c75ff8eaeeeebd572193bec38", true},
		{"e4b4a35bee3d92a0b07f16e3253ae8459e817305514dcd0ed0c64342312b41d8", false},
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
		{"89b9c1cfc8ab8937cfda96393d4cf2f9789b824c75ff8eaeeeebd572193bec38", true, "0x0102d53a8080379bb6499966687a9fccd3ac0641010eb53c983b9dd7f6a0c860b1665275b26d616eecee10d7bd03755c31c4e1ab7ca45e3b7b266442f7f64efa03"},
		{"e4b4a35bee3d92a0b07f16e3253ae8459e817305514dcd0ed0c64342312b41d8", false, "0x043ee1708e4b431e71b1cc596c15425b8e889b80ec120840b6dd998a3a6397142405875eebe6b3488723e6ad3c5c7397c42c57696ac1e2fa925c0a1f6a61fc20a7"},
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
		// {"./keystore/UTC--2020-08-19T05-48-44.625362500Z--did-bid-ZFT4CziA2ktCNgfQPqSm1GpQxSck5q4", "teleinfo", true, "0x0102d53a8080379bb6499966687a9fccd3ac0641010eb53c983b9dd7f6a0c860b1665275b26d616eecee10d7bd03755c31c4e1ab7ca45e3b7b266442f7f64efa03"},
		// {"./keystore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb", "teleinfo", false, "0x043ee1708e4b431e71b1cc596c15425b8e889b80ec120840b6dd998a3a6397142405875eebe6b3488723e6ad3c5c7397c42c57696ac1e2fa925c0a1f6a61fc20a7"},
		{"./keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs", "teleinfo", false, "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"},
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
func TestGenerateNodeUrl(t *testing.T) {
	for _, test := range []struct {
		storeKeyDir string
		isSm2       bool
		password    string
		host        string
		port        uint64
	}{
		{"./keystore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb", true, "teleinfo", "127.0.0.1", 55555},
		{"./keystore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb", false, "teleinfo", "127.0.0.1", 55555},
	} {
		nodeUrl, err := account.GenerateNodeUrl(test.storeKeyDir, test.password, test.host, test.isSm2, test.port)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("nodeUrl is : %s \n",nodeUrl)
	}

}
