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

// Package keystore 实现secp256k1私钥的加密存储
//
// 根据Bif Secret Storage规范，将密钥存储为加密的JSON文件。
package keystore

import (
	"crypto/ecdsa"
	"errors"
	"github.com/tchain/go-tchain-sdk/account/types"
	"path/filepath"
)

var (
	ErrDecrypt = errors.New("could not decrypt key with given passphrase")
)

// KeyStoreScheme is the protocol scheme prefixing account and wallet URLs.
const KeyStoreScheme = "keystore"

// KeyStore manages a key storage directory on disk.
type KeyStore struct {
	storage keyStore // Storage backend, might be cleartext or encrypted
}

// zeroKey zeroes a private key in memory.
func zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

// NewKeyStore creates a keystore for the given directory.
func NewKeyStore(keydir string, scryptN, scryptP int) *KeyStore {
	keydir, _ = filepath.Abs(keydir)
	ks := &KeyStore{storage: &keyStorePassphrase{keydir, scryptN, scryptP}}
	return ks
}

// ImportECDSA stores the given key into the key directory, encrypting it with the passphrase.
func (ks *KeyStore) ImportECDSA(priv *ecdsa.PrivateKey, passphrase, chainCode string) (types.Account, error) {
	key := NewKeyFromECDSA(priv)
	return ks.importKey(key, passphrase, chainCode)
}

func (ks *KeyStore) importKey(key *Key, passphrase, chainCode string) (types.Account, error) {
	a := types.Account{Address: key.Address, URL: types.URL{Scheme: KeyStoreScheme, Path: ks.storage.JoinPath(keyFileName(key.Address, chainCode))}}
	if err := ks.storage.StoreKey(a.URL.Path, key, passphrase, chainCode); err != nil {
		return types.Account{}, err
	}
	return a, nil
}
