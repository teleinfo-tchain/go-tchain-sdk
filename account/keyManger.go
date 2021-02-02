package account

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/account/keystore"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/crypto/config"
	"github.com/bif/bif-sdk-go/utils"
	libp2pcorecrypto "github.com/libp2p/go-libp2p-core/crypto"
	libp2pcorepeer "github.com/libp2p/go-libp2p-core/peer"
	"io/ioutil"
	"regexp"
)

func preCheck(keyStorePath string, password string, UseLightweightKDF bool) (int, int, error) {
	if keyStorePath == "" {
		return -1, -1, errors.New("empty keyStorePath, please check")
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
  	- keyStorePath string          私钥文件生成的存储地址
  	- isSM2   bool            是否采用国密生成私钥，true为是国密，false为否
  	- password string        私钥文件密码，用于加密私钥
  	- UseLightweightKDF bool 一般选择false；如果是true会降低密钥库的内存和CPU要求,不过是以牺牲安全性为代价

  Returns:
  	- string  账户地址
	- error

  Call permissions: Anyone
*/
func GenerateKeyStore(keyStorePath string, isSM2 bool, password, chainCode string, UseLightweightKDF bool) (string, error) {

	scryptN, scryptP, err := preCheck(keyStorePath, password, UseLightweightKDF)
	if err != nil {
		return "", err
	}
	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	address, err := keystore.StoreKey(keyStorePath, password, chainCode, scryptN, scryptP, cryptoType)
	if err != nil {
		return "", err
	}
	return address.String(chainCode), nil
}

/*
  描述:
   	EN - Get the private key and account from the keyStore file
	CN - 从文件获取私钥和账户地址
  Params:
  	- address string         要解析的账户地址
  	- keyStorePath string     keyStore文件地址
  	- password string        keyStore加密密码

  Returns:
  	- string privateKey      私钥
	- error

  Call permissions: Anyone
*/
func GetPrivateKeyFromFile(address string, keyStorePath, password string) (string, error) {
	keyJson, err := ioutil.ReadFile(keyStorePath)
	if err != nil {
		return "", err
	}
	var key *keystore.Key

	addr := utils.StringToAddress(address)
	key, _, err = keystore.DecryptKey(keyJson, password, addr.CryptoType())

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
  	- keyStorePath string      keyStore文件地址
  	- isSM2   bool       是否采用国密生成私钥，true为是国密，false为否
  	- privateKey string  私钥
  	- password string    keyStore加密密码

  Returns:
  	- bool     true为转换成功,false为转换失败
	- error

  Call permissions: Anyone
*/
func PrivateKeyToKeyStoreFile(keyStorePath string, isSM2 bool, privateKey string, password, chainCode string) (bool, error) {
	if utils.Has0xPrefix(privateKey) {
		privateKey = privateKey[2:]
	}

	if !utils.IsHex(privateKey) || len(privateKey) != 64 {
		return false, errors.New("privateKey is not hex string or not 32 Bytes")
	}

	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return false, err
	}

	scryptN, scryptP, err := preCheck(keyStorePath, password, false)
	if err != nil {
		return false, err
	}

	key := keystore.NewKeyStore(keyStorePath, scryptN, scryptP)
	_, err = key.ImportECDSA(privateKeyN, password, chainCode)
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
func GetAddressFromPrivate(privateKey string, isSM2 bool) (utils.Address, error) {
	if utils.Has0xPrefix(privateKey) {
		privateKey = privateKey[2:]
	}

	if !utils.IsHex(privateKey) || len(privateKey) != 64 {
		return utils.Address{}, errors.New("privateKey is not hex string or not 32 Bytes")
	}

	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return utils.Address{}, err
	}
	// 转换成地址
	return crypto.PubkeyToAddress(*privateKeyN.Public().(*ecdsa.PublicKey)), nil
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
	if utils.Has0xPrefix(privateKey) {
		privateKey = privateKey[2:]
	}

	if !utils.IsHex(privateKey) || len(privateKey) != 64 {
		return "", errors.New("privateKey is not hex string or not 32 Bytes")
	}

	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	pubBytes := crypto.FromECDSAPub(privateKeyN.Public().(*ecdsa.PublicKey))
	if isSM2 {
		return "0x" + utils.Bytes2Hex(pubBytes) + utils.Bytes2Hex(privateKeyN.Public().(*ecdsa.PublicKey).Y.Bytes()), nil
	} else {
		return "0x" + utils.Bytes2Hex(pubBytes), nil
	}
}

/*
  GetPublicKeyFromFile:
   	EN -
	CN - 根据keyStore获取公钥
  Params:
  	- keyStorePath       string   keyStore的存储路径
  	- password           string   keyStore的加密密码
	- isSM2              bool     是否采用国密生成私钥，true为是国密，false为否

  Returns:
  	- string  公钥（65字节）
	- error

  Call permissions: Anyone
*/
func GetPublicKeyFromFile(keyStorePath, password string, isSM2 bool) (string, error) {
	keyJson, err := ioutil.ReadFile(keyStorePath)
	if err != nil {
		return "", errors.New("not find privateKeyFile")
	}

	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	key, _, err := keystore.DecryptKey(keyJson, password, cryptoType)
	if err != nil {
		return "", err
	}
	privateKey := hex.EncodeToString(key.PrivateKey.D.Bytes())
	privateKeyN, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	pubBytes := crypto.FromECDSAPub(privateKeyN.Public().(*ecdsa.PublicKey))
	if isSM2 {
		return "0x" + utils.Bytes2Hex(pubBytes) + utils.Bytes2Hex(privateKeyN.Public().(*ecdsa.PublicKey).Y.Bytes()), nil
	} else {
		return "0x" + utils.Bytes2Hex(pubBytes), nil
	}
}

/*
  CheckPublicKeyToAccount:
   	EN -
	CN - 判断给定的公钥是否与账户地址相匹配
  Params:
  	-
  	-

  Returns:
  	-
     - error

  Call permissions: Anyone
*/
func CheckPublicKeyToAccount(account, publicKey string) (bool, error) {

	addr := utils.StringToAddress(account)

	// 检测是否带有前缀
	if utils.Has0xPrefix(publicKey) {
		publicKey = publicKey[2:]
	}

	// 检测公钥合法性
	if !(utils.IsHex(publicKey) && len(publicKey) == 130) {
		return false, errors.New("publicKey is not a hexadecimal string or the length is less than 130(132 with prefix '0x'")
	}

	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, errors.New("publicKey can't be decode to bytes")
	}
	key, err := crypto.UnmarshalPubkey(pubBytes)
	if err != nil {
		return false, errors.New("publicKey is not valid")
	}
	addrParse := crypto.PubkeyToAddress(*key)
	if addr != addrParse {
		return false, errors.New("publicKey does not match account")
	}

	return true, nil
}

// 节点的私钥全部是采用SECP256K1的方式生成的，解密也是按照这个方式
func GenerateNodeUrlFromKeyStore(nodePrivateKeyPath, password, host string, port uint64) (string, error) {
	if !isLegalIP(host) {
		return "", errors.New("host is illegal")
	}
	if port > 65535 {
		return "", errors.New("port should be in range 0 to 65535")
	}

	keyJson, err := ioutil.ReadFile(nodePrivateKeyPath)
	if err != nil {
		return "", err
	}

	key, _, err := keystore.DecryptKey(keyJson, password, config.SECP256K1)

	if err != nil {
		return "", err
	}

	var peerId libp2pcorepeer.ID
	if peerId, err = libp2pcorepeer.IDFromPrivateKey((*libp2pcorecrypto.Secp256k1PrivateKey)(key.PrivateKey)); err != nil {
		return "", err
	}

	return fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", host, port, peerId.String()), nil
}

func isLegalIP(ip string) bool {
	// ip地址范围：(1~255).(0~255).(0~255).(0~255)
	// ipRegEx := "^([1-9]|([1-9][0-9])|(1[0-9][0-9])|(2[0-4][0-9])|(25[0-5]))(\\.([0-9]|([1-9][0-9])|(1[0-9][0-9])|(2[0-4][0-9])|(25[0-5]))){3}$"
	ipRegEx := "^([1-9]|([1-9]\\d)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))(\\.(\\d|([1-9]\\d)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))){3}$"
	// ipRegEx := "^(([1-9]\\d?)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))(\\.(0|([1-9]\\d?)|(1\\d{2})|(2[0-4]\\d)|(25[0-5]))){3}$"
	// Pattern
	reg, _ := regexp.Compile(ipRegEx)
	// Matcher
	return reg.MatchString(ip)
}

func GenerateNodeUrlFromPrivateKey(privateKey, host string, port uint64) (string, error) {
	if !isLegalIP(host) {
		return "", errors.New("host is illegal")
	}
	if port > 65535 {
		return "", errors.New("port should be in range 0 to 65535")
	}

	if utils.Has0xPrefix(privateKey) {
		privateKey = privateKey[2:]
	}

	if !utils.IsHex(privateKey) || len(privateKey) != 64 {
		return "", errors.New("privateKey is not hex string or not 32 Bytes")
	}

	privateKeyN, err := crypto.HexToECDSA(privateKey, config.SECP256K1)
	if err != nil {
		return "", err
	}
	var peerId libp2pcorepeer.ID
	if peerId, err = libp2pcorepeer.IDFromPrivateKey((*libp2pcorecrypto.Secp256k1PrivateKey)(privateKeyN)); err != nil {
		return "", err
	}

	return fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", host, port, peerId.String()), nil
}
