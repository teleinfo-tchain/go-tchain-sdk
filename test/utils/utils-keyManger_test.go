package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go/utils"
	"testing"
)

func TestGenerateKeyStore(t *testing.T) {
	addr, err := utils.GenerateKeyStore("./keystore", 1, "123456", false)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("addr:", addr)
}

func TestGetPrivateKeyFromFile(t *testing.T) {
	pri, pub, err := utils.GetPrivateKeyFromFile("did:bid:c6a0fa74323e4fd57e78e19a", "./keystore/UTC--2020-04-28T10-29-33.363515000Z--did-bid-c6a0fa74323e4fd57e78e19a", "123456")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("pri:", pri)
	fmt.Println("pub:", pub)
}

func TestPrivateKeyToKeyStoreFile(t *testing.T) {
	pub, err := utils.PrivateKeyToKeyStoreFile("./keystore", "did:bid:73f6a70d05af2141dd4ad995", "eea9354b98fd51d7b962cb4c7e61d691e4d540951e1cf277dd72f2a37544c1da", "123")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("pub:", pub)
}

func TestGetAddressFromPrivate(t *testing.T) {
	pub, err := utils.GetAddressFromPrivate("eea9354b98fd51d7b962cb4c7e61d691e4d540951e1cf277dd72f2a37544c1da", 0)
	if err != nil {
		fmt.Println("getaddress from private error")
		return
	}
	fmt.Println("pub:", pub)
}
