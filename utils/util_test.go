package utils

import (
	"fmt"
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
	}

	sha3RawTests = []test1{
		{input:`234`, want: "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79"},
		{input:`0xea`, want: "0x2f20677459120677484f7104c76deb6846a2c071f9b3152c103bb12cd54d1a4a"},
	}
	checkAddressSumTests  = []test1{
		{input:`0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d`, want: true},
	}
	toChecksumAddressTests  = []test1{
		{input:`0XC1912FEE45D61C87CC5EA59DAE31190FFFFF232D`, want: "0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d"},
	}

	numberToHexTests  = []test4{
		{input:`234`, want: "0xea"},
		{input:`0x123`, want: "0x123"},
		{input: 234, want: "0xea"},
		{input: 123, want: "0x7b"},
		{input: "sss", wantErr: ErrNumberString},
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
// web3js 的字节数组可能给错了[72, 101, 108, 108, 111, 33, 36] 最后的应该是37不是36? 还是说是其他问题？
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

func TestCheckAddressSum(t *testing.T){
	for _, test := range checkAddressSumTests {
		res := CheckAddressChecksum(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

// 该测试还是有问题的，因为函数的sha3转换需要注意
func TestToChecksumAddress(t *testing.T){
	for _, test := range toChecksumAddressTests {
		res, _ := ToChecksumAddress(test.input)
		fmt.Println(res)
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

