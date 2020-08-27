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

package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"
	"math/big"
	"math/rand"
	"reflect"
	"strings"

	"github.com/bif/bif-sdk-go/utils/hexutil"
	// "golang.org/x/crypto/sha3"
)

// todo: 需要置于加密模块中
type (
	AddressPrefix  string
	HashLengthType string
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength       = 31
	AddressPrefixString = "did:bid:"
)

var (
	hashT             = reflect.TypeOf(Hash{})
	addressT          = reflect.TypeOf(Address{})
	AddressPrefixByte = []byte(AddressPrefixString)
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func StringToHash(s string) Hash { return BytesToHash([]byte(s)) } // dep: Istanbul

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

func HexToHashWithOutPre(s string) Hash { return BytesToHash(FromHexWithoutPre(s)) }

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

// String implements the stringer common and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer common used for logging.
func (h Hash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h[:])
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func EmptyHash(h Hash) bool {
	return h == Hash{}
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != HashLength {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), HashLength)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (_ Hash) ImplementsGraphQLType(name string) bool { return name == "Bytes32" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (h *Hash) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		*h = HexToHash(input)
	default:
		err = fmt.Errorf("Unexpected type for Bytes32: %v", input)
	}
	return err
}

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	if !bytes.HasPrefix(b, AddressPrefixByte) {
		return Address{}
	}
	var a Address
	a.SetBytes(b)
	return a
}

func BytesToAddressWithoutPre(b []byte) Address {
	var a Address
	a.SetBytesWithoutPre(b)
	return a
}

func StringToAddressWithoutPre(s string) Address { return BytesToAddressWithoutPre([]byte(s)) }

func StringToAddress(s string) Address {
	s = strings.TrimSpace(s)
	if s == "" {
		return Address{}
	}
	// todo:Base58的长度在一定范围内，参展base58的实现修改
	if len(s) < 32 {
		return Address{}
	}

	if !strings.HasPrefix(s, AddressPrefixString) {
		return Address{}
	}
	var b bytes.Buffer
	b.Write([]byte(s[:11]))
	// todo: is always Base58??
	b.Write(base58.Decode(s[11:]))

	return BytesToAddress(b.Bytes())
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsHexAddress verifies whether a string can represent a valid hex-encoded address or not.
func IsHexAddress(s string) bool {
	// 暂时值判断前缀，后面等格式确定，再增加额外的判断
	if HasDidBidPrefix(s) {
		return true
	}
	return false
}

// ///////// Address

// Address represents the 31 byte address of an Bif account.
type Address [AddressLength]byte

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// 判断加密是否采用国密
func (a Address) IsSM2() bool {
	if bytes.HasPrefix(a.Bytes(), []byte(AddressPrefixString)) && a[8] == 90 {
		return true
	} else {
		return false
	}
}

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

// String implements fmt.Stringer.
func (a Address) String() string {
	// return a.Hex()
	tmpAddress := Address{}
	if a == tmpAddress {
		return ""
	}
	addr := a.Bytes()
	if !bytes.Equal(addr[:8], AddressPrefixByte) {
		return ""
	}
	prefix := fmt.Sprintf("%s", addr[:11])
	suffix := base58.Encode(addr[11:])
	return prefix + suffix
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer common used for logging.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

func (a *Address) SetBytesWithoutPre(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}
