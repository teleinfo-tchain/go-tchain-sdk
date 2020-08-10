package account

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/bif/bif-sdk-go/account/keystore"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/utils"
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

/*
  GenerateKeyStore:
   	EN - Generate private key file
	CN - 生成私钥文件
  Params:
  	- keyDir string          私钥文件生成的存储地址
  	- isSM2   bool            是否采用国密生成私钥，true为是国密，false为否
  	- password string        私钥文件密码，用于加密私钥
  	- UseLightweightKDF bool 一般选择false；如果是true会降低密钥库的内存和CPU要求,不过是以牺牲安全性为代价

  Returns:
  	- string  账户地址
	- error

  Call permissions: Anyone
*/
func GenerateKeyStore(keyDir string, isSM2 bool, password string, UseLightweightKDF bool) (string, error) {

	scryptN, scryptP, err := preCheck(keyDir, password, UseLightweightKDF)
	if err != nil {
		return "", err
	}
	var cryptoType crypto.CryptoType
	if isSM2 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}

	address, err := keystore.StoreKey(keyDir, password, scryptN, scryptP, cryptoType)
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}

/*
  描述:
   	EN - Get the private key and account from the keyStore file
	CN - 从文件获取私钥和账户地址
  Params:
  	- address string         要解析的账户地址
  	- privateKeyFile string  keyStore文件地址
  	- password string        keyStore加密密码

  Returns:
  	- string privateKey      私钥
	- error

  Call permissions: Anyone
*/
func GetPrivateKeyFromFile(address string, privateKeyFile, password string) (string, error) {
	keyJson, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return "", err
	}
	var key *keystore.Key

	addr := utils.StringToAddress(address)
	if bytes.HasPrefix(addr.Bytes(), []byte("did:bid:")) && addr[8] == 115 {
		key, err = keystore.DecryptKey(keyJson, password, crypto.SM2)
	} else {
		key, err = keystore.DecryptKey(keyJson, password, crypto.SECP256K1)
	}
	if err != nil {
		return "", err
	}
	privateKey := hex.EncodeToString(key.PrivateKey.D.Bytes())
	addrRes := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)
	if addr != addrRes {
		return "", errors.New("address Not Match keyStoreFile")
	}
	return privateKey, nil
}

/*
  PrivateKeyToKeyStoreFile(默认私钥文件生成采用UseLightweightKDF为false即安全模式):
   	EN - Transfer private key to keyStore file
	CN - 私钥转keyStore文件
  Params:
  	- keyDir string      keyStore文件地址
  	- isSM2   bool       是否采用国密生成私钥，true为是国密，false为否
  	- privateKey string  私钥
  	- password string    keyStore加密密码

  Returns:
  	- bool     true为转换成功,false为转换失败
	- error

  Call permissions: Anyone
*/
func PrivateKeyToKeyStoreFile(keyDir string, isSM2 bool, privateKey string, password string) (bool, error) {
	if utils.Has0xPrefix(privateKey){
		privateKey = privateKey[2:]
	}

	if !utils.IsHex(privateKey) || len(privateKey) != 64{
		return false, errors.New("privateKey is not hex string or not 32 Bytes")
	}

	var cryptoType crypto.CryptoType
	if isSM2 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}

	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return false, err
	}

	scryptN, scryptP, err := preCheck(keyDir, password, false)
	if err != nil {
		return false, err
	}

	key := keystore.NewKeyStore(keyDir, scryptN, scryptP)
	_, err = key.ImportECDSA(privateKeyN, password)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
  GetAddressFromPrivate:
   	EN - Get address based on private key
	CN - 根据私钥获取账户地址
  Params:
  	- privateKey string  私钥
  	- isSM2      bool    是否采用国密生成私钥，true为是国密，false为否

  Returns:
  	- string  账户地址
	- error

  Call permissions: Anyone
*/
func GetAddressFromPrivate(privateKey string, isSM2 bool) (string, error) {
	if utils.Has0xPrefix(privateKey){
		privateKey = privateKey[2:]
	}

	if !utils.IsHex(privateKey) || len(privateKey) != 64{
		return "", errors.New("privateKey is not hex string or not 32 Bytes")
	}

	var cryptoType crypto.CryptoType
	if isSM2 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}

	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	// 转换成地址
	return crypto.PubkeyToAddress(*privateKeyN.Public().(*ecdsa.PublicKey)).Hex(), nil
}

/*
  GetPublicKeyFromPrivate:
   	EN - Get the public key based on the private key
	CN - 根据私钥获取公钥
  Params:
  	- privateKey string   私钥
  	- isSM2      bool     是否采用国密生成私钥，true为是国密，false为否

  Returns:
  	- string  公钥（65字节）
	- error

  Call permissions: Anyone
*/
func GetPublicKeyFromPrivate(privateKey string, isSM2 bool) (string, error) {
	if utils.Has0xPrefix(privateKey){
		privateKey = privateKey[2:]
	}

	if !utils.IsHex(privateKey) || len(privateKey) != 64{
		return "", errors.New("privateKey is not hex string or not 32 Bytes")
	}

	var cryptoType crypto.CryptoType
	if isSM2 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}

	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	pubBytes := crypto.FromECDSAPub(privateKeyN.Public().(*ecdsa.PublicKey))
	return "0x" + utils.Bytes2Hex(pubBytes), nil
}

/*
  GetPublicKeyFromFile:
   	EN -
	CN - 根据keyStore获取公钥
  Params:
  	- privateKeyFilePath string   keyStore的存储路径
  	- password           string   keyStore的加密密码
	- isSM2              bool     是否采用国密生成私钥，true为是国密，false为否

  Returns:
  	- string  公钥（65字节）
	- error

  Call permissions: Anyone
*/
func GetPublicKeyFromFile(privateKeyFilePath, password string, isSM2 bool) (string, error) {
	keyJson, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return "", errors.New("not find privateKeyFile")
	}

	var cryptoType crypto.CryptoType
	if isSM2 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}

	key, err := keystore.DecryptKey(keyJson, password, cryptoType)
	if err != nil {
		return "", err
	}
	privateKey := hex.EncodeToString(key.PrivateKey.D.Bytes())
	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	pubBytes := crypto.FromECDSAPub(privateKeyN.Public().(*ecdsa.PublicKey))
	if isSM2{
		return "0x" + utils.Bytes2Hex(pubBytes)+utils.Bytes2Hex(privateKeyN.Public().(*ecdsa.PublicKey).Y.Bytes()), nil
	}else{
		return "0x" + utils.Bytes2Hex(pubBytes), nil
	}
}
