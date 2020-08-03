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
	"encoding/hex"
	"errors"
	"github.com/bif/bif-sdk-go/utils"
	"io"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/bif/bif-sdk-go/utils/math"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(cryptoType CryptoType, data ...[]byte) []byte {
	switch cryptoType {
	case SM2:
		return Keccak256Sm2(data...)
	case SECP256K1:
		return Keccak256Btc(data...)
	default:
		return Keccak256Sm2(data...)
	}
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(cryptoType CryptoType, data ...[]byte) (h utils.Hash) {
	switch cryptoType {
	case SM2:
		return Keccak256HashSm2(data...)
	case SECP256K1:
		return Keccak256HashBtc(data...)
	default:
		return Keccak256HashSm2(data...)
	}
}

// ToECDSA creates a private key with the given D value.
func ToECDSA(d []byte, cryptoType CryptoType) (*ecdsa.PrivateKey, error) {
	switch cryptoType {
	case SM2:
		return toECDSASm2(d, true)
	case SECP256K1:
		return toECDSABtc(d, true)
	default:
		return toECDSASm2(d, true)
	}
}

// ToECDSAUnsafe blindly converts a binary blob to a private key. It should almost
// never be used unless you are sure the input is valid and want to avoid hitting
// errors due to bad origin encoding (0 prefixes cut off).
func ToECDSAUnsafe(d []byte, cryptoType CryptoType) *ecdsa.PrivateKey {
	var priv *ecdsa.PrivateKey
	switch cryptoType {
	case SM2:
		priv, _ = toECDSASm2(d, false)
	case SECP256K1:
		priv, _ = toECDSABtc(d, false)
	default:
		priv, _ = toECDSASm2(d, false)
	}
	return priv
}

// FromECDSA exports a private key into a binary dump.
func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

// UnmarshalPubkey converts bytes to a secp256k1 public key.
func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	if nil == pub{
		return nil, errors.New("nil pub byte")
	}
	if pub[0] != 4 { // uncompressed form
		return nil, errors.New("invalid pub byte")
	}
	x := new(big.Int).SetBytes(pub[1:33])
	y := new(big.Int).SetBytes(pub[33:])

	if S256Sm2().IsOnCurve(x, y) {
		return UnmarshalPubkeySm2(pub)
	}
	if S256Btc().IsOnCurve(x, y) {
		return UnmarshalPubkeyBtc(pub)
	}
	return UnmarshalPubkeySm2(pub)
}

func FromECDSAPub(p *ecdsa.PublicKey) []byte {
	if S256Sm2().IsOnCurve(p.X, p.Y) {
		return FromECDSAPubSm2(p)
	}
	if S256Btc().IsOnCurve(p.X, p.Y) {
		return FromECDSAPubBtc(p)
	}
	return FromECDSAPubSm2(p)
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA(hexkey string, cryptoType CryptoType) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToECDSA(b, cryptoType)
}

// LoadECDSA loads a secp256k1 private key from the given file.
func LoadECDSA(file string, cryptoType CryptoType) (*ecdsa.PrivateKey, error) {
	buf := make([]byte, 64)
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err := io.ReadFull(fd, buf); err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	return ToECDSA(key, cryptoType)
}

// SaveECDSA saves a secp256k1 private key to the given file with
// restrictive permissions. The key data is saved hex-encoded.
func SaveECDSA(file string, key *ecdsa.PrivateKey) error {
	k := hex.EncodeToString(FromECDSA(key))
	return ioutil.WriteFile(file, []byte(k), 0600)
}

func GenerateKey(cryptoType CryptoType) (*ecdsa.PrivateKey, error) {
	switch cryptoType {
	case SM2:
		return GenerateKeySm2()
	case SECP256K1:
		return GenerateKeyBtc()
	default:
		return GenerateKeySm2()
	}
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	//if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
	//	return false
	//}
	//// reject upper range of s values (ECDSA malleability)
	//// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
	//if homestead && s.Cmp(secp256k1halfN) > 0 {
	//	return false
	//}
	//// Frontier: allow s to be in full N range
	//return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1)
	return true
}

func PubkeyToAddress(p ecdsa.PublicKey) utils.Address {
	if S256Sm2().IsOnCurve(p.X, p.Y) {
		return PubkeyToAddressSm2(p)
	}
	if S256Btc().IsOnCurve(p.X, p.Y) {
		return PubkeyToAddressBtc(p)
	}
	return PubkeyToAddressSm2(p)
}
