package account

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/crypto/config"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"github.com/bif/bif-sdk-go/utils/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
)

// SignatureLength indicates the byte length required to carry a signature with recovery id.
// deprecated: 注意由于已经经过国密改造，这个的意义是否有改变？？
const SignatureLength = 64 + 1 // 64 bytes ECDSA signature + 1 byte recovery id

type SignTxParams struct {
	ChainId   uint64
	Nonce     *big.Int
	GasPrice  *big.Int
	GasLimit  uint64
	Sender    *utils.Address
	Recipient *utils.Address
	Amount    *big.Int
	Payload   []byte
}

type txData struct {
	ChainId   uint64         `json:"chainId"    gencodec:"required"`
	Nonce     uint64         `json:"nonce"    gencodec:"required"`
	GasPrice  *big.Int       `json:"gasPrice" gencodec:"required"`
	GasLimit  uint64         `json:"gas"      gencodec:"required"`
	Sender    *utils.Address `json:"sender"     rlp:"nil"`      // nil means contract creation
	Recipient *utils.Address `json:"recipient"       rlp:"nil"` // nil means contract creation
	Amount    *big.Int       `json:"amount"    gencodec:"required"`
	Payload   []byte         `json:"payload"    gencodec:"required"`

	// account Signature values
	SignUser []byte `json:"signUser"    gencodec:"required"` // 33 + 1 + 32 + 32  即  公钥 + 类型 + r + s

	// This is only used when marshaling to JSON.
	Hash *utils.Hash `json:"hash" rlp:"-"`
}

type SignTransactionResult struct {
	Raw hexutil.Bytes `json:"raw"`
	Tx  *txData       `json:"tx"`
}

type BIFSigner struct{}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (bfs BIFSigner) Hash(tx *txData) utils.Hash {
	return rlpHash([]interface{}{
		tx.ChainId,
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.Sender,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
	})
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *txData, s BIFSigner, prv *ecdsa.PrivateKey, cryptoType config.CryptoType) (*txData, error) {
	h := s.Hash(tx)
	var err error
	Signature, err := crypto.GenSignature(h[:], prv, cryptoType)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	b.Write(Signature.PublicKey)
	b.Write(Signature.CryptoType)
	b.Write(Signature.Signature)
	tx.SignUser = b.Bytes()
	return tx, nil
}

/*
  SignTransaction:
   	EN -
 	CN - 使用地址私钥给指定的交易签名，返回签名结果
  Params:
  	- transaction: *SignTxParams 指定的交易信息
		- Recipient        string    （可选）交易的接收方，如果是部署合约，则为空
		- AccountNonce     uint64    （可选）整数，可以允许你覆盖你自己的相同nonce的，待pending中的交易；默认是Core.GetTransactionCount()
		- GasPrice       uint64     交易可使用的gas，未使用的gas会退回。
		- GasLimit  *big.Int  （可选）默认是自动确定，交易的gas价格，默认是 Core.GetGasPrice()
		- Amount     *big.Int  （可选）交易转移的bifer，以wei为单位
		- Payload      []byte    （可选）合约函数交互中调用的数据的ABI字节字符串或者合约创建时初始的字节码
		- ChainId   *big.Int   签署此交易时要使用的链ID，默认是Core.GetChainId
 	- privateKey: string, 私钥（transaction中的from地址对应的私钥）
 	- isSM2: bool,  私钥生成是否采用国密，是的话为True，否则为false

  Returns:
  	- *SignTransactionResult
 	- error

  Call permissions: Anyone
  TODO: 签署交易在构造交易时，*txData中的NT， NV， NR，NS暂未处理，后期需要加上(现在未加，因为加上后，使用Send_rawTransaction)报错。rlp: input list has too many elements for types.Txdata
*/
func SignTransaction(signData *SignTxParams, privateKey string, isSM2 bool) (*SignTransactionResult, error) {
	// 1 check input
	tx, err := basePreCheck(signData, privateKey, isSM2)
	if err != nil {
		return nil, err
	}

	// 2 Get signature type based on Sender type
	var cryptoType config.CryptoType
	if isSM2 {
		cryptoType = config.SM2
	} else {
		cryptoType = config.SECP256K1
	}

	privKey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		return nil, err
	}

	// 3 New Signer
	signer := &BIFSigner{}
	signed, err := SignTx(tx, *signer, privKey, cryptoType)
	if err != nil {
		return nil, err
	}

	data, err := rlp.EncodeToBytes(signed)

	if err != nil {
		return nil, err
	}
	// fmt.Printf("signed  %#v \n", signed)
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
  TODO：由于解密交易和签名交易的还没统一，后期需要修改，将rlp.DecodeBytes(rawTx, &tx)加上错误处理（即上文注释的部分）,暂时不做这个接口

