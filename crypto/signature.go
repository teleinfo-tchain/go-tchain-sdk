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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/crypto/config"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/prometheus/common/log"
	"time"
)

type Signature struct {
	PublicKey  []byte `json:"publicKey"    gencodec:"required"`  // 公钥，33字节，第1个字节是类型0, 1, 2，3，后32字节是公钥的x
	CryptoType []byte `json:"cryptoType"    gencodec:"required"` // 签名类型，0是sm2，1是secp
	Signature  []byte `json:"signature"    gencodec:"required"`  // 签名，64字字，前32字节是r，后32字节是s
}

func GenSignature(hash []byte, prv *ecdsa.PrivateKey, cryptoType config.CryptoType) (*Signature, error) {
	signature := &Signature{
		PublicKey:  make([]byte, 0, 33),
		CryptoType: make([]byte, 0, 1),
		Signature:  make([]byte, 0, 64),
	}
	sig, err := sign(hash, prv, cryptoType)
	if err != nil {
		return signature, err
	}

	var pubkey []byte
	var ct []byte

	switch cryptoType {
	case config.SM2:
		pubkey = CompressPubkeySm2(&prv.PublicKey)
		ct = []byte{0}
	case config.SECP256K1:
		pubkey = CompressPubkeyBtc(&prv.PublicKey)
		ct = []byte{1}
	default:
		pubkey = CompressPubkeySm2(&prv.PublicKey)
		ct = []byte{0}
	}

	signature.PublicKey = pubkey
	signature.CryptoType = ct
	signature.Signature = sig[:64]

	return signature, nil
}

func sign(hash []byte, prv *ecdsa.PrivateKey, cryptoType config.CryptoType) (sig []byte, err error) {
	switch cryptoType {
	case config.SM2:
		// start := time.Now()
		sig, err = SignSm2(hash, prv)
		// log.Debug("Sign SM2", "nanoseconds", common.PrettyDuration(time.Since(start)))
	case config.SECP256K1:
		// start := time.Now()
		sig, err = SignBtc(hash, prv)
		// log.Debug("Sign SECP256K1", "nanoseconds", common.PrettyDuration(time.Since(start)))
	default:
		// start := time.Now()
		sig, err = SignBtc(hash, prv)
		// log.Debug("Sign SECP256K1", "nanoseconds", common.PrettyDuration(time.Since(start)))
	}

	if bytes.Count(sig[:32], []byte{0}) == 32 {
		log.Error("Sign, r is nil", "sig", sig, "err", err)
		return nil, errors.New("r of sig of Sign is nil")
	}

	if bytes.Count(sig[32:64], []byte{0}) == 32 {
		log.Error("Sign, s is nil", "sig", sig, "err", err)
		return nil, errors.New("s of sig of Sign is nil")
	}

	return sig, err
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	if bytes.Count(sig, []byte{0}) == 66 {
		fmt.Println("Ecrecover", "sig", sig, "hash", hash)
		return nil, errors.New("sig of Ecrecover is nil")
	}

	switch sig[65] {
	case 0:
		start := time.Now()
		sig = sig[:len(sig)-1]
		bytes, err := EcrecoverSm2(hash, sig)
		fmt.Println("Ecrecover SM2", "nanoseconds", utils.PrettyDuration(time.Since(start)))
		return bytes, err
	case 1:
		start := time.Now()
		sig = sig[:len(sig)-1]
		bytes, err := EcrecoverBtc(hash, sig)
		fmt.Println("Ecrecover SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
		return bytes, err
	default:
		start := time.Now()
		sig = sig[:len(sig)-1]
		bytes, err := EcrecoverBtc(hash, sig)
		fmt.Println("Ecrecover SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
		return bytes, err
	}
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	if bytes.Count(sig, []byte{0}) == 66 {
		fmt.Println("SigToPub", "sig", sig, "hash", hash)
		return nil, errors.New("sig of SigToPub is nil")
	}

	switch sig[65] {
	case 0:
		start := time.Now()
		sig = sig[:len(sig)-1]
		bytes, err := SigToPubSm2(hash, sig)
		fmt.Println("SigToPub SM2", "nanoseconds", utils.PrettyDuration(time.Since(start)))
		return bytes, err
	case 1:
		start := time.Now()
		sig = sig[:len(sig)-1]
		bytes, err := SigToPubBtc(hash, sig)
		fmt.Println("SigToPub SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
		return bytes, err
	default:
		start := time.Now()
		sig = sig[:len(sig)-1]
		bytes, err := SigToPubBtc(hash, sig)
		fmt.Println("SigToPub SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
		return bytes, err
	}
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given hash cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
// hash = Hash(M) 32  hash = Hash(M) + uid + 标志位
func Sign(hash []byte, prv *ecdsa.PrivateKey, cryptoType config.CryptoType) (sig []byte, err error) {
	switch cryptoType {
	case config.SM2:
		start := time.Now()
		sig, err = SignSm2(hash, prv)
		fmt.Println("Sign SM2", "nanoseconds", utils.PrettyDuration(time.Since(start)))
	case config.SECP256K1:
		start := time.Now()
		sig, err = SignBtc(hash, prv)
		fmt.Println("Sign SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
	default:
		start := time.Now()
		sig, err = SignBtc(hash, prv)
		fmt.Println("Sign SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
	}

	if bytes.Count(sig[:32], []byte{0}) == 32 {
		fmt.Println("Sign, r is nil", "sig", sig, "err", err)
		return nil, errors.New("r of sig of Sign is nil")
	}

	if bytes.Count(sig[32:64], []byte{0}) == 32 {
		fmt.Println("Sign, s is nil", "sig", sig, "err", err)
		return nil, errors.New("s of sig of Sign is nil")
	}

	return sig, err
}

// S256 returns an instance of the secp256k1 curve.
func S256(cryptoType config.CryptoType) elliptic.Curve {
	switch cryptoType {
	case config.SM2:
		return S256Sm2()
	case config.SECP256K1:
		return S256Btc()
	default:
		return S256Sm2()
	}
}
