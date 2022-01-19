package account

import (
	"bytes"
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
	"strconv"
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

/*
  Create:
   	EN - Generate public and private key pair
 	CN - 生成公私钥对
  Params:
  	- isSM2,bool, 私钥生成方式，采用国密还是非国密

  Returns:
  	- accountAddress:string, privateKey:string  返回账户地址和私钥，类型为hex String
  	- error

  Call permissions: Anyone
*/
func Create(isSM2 bool, chainCode string) (string, string, error) {
	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}
	privateKeyECDSA, err := crypto.GenerateKey(cryptoType)
	if err != nil {
		return "", "", err
	}
	account := crypto.PubkeyToAddress(*privateKeyECDSA.Public().(*ecdsa.PublicKey)).String(chainCode)
	privateKey := hex.EncodeToString(privateKeyECDSA.D.Bytes())
	return account, privateKey, nil
}

/*
  GenKeyStore:
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
func GenKeyStore(keyStorePath string, isSM2 bool, password, chainCode string, UseLightweightKDF bool) (string, error) {

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
func PriKeyFromKeyStore(address string, keyStorePath, password string) (string, error) {
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
  PriKeyToKeyStore(默认私钥文件生成采用UseLightweightKDF为false即安全模式):
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
func PriKeyToKeyStore(keyStorePath string, isSM2 bool, privateKey string, password, chainCode string) (bool, error) {
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
  PriKeyToAccount:
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
func PriKeyToAccount(privateKey string, isSM2 bool, chainCode string) (string, error) {
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
	// 转换成地址
	return crypto.PubkeyToAddress(*privateKeyN.Public().(*ecdsa.PublicKey)).String(chainCode), nil
}

/*
  PriKeyToPublicKey:
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
func PriKeyToPublicKey(privateKey string, isSM2 bool) (string, error) {
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
  PublicKeyFromFile:
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
func PublicKeyFromFile(keyStorePath, password string, isSM2 bool) (string, error) {
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

func PublicKeyToAccount(pubBytes []byte, isSM2 bool, chainCode string) (string, error) {
	var addr []byte
	if isSM2 {
		addr = crypto.Keccak256(config.SECP256K1, pubBytes[1:])[12:]
		addr[8] = 115
	} else {
		addr = crypto.Keccak256(config.SECP256K1, pubBytes[1:])[12:]
		if addr[8] == 115 {
			addr[8] = 103
		}
	}
	return utils.BytesToAddress(addr).String(chainCode), nil
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

/*
  Encrypt:
   	EN - Encrypts a private key to the keystore v3 standard.
 	CN - 将私钥加密为keystore v3标准
  Params:
   	- privateKey, string, 私钥的hex string
   	- cryType, uint， 加密的类型，0为国密加密；1为非国密；输入其他数值则默认为国密
 	- password, string， 用于加密的密码
 	- UseLightweightKDF, bool，密钥生成时，是否采用低标准的密钥库的内存和CPU，一般为false

  Returns:
  	- []byte,加密的密钥库v3
 	- error

  Call permissions: Anyone
*/
func Encrypt(privateKey string, isSM2 bool, password, chainCode string, UseLightweightKDF bool) ([]byte, error) {
	if password == "" {
		return nil, errors.New("empty password, please check")
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}

	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	var privkey *ecdsa.PrivateKey
	privkey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return nil, err
	}

	return keystore.EncryptKey(keystore.NewKeyFromECDSA(privkey), password, chainCode, scryptN, scryptP)
}

/*
  Decrypt:
   	EN - Decrypts a keystore v3
 	CN - 解密密钥库v3
  Params:
   	- keystoreJson, []byte, 私钥的hex string
   	- cryType, uint， 加密的类型，0为国密加密；1为非国密；输入其他数值则默认为国密
 	- password, string， 用于加密的密码

  Returns:
  	- string, 账户地址(hexString)
  	- string, 私钥(hexString)
 	- error

  Call permissions: Anyone
*/
func Decrypt(keystoreJson []byte, isSM2 bool, password string) (string, string, error) {
	if len(keystoreJson) == 0 {
		return "", "", errors.New("keystoreJson is empty")
	}

	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	var addressString string
	key, addressString, err := keystore.DecryptKey(keystoreJson, password, cryptoType)
	if err != nil {
		return "", "", err
	}

	return addressString, hex.EncodeToString(key.PrivateKey.D.Bytes()), nil
}

func MessageSignatureBtc(message, password string, keyFileData []byte) (string, string, error) {
	messageSha3 := utils.Sha3Raw(message)

	_, privateKey, err := Decrypt(keyFileData, false, password)
	privKey, err := crypto.HexToECDSA(privateKey, config.SECP256K1)
	if err != nil {
		return "", "", err
	}

	var cryptoType config.CryptoType
	cryptoType = config.SECP256K1

	messageSha3Bytes := utils.Hex2Bytes(messageSha3[2:])
	sig, err := crypto.NewSignature(messageSha3Bytes, privKey, cryptoType)
	if err != nil {
		return "", "", err
	}
	// r := new(big.Int).SetBytes(sig[:32])
	// s := new(big.Int).SetBytes(sig[32:64])
	// v := new(big.Int).SetBytes([]byte{sig[64] + 27})
	// fmt.Printf("r %x \n", r)
	// fmt.Printf("s %x \n", s)
	// fmt.Printf("v %x \n", v)
	// fmt.Printf("sig len is  %x \n", len(sig))
	var buf bytes.Buffer
	// 链的代码已经撤销这个
	// buf.Write([]byte{sig[64] + 27})
	buf.Write(sig.Signature[:64])
	// fmt.Printf("sig is  %s \n", t+utils.Bytes2Hex(buf.Bytes()))
	return messageSha3, utils.Bytes2Hex(buf.Bytes()), err
}

/*
  HashMessage:
   	EN - Hashes the given message，The data will be UTF-8 HEX decoded and enveloped
 	CN - 对给定消息进行哈希处理，数据将按以UTF-8 HEX解码和封装
  Params:
  	- message, string, 散列消息，如果为十六进制，则先对其进行UTF8解码

  Returns:
  	- string: 散列消息
 	- error

  Call permissions: Anyone
*/
func HashMessage(message string, isSm2 bool) string {
	var messageHex string
	if !utils.IsHexStrict(message) {
		messageHex = utils.Utf8ToHex(message)
	} else {
		messageHex = message
	}
	messageBytes := utils.Hex2Bytes(messageHex[2:])
	preamble := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(messageBytes))
	var buffer bytes.Buffer
	buffer.Write([]byte(preamble))
	buffer.Write(messageBytes)

	var cryptoType config.CryptoType
	if isSm2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}
	hashBytes := crypto.Keccak256(cryptoType, buffer.Bytes())
	return "0x" + utils.Bytes2Hex(hashBytes)
}

// 节点的私钥全部是采用Secp256k1的方式生成的，解密也是按照这个方式
func GenNodeUrlFromKeyStore(nodePrivateKeyPath, password, host string, port uint64) (string, error) {
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

func GenNodeUrlFromPriKey(privateKey, host string, port uint64) (string, error) {
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