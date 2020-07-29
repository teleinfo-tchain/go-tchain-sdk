package account

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/common"
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
		// {"国密交易签名解析", resources.AddressSM2, true, "0xf87f82936d23831e8480946469643a6269643a7368eb1b26408f44721e20bd946469643a6269643ac935bd29a90fbeea87badf3e850ba43b740080808202bda0283d1cb7b8800b8ad4af0251ac8add797fddb0960f3a6190420542c58bd3c283a07c0dde04349f05737df8f67562e2babfdfe3b64a22e54c18e0e65e935748429a"},
		{"非国密交易签名解析", resources.CoinBase, false, "0xf87f82936d23831e8480946469643a6269643ac935bd29a90fbeea87badf3e946469643a6269643a7368eb1b26408f44721e20bd850ba43b740080018202bda0578eaf57135aa31984dcc8fb6aa825550dd4ae56f69bfc3595636df2a5d98e38a05f9acfaa1cf9bd6a32be87c8b41eecfacf25fc5cffa18e6c803d0cd893a93466"},
		// 这个为什么会抛出panic？？？？  不应该提示错误吗？
		// E:\code\go\pkg\mod\github.com\teleinfo-bif\bit-gmsm@v1.0.5\sm2\sm2.go
		// E:\code\go\pkg\mod\github.com\teleinfo-bif\bit-gmsm@v1.0.5\sm2\p256.go
		// {resources.AddressSM2, true, "0xf87f82936d23831e8480946469643a6269643ac935bd29a90fbeea87badf3e946469643a6269643a7368eb1b26408f44721e20bd850ba43b740080018202bda0578eaf57135aa31984dcc8fb6aa825550dd4ae56f69bfc3595636df2a5d98e38a05f9acfaa1cf9bd6a32be87c8b41eecfacf25fc5cffa18e6c803d0cd893a93466"},
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	nonce, err := connection.Core.GetTransactionCount(resources.CoinBase, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Errorf("Failed getChainId")
		t.FailNow()
	}

	for _, test := range []struct {
		id         string
		cryType    uint
		sender     string
		privateKey string
		to         string
	}{
		{"国密签名", 0, resources.AddressSM2, resources.AddressPriKey, resources.CoinBase},
		// {"非国密签名", 1, resources.CoinBase, resources.CoinBasePriKey, resources.AddressSM2},
	} {
		sender, err := utils.GetAddressFromPrivate(test.privateKey, test.cryType)
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

	// rawTxString := "0xf87f82936d23831e8480946469643a6269643a7368eb1b26408f44721e20bd946469643a6269643ac935bd29a90fbeea87badf3e850ba43b740080808202bda0dd9fa03c2c67816bbac4b3ca3e29ed01116e3ab7a9e0afc0d3481853ca5ed2b5a017630ed8aef1b7be5fcd69674bab68b17698ace7a1a521ba45628caedf800a08"
	// if rawTxString[:2] == "0x" || rawTxString[:2] == "0X" {
	// 	rawTxString = rawTxString[2:]
	// }
	//
	// rawTx, _ := hex.DecodeString(rawTxString)
	//
	// var tx account.TxData
	// rawTx, _ = hex.DecodeString(rawTxString)
	//
	// err := rlp.DecodeBytes(rawTx, &tx)
	//
	// sigHash := rlpHash([]interface{}{
	// 	tx.AccountNonce,
	// 	tx.Price,
	// 	tx.GasLimit,
	// 	tx.Sender,
	// 	tx.Recipient,
	// 	tx.Amount,
	// 	tx.Payload,
	// 	big.NewInt(333), uint(0), uint(0),
	// })
	// privKey, err := crypto.HexToECDSA(resources.AddressPriKey, crypto.SM2)
	// if err != nil{
	// 	fmt.Println("err is ", err)
	// }
	// msg := sigHash.Bytes()
	// // sig := []byte{13,105,24,131,219,111,82,44,204,35,60,54,188,208,54,65,94,63,27,124,146,108,217,126,70,247,57,195,131,0,15,21,235,161,66,96,46,92,196,200,216,155,98,103,108,87,15,73 ,37, 52 ,137, 243, 65, 184, 94 ,91 ,134, 38, 81 ,107 ,63 ,120, 69 ,208 ,0 ,0}
	// sig := []byte{221,159 ,160, 60, 44, 103, 129 ,107 ,186 ,196 ,179, 202 ,62 ,41 ,237, 1 ,17 ,110 ,58, 183, 169, 224, 175 ,192, 211, 72 ,24 ,83, 202 ,94 ,210 ,181 ,23 ,99, 14, 216, 174, 241, 183, 190 ,95, 205, 105, 103, 75, 171 ,104 ,177, 118, 152 ,172, 231, 161, 165, 33, 186, 69,98, 140, 174, 223, 128, 10 ,8 ,0, 0}
	// pubkey, err := sm2.RecoverPubKey(msg, sig[:65])
	// if err != nil {
	// 	fmt.Println("err is ", err)
	// }
	// fmt.Println(pubkey)

}

// func rlpHash(x interface{}) (h common.Hash) {
// 	hw := sha3.NewLegacyKeccak256()
// 	rlp.Encode(hw, x)
// 	hw.Sum(h[:0])
// 	return h
// }
