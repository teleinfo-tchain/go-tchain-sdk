package sm2

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	gmsm2 "github.com/teleinfo-bif/bit-gmsm/sm2"
	"math/big"
)

// Ecrecover returns the uncompressed public key that created the given signature.
func EcrecoverSm2(hash, sig []byte) ([]byte, error) {
	bytes, err := gmsm2.RecoverPubKey(hash, sig)
	return bytes, err
}

// SigToPub returns the public key that created the given signature.
func SigToPubSm2(hash, sig []byte) (*ecdsa.PublicKey, error) {
	// Convert to btcec input format with 'recovery id' v at the beginning.
	s, err := EcrecoverSm2(hash, sig)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(S256Sm2(), s)
	return &ecdsa.PublicKey{Curve: gmsm2.P256Sm2(), X: x, Y: y}, nil
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given hash cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func SignSm2(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	prvKey := gmsm2.PrivateKey{
		PublicKey: gmsm2.PublicKey{
			Curve: S256Sm2(),
			X:     prv.X,
			Y:     prv.Y,
		},
		D: prv.D,
	}
	var uid []byte

	return gmsm2.Sm2Sign(&prvKey, hash, uid)
}

func VerifySignatureSm2(pubKey *ecdsa.PublicKey, hash, signature []byte) bool {

	publicKey := &gmsm2.PublicKey{
		Curve: S256Sm2(),
		X:     pubKey.X,
		Y:     pubKey.Y,
	}

	var uid []byte

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	return gmsm2.Sm2Verify(publicKey, hash, uid, r, s)
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkeySm2(pubkey []byte) (*ecdsa.PublicKey, error) {
	public := gmsm2.Decompress(pubkey)
	if public.X == nil || public.Y == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{X: public.X, Y: public.Y, Curve: S256Sm2()}, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkeySm2(pubkey *ecdsa.PublicKey) []byte {
	key := &gmsm2.PublicKey{
		Curve: S256Sm2(),
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	return gmsm2.Compress(key)
}

// S256 returns an instance of the secp256k1 curve.
func S256Sm2() elliptic.Curve {
	return gmsm2.P256Sm2()
}
