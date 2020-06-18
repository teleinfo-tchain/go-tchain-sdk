package hexutil

import (
	"bytes"
	"testing"
)

type test1 struct {
	input string
	want  interface{}
}

type test2 struct {
	input        string
	want         interface{}
	wantErr      error // if set, decoding must fail on any platform
	wantErr32bit error // if set, decoding must fail on 32bit platforms (used for Uint tests)
}

type test3 struct {
	input1 string
	input2 int
	input3 []string
	want  interface{}
}

var (
	isHexTests = []unmarshalTest{
		// false
		{input: `0xZ1912`, want: false},
		{input: `Hello`, want: false},
		// true
		{input: `0xc1912`, want: true},
		{input: `c1912`, want: true},
	}

	isHexStrictTests = []unmarshalTest{
		// false
		{input: `0xZ1912`, want: false},
		{input: `Hello`, want: false},
		{input: `c1912`, want: false},
		// true
		{input: `0xc1912`, want: true},
	}

	hexToBytesTests = []test1{
		{input: `0x000000ea`, want: []byte{0, 0, 0, 234}},
		{input: `0x48656c6c6f2125`, want: []byte{72,101,108,108,111,33,37}},
	}

	hexToUtf8Tests = []test1{
		{input: `0x49206861766520313030e282ac`, want: "I have 100€"},
	}

	hexToAsciiTests = []test1{
		{input: `0x4920686176652031303021`, want: "I have 100!"},
	}

	hexToNumberStringTests = []unmarshalTest{
		{input: `0xea`, want: "234"},
	}
	hexToNumberTests = []test2{
		{input: `0xea`, want: int64(234)},
		{input: `0xc1912fee45d61c8`, want: int64(871748888835219912)},
		{input: `0xc1912fee45d61c8111111`, wantErr: ErrUint64Range},
	}
	asciiToHexTests = []test1{
		{input: `I have 100!`, want: "0x4920686176652031303021"},
	}
	utf8ToHexTests = []test1{
		{input: `I have 100€`, want: "0x49206861766520313030e282ac"},
	}

	padLeftTests = []test3{
		{input1: `0x3456ff`, input2: 20, input3: []string{}, want: "0x000000000000003456ff"},
		{input1: `0x3456ff`, input2: 20, input3: []string{"x"}, want: "0xxxxxxxxxxxxxxx3456ff"},
		{input1: `Hello`, input2: 20, input3: []string{"x"}, want: "xxxxxxxxxxxxxxxHello"},
	}

	padRightTests = []test3{
		{input1: `0x3456ff`, input2: 20, input3: []string{}, want: "0x3456ff00000000000000"},
		{input1: `0x3456ff`, input2: 20, input3: []string{"x"}, want: "0x3456ffxxxxxxxxxxxxxx"},
		{input1: `0x3456ff`, input2: 20, input3: []string{"xd"}, want: "0x3456ffxdxdxdxdxdxdxd"},
		{input1: `Hello`, input2: 20, input3: []string{"x"}, want: "Helloxxxxxxxxxxxxxxx"},
	}

)
func TestIsHex(t *testing.T){
	for _, test := range isHexTests {
		res := IsHex(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestIsHexStrict(t *testing.T){
	for _, test := range isHexStrictTests {
		res := IsHexStrict(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToBytes(t *testing.T){
	for _, test := range hexToBytesTests {
		res, _:= HexToBytes(test.input)
		if !bytes.Equal(test.want.([]byte), res) {
			t.Errorf("input %s: value mismatch: got %v, want %v", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToUtf8(t *testing.T){
	for _, test := range hexToUtf8Tests {
		res, _:= HexToUtf8(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToAscii(t *testing.T){
	for _, test := range hexToAsciiTests {
		res, _:= HexToAscii(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToNumberString(t *testing.T){
	for _, test := range hexToNumberStringTests {
		res, _:= HexToNumberString(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestHexToNumber(t *testing.T){
	for _, test := range hexToNumberTests {
		res,err := HexToNumber(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(int64) {
			t.Errorf("input %s: value mismatch: got %d, want %d", test.input, res, test.want)
			continue
		}
	}
}

func TestAsciiToHex(t *testing.T){
	for _, test := range asciiToHexTests {
		res := AsciiToHex(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestUtf8ToHex(t *testing.T){
	for _, test := range utf8ToHexTests {
		res := Utf8ToHex(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestPadLeft(t *testing.T){
	for _, test := range padLeftTests {
		res := PadLeft(test.input1, test.input2, test.input3...)
		if res != test.want.(string) {
			t.Errorf("input %s:  number %d  value mismatch: got %s, want %s", test.input1, test.input2, res, test.want)
			continue
		}
	}
}

func TestPadRight(t *testing.T){
	for _, test := range padRightTests {
		res := PadRight(test.input1, test.input2, test.input3...)
		if res != test.want.(string) {
			t.Errorf("input %s:  number %d  value mismatch: got %s, want %s", test.input1, test.input2, res, test.want)
			continue
		}
	}
}

//Converts a negative numer into a two’s complement.
//func ToTwosComplement() string{
//
//}