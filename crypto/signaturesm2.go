// Copyright 2019 The go-bif Authors
// This file is part of the go-bif library.
//
// The go-bif library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-bif library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-bif library. If not, see <http://www.gnu.org/licenses/>.

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/teleinfo-bif/bit-gmsm/sm2"
	"math/big"
)

// Ecrecover returns the uncompressed public key that created the given signature.
func EcrecoverSm2(hash, sig []byte) ([]byte, error) {
	bytes, err := sm2.RecoverPubKey(hash, sig)
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
	return &ecdsa.PublicKey{Curve: sm2.P256Sm2(), X: x, Y: y}, nil
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
	prvKey := sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: S256Sm2(),
			X:     prv.X,
			Y:     prv.Y,
		},
		D: prv.D,
	}
	var uid []byte
	return sm2.Sm2Sign(&prvKey, hash, uid)
}

func VerifySignatureSm2(pubkey ecdsa.PublicKey, hash, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}
	R := new(big.Int).SetBytes(signature[:32])
	S := new(big.Int).SetBytes(signature[32:])
	uid := PubkeyToAddress(pubkey).Bytes()
	p := sm2.PublicKey{X:pubkey.X, Y:pubkey.Y,Curve:S256Sm2(),}

	// Reject malleable signatures. libsecp256k1 does this check but btcec doesn't.
	return sm2.Sm2Verify(&p, hash, uid, R, S)
}

// S256 returns an instance of the secp256k1 curve.
func S256Sm2() elliptic.Curve {
	return sm2.P256Sm2()
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkeySm2(pubkey *ecdsa.PublicKey) []byte {
	key := &sm2.PublicKey{
		Curve: S256Sm2(),
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	return sm2.Compress(key)
}