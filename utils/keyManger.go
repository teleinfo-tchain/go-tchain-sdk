package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/utils/keystore"
	"io/ioutil"
)

func preCheck(keyDir string, password string, UseLightweightKDF bool) (int, int, error) {
	if keyDir == "" {
		return -1, -1, errors.New("empty keyDir, please check")
	}
	if password == "" {
		return -1, -1, errors.New("empty password, please check")
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}
	return scryptN, scryptP, nil
}

// 1. 私钥文件生成
func GenerateKeyStore(keyDir string, cryType uint, password string, UseLightweightKDF bool) (string, error) {

	scryptN, scryptP, err := preCheck(keyDir, password, UseLightweightKDF)
	if err != nil {
		return "", err
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
	address, err := keystore.StoreKey(keyDir, password, scryptN, scryptP, cryptoType)
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}

// 2. 从文件获取私钥和地址
func GetPrivateKeyFromFile(addrParse string, privateKeyFile, password string) (string, string, error) {
	keyJson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return "", "", err
	}
	var key *keystore.Key

	addr := common.StringToAddress(addrParse)
	if bytes.HasPrefix(addr.Bytes(), []byte("did:bid:")) && addr[8] == 115 {
		key, err = keystore.DecryptKey(keyJson, password, crypto.SM2)
	} else {
		key, err = keystore.DecryptKey(keyJson, password, crypto.SECP256K1)
	}
	if err != nil {
		return "", "", err
	}
	privateKey := hex.EncodeToString(key.PrivateKey.D.Bytes())
	addrRes := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)
	if addrParse != common.ByteAddressToString(addrRes.Bytes()) {
		return "", "", errors.New("addrParse Not Match keyStoreFile")
	}
	return privateKey, addrRes.String(), nil
}

// 3. 私钥转文件
func PrivateKeyToKeyStoreFile(keyDir string, addrParse string, privateKey string, password string) (string, error) {
	addr := common.StringToAddress(addrParse)
	var cryptoType crypto.CryptoType
	if bytes.HasPrefix(addr.Bytes(), []byte("did:bid:")) && addr[8] == 115 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}
	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	addrRes := crypto.PubkeyToAddress(privateKeyN.PublicKey)
	if addrParse != common.ByteAddressToString(addrRes.Bytes()) {
		return "", errors.New("addrParse Not Match privateKey")
	}

	scryptN, scryptP, err := preCheck(keyDir, password, false)
	if err != nil {
		return "", err
	}
	key := keystore.NewKeyStore(keyDir, scryptN, scryptP)
	account, err := key.ImportECDSA(privateKeyN, password)
	if err != nil {
		return "", err
	}
	return account.Address.Hex(), nil
}

// 4.根据私钥获取地址
func GetAddressFromPrivate(privateKey string, cryType uint) (string, error) {
	var cryptoType crypto.CryptoType
	switch cryType {
	case 0:
		cryptoType = crypto.SM2
	case 1:
		cryptoType = crypto.SECP256K1
	default:
		cryptoType = crypto.SM2
	}
	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	// 转换成地址
	return crypto.PubkeyToAddress(*privateKeyN.Public().(*ecdsa.PublicKey)).Hex(), nil
}
