package test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"testing"
)

func checkError(t *testing.T, input string, got, want error) bool {
	if got == nil {
		if want != nil {
			t.Errorf("input %s: got no error, want %q", input, want)
			return false
		}
		return true
	}
	if want == nil {
		t.Errorf("input %s: unexpected error %q", input, got)
	} else if got.Error() != want.Error() {
		t.Errorf("input %s: got error %q, want %q", input, got, want)
	}
	return false
}

func checkNumberError(t *testing.T, input interface{}, got, want error) bool {
	if got == nil {
		if want != nil {
			t.Errorf("input %v: got no error, want %q", input, want)
			return false
		}
		return true
	}
	if want == nil {
		t.Errorf("input %v: unexpected error %q", input, got)
	} else if got.Error() != want.Error() {
		t.Errorf("input %v: got error %q, want %q", input, got, want)
	}
	return false
}

var (
	kbifer, _ = new(big.Int).SetString("1000000000000000000000", 10)
	grand, _  = new(big.Int).SetString("1000000000000000000000", 10)
	mbifer, _ = new(big.Int).SetString("1000000000000000000000000", 10)
	gbifer, _ = new(big.Int).SetString("1000000000000000000000000000", 10)
	tbifer, _ = new(big.Int).SetString("1011100000000000000000000000000", 10)
)

func TestToWei(t *testing.T) {
	for _, test := range []struct {
		balance string
		uint    []string
		expect  *big.Int
		wantErr error
	}{
		{"001.00000001", []string{"kwei"}, nil, errors.New(fmt.Sprintf("while converting number %s to wei,  too many decimal points not met", "001.00000001"))},
		{"-001.001", []string{"kwei"}, big.NewInt(-1001), nil},
		{"1.10", []string{"kwei"}, big.NewInt(1100), nil},
		{"1.000000000555599001", nil, big.NewInt(1000000000555599001), nil},
		{"1", nil, big.NewInt(1000000000000000000), nil},
		{"1", []string{"wei"}, big.NewInt(1), nil},
		{"1", []string{"kwei"}, big.NewInt(1000), nil},
		{"1", []string{"mwei"}, big.NewInt(1000000), nil},
		{"1", []string{"gwei"}, big.NewInt(1000000000), nil},
		{"1", []string{"microbif"}, big.NewInt(1000000000000), nil},
		{"1", []string{"micro"}, big.NewInt(1000000000000), nil},
		{"1", []string{"millibif"}, big.NewInt(1000000000000000), nil},
		{"1", []string{"milli"}, big.NewInt(1000000000000000), nil},
		{"1", []string{"bif"}, big.NewInt(1000000000000000000), nil},
	} {
		res, err := utils.ToWei(test.balance, test.uint...)
		if !checkNumberError(t, test.balance, err, test.wantErr) {
			continue
		}
		if res.Cmp(test.expect) != 0 {
			t.Errorf("balance is  %v, uint is %v, value mismatch: got %v, want %v", test.balance, test.uint, res, test.expect)
			continue
		}
	}
}

func TestFromWei(t *testing.T) {
	for _, test := range []struct {
		balance *big.Int
		uint    []string
		expect  string
		wantErr error
	}{
		{big.NewInt(-1000000000555599118), nil, "-1.000000000555599118", nil},
		{big.NewInt(1000000000000000000), nil, "1", nil},
		{big.NewInt(1), []string{"wei"}, "1", nil},
		{big.NewInt(1000), []string{"kwei"}, "1", nil},
		{big.NewInt(1000000), []string{"mwei"}, "1", nil},
		{big.NewInt(1000000000), []string{"gwei"}, "1", nil},
		{big.NewInt(1000000000000), []string{"microbif"}, "1", nil},
		{big.NewInt(1000000000000), []string{"micro"}, "1", nil},
		{big.NewInt(1000000000000000), []string{"millibif"}, "1", nil},
		{big.NewInt(1000000000000000), []string{"milli"}, "1", nil},
		{big.NewInt(1000000000000000000), []string{"bif"}, "1", nil},
	} {
		res, err := utils.FromWei(test.balance, test.uint...)
		if !checkNumberError(t, test.balance, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("balance is  %v, uint is %v, value mismatch: got %v, want %s", test.balance, test.uint, res, test.expect)
			continue
		}
	}
}

