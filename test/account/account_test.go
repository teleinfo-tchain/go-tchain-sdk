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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	accountAddress, privateKey, err := connection.Account.Create(false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("account is %s, privateKey is %s", accountAddress, privateKey)
}

func TestPrivateKeyToAccount(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	accountAddress, privateKey, err := connection.Account.Create(false)
	parsedAccount, err := connection.Account.PrivateKeyToAccount(privateKey, false)
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	privateKey := "683179891753cad00325ae51a7c274ad1f39563ca56e07e4d86c2ef5e9ab82b3"
	password := "teleinfo"
	keyJson, err := connection.Account.Encrypt(privateKey, false, password, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("keyJson is ", keyJson)
}

func TestDecrypt(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	privateKeyFile := "../resources/superNodeKeyStore/UTC--172.17.6.51--did-bid-590ed37615bdfefa496224c7"
	password := "teleinfo"
	keyJson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	address, privateKey, err := connection.Account.Decrypt(keyJson, false, password)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("PublicAddress is ", address)
	t.Log("PrivateKey is ", privateKey)
}

func TestRecoverTransaction(t *testing.T) {
	for _, test := range []struct {
		id          string
		signAddress string
		isSM2       bool
		rawTx       string
	}{
		// {"国密交易签名解析", resources.AddressSM2, true, "0xf87d0623831e8480946469643a6269643a7368eb1b26408f44721e20bd946469643a6269643ac935bd29a90fbeea87badf3e850ba43b740080808202bea07b5228d865d1ee9dc1dd33297f3a0e4442a31430dd0c22538f0137ee37fdc8e9a09a99f1cb34c437a4640a22d9cec6f57ead94d6734fd417d56a8a7d66bdc7bd6b"},
		// 这个为什么会抛出panic？？？？  不应该提示错误吗？ 是国密解密的初始化问题，国密那里代码需要修改
		// E:\code\go\pkg\mod\github.com\teleinfo-bif\bit-gmsm@v1.0.5\sm2\sm2.go
		// E:\code\go\pkg\mod\github.com\teleinfo-bif\bit-gmsm@v1.0.5\sm2\p256.go
		// 注意使用不同的chainId签名会得到不同的签名交易，因此解析可能报错（因为签名的chainId变了）。
		{"非国密交易签名解析", resources.CoinBase, false, "0xf87d8023831e8480946469643a6269643ac935bd29a90fbeea87badf3e946469643a6269643a7368eb1b26408f44721e20bd850ba43b74008001820558a018eaaade5d45e1d461ddb8d763ec7d5f5b12bf17f36daa4bc3cdefae3a395fdba00d95dc9143c1d3cde74b3ddcb7a87e562ca89fca9dbe1530a73ba38e2a15cfb7"},
		// 这里试着对从GetRawTransactionByHash中获得hash获取；
		{"非国密交易签名解析", resources.CoinBase, false, "0xf8bf8292fe19829c40946469643a6269643ac935bd29a90fbeea87badf3e946469643a6269643ac935bd29a90fbeea87badf3e0a80018202bda0e96150ad4fd738be89a982dd4a0d739fc3992ebfe48c63ce07e0b5249f03f57ca0481e40b26a7bd13a5b9a5a1429c761152c2a5b44bb7d96ec7e6352cd0952908b018202bda0e96150ad4fd738be89a982dd4a0d739fc3992ebfe48c63ce07e0b5249f03f57ca0481e40b26a7bd13a5b9a5a1429c761152c2a5b44bb7d96ec7e6352cd0952908b"},
	} {
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
		recoverAddress, err := connection.Account.RecoverTransaction(test.rawTx, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// fmt.Println(common.HexToAddress(test.signAddress).String())
		if utils.HexToAddress(test.signAddress) != utils.HexToAddress(recoverAddress) {
			t.Error(test.id + ": signAddress is not match recoverAddress")
			t.FailNow()
		}

		t.Logf("%s， 交易签名地址： %s", test.id, recoverAddress)
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
		// {"国密签名转账", resources.AddressSM2, resources.AddressPriKey, resources.CoinBase, true},
		{"非国密签名转账", resources.CoinBase, resources.CoinBasePriKey, resources.AddressSM2, false},
	} {
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
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

		// tx := &account.SignTxParams{
		// 	To:       test.to,
		// 	Nonce:    nonce.Uint64(),
		// 	Gas:      2000000,
		// 	GasPrice: big.NewInt(30),
		// 	Value:    big.NewInt(50000000000),
		// 	Data:     nil,
		// 	ChainId:  big.NewInt(0).SetUint64(chainId),
		// }
		t.Log("nonce is ", nonce, " chainId is ", chainId)
		tx := &account.SignTxParams{
			To:       test.to,
			Nonce:    0,
			Gas:      21000,
			GasPrice: nil,
			Value:    big.NewInt(50000000000),
			Data:     nil,
			ChainId:  nil,
		}
		res, err := connection.Account.SignTransaction(tx, test.privateKey, test.isSM2)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("%s : %v \n", test.id, res.Raw)
		// // t.Logf("%#v \n", res.Tx.Hash)
		//
		// txHash, err := connection.Core.SendRawTransaction(res.Raw.String())
		// if err != nil {
		// 	t.Error(err)
		// 	t.FailNow()
		// }
		// t.Logf("txHash is %s", txHash)

	}
}

func TestHashMessage(t *testing.T) {
	for _, test := range []struct {
		message     string
		isSm2       bool
		hashMessage string
	}{
		{"Hello World", false, "0xa1de988600a42c4b4ab089b619297c17d53cffae5d5120d82d8a92d0bb3b78f2"},
		{utils.NewUtils().Utf8ToHex("Hello World"), false, "0xa1de988600a42c4b4ab089b619297c17d53cffae5d5120d82d8a92d0bb3b78f2"},
	} {
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
		res := connection.Account.HashMessage(test.message, test.isSm2)
		if res != test.hashMessage {
			t.Errorf("hash error, input is %s, result is %s, expect is %s \n", test.message, res, test.hashMessage)
			t.FailNow()
		}

	}
}

func TestSign(t *testing.T) {
	for _, test := range []struct {
		id         string
		isSm2      bool
		privateKey string
		message    string
	}{
		// {"国密签名", true, resources.AddressPriKey, "Some data"},
		{"非国密签名", false, resources.CoinBasePriKey, "Some data"},
	} {
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

		chainId, err := connection.Core.GetChainId()
		if err != nil {
			t.Errorf("Failed getChainId")
			t.FailNow()
		}

		messageHash := connection.Account.HashMessage(test.message, test.isSm2)
		sign := &account.SignData{
			Message:     test.message,
			MessageHash: messageHash,
			V:           new(big.Int),
			R:           new(big.Int),
			S:           new(big.Int),
			T:           big.NewInt(0),
			// NT:           new(big.Int),
			// NV:           new(big.Int),
			// NR:           new(big.Int),
			// NS:           new(big.Int),
		}

		res, err := connection.Account.Sign(sign, test.privateKey, test.isSm2, big.NewInt(0).SetUint64(chainId))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("%s : %v \n", test.id, res)
	}
}
