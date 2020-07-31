package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bif/bif-sdk-go/utils"
	"io/ioutil"
	"math/big"
	"testing"
)

var util = utils.NewUtils()

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

func TestToWei(t *testing.T) {
	for _, test := range []struct {
		balance string
		uint    []string
		expect  *big.Int
		wantErr error
	}{
		// {"1", nil, big.NewInt(1000000000000000000), nil},
		{"1.000000000555599001", nil, big.NewInt(1000000000555599001), nil},
		// {"1", []string{"bifer"}, big.NewInt(1000000000000000000), nil},
		// {"1", []string{"finney"}, big.NewInt(1000000000000000), nil},
		// {"1", []string{"szabo"}, big.NewInt(1000000000000), nil},
		// {"1", []string{"shannon"}, big.NewInt(1000000000), nil},
		// {"1", []string{"shannon", "bifer"}, big.NewInt(1000000000), nil},
		// {"1", []string{"shan"}, nil, utils.ErrUintNoExist},
	} {
		res, err := util.ToWei(test.balance, test.uint...)
		if !checkNumberError(t, test.balance, err, test.wantErr) {
			continue
		}
		if res.Cmp(test.expect) != 0 {
			t.Errorf("balance is  %v, uint is %v, value mismatch: got %v, want %s", test.balance, test.uint, res, test.expect)
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
		// {big.NewInt(1000000000000000000), nil, "1", nil},
		{big.NewInt(1000000001000000013), nil, "1.0000002", nil},
		{big.NewInt(10000001), []string{"tbifer"}, "1.0000001e-23", nil},
		// {big.NewInt(1000000000000000000), []string{"bifer"}, "1", nil},
		// {big.NewInt(1000000000000000), []string{"finney"}, "1", nil},
		// {big.NewInt(1000000000000), []string{"szabo"}, "1", nil},
		// {big.NewInt(1000000000), []string{"shannon"}, "1", nil},
		// {big.NewInt(1000000000), []string{"shannon", "bifer"}, "1", nil},
		// {big.NewInt(1), []string{"shan"}, "", utils.ErrUintNoExist},
	} {
		res, err := util.FromWei(test.balance, test.uint...)
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
		res := util.ByteToHex(test.byteArr)

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
		res, err := util.Sm3(test.input)
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
		res, err := util.Sha3(test.input)
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
		res := util.Sha3Raw(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}
}

func TestCheckBidChecksum(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect bool
	}{
		{"0x6469643a6269643ab9A09F33669435E7ef1bEaEd", true},
		{"0x6469643a6269643ab9A09F33669435E7ef1bEaEd", true},
		{"0x6469643a6269643AC117c1794Fc7a27Bd301AE52", true},
		{"0x6469643a6269643AC117c1794Fc7a27Bd301Ae52", false},
	} {
		res := util.CheckBidChecksum(test.input)
		if res != test.expect {
			t.Errorf("input %s: value mismatch: got %t, want %t", test.input, res, test.expect)
			continue
		}
	}
}

func TestIsBid(t *testing.T) {
	for _, test := range []struct {
		bidAddress string
		expect     bool
	}{
		{"did:bid:590ed37615bdfefa496224c7", true},
		{"0x6469643a6269643a590ed37615bdfefa496224c7", true},
		{"6469643a6269643a590ed37615bdfefa496224c7", false},
		{"6469643a6269643A590ed37615bdfefa496224c7", false},
	} {
		isBid := util.IsBid(test.bidAddress)

		if isBid != test.expect {
			t.Errorf("bidAddress is  %s, got %t, want %t", test.bidAddress, isBid, test.expect)
			continue
		}
	}
}

func TestToChecksumBid(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{"did:bid:590ed37615bdfefa496224c7", "0x6469643A6269643A590eD37615BDFeFA496224C7"},
		{"did:bid:c935bd29a90fbeea87badf3e", "0x6469643a6269643AC935BD29A90fbEea87baDF3e"},
		{"0x6469643a6269643AC935BD29A90fbEea87baDF3e", "0x6469643a6269643AC935BD29A90fbEea87baDF3e"},
		{"6469643a6269643aC935BD29A90fbEea87baDF3e", "0x6469643a6269643AC935BD29A90fbEea87baDF3e"},
	} {
		res := util.ToChecksumBid(test.input)
		if res != test.expect {
			t.Errorf("input %v: value mismatch: got %s, want %s", test.input, res, test.expect)
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
		res := util.LeftPadBytes(test.slice, test.length)
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
		res := util.RightPadBytes(test.slice, test.length)
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
		res, err := util.ToTwosComplement(test.input)
		if !checkNumberError(t, test.input, err, test.wantErr) {
			continue
		}
		if res != test.expect {
			t.Errorf("input is %v: value mismatch: got %s, want %s", test.input, res, test.expect)
			continue
		}
	}

}

func TestByteCodeDeploy(t *testing.T) {
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
	fmt.Println(res, err)
}

func TestByteCodeInteract(t *testing.T) {
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
	fmt.Println(res, err)
}
