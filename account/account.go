package account

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/bif/bif-sdk-go/account/keystore"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/rlp"
	"math/big"
	"strconv"
)

// SignatureLength indicates the byte length required to carry a signature with recovery id.
// deprecated: 注意由于已经经过国密改造，这个的意义是否有改变？？
const SignatureLength = 64 + 1 // 64 bytes ECDSA signature + 1 byte recovery id

// Account - The Account Module
type Account struct {
	provider providers.ProviderInterface
}

//  NewAccount - 初始化Account
func NewAccount(provider providers.ProviderInterface) *Account {
	account := new(Account)
	account.provider = provider
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
	return utils.BytesToAddress(addr).String(), nil
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

// - To        string    （可选）交易的接收方，如果是部署合约，则为空
// - Nonce     uint64    （可选）整数，可以允许你覆盖你自己的相同nonce的，待pending中的交易；默认是Core.GetTransactionCount()
// - Gas       uint64     交易可使用的gas，未使用的gas会退回。
// - GasPrice  *big.Int  （可选）默认是自动确定，交易的gas价格，默认是 Core.GetGasPrice()
// - Value     *big.Int  （可选）交易转移的bifer，以wei为单位
// - Data      []byte    （可选）合约函数交互中调用的数据的ABI字节字符串或者合约创建时初始的字节码
// - ChainId   *big.Int   签署此交易时要使用的链ID，默认是Core.GetChainId
func (account *Account) preCheckTx(signData *SignTxParams, privateKey string, isSM2 bool) (*txData, error) {
	// 暂缓判断
	if signData.Gas < 0 {
		return nil, errors.New("params(GasPrice|Value|ChainId) should not be less than 0")
	}

	if signData.Gas == 0 {
		return nil, errors.New("gas should be greater than 0")
	}

	var cryptoType uint

	if isSM2 {
		cryptoType = 0
	} else {
		cryptoType = 1
	}

	publicAddr, err := GetAddressFromPrivate(privateKey, cryptoType)
	if err != nil {
		return nil, errors.New("not invalid privateKey")
	}

	var recipient utils.Address

	if signData.To != "" {
		recipient = utils.StringToAddress(signData.To)
	}

	if signData.Nonce == 0 {
		pointer := &dto.CoreRequestResult{}
		err := account.provider.SendRequest(pointer, "core_getTransactionCount", []string{publicAddr, "latest"})

		if err != nil {
			return nil, err
		}

		signData.Nonce, err = pointer.ToUint64()
		if err != nil {
			return nil, err
		}
	}

	if signData.GasPrice == nil {
		pointer := &dto.CoreRequestResult{}

		err := account.provider.SendRequest(pointer, "core_gasPrice", nil)

		if err != nil {
			return nil, err
		}

		signData.GasPrice, err = pointer.ToBigInt()
		if err != nil {
			return nil, err
		}
	}

	if signData.ChainId == nil {
		pointer := &dto.CoreRequestResult{}

		err := account.provider.SendRequest(pointer, "core_chainId", nil)

		if err != nil {
			return nil, err
		}

		signData.ChainId, err = pointer.ToBigInt()
		if err != nil {
			return nil, err
		}
	}

	sender := utils.StringToAddress(publicAddr)
	tx := &txData{
		AccountNonce: signData.Nonce,
		Price:        signData.GasPrice,
		GasLimit:     signData.Gas,
		Sender:       &sender,
		Recipient:    &recipient,
		Amount:       signData.Value,
		Payload:      signData.Data,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
		T:            new(big.Int),
		// NT:           new(big.Int),
		// NV:           new(big.Int),
		// NR:           new(big.Int),
		// NS:           new(big.Int),
	}
	return tx, nil
}

/*
  SignTransaction:
   	EN -
 	CN - 使用地址私钥给指定的交易签名，返回签名结果
  Params:
  	- transaction: *SignTxParams 指定的交易信息
		- To        string    （可选）交易的接收方，如果是部署合约，则为空
		- Nonce     uint64    （可选）整数，可以允许你覆盖你自己的相同nonce的，待pending中的交易；默认是Core.GetTransactionCount()
		- Gas       uint64     交易可使用的gas，未使用的gas会退回。
		- GasPrice  *big.Int  （可选）默认是自动确定，交易的gas价格，默认是 Core.GetGasPrice()
		- Value     *big.Int  （可选）交易转移的bifer，以wei为单位
		- Data      []byte    （可选）合约函数交互中调用的数据的ABI字节字符串或者合约创建时初始的字节码
		- ChainId   *big.Int   签署此交易时要使用的链ID，默认是Core.GetChainId
 	- privateKey: string, 私钥（transaction中的from地址对应的私钥）
 	- isSM2: bool,  私钥生成是否采用国密，是的话为True，否则为false

  Returns:
  	- *SignTransactionResult
 	- error

  Call permissions: Anyone
  TODO: 签署交易在构造交易时，*txData中的NT， NV， NR，NS暂未处理，后期需要加上(现在未加，因为加上后，使用Send_rawTransaction)报错。rlp: input list has too many elements for types.Txdata
*/
func (account *Account) SignTransaction(signData *SignTxParams, privateKey string, isSM2 bool) (*SignTransactionResult, error) {
	// 1 check input
	tx, err := account.preCheckTx(signData, privateKey, isSM2)
	if err!=nil {
		return nil, err
	}

	// 2 Get signature type based on Sender type
	var cryptoType crypto.CryptoType
	if isSM2 {
		cryptoType = crypto.SM2
	} else {
		tx.T = utils.Big1
		cryptoType = crypto.SECP256K1
	}

	privKey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return nil, err
	}

	// 3 New Signer
	signer := &BIFSigner{
		chainId:    signData.ChainId,
		chainIdMul: new(big.Int).Mul(signData.ChainId, big.NewInt(2)),
	}

	signed, err := SignTx(tx, *signer, privKey)

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
  BUG:国密解密有问题(SM2初始化的问题)
  TODO：由于解密交易和签名交易的还没统一，后期需要修改，将rlp.DecodeBytes(rawTx, &tx)加上错误处理（即上文注释的部分）
*/
func (account *Account) RecoverTransaction(rawTxString string, isSM2 bool) (string, error) {
	if rawTxString[:2] == "0x" || rawTxString[:2] == "0X" {
		rawTxString = rawTxString[2:]
	}

	rawTx, err := hex.DecodeString(rawTxString)
	if err != nil {
		return "", err
	}

	var tx txData
	rawTx, err = hex.DecodeString(rawTxString)
	if err != nil {
		return "", err
	}

	// err = rlp.DecodeBytes(rawTx, &tx)
	// // fmt.Printf("tx is %v \n", tx)
	// if err != nil {
	// 	return "", err
	// }
	rlp.DecodeBytes(rawTx, &tx)

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
func (account *Account) HashMessage(message string, isSm2 bool) string {
	util := utils.NewUtils()
	var messageHex string
	if !util.IsHexStrict(message) {
		messageHex = util.Utf8ToHex(message)
	} else {
		messageHex = message
	}
	messageBytes := utils.Hex2Bytes(messageHex[2:])
	preamble := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(messageBytes))
	var buffer bytes.Buffer
	buffer.Write([]byte(preamble))
	buffer.Write(messageBytes)

	var cryptoType crypto.CryptoType
	if isSm2 {
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
	}
	hashBytes := crypto.Keccak256(cryptoType, buffer.Bytes())
	return "0x" + utils.Bytes2Hex(hashBytes)
}

/*
  Sign:
   	EN - Signs arbitrary data
 	CN - 根据给定数据进行签名
  Params:
  	- signData: *SignData 指定的签名数据
 	- privateKey: string, 私钥（transaction中的from地址对应的私钥）
 	- isSm2: bool, 私钥生成的类型是否采用国密，如果是则为true，否则为false
 	- chainId: int64, 链的ChainId

  Returns:
  	- *Signature
 	- error

  Call permissions: Anyone
  Deprecated: 这个接口暂时没有用
*/
func (account *Account) Sign(signData *SignData, privateKey string, isSm2 bool, chainId *big.Int) (*SignData, error) {
	// 1 Get signature type based on Sender type
	var cryptoType crypto.CryptoType
	if isSm2 {
		signData.T = utils.Big0
		cryptoType = crypto.SM2
	} else {
		signData.T = utils.Big1
		cryptoType = crypto.SECP256K1
	}

	privKey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return nil, err
	}

	// 2 New Signer
	signer := &BIFSigner{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}

	signed, err := SignDt(signData, *signer, privKey)

	if err != nil {
		return nil, err
	}

	return signed, nil
}

/*
  Recover:
   	EN - Recovers the Bif address which was used to sign the given data
 	CN - 恢复用于签名给定数据的Bif地址
  Params:
  	- rlpTransaction, string, RLP编码的交易

  Returns:
  	- string, 签名该交易的地址，为hexString
 	- error

  Call permissions: Anyone
  Bug:国密解密有问题
*/
func (account *Account) Recover(rawTxString string, isSM2 bool) (string, error) {
	if rawTxString[:2] == "0x" || rawTxString[:2] == "0X" {
		rawTxString = rawTxString[2:]
	}

	rawTx, err := hex.DecodeString(rawTxString)
	if err != nil {
		return "", err
	}

	var tx txData
	rawTx, err = hex.DecodeString(rawTxString)
	if err != nil {
		return "", err
	}

	// err = rlp.DecodeBytes(rawTx, &tx)
	// // fmt.Printf("tx is %v \n", tx)
	// if err != nil {
	// 	return "", err
	// }
	rlp.DecodeBytes(rawTx, &tx)

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
