package account

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/account/keystore"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/rlp"
	"github.com/bif/bif-sdk-go/crypto"
	"math/big"
)

// SignatureLength indicates the byte length required to carry a signature with recovery id.
// deprecated: 注意由于已经经过国密改造，这个的意义是否有改变？？
const SignatureLength = 64 + 1 // 64 bytes ECDSA signature + 1 byte recovery id

// Account - The Account Module
type Account struct{}

//  NewAccount - 初始化Account
func NewAccount() *Account {
	account := new(Account)
	return account
}

func cryptoType(isSM2 bool) crypto.CryptoType {
	if isSM2 {
		return crypto.SM2
	}
	return crypto.SECP256K1
}

func publicKeyStrToAddress(pubBytes []byte, isSM2 bool) (string, error) {
	var addr []byte
	if isSM2 {
		addr = crypto.Keccak256(crypto.SECP256K1, pubBytes[1:])[12:]
		addr[8] = 115
	} else {
		addr = crypto.Keccak256(crypto.SECP256K1, pubBytes[1:])[12:]
		if addr[8] == 115 {
			addr[8] = 103
		}
	}
	return common.BytesToAddress(addr).String(), nil
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
func (account *Account) Create(isSM2 bool) (string, string, error) {
	cryptoType := cryptoType(isSM2)
	privateKeyECDSA, err := crypto.GenerateKey(cryptoType)
	if err != nil {
		return "", "", err
	}
	accountAddress := crypto.PubkeyToAddress(*privateKeyECDSA.Public().(*ecdsa.PublicKey)).Hex()
	privateKey := hex.EncodeToString(privateKeyECDSA.D.Bytes())
	return accountAddress, privateKey, nil
}

/*
  PrivateKeyToAccount:
   	EN - Generate account address based on private key
 	CN - 根据私钥生成账户地址
  Params:
	- privateKey, string, 私钥的hex string
	- cryType, uint， 加密的类型，0为国密加密；1为非国密；输入其他数值则默认为国密

  Returns:
  	- string, 账户地址(hexString)
 	- error

  Call permissions: Anyone
*/
func (account *Account) PrivateKeyToAccount(privateKey string, isSM2 bool) (string, error) {
	cryptoType := cryptoType(isSM2)
	privKey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return "", err
	}
	// 转换成地址
	return crypto.PubkeyToAddress(*privKey.Public().(*ecdsa.PublicKey)).Hex(), nil
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
func (account *Account) Encrypt(privateKey string, isSM2 bool, password string, UseLightweightKDF bool) ([]byte, error) {
	if password == "" {
		return nil, errors.New("empty password, please check")
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}

	cryptoType := cryptoType(isSM2)

	var privkey *ecdsa.PrivateKey
	privkey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return nil, err
	}

	return keystore.EncryptKey(keystore.NewKeyFromECDSA(privkey), password, scryptN, scryptP)
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
func (account *Account) Decrypt(keystoreJson []byte, isSM2 bool, password string) (string, string, error) {
	cryptoType := cryptoType(isSM2)

	key, err := keystore.DecryptKey(keystoreJson, password, cryptoType)
	if err != nil {
		return "", "", err
	}

	return key.Address.String(), hex.EncodeToString(key.PrivateKey.D.Bytes()), nil
}

/*
  SignTransaction:
   	EN -
 	CN - 使用地址私钥给指定的交易签名，返回签名结果
  Params:
  	- transaction: *TxData 指定的交易信息
 	- privateKey: string, 私钥（transaction中的from地址对应的私钥）
 	- chainId: int64, 链的ChainId

  Returns:
  	- *SignTransactionResult
 	- error

  Call permissions: Anyone
*/
func (account *Account) SignTransaction(transaction *TxData, privateKey string, chainId *big.Int) (*SignTransactionResult, error) {
	// 1 check input
	ret, err := transaction.PreCheck(privateKey)
	if !ret {
		return nil, err
	}

	// 2 Get signature type based on Sender type
	var cryptoType crypto.CryptoType
	if transaction.Sender[8] == 115 {
		cryptoType = crypto.SM2
	} else {
		transaction.T = common.Big1
		cryptoType = crypto.SECP256K1
	}

	privKey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return nil, err
	}

	// 3 New Signer
	signer := &BIFSigner{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}

	signed, err := SignTx(transaction, *signer, privKey)

	if err != nil {
		return nil, err
	}

	data, err := rlp.EncodeToBytes(signed)

	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
}

/*
  RecoverTransaction:
   	EN - Recovers the Bif address which was used to sign the given RLP encoded transaction.
 	CN - 恢复用于签名给定RLP编码交易的Bif地址。
  Params:
  	- rlpTransaction, string, RLP编码的交易

  Returns:
  	- string, 签名该交易的地址，为hexString
 	- error

  Call permissions: Anyone
  Bug:国密解密有问题
*/
func (account *Account) RecoverTransaction(rawTxString string, isSM2 bool) (string, error) {
	if rawTxString[:2] == "0x" || rawTxString[:2] == "0X" {
		rawTxString = rawTxString[2:]
	}

	rawTx, err := hex.DecodeString(rawTxString)
	if err != nil {
		return "", err
	}

	var tx TxData
	rawTx, err = hex.DecodeString(rawTxString)
	if err != nil {
		return "", err
	}

	err = rlp.DecodeBytes(rawTx, &tx)
	fmt.Printf("tx is %v \n", tx)
	if err != nil {
		return "", err
	}

	// fmt.Printf("%v \n ", tx)
	// deprecated: 应该可以从txSign中拿到或者其他的方法
	signer := &BIFSigner{
		chainId:    big.NewInt(333),
		chainIdMul: new(big.Int).Mul(big.NewInt(333), big.NewInt(2)),
	}
	sigHash := signer.Hash(&tx)

	r, s := tx.R.Bytes(), tx.S.Bytes()
	// deprecated: 国密加密不知道为啥会变长
	sig := make([]byte, SignatureLength+1)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	var v byte
	if signer.chainId.Sign() != 0 {
		tx.V.Sub(tx.V, signer.chainIdMul)
		v = byte(tx.V.Uint64() - 35)
	} else {
		v = byte(tx.V.Uint64() - 27)

	}

	sig[64] = v
	if isSM2 {
		sig[65] = byte(0)
	} else {
		sig[65] = byte(1)
	}
	// _, err = crypto.HexToECDSA(resources.AddressPriKey, crypto.SM2)
	pubBytes, err := crypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		return "", err
	}
	return publicKeyStrToAddress(pubBytes, isSM2)
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
func (account *Account) HashMessage(message string) {

}

func (account *Account) Sign() {

}

func (account *Account) Recover() {

}
