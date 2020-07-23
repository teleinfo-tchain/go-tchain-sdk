package core

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/common/rlp"
	"github.com/bif/bif-sdk-go/crypto"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type Txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Sender       *common.Address `json:"from"     gencodec:"required"`
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

type SignTransactionResult struct {
	Raw hexutil.Bytes `json:"raw"`
	Tx  *Txdata       `json:"tx"`
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type BIFSigner struct {
	chainId, chainIdMul *big.Int
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (bfs BIFSigner) Hash(tx *Txdata) common.Hash {
	return rlpHash([]interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Sender,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
		bfs.chainId, uint(0), uint(0),
	})
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (bfs BIFSigner) SignatureValue(tx *Txdata, sig []byte) (r, s, v *big.Int, err error) {
	//if len(sig) != SignatureLength {
	//	panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), SignatureLength))
	//}
	if len(sig) < 65 {
		fmt.Println("len sig is ", len(sig))
		return nil, nil, nil, fmt.Errorf("wrong size for signature: got %d, want %d ", len(sig), 65)
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v, nil
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (bfs BIFSigner) SignatureValues(tx *Txdata, sig []byte) (R, S, V *big.Int, err error) {
	R, S, V, err = bfs.SignatureValue(tx, sig)
	if err != nil {
		return nil, nil, nil, err
	}
	if bfs.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, bfs.chainIdMul)
	}
	return R, S, V, nil
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (tx *Txdata) WithSignature(signer BIFSigner, sig []byte) (*Txdata, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := tx
	cpy.R, cpy.S, cpy.V = r, s, v
	return cpy, nil
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *Txdata, s BIFSigner, prv *ecdsa.PrivateKey) (*Txdata, error) {
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

func (tx *Txdata) PreCheck() (bool, error) {
	if tx.Sender == nil {
		return false, fmt.Errorf("sender not specified")
	}
	// 判断地址格式的有效性
	if !(bytes.HasPrefix(tx.Sender.Bytes(), []byte("did:bid:"))) {
		return false, fmt.Errorf("not invalid sender address")
	}

	if tx.Price == nil {
		return false, fmt.Errorf("gasPrice not specified")
	}

	if tx.GasLimit == 0 {
		return false, fmt.Errorf("gasLimit not specified")
	}

	return true, nil
}

/*
SignTransaction: 使用地址私钥给指定的交易签名，返回签名结果

Params:
	- transaction: *Txdata 指定的交易信息
	- privKey: string, 私钥（transaction中的from地址对应的私钥）
	- chainId: int64, 链的ChainId

Returns:
	- *SignTransactionResult
	- error
*/
func SignTransaction(transaction *Txdata, privKey string, chainId int64) (*SignTransactionResult, error) {
	// 1 check input
	ret, err := transaction.PreCheck()
	fmt.Println("t is ", transaction.T)
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

	privateKey, err := crypto.HexToECDSA(privKey, cryptoType)
	if err != nil {
		return nil, err
	}

	// 3 New Signer
	id := big.NewInt(chainId)
	signer := &BIFSigner{
		chainId:    id,
		chainIdMul: new(big.Int).Mul(id, big.NewInt(2)),
	}

	signed, err := SignTx(transaction, *signer, privateKey)

	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(signed)

	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signed}, nil
}
