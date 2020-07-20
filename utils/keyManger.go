package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/utils/keystore"
	"io/ioutil"
)

func preCheck(keydir string, password string, UseLightweightKDF bool) (bool, int, int) {
	if keydir == "" {
		fmt.Println("empty keydir, please check")
		return false, -1, -1
	}
	if password == "" {
		fmt.Println("empty password, please check")
		return false, -1, -1
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}
	return true, scryptN, scryptP
}

//1. 私钥文件生成
func GenerateKeyStore(keydir string, cryType uint, password string, UseLightweightKDF bool) (string, error) {

	check, scryptN, scryptP := preCheck(keydir, password, UseLightweightKDF)
	if !check {
		return "", errors.New("preCheck err")
	}
	var cryptoType crypto.CryptoType
	switch cryType {
	case 0:
		cryptoType = crypto.SM2
	case 1:
		cryptoType = crypto.SECP256K1
	default:
		cryptoType = crypto.SM2
	}
	address, err := keystore.StoreKey(keydir, password, scryptN, scryptP, cryptoType)
	if err != nil {
		fmt.Println("Failed to create account:", err)
		return "", err
	}
	fmt.Printf("Address: {%s%x}\n", address[:8], address[8:])
	return address.Hex(), nil
}

//2. 从文件获取私钥和地址
func GetPrivateKeyFromFile(addrParse string, privateKeyFile, password string) (string, string, error) {
	keyjson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		fmt.Println("read keyjson file failed: ", err)
		return "", "", err
	}
	var key *keystore.Key

	addr := common.StringToAddress(addrParse)
	if bytes.HasPrefix(addr.Bytes(), []byte("did:bid:")) && addr[8] == 115 {
		key, err = keystore.DecryptKey(keyjson, password, crypto.SM2)
	} else {
		key, err = keystore.DecryptKey(keyjson, password, crypto.SECP256K1)
	}
	if err != nil {
		return "", "", err
	}
	privKey := hex.EncodeToString(key.PrivateKey.D.Bytes())
	addrRes := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)
	if addrParse != common.ByteAddressToString(addrRes.Bytes()) {
		fmt.Println("addrParse Not Match keyStoreFile")
		return "", "", errors.New("addrParse Not Match keyStoreFile")
	}
	return privKey, addrRes.String(), nil
}

//3. 私钥转文件
func PrivateKeyToKeyStoreFile(keydir string, addrParse string, privKey string, password string) (string, error) {
	addr := common.StringToAddress(addrParse)
	var cryptoType crypto.CryptoType
	if bytes.HasPrefix(addr.Bytes(), []byte("did:bid:")) && addr[8] == 115 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}
	privateKey, err := crypto.HexToECDSA(privKey, cryptoType)
	if err != nil {
		return "", err
	}
	addrRes := crypto.PubkeyToAddress(privateKey.PublicKey)
	if addrParse != common.ByteAddressToString(addrRes.Bytes()) {
		fmt.Println("addrParse Not Match privKey")
		return "", errors.New("addrParse Not Match privKey")
	}

	check, scryptN, scryptP := preCheck(keydir, password, false)
	if !check {
		return "", errors.New("preCheck err")
	}
	key := keystore.NewKeyStore(keydir, scryptN, scryptP)
	account, err := key.ImportECDSA(privateKey, password)
	if err != nil {
		return "", err
	}
	fmt.Printf("Address: {%s%x}\n", account.Address[:8], account.Address[8:])
	return account.Address.Hex(), nil
}

//4.根据私钥获取地址
func GetAddressFromPrivate(privKey string, cryType uint) (string, error) {
	var cryptoType crypto.CryptoType
	switch cryType {
	case 0:
		cryptoType = crypto.SM2
	case 1:
		cryptoType = crypto.SECP256K1
	default:
		cryptoType = crypto.SM2
	}
	privateKey, err := crypto.HexToECDSA(privKey, cryptoType)
	if err != nil {
		return "", err
	}
	//转换成地址
	return crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey)).Hex(), nil
}
