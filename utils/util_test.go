package utils

import (
	"math/big"
	"testing"
)

type test3 struct {
	input        [] byte
	want         interface{}
	wantErr      error // if set, decoding must fail on any platform
	wantErr32bit error // if set, decoding must fail on 32bit platforms (used for Uint tests)
}

type test4 struct {
	input        interface{}
	want         interface{}
	wantErr      error // if set, decoding must fail on any platform
	wantErr32bit error // if set, decoding must fail on 32bit platforms (used for Uint tests)
}
var (
	byteToHexTests = []test3{
		{input: []byte{72, 101, 108, 108, 111, 33, 37} , want: "0x48656c6c6f2125"},
	}
	sha3Tests = []test1{
		{input:`234`, want: "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79"},
		{input:`0xea`, want: "0x2f20677459120677484f7104c76deb6846a2c071f9b3152c103bb12cd54d1a4a"},
		{input:`c1912fee45d61c87cc5ea59dae31190fffff232d`, want: "0x4fb647abf5735d02e3a8a6c94c29977abed5bcc26e646c8e079e46759c1e0b04"},
	}

	sha3RawTests = []test1{
		{input:`234`, want: "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79"},
		{input:`0xea`, want: "0x2f20677459120677484f7104c76deb6846a2c071f9b3152c103bb12cd54d1a4a"},
	}

	checkAddressChecksumTests  = []test1{
		{input:`0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d`, want: true},
	}

	isAddressTests  = []test1{
		{input:`0xc1912fee45d61c87cc5ea59dae31190fffff232d`,   want: true},
		{input:`c1912fee45d61c87cc5ea59dae31190fffff232d`,     want: true},
		{input:`0XC1912FEE45D61C87CC5EA59DAE31190FFFFF232D`,   want: true},
		{input:`0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d`,   want: true},
		{input:`0xC1912fEE45d61C87Cc5EA59DaE31190FFFFf232d`,   want: false },
	}

	toChecksumAddressTests  = []test1{
		{input:`0xc1912fee45d61c87cc5ea59dae31190fffff232d`, want: "0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d"},
		{input:`0XC1912FEE45D61C87CC5EA59DAE31190FFFFF232D`, want: "0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d"},
		{input:`0Xaaa`,                                      wantErr: ErrInvalidAddress},
		{input:`0xcg912fee45d61c87cc5ea59dae31190fffff232d`, wantErr: ErrInvalidAddress},
	}

	numberToHexTests  = []test4{
		{input:`234`,     want: "0xea"},
		{input:`0x123`,   want: "0x123"},
		{input: 234,      want: "0xea"},
		{input: 123,      want: "0x7b"},
		{input: "sss",    wantErr: ErrNumberString},
		{input: -1,       wantErr: ErrNegInt},
		{input: "-1",     wantErr: ErrNumberString},
	}

	toHexTests  = []test4{
		{input:`234`,               want: "0xea"},
		{input: 234,                want: "0xea"},
		{input: "I have 100€",      want: "0x49206861766520313030e282ac"},
		{input: "sss",              want: "0x737373"},
		{input: -1,                 wantErr: ErrNegInt},
		{input: "-1",               want: "0x2d31"},
		{input: big.NewInt(-2),  wantErr: ErrBigNegInt},
		{input: 1.2,                wantErr: ErrNumberInput},
	}

	toBNTests = []test4{
		{input:`1234`,                  want: big.NewInt(1234)},
		{input:`0xea`,                  want: big.NewInt(234)},
		{input: 234,                    want: big.NewInt(234)},
		{input: 0x2,                    want: big.NewInt(2)},
		{input: -1,                     want: big.NewInt(-1)},
		{input: big.NewRat(2,1),  wantErr:  ErrNumberInput},
		{input: "sss",                  wantErr: ErrNumberString},
	}

	isBNTests   = []test4{
		{input:`-1`,                    want:  false},
		{input:`-asd`,                  want:  false},
		{input: 0x1,                    want:  false},
		{input: big.NewInt(2),       want:  true},
		{input: big.NewFloat(2),     want:  false},
		{input: big.NewRat(2,1),  want:  false},
	}
	toTwosComplementTests  = []test4{
		{input:`-1`, want: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{input: -1, want: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{input: 0x1, want: "0x0000000000000000000000000000000000000000000000000000000000000001"},
		{input: -15, want: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1"},
		{input: -0x1, want: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{input: "asd", wantErr: ErrNumberInput},
	}
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

func TestByteToHex(t *testing.T) {
	for _, test := range byteToHexTests {
		res := ByteToHex(test.input)

		if res != test.want.(string) {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestSha3(t *testing.T){
	for _, test := range sha3Tests {
		res, err := Sha3(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}


func TestSha3Raw(t *testing.T){
	for _, test := range sha3RawTests {
		res := Sha3Raw(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestCheckAddressChecksum(t *testing.T){
	for _, test := range checkAddressChecksumTests {
		res := CheckAddressChecksum(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestIsAddress(t *testing.T){
	for _, test := range isAddressTests {
		res := IsAddress(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestToChecksumAddress(t *testing.T){
	for _, test := range toChecksumAddressTests {
		res, err := ToChecksumAddress(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestNumberToHex(t *testing.T){
	for _, test := range numberToHexTests {
		res, err := NumberToHex(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestToHex(t *testing.T){
	for _, test := range toHexTests {
		res, err := ToHex(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestToBN(t *testing.T){
	for _, test := range toBNTests {
		res, err := ToBN(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res.Cmp(test.want.(*big.Int)) != 0 {
			t.Errorf("input %v: value mismatch: got %v, want %v", test.input, res, test.want)
			continue
		}
	}
}

func TestIsBN(t *testing.T){
	for _, test := range isBNTests {
		res := IsBN(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %v: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestToTwosComplement(t *testing.T){
	for _, test := range toTwosComplementTests {
		res, err := ToTwosComplement(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}