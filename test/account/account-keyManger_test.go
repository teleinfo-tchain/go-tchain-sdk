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
		{"./keystore", true, "123456", false},
		{"./keystore", false, "123456", false},
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
		{"123456", "0x6469643A6269643a738cB17e4DB903d1bF802fEC", "./keystore/UTC--2020-08-08T16-09-37.005318300Z--did-bid-738cb17e4db903d1bf802fec"},
		{"123456", "0x6469643A6269643a46727D4aA2Bf54976C22b6d4", "./keystore/UTC--2020-08-08T16-09-38.381643000Z--did-bid-46727d4aa2bf54976c22b6d4"},
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
		{"123456", "d64912c32615ed55b9d9602539714fdf22bbe5c32608d5f1082a60a2bc22ac0a", true, "./keystore"},
		{"123456", "dc06dc71470401dbbd7e064a82df5b714d6e60d5dfecd694d19b40d6f6398137", false, "./keystore"},
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
		{"d64912c32615ed55b9d9602539714fdf22bbe5c32608d5f1082a60a2bc22ac0a", true},
		{"dc06dc71470401dbbd7e064a82df5b714d6e60d5dfecd694d19b40d6f6398137", false},
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
	}{
		{"d64912c32615ed55b9d9602539714fdf22bbe5c32608d5f1082a60a2bc22ac0a", true},
		{"dc06dc71470401dbbd7e064a82df5b714d6e60d5dfecd694d19b40d6f6398137", false},
	} {
		publicKey, err := account.GetPublicKeyFromPrivate(test.privateKey, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("len is %d , publicKey is : %s \n", len(publicKey[2:]), publicKey)
	}

}

func TestGetGetPublicKeyFromFile(t *testing.T) {
	for _, test := range []struct {
		privateKeyFilePath string
		password           string
		isSM2              bool
	}{
		{"./keystore/UTC--2020-08-08T16-09-37.005318300Z--did-bid-738cb17e4db903d1bf802fec", "123456", true},
		{"./keystore/UTC--2020-08-08T16-09-38.381643000Z--did-bid-46727d4aa2bf54976c22b6d4", "123456", false},
	} {
		publicKey, err := account.GetPublicKeyFromFile(test.privateKeyFilePath, test.password, test.isSM2)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// 0x01cea5a1d9b901bb2b3306e8cf28762dbed3c9a886c40b5c3860582a9bdf51c9ed
		// 0x042466cb80265bcf5669f65f041c31e433cf8328a3070bc7d20279015c61b27fe6826f235eda0baacbff7e5e19567995410ecb3efcc280d7daaa44331b5c5027b9
		t.Logf("len is %d , publicKey is : %s \n", len(publicKey[2:]), publicKey)
	}

}
