package account

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/crypto/config"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
)

type SignData struct {
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

// SignDt signs the data using the given signer and private key
func SignDt(signData *SignData, s BIFSigner, prv *ecdsa.PrivateKey) (*SignData, error) {
	h := utils.Hex2Bytes(signData.MessageHash[2:])
	var sig []byte
	var err error
	t := signData.T.Uint64()
	switch t {
	case 0:
		sig, err = crypto.Sign(h[:], prv, config.SM2)
	case 1:
		sig, err = crypto.Sign(h[:], prv, config.SECP256K1)
	default:
		sig, err = crypto.Sign(h[:], prv, config.SM2)
	}
	if err != nil {
		return nil, err
	}
	fmt.Println("sig ", sig)
	return nil, errors.New("接口待修改")
}
