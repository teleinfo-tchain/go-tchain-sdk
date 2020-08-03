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
	acc := account.NewAccount()
	accountAddress, privateKey, err := acc.Create(false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("account is %s, privateKey is %s", accountAddress, privateKey)
}

func TestPrivateKeyToAccount(t *testing.T) {
	acc := account.NewAccount()
	accountAddress, privateKey, err := acc.Create(false)
	parsedAccount, err := acc.PrivateKeyToAccount(privateKey, false)
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
	acc := account.NewAccount()
	privateKey := "683179891753cad00325ae51a7c274ad1f39563ca56e07e4d86c2ef5e9ab82b3"
	password := "teleinfo"
	keyJson, err := acc.Encrypt(privateKey, false, password, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("keyJson is ", keyJson)
}

func TestDecrypt(t *testing.T) {
	acc := account.NewAccount()
	privateKeyFile := "../resources/superNodeKeyStore/UTC--172.17.6.51--did-bid-590ed37615bdfefa496224c7"
	password := "teleinfo"
	keyJson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	address, privateKey, err := acc.Decrypt(keyJson, false, password)
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
		{"非国密交易签名解析", resources.CoinBase, false, "0xf87f82938f23831e8480946469643a6269643ac935bd29a90fbeea87badf3e946469643a6269643a7368eb1b26408f44721e20bd850ba43b740080018202bea0368f6865897245b3d2e3148bf110e81ef7a8fc8d0284700342e4f289319ac485a00a7463274634171a32d412a620df4db3a0f587ee879c1ea27353239b1d6c05c6"},
		// 这里试着对从GetRawTransactionByHash中获得hash获取；
		{"非国密交易签名解析", resources.CoinBase, false, "0xf8bf8292fe19829c40946469643a6269643ac935bd29a90fbeea87badf3e946469643a6269643ac935bd29a90fbeea87badf3e0a80018202bda0e96150ad4fd738be89a982dd4a0d739fc3992ebfe48c63ce07e0b5249f03f57ca0481e40b26a7bd13a5b9a5a1429c761152c2a5b44bb7d96ec7e6352cd0952908b018202bda0e96150ad4fd738be89a982dd4a0d739fc3992ebfe48c63ce07e0b5249f03f57ca0481e40b26a7bd13a5b9a5a1429c761152c2a5b44bb7d96ec7e6352cd0952908b"},
	} {
		acc := account.NewAccount()
		recoverAddress, err := acc.RecoverTransaction(test.rawTx, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// fmt.Println(common.HexToAddress(test.signAddress).String())
		if utils.HexToAddress(test.signAddress) != utils.HexToAddress(recoverAddress) {
			t.Error(test.id+": signAddress is not match recoverAddress")
			t.FailNow()
		}

		t.Logf("%s， 交易签名地址： %s", test.id, recoverAddress)
	}
}

func TestSignTransaction(t *testing.T) {
	for _, test := range []struct {
		id         string
		cryType    uint
		sender     string
		privateKey string
		to         string
	}{
		{"国密签名", 0, resources.AddressSM2, resources.AddressPriKey, resources.CoinBase},
		{"非国密签名", 1, resources.CoinBase, resources.CoinBasePriKey, resources.AddressSM2},
	} {
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
		nonce, err := connection.Core.GetTransactionCount(test.sender, block.LATEST)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		chainId, err := connection.Core.GetChainId()
		if err != nil {
			t.Errorf("Failed getChainId")
			t.FailNow()
		}


		sender, err := account.GetAddressFromPrivate(test.privateKey, test.cryType)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		from := utils.StringToAddress(sender)
		if from != utils.StringToAddress(test.sender) {
			t.Errorf("Address and private key do not match")
			t.FailNow()
		}

		to := utils.StringToAddress(test.to)
		tx := &account.TxData{
			AccountNonce: nonce.Uint64(),
			Price:        big.NewInt(35),
			GasLimit:     2000000,
			Sender:       &from,
			Recipient:    &to,
			Amount:       big.NewInt(50000000000),
			Payload:      nil,
			V:            new(big.Int),
			R:            new(big.Int),
			S:            new(big.Int),
			T:            big.NewInt(0),
			// NT:           new(big.Int),
			// NV:           new(big.Int),
			// NR:           new(big.Int),
			// NS:           new(big.Int),
		}

		acc := account.NewAccount()
		res, err := acc.SignTransaction(tx, test.privateKey, big.NewInt(0).SetUint64(chainId))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("%s : %v \n", test.id, res.Raw)
		// t.Logf("%#v \n", res.Raw)
	}
}
