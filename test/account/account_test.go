package account

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestCreate(t *testing.T) {
	accountAddress, privateKey, err := account.Create(false, "sqtx")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("account is %s, privateKey is %s", accountAddress, privateKey)
}

func TestPrivateKeyToAccount(t *testing.T) {
	accountAddress, privateKey, err := account.Create(false, resources.ChainCode)
	parsedAccount, err := account.PriKeyToAccount(privateKey, false, resources.ChainCode)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if parsedAccount != accountAddress {
		t.Error("Account parsed error")
		t.FailNow()
	}

	t.Log("account is ", parsedAccount)

}

func TestEncrypt(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		isSM2      bool
		password   string
		chainCode  string
	}{
		{resources.Addr1Pri, true, resources.PassWord, resources.ChainCode},
		{resources.Addr2Pri, false, resources.PassWord, resources.ChainCode},
	} {
		keyJson, err := account.Encrypt(test.privateKey, test.isSM2, test.password, resources.ChainCode, false)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log("keyJson is ", keyJson)
	}
}

func TestDecrypt(t *testing.T) {
	for _, test := range []struct {
		isSM2       bool
		password    string
		keyDir      string
		wantAddress string
	}{
		{false, resources.PassWord, "UTC--2021-10-19T05-33-49.419105162Z--did_bid_qwer_sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"},
	} {
		file := bif.GetCurrentAbPath()+resources.KeyStoreFile+test.keyDir
		keyJson, err := ioutil.ReadFile(file)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		address, _, err := account.Decrypt(keyJson, test.isSM2, test.password)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if address != test.wantAddress {
			t.Logf("result address is %s, want address is %s \n", address, test.wantAddress)
		}
	}
}

func TestSignTransaction(t *testing.T) {
	for _, test := range []struct {
		id         string
		address    string
		privateKey string
		to         string
		isSM2      bool
	}{
		{"非国密签名转账", "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a", "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY", false},
	} {
		connection := bif.NewBif(providers.NewHTTPProvider("39.99.132.122:44002", 10, false))
		nonce, err := connection.Core.GetTransactionCount(test.address, block.LATEST)
		if err != nil {
			t.Errorf("Failed connection")
			t.FailNow()
		}

		chainId, err := connection.Core.GetChainId()
		if err != nil {
			t.Errorf("Failed getChainId")
			t.FailNow()
		}

		t.Log("nonce is ", nonce, " chainId is ", chainId)
		var recipient utils.Address
		recipient = utils.StringToAddress(test.to)
		tx := &account.SignTxParams{
			Recipient: &recipient,
			Nonce:     nonce,
			GasPrice:  big.NewInt(0),
			GasLimit:  21000,
			Amount:    big.NewInt(50000000000),
			Payload:   nil,
			ChainId:   chainId,
		}
		res, err := account.SignTransaction(tx, test.privateKey, test.isSM2)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("%s : %v \n", test.id, res.Raw)

		txHash, err := connection.Core.SendRawTransaction(res.Raw.String())
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("txHash is %s", txHash)
	}
}

func TestHashMessage(t *testing.T) {
	for _, test := range []struct {
		message     string
		isSm2       bool
		hashMessage string
	}{
		{"Hello World", false, "0xa1de988600a42c4b4ab089b619297c17d53cffae5d5120d82d8a92d0bb3b78f2"},
		{utils.Utf8ToHex("Hello World"), false, "0xa1de988600a42c4b4ab089b619297c17d53cffae5d5120d82d8a92d0bb3b78f2"},
	} {
		res := account.HashMessage(test.message, test.isSm2)
		if res != test.hashMessage {
			t.Errorf("hash error, input is %s, result is %s, expect is %s \n", test.message, res, test.hashMessage)
			t.FailNow()
		}

	}
}

func TestGenKeyStore(t *testing.T) {
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
		addr, err := account.GenKeyStore(test.storeKeyDir, test.isSm2, test.password, test.chainCode, test.useLightweightKDF)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("是否选择国密: %t , 生成的地址是: %s \n", test.isSm2, addr)
	}

}

func TestPriKeyFromKeyStore(t *testing.T) {
	for _, test := range []struct {
		password string
		address  string
		keyDir   string
	}{
		{"teleinfo", "did:bid:sqtx:sfBUW9rkT9VCVLVwvHd8G8qEUaSFr5GN", "UTC--2021-11-17T09-20-40.402116724Z--did_bid_sqtx_sfBUW9rkT9VCVLVwvHd8G8qEUaSFr5GN"},
	} {
		path := bif.GetCurrentAbPath()+resources.KeyStoreFile+test.keyDir
		pri, err := account.PriKeyFromKeyStore(test.address, path, test.password)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log("pri:", pri)
	}
}

func TestPriKeyToKeyStore(t *testing.T) {
	for _, test := range []struct {
		keyDir     string
		isSM2      bool
		privateKey string
		password   string
		chainCode  string
	}{
		{"../resources/keystore", false, resources.Addr1Pri, resources.PassWord, resources.ChainCode},
	} {
		isSuccess, err := account.PriKeyToKeyStore(test.keyDir, test.isSM2, test.privateKey, test.password, test.chainCode)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if isSuccess {
			t.Log("privateKey Recipient KeyStoreFile success")
		}
	}

}

func TestPriKeyToAccount(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		isSM2      bool
		chainCode  string
	}{
		{resources.Addr1Pri, false, resources.ChainCode},
		{resources.Addr2Pri, false, resources.ChainCode},
	} {
		accountAddress, err := account.PriKeyToAccount(test.privateKey, test.isSM2, test.chainCode)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log("accountAddress: ", accountAddress)
	}
}

func TestPriKeyToPublicKey(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		isSM2      bool
		want       string
	}{
		{resources.Addr1Pri, false, resources.Addr1Hex},
		{resources.Addr2Pri, false, resources.Addr2Hex},
	} {
		publicKey, err := account.PriKeyToPublicKey(test.privateKey, test.isSM2)
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

func TestGetPublicKeyFromFile(t *testing.T) {
	for _, test := range []struct {
		privateKeyFilePath string
		password           string
		isSM2              bool
		want               string
	}{
		{"UTC--2021-10-19T05-33-49.419105162Z--did_bid_qwer_sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", "tele", false, "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"},
	} {
		file := bif.GetCurrentAbPath()+resources.KeyStoreFile+test.privateKeyFilePath
		publicKey, err := account.PublicKeyFromFile(file, test.password, test.isSM2)
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
func TestGenNodeUrlFromKeyStore(t *testing.T) {
	for _, test := range []struct {
		nodePrivateKeyPath string
		password           string
		host               string
		port               uint64
	}{
		{"UTC--2021-10-19T05-33-49.419105162Z--did_bid_qwer_sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", "tele", "127.0.0.1", 55555},
	} {
		file := bif.GetCurrentAbPath()+resources.KeyStoreFile+test.nodePrivateKeyPath
		nodeUrl, err := account.GenNodeUrlFromKeyStore(file, test.password, test.host, test.port)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("nodeUrl is : %s \n", nodeUrl)
	}

}

func TestGenNodeUrlFromPriKey(t *testing.T) {
	for _, test := range []struct {
		privateKey string
		host       string
		port       uint64
	}{
		{resources.Addr1Pri, resources.IP00, resources.Port},
		{resources.Addr2Pri, resources.IP00, resources.Port},
	} {
		nodeUrl, err := account.GenNodeUrlFromPriKey(test.privateKey, test.host, test.port)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("nodeUrl is : %s \n", nodeUrl)
	}

}
