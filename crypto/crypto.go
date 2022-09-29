package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/tchain/go-tchain-sdk/crypto/config"
	"github.com/tchain/go-tchain-sdk/crypto/secp"
	"github.com/tchain/go-tchain-sdk/crypto/sm2"
	"github.com/tchain/go-tchain-sdk/utils"
	"github.com/tchain/go-tchain-sdk/utils/math"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(cryptoType config.CryptoType, data ...[]byte) []byte {
	switch cryptoType {
	case config.SM2:
		return sm2.Keccak256Sm2(data...)
	case config.SECP256K1:
		return secp.Keccak256Btc(data...)
	default:
		return sm2.Keccak256Sm2(data...)
	}
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(cryptoType config.CryptoType, data ...[]byte) (h utils.Hash) {
	switch cryptoType {
	case config.SM2:
		return sm2.Keccak256HashSm2(data...)
	case config.SECP256K1:
		return secp.Keccak256HashBtc(data...)
	default:
		return sm2.Keccak256HashSm2(data...)
	}
}

// ToECDSA creates a private key with the given D value.
func ToECDSA(d []byte, cryptoType config.CryptoType) (*ecdsa.PrivateKey, error) {
	switch cryptoType {
	case config.SM2:
		return sm2.ToECDSASm2(d, true)
	case config.SECP256K1:
		return secp.ToECDSABtc(d, true)
	default:
		return sm2.ToECDSASm2(d, true)
	}
}

// ToECDSAUnsafe blindly converts a binary blob to a private key. It should almost
// never be used unless you are sure the input is valid and want to avoid hitting
// errors due to bad origin encoding (0 prefixes cut off).
func ToECDSAUnsafe(d []byte, cryptoType config.CryptoType) *ecdsa.PrivateKey {
	var priv *ecdsa.PrivateKey
	switch cryptoType {
	case config.SM2:
		priv, _ = sm2.ToECDSASm2(d, false)
	case config.SECP256K1:
		priv, _ = secp.ToECDSABtc(d, false)
	default:
		priv, _ = sm2.ToECDSASm2(d, false)
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

	if pub[0] != 4 { // uncompressed form
		return nil, errors.New("invalid pub byte")
	}
	x := new(big.Int).SetBytes(pub[1:33])
	y := new(big.Int).SetBytes(pub[33:])

	if sm2.S256Sm2().IsOnCurve(x, y) {
		return sm2.UnmarshalPubkeySm2(pub)
	}
	if secp.S256Btc().IsOnCurve(x, y) {
		return secp.UnmarshalPubkeyBtc(pub)
	}
	return sm2.UnmarshalPubkeySm2(pub)
}

func FromECDSAPub(p *ecdsa.PublicKey) []byte {
	if sm2.S256Sm2().IsOnCurve(p.X, p.Y) {
		return sm2.FromECDSAPubSm2(p)
	}
	if secp.S256Btc().IsOnCurve(p.X, p.Y) {
		return secp.FromECDSAPubBtc(p)
	}
	return sm2.FromECDSAPubSm2(p)
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA(hexkey string, cryptoType config.CryptoType) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToECDSA(b, cryptoType)
}

// LoadECDSA loads a secp256k1 private key from the given file.
func LoadECDSA(file string, cryptoType config.CryptoType) (*ecdsa.PrivateKey, error) {
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

func GenerateKey(cryptoType config.CryptoType) (*ecdsa.PrivateKey, error) {
	switch cryptoType {
	case config.SM2:
		return sm2.GenerateKeySm2()
	case config.SECP256K1:
		return secp.GenerateKeyBtc()
	default:
		return sm2.GenerateKeySm2()
	}
}

func PubkeyToAddress(p ecdsa.PublicKey) utils.Address {
	length := config.HashLength
	var publicKey []byte
	var hash []byte
	var prefix strings.Builder

	if sm2.S256Sm2().IsOnCurve(p.X, p.Y) {
		publicKey = sm2.CompressPubkeySm2(&p)
		hash = sm2.Keccak256Sm2(publicKey)
		prefix.WriteString(string(config.SM2_Prefix))
	} else if secp.S256Btc().IsOnCurve(p.X, p.Y) {
		publicKey = secp.CompressPubkeyBtc(&p)
		hash = secp.Keccak256Btc(publicKey)
		prefix.WriteString(string(config.SECP256K1_Prefix))
	} else {
		publicKey = sm2.CompressPubkeySm2(&p)
		hash = sm2.Keccak256Sm2(publicKey)
		prefix.WriteString(string(config.SM2_Prefix))
	}

	hashLength := len(hash)
	if hashLength < length {
		length = hashLength
	}
	h := hash[hashLength-length:]

	var encode string
	encode = base58.Encode(h)
	prefix.WriteString(string(config.BASE58_Prefix))

	return utils.StringToAddress(prefix.String() + encode)
}
