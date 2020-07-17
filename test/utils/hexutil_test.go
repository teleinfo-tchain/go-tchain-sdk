package test

import (
	"bytes"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"math/big"
	"testing"
)

var (
	isHexTests = []test1{
		// false
		{input: `0xZ1912`, want: false},
		{input: `Hello`, want: false},
		// true
		{input: `0xc1912`, want: true},
		{input: `c1912`, want: true},
	}

	isHexStrictTests = []test1{
		// false
		{input: `0xZ1912`, want: false},
		{input: `Hello`, want: false},
		{input: `c1912`, want: false},
		// true
		{input: `0xc1912`, want: true},
	}

	hexToBytesTests = []test1{
		{input: `0x000000ea`, want: []byte{0, 0, 0, 234}},
		{input: `0x48656c6c6f2125`, want: []byte{72, 101, 108, 108, 111, 33, 37}},
	}

	hexToUtf8Tests = []test1{
		{input: `0x49206861766520313030e282ac`, want: "I have 100€"},
	}

	hexToAsciiTests = []test1{
		{input: `0x4920686176652031303021`, want: "I have 100!"},
	}

	hexToNumberStringTests = []test1{
		{input: `0xea`, want: "234"},
	}

	hexToUint64NumberTests = []test1{
		{input: `0xea`, want: uint64(234)},
		{input: `0xc1912fee45d61c8`, want: uint64(871748888835219912)},
		{input: `0xc1912fee45d61c8111111`, wantErr: hexutil.ErrUint64Range},
	}

	hexToBigNumberTests = []test1{
		{input: `0xea`, want: big.NewInt(234)},
		{input: `0xc1912fee45d61c8`, want: big.NewInt(871748888835219912)},
		{input: `0xc1912fee45d61c8000000000000000000000000000000000000000000000000000`, wantErr: hexutil.ErrBig256Range},
	}

	asciiToHexTests = []test1{
		{input: `I have 100!`, want: "0x4920686176652031303021"},
	}
	utf8ToHexTests = []test1{
		{input: `I have 100€`, want: "0x49206861766520313030e282ac"},
	}
)

func TestIsHex(t *testing.T) {
	for _, test := range isHexTests {
		res := util.IsHex(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestIsHexStrict(t *testing.T) {
	for _, test := range isHexStrictTests {
		res := util.IsHexStrict(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToBytes(t *testing.T) {
	for _, test := range hexToBytesTests {
		res, _ := util.HexToBytes(test.input)
		if !bytes.Equal(test.want.([]byte), res) {
			t.Errorf("input %s: value mismatch: got %v, want %v", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToUtf8(t *testing.T) {
	for _, test := range hexToUtf8Tests {
		res, _ := util.HexToUtf8(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToAscii(t *testing.T) {
	for _, test := range hexToAsciiTests {
		res, _ := util.HexToAscii(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToNumberString(t *testing.T) {
	for _, test := range hexToNumberStringTests {
		res, _ := util.HexToNumberString(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToUint64Number(t *testing.T) {
	for _, test := range hexToUint64NumberTests {
		res, err := util.HexToUint64Number(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(uint64) {
			t.Errorf("input %s: value mismatch: got %d, want %d", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToBigNumber(t *testing.T) {
	for _, test := range hexToBigNumberTests {
		res, err := util.HexToBigNumber(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res.Cmp(test.want.(*big.Int)) != 0 {
			t.Errorf("input %s: value mismatch: got %v, want %v", test.input, res, test.want)
			continue
		}
	}
}

func TestAsciiToHex(t *testing.T) {
	for _, test := range asciiToHexTests {
		res := util.AsciiToHex(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestUtf8ToHex(t *testing.T) {
	for _, test := range utf8ToHexTests {
		res := util.Utf8ToHex(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}
