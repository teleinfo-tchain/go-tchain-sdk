package account

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"github.com/bif/bif-sdk-go/utils/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type SignTxParams struct {
	To       string
	Nonce    uint64
	Gas      uint64
	GasPrice *big.Int
	Value    *big.Int
	Data     []byte
	ChainId  *big.Int
}

type txData struct {
	AccountNonce uint64         `json:"nonce"    gencodec:"required"`
	Price        *big.Int       `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64         `json:"gas"      gencodec:"required"`
	Sender       *utils.Address `json:"from"     rlp:"nil"`
	Recipient    *utils.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int       `json:"value"    gencodec:"required"`
	Payload      []byte         `json:"input"    gencodec:"required"`

	// account Signature values
	T *big.Int `json:"t" gencodec:"required"`
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// // node Signature values
	// NT *big.Int `json:"nt" gencodec:"required"`
	// NV *big.Int `json:"nv" gencodec:"required"`
	// NR *big.Int `json:"nr" gencodec:"required"`
	// NS *big.Int `json:"ns" gencodec:"required"`
	// This is only used when marshaling to JSON.
	Hash *utils.Hash `json:"hash" rlp:"-"`
}

type SignTransactionResult struct {
	Raw hexutil.Bytes `json:"raw"`
	Tx  *txData       `json:"tx"`
}

func rlpHash(x interface{}) (h utils.Hash) {
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
func (bfs BIFSigner) Hash(tx *txData) utils.Hash {
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
func (bfs BIFSigner) SignatureValue(sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) < SignatureLength {
		return nil, nil, nil, fmt.Errorf("wrong size for signature: got %d, want %d ", len(sig), 65)
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v, nil
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (bfs BIFSigner) SignatureValues(sig []byte) (R, S, V *big.Int, err error) {
	R, S, V, err = bfs.SignatureValue(sig)
	// fmt.Println("v is ", V)
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
func (tx *txData) WithSignature(signer BIFSigner, sig []byte) (*txData, error) {
	r, s, v, err := signer.SignatureValues(sig)
	if err != nil {
		return nil, err
	}
	cpy := tx
	cpy.R, cpy.S, cpy.V = r, s, v
	return cpy, nil
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *txData, s BIFSigner, prv *ecdsa.PrivateKey) (*txData, error) {
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
	// fmt.Println("Hash is ", h[:])
	// fmt.Println("sig is ", sig)
	// fmt.Println("sig len is ", len(sig))
	return tx.WithSignature(s, sig)
}