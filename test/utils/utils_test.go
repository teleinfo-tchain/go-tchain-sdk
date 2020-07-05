package test

import (
	"encoding/json"
	"fmt"
	"github.com/bif/bif-sdk-go/utils"
	"io/ioutil"
	"math/big"
	"testing"
)

var (
	util = utils.NewUtils()

	toWeiTests = []test5{
		{inputOne: `1`, inputTwo: []string{}, want: "1e+18"},
		{inputOne: `1`, inputTwo: []string{"bifer"}, want: "1e+18"},
		{inputOne: `1`, inputTwo: []string{"finney"}, want: "1e+15"},
		{inputOne: `1`, inputTwo: []string{"szabo"}, want: "1e+12"},
		{inputOne: `1`, inputTwo: []string{"shannon"}, want: "1000000000"},
		{inputOne: `1`, inputTwo: []string{"shannon", "bifer"}, wantErr: utils.ErrParameter},
		{inputOne: `1`, inputTwo: []string{"shan"}, wantErr: utils.ErrUintNoExist},
		{inputOne: `1asd`, inputTwo: []string{}, wantErr: utils.ErrNumberString},
	}

	fromWeiTests = []test5{
		{inputOne: `1`, inputTwo: []string{}, want: "1e-18"},
		{inputOne: `1`, inputTwo: []string{"bifer"}, want: "1e-18"},
		{inputOne: `1`, inputTwo: []string{"finney"}, want: "1e-15"},
		{inputOne: `1`, inputTwo: []string{"szabo"}, want: "1e-12"},
		{inputOne: `1`, inputTwo: []string{"shannon"}, want: "1e-09"},
		{inputOne: `1`, inputTwo: []string{"shannon", "bifer"}, wantErr: utils.ErrParameter},
		{inputOne: `1`, inputTwo: []string{"shan"}, wantErr: utils.ErrUintNoExist},
		{inputOne: `1asd`, inputTwo: []string{}, wantErr: utils.ErrNumberString},
	}

	byteToHexTests = []test3{
		{input: []byte{72, 101, 108, 108, 111, 33, 37} , want: "0x48656c6c6f2125"},
	}

	sm3Tests = []test1{
		{input:`234`, want: "0x4910a0057a0cf8c5297f47bc650a1e080c2f33cc5d6fd6a7428ef63d3e3a6e29"},
		{input:`0xea`, want: "0x04991392012bc6618739cf856ca07878283b310b45aada09bcfaab1420aa72c0"},
		{input:``, wantErr: utils.ErrInvalidSm3},
	}

	sha3Tests = []test1{
		{input:`234`, want: "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79"},
		{input:`0xea`, want: "0x2f20677459120677484f7104c76deb6846a2c071f9b3152c103bb12cd54d1a4a"},
		{input:`c1912fee45d61c87cc5ea59dae31190fffff232d`, want: "0x4fb647abf5735d02e3a8a6c94c29977abed5bcc26e646c8e079e46759c1e0b04"},
		{input:``, wantErr: utils.ErrInvalidSha3},
	}

	sha3RawTests = []test1{
		{input:`234`, want: "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79"},
		{input:`0xea`, want: "0x2f20677459120677484f7104c76deb6846a2c071f9b3152c103bb12cd54d1a4a"},
		{input:``, want: "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},
	}

	checkBidChecksumTests  = []test1{
		{input:`0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d`, want: true},
	}

	isBidTests  = []test1{
		{input:`0xc1912fee45d61c87cc5ea59dae31190fffff232d`,   want: true},
		{input:`c1912fee45d61c87cc5ea59dae31190fffff232d`,     want: true},
		{input:`0XC1912FEE45D61C87CC5EA59DAE31190FFFFF232D`,   want: true},
		{input:`0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d`,   want: true},
		{input:`0xC1912fEE45d61C87Cc5EA59DaE31190FFFFf232d`,   want: false },
	}

	toChecksumBidTests  = []test1{
		{input:`0xc1912fee45d61c87cc5ea59dae31190fffff232d`, want: "0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d"},
		{input:`0XC1912FEE45D61C87CC5EA59DAE31190FFFFF232D`, want: "0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d"},
		{input:`0Xaaa`,                                      wantErr: utils.ErrInvalidBid},
		{input:`0xcg912fee45d61c87cc5ea59dae31190fffff232d`, wantErr: utils.ErrInvalidBid},
	}

	numberToHexTests  = []test4{
		{input:`234`,     want: "0xea"},
		{input:`0x123`,   want: "0x123"},
		{input: 234,      want: "0xea"},
		{input: 123,      want: "0x7b"},
		{input: -1,       want: "-0x1"},
		{input: "-1",     want: "-0x1"},
		{input: "sss",    wantErr: utils.ErrNumberString},
	}

	toHexTests  = []test4{
		{input:`234`,               want: "0xea"},
		{input: 234,                want: "0xea"},
		{input: "I have 100€",      want: "0x49206861766520313030e282ac"},
		{input: "sss",              want: "0x737373"},
		{input: -1,                 want: "-0x1"},
		{input: "-1",               want: "-0x1"},
		{input: big.NewInt(-2),  want: "-0x2"},
		{input: 1.2,                wantErr: utils.ErrNumberInput},
	}

	toBNTests = []test4{
		{input:`1234`,                  want: big.NewInt(1234)},
		{input:`0xea`,                  want: big.NewInt(234)},
		{input: 234,                    want: big.NewInt(234)},
		{input: 0x2,                    want: big.NewInt(2)},
		{input: -1,                     want: big.NewInt(-1)},
		{input: big.NewRat(2,1),  wantErr: utils.ErrNumberInput},
		{input: "sss",                  wantErr: utils.ErrNumberString},
	}

	isBNTests   = []test4{
		{input:`-1`,                    want:  false},
		{input:`-asd`,                  want:  false},
		{input: 0x1,                    want:  false},
		{input: big.NewInt(2),       want:  true},
		{input: big.NewFloat(2),     want:  false},
		{input: big.NewRat(2,1),  want:  false},
	}

	padLeftTests = []test2{
		{input1: `0x3456ff`, input2: 20, input3: []string{}, want: "0x000000000000003456ff"},
		{input1: `0x3456ff`, input2: 20, input3: []string{"x"}, want: "0xxxxxxxxxxxxxxx3456ff"},
		{input1: `Hello`, input2: 20, input3: []string{"x"}, want: "xxxxxxxxxxxxxxxHello"},
	}

	padRightTests = []test2{
		{input1: `0x3456ff`, input2: 20, input3: []string{}, want: "0x3456ff00000000000000"},
		{input1: `0x3456ff`, input2: 20, input3: []string{"x"}, want: "0x3456ffxxxxxxxxxxxxxx"},
		{input1: `0x3456ff`, input2: 20, input3: []string{"xd"}, want: "0x3456ffxdxdxdxdxdxdxd"},
		{input1: `Hello`, input2: 20, input3: []string{"x"}, want: "Helloxxxxxxxxxxxxxxx"},
	}

	toTwosComplementTests  = []test4{
		{input:`-1`, want: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{input: -1, want: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{input: 0x1, want: "0x0000000000000000000000000000000000000000000000000000000000000001"},
		{input: -15, want: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1"},
		{input: -0x1, want: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{input: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", wantErr: utils.ErrBigInt},
		{input: "asd", wantErr: utils.ErrNumberInput},
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

func TestToWei(t *testing.T){
	for _, test := range toWeiTests {
		res, err := util.ToWei(test.inputOne, test.inputTwo...)
		if !checkError(t, test.inputOne, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %s %v: value mismatch: got %s, want %s", test.inputOne, test.inputTwo, res, test.want)
			continue
		}
	}
}

func TestFromWei(t *testing.T){
	for _, test := range fromWeiTests {
		res, err := util.FromWei(test.inputOne, test.inputTwo...)
		if !checkError(t, test.inputOne, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %s %v: value mismatch: got %s, want %s", test.inputOne, test.inputTwo, res, test.want)
			continue
		}
	}
}

func TestByteToHex(t *testing.T) {
	for _, test := range byteToHexTests {
		res := util.ByteToHex(test.input)

		if res != test.want.(string) {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestSm3(t *testing.T){
	for _, test := range sm3Tests {
		res, err := util.Sm3(test.input)
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestSha3(t *testing.T){
	for _, test := range sha3Tests {
		res, err := util.Sha3(test.input)
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
		res := util.Sha3Raw(test.input)
		if res != test.want.(string) {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestCheckBidChecksum(t *testing.T){
	for _, test := range checkBidChecksumTests {
		res := util.CheckBidChecksum(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestIsBid(t *testing.T){
	for _, test := range isBidTests {
		res := util.IsBid(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestToChecksumBid(t *testing.T){
	for _, test := range toChecksumBidTests {
		res, err := util.ToChecksumBid(test.input)
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
		res, err := util.NumberToHex(test.input)
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
		res, err := util.ToHex(test.input)
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
		res, err := util.ToBN(test.input)
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
		res := util.IsBN(test.input)
		if res != test.want.(bool) {
			t.Errorf("input %v: value mismatch: got %t, want %t", test.input, res, test.want)
			continue
		}
	}
}

func TestPadLeft(t *testing.T){
	for _, test := range padLeftTests {
		res := util.PadLeft(test.input1, test.input2, test.input3...)
		if res != test.want.(string) {
			t.Errorf("input %s:  number %d  value mismatch: got %s, want %s", test.input1, test.input2, res, test.want)
			continue
		}
	}
}

func TestPadRight(t *testing.T){
	for _, test := range padRightTests {
		res := util.PadRight(test.input1, test.input2, test.input3...)
		if res != test.want.(string) {
			t.Errorf("input %s:  number %d  value mismatch: got %s, want %s", test.input1, test.input2, res, test.want)
			continue
		}
	}
}

func TestToTwosComplement(t *testing.T){
	for _, test := range toTwosComplementTests {
		res, err := util.ToTwosComplement(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.want)
			continue
		}
	}
}

func TestByteCodeDeploy(t *testing.T){
	content, err := ioutil.ReadFile("../resources/simple-contract.json")
	if err != nil {
		t.Errorf("File not exist")
		t.Error(err)
		t.FailNow()
	}

	type Contract struct {
		Abi      string `json:"abi"`
		Bytecode string `json:"bytecode"`
	}

	var unmarshalResponse Contract

	err = json.Unmarshal(content, &unmarshalResponse)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	res, err := util.ByteCodeDeploy(unmarshalResponse.Abi, unmarshalResponse.Bytecode, big.NewInt(2))
	fmt.Println(res,err)
}

func TestByteCodeInteract(t *testing.T){
	content, err := ioutil.ReadFile("../resources/simple-contract.json")
	if err != nil {
		t.Errorf("File not exist")
		t.Error(err)
		t.FailNow()
	}

	type Contract struct {
		Abi      string `json:"abi"`
		Bytecode string `json:"bytecode"`
	}

	var unmarshalResponse Contract

	err = json.Unmarshal(content, &unmarshalResponse)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	res, err := util.ByteCodeInteract(unmarshalResponse.Abi, "multiply", big.NewInt(2))
	fmt.Println(res,err)
}