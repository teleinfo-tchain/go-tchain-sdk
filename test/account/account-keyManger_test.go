package account

import (
	"github.com/bif/bif-sdk-go/account"
	"testing"
)

func TestGenerateKeyStore(t *testing.T) {
	addr, err := account.GenerateKeyStore("./keystore", 1, "123456", false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("addr:", addr)
}

func TestGetPrivateKeyFromFile(t *testing.T) {
	pri, err := account.GetPrivateKeyFromFile("did:bid:c6a0fa74323e4fd57e78e19a", "./keystore/UTC--2020-04-28T10-29-33.363515000Z--did-bid-c6a0fa74323e4fd57e78e19a", "123456")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("pri:", pri)
}

func TestPrivateKeyToKeyStoreFile(t *testing.T) {
	isSuccess, err := account.PrivateKeyToKeyStoreFile("./keystore", 0, "eea9354b98fd51d7b962cb4c7e61d691e4d540951e1cf277dd72f2a37544c1da", "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if isSuccess{
		t.Log("privateKey To KeyStoreFile success")
	}

}

func TestGetAddressFromPrivate(t *testing.T) {
	accountAddress, err := account.GetAddressFromPrivate("eea9354b98fd51d7b962cb4c7e61d691e4d540951e1cf277dd72f2a37544c1da", 0)
	if err != nil {
		t.Error("get address from private error")
		t.FailNow()
	}
	t.Log("accountAddress: ", accountAddress)
}

func TestGetPublicKeyFromPrivate(t *testing.T) {
	publicKey, err := account.GetPublicKeyFromPrivate("eea9354b98fd51d7b962cb4c7e61d691e4d540951e1cf277dd72f2a37544c1da", 0)
	if err != nil {
		t.Error("get publicKey from private error")
		t.FailNow()
	}
	t.Logf("publicKey is : %s, len is %d", publicKey, len(publicKey[2:]))
}
