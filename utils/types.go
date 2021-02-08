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
	"fmt"
	"github.com/bif/bif-sdk-go/crypto/config"
	"github.com/btcsuite/btcutil/base58"
	"math/big"
	"reflect"
	"regexp"
	"strings"

	"github.com/bif/bif-sdk-go/utils/hexutil"
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
	AddressLength       = 24
	AddressPrefixString = "did:bid:hufa:"
	AddressSplit        = ":"
)

var (
	hashT             = reflect.TypeOf(Hash{})
	addressT          = reflect.TypeOf(Address{})
	AddressPrefixByte = []byte(AddressPrefixString)
	AddressRegexp     = regexp.MustCompile(`^(did:bid:hufa:(([a-z0-9]{4}):)?)?([zes][sft])([1-9a-km-zA-HJ-NP-Z]+)$`)
	ChainCodeRegexp   = regexp.MustCompile(`^[a-z0-9]{4}$`)
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

/////////// Address

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

func (a *Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
	//return a.Bytes(), nil
}

func (a *Address) UnmarshalText(input []byte) error {
	return a.unmarshalAddress(input)
}

func (a *Address) UnmarshalJSON(input []byte) error {
	return a.unmarshalAddress(input)
}

func (a *Address) unmarshalAddress(input []byte) error {
	inputLength := len(input)
	if inputLength < 32 {
		return fmt.Errorf("input address's length is short")
	}

	// 删除前后引号
	if input[inputLength-1] == 34 {
		input = input[:inputLength-1]
	}
	if input[0] == 34 {
		input = input[1:]
	}

	s := string(input)
	if !AddressRegexp.MatchString(s) {
		return fmt.Errorf("input address is error, addr=%s", s)
	}

	subString := AddressRegexp.FindStringSubmatch(s)

	save := subString[4]
	suf := subString[5]

	decode := base58.Decode(suf)
	if len(decode) != config.HashLength {
		return fmt.Errorf("input address's length is error, len=%d", len(decode))
	}

	var addr bytes.Buffer
	addr.WriteString(save)
	addr.Write(decode)

	copy(a[:], addr.Bytes())

	return nil
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	//if !bytes.HasPrefix(b, AddressPrefixByte) {
	//	return Address{}
	//}

	var a Address
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[:], b)
	return a
}

func StringToAddress(s string) Address {
	s = strings.TrimSpace(s)
	if s == "" {
		return Address{}
	}
	if len(s) < 24 {
		return Address{}
	}

	if !AddressRegexp.MatchString(s) {
		return Address{}
	}

	subString := AddressRegexp.FindStringSubmatch(s)

	//code := subString[2]
	save := subString[4]
	suf := subString[5]

	decode := base58.Decode(suf)
	if len(decode) != config.HashLength {
		return Address{}
	}

	var addr bytes.Buffer
	addr.WriteString(save)
	addr.Write(decode)

	return BytesToAddress(addr.Bytes())
}

func HashToAddress(hash []byte) Address {
	length := config.HashLength
	var prefix strings.Builder
	//prefix.WriteString(AddressPrefixString)
	prefix.WriteString(string(config.SECP256K1_Prefix))

	hashLength := len(hash)
	if hashLength < length {
		length = hashLength
	}
	h := hash[hashLength-length:]

	var encode string
	encode = base58.Encode(h)
	prefix.WriteString(string(config.BASE58_Prefix))
	//prefix.WriteString(string(config.HashLength20))

	return StringToAddress(prefix.String() + encode)
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

func (a Address) ValidBytes() []byte { return a[0:] }

func (a Address) CryptoType() config.CryptoType {
	if bytes.HasPrefix(a.Bytes(), []byte(config.SM2_Prefix)) {
		return config.SM2
	} else if bytes.HasPrefix(a.Bytes(), []byte(config.SECP256K1_Prefix)) {
		return config.SECP256K1
	} else {
		return config.SECP256K1
	}
}

// String implements fmt.Stringer.
func (a Address) String(chainCode string) string {
	var addressString strings.Builder
	if a == (Address{}) {
		return ""
	}
	addr := a.Bytes()
	addrLen := len(addr)

	addressString.WriteString(AddressPrefixString)
	if chainCode != "" {
		addressString.WriteString(chainCode)
		addressString.WriteString(AddressSplit)
	}
	addressString.Write(addr[:addrLen-config.HashLength])
	addressString.WriteString(base58.Encode(addr[addrLen-config.HashLength:]))

	return addressString.String()
}

func (a Address) Equal(addr Address) bool {
	if a.String("") == addr.String("") {
		return true
	}
	return false
}

func (a Address) EqualString(addr string) bool {
	if a.String("") == StringToAddress(addr).String("") {
		return true
	}
	return false
}
