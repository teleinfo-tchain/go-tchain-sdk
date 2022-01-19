package sm2

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/prometheus/common/log"
	gmsm2 "github.com/teleinfo-bif/bit-gmsm/sm2"
	gmsm3 "github.com/teleinfo-bif/bit-gmsm/sm3"
)

var errInvalidPubkey = errors.New("invalid sm2 public key")

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256Sm2(data ...[]byte) []byte {
	d := gmsm3.New()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256HashSm2(data ...[]byte) (h utils.Hash) {
	d := gmsm3.New()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func ToECDSASm2(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	privateKey, err := gmsm2.GenerateKeyBySeed(d, true)
	if err != nil {
		log.Error("国密版公私钥创建错误，请重新创建，错误：%s", err)
	}
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: S256Sm2(),
			X:     privateKey.X,
			Y:     privateKey.Y,
		},
		D: privateKey.D,
	}, nil
}

// UnmarshalPubkey converts bytes to a secp256k1 public key.
func UnmarshalPubkeySm2(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(S256Sm2(), pub)
	if x == nil {
		return nil, errInvalidPubkey
	}
	return &ecdsa.PublicKey{Curve: S256Sm2(), X: x, Y: y}, nil
}

func FromECDSAPubSm2(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	key := gmsm2.PublicKey{
		Curve: S256Sm2(),
		X:     pub.X,
		Y:     pub.Y,
	}
	return gmsm2.Compress(&key)
}

func GenerateKeySm2() (*ecdsa.PrivateKey, error) {
	privateKey, err := gmsm2.GenerateKey()
	if err != nil {
		log.Error("国密版公私钥创建错误，请重新创建，错误：%s", err)
	}
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: S256Sm2(),
			X:     privateKey.X,
			Y:     privateKey.Y,
		},
		D: privateKey.D,
	}, nil
}
