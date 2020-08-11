package account

import (
	"crypto/ecdsa"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
)

type SignData  struct {
	Message     string `json:"message"    gencodec:"required"`
	MessageHash string `json:"messageHash" gencodec:"required"`
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

// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (sign *SignData) WithSignature(signer BIFSigner, sig []byte) (*SignData, error) {
	r, s, v, err := signer.SignatureValues(sig)
	if err != nil {
		return nil, err
	}
	cpy := sign
	cpy.R, cpy.S, cpy.V = r, s, v
	return cpy, nil
}

// SignDt signs the data using the given signer and private key
func SignDt(signData *SignData, s BIFSigner, prv *ecdsa.PrivateKey) (*SignData, error) {
	h := utils.Hex2Bytes(signData.MessageHash[2:])
	var sig []byte
	var err error
	t := signData.T.Uint64()
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
	return signData.WithSignature(s, sig)
}