*/
func RecoverTransaction(rawTxString string, isSM2 bool) (string, error) {
	if utils.Has0xPrefix(rawTxString) {
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

	// // fmt.Printf("%v \n ", tx)
	// chainId, err := account.getChainId()
	// if err != nil {
	// 	return "", err
	// }
	//
	// signer := &BIFSigner{}
	// sigHash := signer.Hash(&tx)
	//
	// r, s := tx.R.Bytes(), tx.S.Bytes()
	// // deprecated: 国密加密不知道为啥会变长
	// sig := make([]byte, SignatureLength+1)
	// copy(sig[32-len(r):32], r)
	// copy(sig[64-len(s):64], s)
	// var v byte
	// if signer.chainId.Sign() != 0 {
	// 	tx.V.Sub(tx.V, signer.chainIdMul)
	// 	v = byte(tx.V.Uint64() - 35)
	// } else {
	// 	v = byte(tx.V.Uint64() - 27)
	//
	// }
	//
	// sig[64] = v
	// if isSM2 {
	// 	sig[65] = byte(0)
	// } else {
	// 	sig[65] = byte(1)
	// }
	// // _, err = crypto.HexToECDSA(resources.AddressPriKey, crypto.SM2)
	// pubBytes, err := crypto.Ecrecover(sigHash[:], sig)
	// if err != nil {
	// 	return "", err
	// }
	// return publicKeyStrToAddress(pubBytes, isSM2)
	return "", err
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
todo:暂时注释这个接口
*/
func Recover(rawTxString string, isSM2 bool) (string, error) {
	// if utils.Has0xPrefix(rawTxString) {
	// 	rawTxString = rawTxString[2:]
	// }
	//
	// rawTx, err := hex.DecodeString(rawTxString)
	// if err != nil {
	// 	return "", err
	// }
	//
	// var tx txData
	// rawTx, err = hex.DecodeString(rawTxString)
	// if err != nil {
	// 	return "", err
	// }
	//
	// // err = rlp.DecodeBytes(rawTx, &tx)
	// // // fmt.Printf("tx is %v \n", tx)
	// // if err != nil {
	// // 	return "", err
	// // }
	// rlp.DecodeBytes(rawTx, &tx)
	//
	// // fmt.Printf("%v \n ", tx)
	// chainId, err := account.getChainId()
	// if err != nil {
	// 	return "", err
	// }
	//
	// signer := &BIFSigner{
	// 	chainId:    chainId,
	// 	chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	// }
	// sigHash := signer.Hash(&tx)
	//
	// r, s := tx.R.Bytes(), tx.S.Bytes()
	// // deprecated: 国密加密不知道为啥会变长
	// sig := make([]byte, SignatureLength+1)
	// copy(sig[32-len(r):32], r)
	// copy(sig[64-len(s):64], s)
	// var v byte
	// if signer.chainId.Sign() != 0 {
	// 	tx.V.Sub(tx.V, signer.chainIdMul)
	// 	v = byte(tx.V.Uint64() - 35)
	// } else {
	// 	v = byte(tx.V.Uint64() - 27)
	//
	// }
	//
	// sig[64] = v
	// if isSM2 {
	// 	sig[65] = byte(0)
	// } else {
	// 	sig[65] = byte(1)
	// }
	// // _, err = crypto.HexToECDSA(resources.AddressPriKey, crypto.SM2)
	// pubBytes, err := crypto.Ecrecover(sigHash[:], sig)
	// if err != nil {
	// 	return "", err
	// }
	// return publicKeyStrToAddress(pubBytes, isSM2)
	return "", errors.New("接口待修改")
}

func rlpHash(x interface{}) (h utils.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

//  只会校验非链上的信息
func basePreCheck(signData *SignTxParams, privateKey string, isSM2 bool) (*txData, error) {

	_, err := PriKeyToAccount(privateKey, isSM2, "")
	if err != nil {
		return nil, errors.New("not invalid privateKey")
	}

	// 校验Nonce是否为0
	if signData.Nonce == nil{
		return nil, errors.New("nonce is nil")
	}

	// 校验gas
	if signData.GasLimit < 21000 {
		return nil, errors.New("gas should be at least 21000")
	}

	// 校验gasPrice
	if signData.GasPrice != nil && signData.GasPrice.Cmp(big.NewInt(0)) == -1 {
		return nil, errors.New("gasPrice can't be negative")
	}

	// 校验Value
	if signData.Amount != nil && signData.Amount.Cmp(big.NewInt(0)) == -1 {
		return nil, errors.New("value can't be negative")
	}

	// 校验ChainId
	if signData.ChainId == 0 {
		return nil, errors.New("chainId can't be 0")
	}

	tx := &txData{
		ChainId:   signData.ChainId,
		Nonce:     signData.Nonce.Uint64(),
		GasPrice:  signData.GasPrice,
		GasLimit:  signData.GasLimit,
		Sender:    signData.Sender,
		Recipient: signData.Recipient,
		Amount:    signData.Amount,
		Payload:   signData.Payload,
		SignUser:  nil,
	}
	return tx, nil
}
