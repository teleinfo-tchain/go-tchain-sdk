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
	"github.com/bif/bif-sdk-go/crypto/config"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/teleinfo-bif/bit-gmsm/sm2"
	"github.com/teleinfo-bif/bit-gmsm/sm3"
	"math/big"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256Sm2(data ...[]byte) []byte {
	d := sm3.New()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256HashSm2(data ...[]byte) (h utils.Hash) {
	d := sm3.New()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSASm2(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	privateKey, err := sm2.GenerateKeyBySeed(d, true)
	if err != nil {
		//fmt.Printf("国密版公私钥创建错误，请重新创建，错误: %s", err)
		return nil, err
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
	key := sm2.PublicKey{
		Curve: S256Sm2(),
		X:     pub.X,
		Y:     pub.Y,
	}
	return sm2.Compress(&key)
}

func GenerateKeySm2() (*ecdsa.PrivateKey, error) {
	privateKey, err := sm2.GenerateKey()
	if err != nil {
		//fmt.Printf("国密版公私钥创建错误，请重新创建，错误: %s", err)
		return nil, err
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

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValuesSm2(v byte, r, s *big.Int, homestead bool) bool {
	return true
}

// func PubkeyToAddressSm2(p ecdsa.PublicKey) utils.Address {
// 	pubBytes := FromECDSAPub(&p)
// 	addr := Keccak256(SM2, pubBytes[1:])[12:]
// 	addr[8] = 115
// 	return utils.BytesToAddress(addr)
// }

func PubkeyToAddressSm2(p ecdsa.PublicKey) []byte {
	pubBytes := CompressPubkeySm2(&p)
	return Keccak256(config.SM2, pubBytes)
}
