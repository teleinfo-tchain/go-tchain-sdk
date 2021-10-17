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
	"strconv"
	"testing"
)

func TestCreate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	accountAddress, privateKey, err := connection.Account.Create(false, resources.ChainCode)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("account is %s, privateKey is %s", accountAddress, privateKey)
}

func TestPrivateKeyToAccount(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	accountAddress, privateKey, err := connection.Account.Create(false, resources.ChainCode)
	parsedAccount, err := connection.Account.PrivateKeyToAccount(privateKey, resources.ChainCode, false)
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
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
		keyJson, err := connection.Account.Encrypt(test.privateKey, test.isSM2, test.password, resources.ChainCode, false)
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
		{true, resources.PassWord, "./keystore/UTC--2020-08-19T05-48-44.625362500Z--did-bid-ZFT4CziA2ktCNgfQPqSm1GpQxSck5q4", "did:bid:ZFT4CziA2ktCNgfQPqSm1GpQxSck5q4"},
		{false, resources.PassWord, "./keystore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb", "did:bid:EFTTQWPMdtghuZByPsfQAUuPkWkWYb"},
	} {
		keyJson, err := ioutil.ReadFile(test.keyDir)
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
		{"国密签名转账", resources.Addr1, resources.Addr1Pri, resources.Addr2, true},
	} {
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
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
		// 	Recipient:       test.to,
		// 	AccountNonce:    nonce.Uint64(),
		// 	GasPrice:      2000000,
		// 	GasLimit: big.NewInt(30),
		// 	Amount:    big.NewInt(50000000000),
		// 	Payload:     nil,
		// 	ChainId:  chainId,
		// }
		t.Log("nonce is ", nonce, " chainId is ", chainId)
		var recipient utils.Address
		recipient = utils.StringToAddress(test.to)
		tx := &account.SignTxParams{
			Recipient:    &recipient,
			AccountNonce: 0,
			GasPrice:     big.NewInt(21000),
			GasLimit:     0,
			Amount:       big.NewInt(50000000000),
			Payload:      nil,
			ChainId:      0,
		}
		res, err := connection.Account.SignTransaction(tx, test.privateKey, test.isSM2)

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
		{utils.NewUtils().Utf8ToHex("Hello World"), false, "0xa1de988600a42c4b4ab089b619297c17d53cffae5d5120d82d8a92d0bb3b78f2"},
	} {
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
		res := connection.Account.HashMessage(test.message, test.isSm2)
		if res != test.hashMessage {
			t.Errorf("hash error, input is %s, result is %s, expect is %s \n", test.message, res, test.hashMessage)
			t.FailNow()
		}

	}
}