func TestByteToHex(t *testing.T) {
	for _, test := range []struct {
		byteArr []byte
		expect  string
		wantErr error
	}{
		{[]byte{72, 101, 108, 108, 111, 33, 37}, "0x48656c6c6f2125", nil},
	} {
		res := utils.ByteToHex(test.byteArr)

		if res != test.expect {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.byteArr, res, test.expect)
			continue
		}
	}
}

func TestSm3(t *testing.T) {
	for _, test := range []struct {
		input   string
		expect  string
		wantErr error
	}{
		{"234", "0x4910a0057a0cf8c5297f47bc650a1e080c2f33cc5d6fd6a7428ef63d3e3a6e29", nil},
		{"0xea", "0x04991392012bc6618739cf856ca07878283b310b45aada09bcfaab1420aa72c0", nil},
		{"", "", utils.ErrInvalidSm3},
	} {
		res, err := utils.Sm3(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}

func TestSha3(t *testing.T) {
	for _, test := range []struct {
		input   string
		expect  string
		wantErr error
	}{
		{"234", "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79", nil},
		{"0xea", "0x2f20677459120677484f7104c76deb6846a2c071f9b3152c103bb12cd54d1a4a", nil},
		{"c1912fee45d61c87cc5ea59dae31190fffff232d", "0x4fb647abf5735d02e3a8a6c94c29977abed5bcc26e646c8e079e46759c1e0b04", nil},
		{"", "", utils.ErrInvalidSha3},
	} {
		res, err := utils.Sha3(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}

func TestSha3Raw(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{"234", "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79"},
		{"0xea", "0x2f20677459120677484f7104c76deb6846a2c071f9b3152c103bb12cd54d1a4a"},
		{"c1912fee45d61c87cc5ea59dae31190fffff232d", "0x4fb647abf5735d02e3a8a6c94c29977abed5bcc26e646c8e079e46759c1e0b04"},
		{"", utils.Sha3Null},
	} {
		res := utils.Sha3Raw(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}

func TestPadLeft(t *testing.T) {
	for _, test := range []struct {
		slice  []byte
		length int
		expect []byte
	}{
		{[]byte{1, 2, 3, 4}, 8, []byte{0, 0, 0, 0, 1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4}, 2, []byte{1, 2, 3, 4}},
	} {
		res := utils.LeftPadBytes(test.slice, test.length)
		if !bytes.Equal(res, test.expect) {
			t.Errorf("slice %v:  length %d  value mismatch: got %v, want %v", test.slice, test.length, res, test.expect)
			continue
		}
	}
}

func TestPadRight(t *testing.T) {
	for _, test := range []struct {
		slice  []byte
		length int
		expect []byte
	}{
		{[]byte{1, 2, 3, 4}, 8, []byte{1, 2, 3, 4, 0, 0, 0, 0}},
		{[]byte{1, 2, 3, 4}, 2, []byte{1, 2, 3, 4}},
	} {
		res := utils.RightPadBytes(test.slice, test.length)
		if !bytes.Equal(res, test.expect) {
			t.Errorf("slice %v:  length %d  value mismatch: got %v, want %v", test.slice, test.length, res, test.expect)
			continue
		}
	}
}

func TestToTwosComplement(t *testing.T) {
	for _, test := range []struct {
		input   *big.Int
		expect  string
		wantErr error
	}{
		{big.NewInt(-1), "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", nil},
		{big.NewInt(0x1), "0x0000000000000000000000000000000000000000000000000000000000000001", nil},
		{big.NewInt(-15), "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1", nil},
		{big.NewInt(-0x1), "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", nil},
	} {
		res, err := utils.ToTwosComplement(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("input is %v: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}

}

func TestIsHex(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect bool
	}{
		{"0xZ1912", false},
		{"Hello", false},
		{"c1912", false},
		{"0xc1912", false},
		{"0x1912", false},
		{"1912", true},
	} {
		res := utils.IsHex(test.input)
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
		res := utils.IsHexStrict(test.input)
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
		res := utils.RandomHex(test.size)
		if len(res) != test.expectLen && utils.IsHex(res) {
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
		res, err := utils.HexToBytes(test.input)
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
		res, err := utils.HexToUint64(test.input)
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
		res, err := utils.HexToBigInt(test.input)
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
		res, err := utils.HexToUtf8(test.input)
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
		res, err := utils.HexToAscii(test.input)
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
		res := utils.Utf8ToHex(test.input)
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
		res := utils.AsciiToHex(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}
