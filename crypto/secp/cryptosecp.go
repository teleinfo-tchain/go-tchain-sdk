package secp

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/tchain/go-tchain-sdk/utils"
	"github.com/tchain/go-tchain-sdk/utils/math"
	"github.com/tchain/go-tchain-sdk/utils/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

var errInvalidPubkey = errors.New("invalid secp256k1 public key")

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256Btc(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256HashBtc(data ...[]byte) (h utils.Hash) {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// Keccak512 calculates and returns the Keccak512 hash of the input data.
func Keccak512(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak512()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b utils.Address, nonce uint64) utils.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return utils.HashToAddress(Keccak256Btc(data))
}

// CreateAddress2 creates an ethereum address given the address bytes, initial
// contract code hash and a salt.
func CreateAddress2(b utils.Address, salt [32]byte, inithash []byte) utils.Address {
	return utils.HashToAddress(Keccak256Btc([]byte{0xff}, b.Bytes(), salt[:], inithash))
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func ToECDSABtc(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256Btc()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

// FromECDSA exports a private key into a binary dump.
func FromECDSABtc(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

// UnmarshalPubkey converts bytes to a secp256k1 public key.
func UnmarshalPubkeyBtc(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(S256Btc(), pub)
	if x == nil {
		return nil, errInvalidPubkey
	}
	return &ecdsa.PublicKey{Curve: S256Btc(), X: x, Y: y}, nil
}

func FromECDSAPubBtc(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}

	return elliptic.Marshal(S256Btc(), pub.X, pub.Y)
}

func GenerateKeyBtc() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(S256Btc(), rand.Reader)
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
