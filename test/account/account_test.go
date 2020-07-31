package account

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
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
		{"国密交易签名解析", resources.AddressSM2, true, "0xf8810423831e8480946469643a6269643a7368eb1b26408f44721e20bd946469643a6269643ac935bd29a90fbeea87badf3e850ba43b740080808202bda05b166e6cd370a403b08b190a4264b9f874058211fc710cc0c9ef183d57ee6368a07a15d5f87b74bb8a2711ed10c305dde979c435bd03ee5cd42fb644e4822c07f880808080"},
		// 这个为什么会抛出panic？？？？  不应该提示错误吗？
		// E:\code\go\pkg\mod\github.com\teleinfo-bif\bit-gmsm@v1.0.5\sm2\sm2.go
		// E:\code\go\pkg\mod\github.com\teleinfo-bif\bit-gmsm@v1.0.5\sm2\p256.go
		{"非国密交易签名解析", resources.CoinBase, false, "0xf88382936e23831e8480946469643a6269643ac935bd29a90fbeea87badf3e946469643a6269643a7368eb1b26408f44721e20bd850ba43b740080018202bea067757f3345009cd5f3edb4c7a3db6bde6142fd374aea15dcb189aa838010401da02493c15c6eab52b82a19fa0963eee6cc44935b2cd482fab9189248332f171afd80808080"},
	} {
		acc := account.NewAccount()
		recoverAddress, err := acc.RecoverTransaction(test.rawTx, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// fmt.Println(common.HexToAddress(test.signAddress).String())
		if common.HexToAddress(test.signAddress) != common.HexToAddress(recoverAddress) {
			t.Error("signAddress is not match recoverAddress")
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

		from := common.StringToAddress(sender)
		if from != common.StringToAddress(test.sender) {
			t.Errorf("Address and private key do not match")
			t.FailNow()
		}

		to := common.StringToAddress(test.to)
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
			NT:           new(big.Int),
			NV:           new(big.Int),
			NR:           new(big.Int),
			NS:           new(big.Int),
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
