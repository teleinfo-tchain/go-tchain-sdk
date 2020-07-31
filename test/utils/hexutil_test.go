package test

import (
	"bytes"
	"errors"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"math/big"
	"testing"
)

func TestIsHex(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect bool
	}{
		{"0xZ1912", false},
		{"Hello", false},
		{"c1912", false},
		{"0xc1912", false},
		{"0x1912", true},
		{"1912", true},
	} {
		res := util.IsHex(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.expect)
			continue
		}
	}
}

func TestIsHexStrict(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect bool
	}{
		{"0xZ1912", false},
		{"Hello", false},
		{"c1912", false},
		{"0xc1912", false},
		{"0x1912", true},
		{"1912", false},
	} {
		res := util.IsHexStrict(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.expect)
			continue
		}
	}
}

func TestRandomHex(t *testing.T) {
	for _, test := range []struct {
		size      int
		expectLen int
	}{
		{2, 4},
		{4, 8},
		{0, 0},
		{-1, 0},
	} {
		res := util.RandomHex(test.size)
		if len(res) != test.expectLen && util.IsHex(res) {
			t.Errorf("size %d: value mismatch: got %s, want %d", test.size, res, test.expectLen)
			continue
		}
	}
}

func TestHexToBytes(t *testing.T) {
	for _, test := range []struct {
		input   string
		expect  []byte
		wantErr error
	}{
		{"0x000000ea", []byte{0, 0, 0, 234}, nil},
		{"0x000000e", []byte{0, 0, 0, 234}, errors.New("hex string of odd length")},
		{"0x48656c6c6f2125", []byte{72, 101, 108, 108, 111, 33, 37}, nil},
		{"48656c6c6f2125", nil, hexutil.ErrMissingPrefix},
		{"", nil, hexutil.ErrEmptyString},
	} {
		res, err := util.HexToBytes(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}

		if !bytes.Equal(test.expect, res) {
			t.Errorf("input %s: value mismatch: got %v, want %v", test.input, res, test.expect)
			continue
		}
	}
}

func TestHexToUint64Number(t *testing.T) {
	for _, test := range []struct {
		input   string
		expect  uint64
		wantErr error
	}{
		{"0xea", 234, nil},
		{"0xc1912fee45d61c8", 871748888835219912, nil},
		{"0x0a", 0, hexutil.ErrLeadingZero},
		{"0xc1912fee45d61c8111111", 0, hexutil.ErrUint64Range},
		{"1a", 0, hexutil.ErrMissingPrefix},
	} {
		res, err := util.HexToUint64(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %d, want %d", test.input, res, test.wantErr)
			continue
		}
	}
}

func TestHexToBigNumber(t *testing.T) {
	for _, test := range []struct {
		input   string
		expect  *big.Int
		wantErr error
	}{
		{"0xea", big.NewInt(234), nil},
		{"0xc1912fee45d61c8", big.NewInt(871748888835219912), nil},
		{"0xc", big.NewInt(12), nil},
		{"0x0c", nil, hexutil.ErrLeadingZero},
		{"0xc1912fee45d61c8000000000000000000000000000000000000000000000000000", nil, hexutil.ErrBig256Range},
	} {
		res, err := util.HexToBigInt(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res.Cmp(test.expect) != 0 {
			t.Errorf("input %s: value mismatch: got %v, want %v", test.input, res, test.expect)
			continue
		}
	}
}

func TestHexToUtf8(t *testing.T) {
	for _, test := range []struct {
		input   string
		expect  string
		wantErr error
	}{
		{"0x4920686176652031303021", "I have 100!", nil},
		{"0xe4b8ad", "中", nil},
		{"0xe282ac", "€", nil},
		{"4e2d", "N-", errors.New("hex string without 0x prefix")},
	} {
		res, err := util.HexToUtf8(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}

func TestHexToAscii(t *testing.T) {
	for _, test := range []struct {
		input   string
		expect  string
		wantErr error
	}{
		{"0x4920686176652031303021", "I have 100!", nil},
		{"0x4e2d", "N-", nil},
		{"4e2d", "", errors.New("hex string without 0x prefix")},
	} {
		res, err := util.HexToAscii(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}

func TestUtf8ToHex(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{"I have 100!", "0x4920686176652031303021"},
		{"中", "0xe4b8ad"},
		{"€", "0xe282ac"},
	} {
		res := util.Utf8ToHex(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}

func TestAsciiToHex(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{"I have 100!", "0x4920686176652031303021"},
		{"中", "0x4e2d"},
		{"€", "0x20ac"},
	} {
		res := util.AsciiToHex(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}
