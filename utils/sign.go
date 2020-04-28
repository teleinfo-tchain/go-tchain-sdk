package utils

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/rlp"
	"github.com/bif/bif-sdk-go/crypto"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type Txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Sender       *common.Address `json:"from"     rlp:"nil"` // nil means contract creation
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	T *big.Int `json:"t" gencodec:"required"`
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

func (tx *Txdata) PreCheck() bool{
	if tx.Sender == nil{
		return false
	}
	if tx.Recipient == nil{
		return false
	}
	if tx.Price == nil{
		return false
	}
	if tx.GasLimit == 0{
		return false
	}
	return true
}

func (tx *Txdata) WithSignature(signer BIFSigner, sig []byte) (*Txdata, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := tx
	cpy.R, cpy.S, cpy.V = r, s, v
	if bytes.HasPrefix(tx.Sender.Bytes(), []byte("did:bid:")) && tx.Sender[8] == 115 {
		cpy.T = common.Big0
	} else {
		cpy.T = common.Big1
	}
	fmt.Println("r:", cpy.R)
	fmt.Println("s:", cpy.S)
	fmt.Println("v:", cpy.V)
	return cpy, nil
}

type BIFSigner struct {
	chainId, chainIdMul *big.Int
}

func (s BIFSigner) Hash(tx *Txdata) common.Hash {
	h := rlpHash([]interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Sender,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
		s.chainId, uint(0), uint(0),
	})
	return h
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (es BIFSigner) SignatureValue(tx *Txdata, sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) < 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v, nil
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s BIFSigner) SignatureValues(tx *Txdata, sig []byte) (R, S, V *big.Int, err error) {
	R, S, V, err = s.SignatureValue(tx, sig)
	if err != nil {
		return nil, nil, nil, err
	}
	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func signTx(tx *Txdata, s BIFSigner, prv *ecdsa.PrivateKey) (*Txdata, error) {
	h := s.Hash(tx)

	var sig []byte
	var err error
	t := tx.T.Uint64()
	switch t {
	case 0:
		sig, err = crypto.Sign(h[:], prv, crypto.SM2)
	case 1:
		sig, err = crypto.Sign(h[:], prv, crypto.SECP256K1)
	default:
		sig, err = crypto.Sign(h[:], prv, crypto.SM2)
	}

	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

//使用自定义的密钥对交易进行签名
func SignTransactionWithKey(transaction *Txdata, privKey string, cryType uint, chainId int64)(string, error){
	//1. 检查交易数据
	ret := transaction.PreCheck()
	if !ret{
		return "", errors.New("parameters error")
	}
	//2. 转换私钥
	var cryptoType crypto.CryptoType
	switch cryType {
	case 0:
		cryptoType = crypto.SM2
	case 1:
		cryptoType = crypto.SECP256K1
	default:
		cryptoType = crypto.SM2
	}
	privateKey,err := crypto.HexToECDSA(privKey, cryptoType)
	if err != nil{
		return "", err
	}
	//3. 新建Signer
	id := big.NewInt(chainId)
	signer := &BIFSigner{
		chainId:    id,
		chainIdMul: new(big.Int).Mul(id, big.NewInt(2)),
	}

	tx, err := signTx(transaction, *signer, privateKey)
	if err != nil{
		return "", err
	}
	signTx, err := rlp.EncodeToBytes(tx)
	if err != nil{
		return "", err
	}
	return common.Bytes2Hex(signTx), nil
}