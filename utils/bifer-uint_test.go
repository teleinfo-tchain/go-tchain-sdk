package utils

import (
	"testing"
)

type test5 struct {
	inputOne     string
	inputTwo     []string
	want         interface{}
	wantErr      error // if set, decoding must fail on any platform
	wantErr32bit error // if set, decoding must fail on 32bit platforms (used for Uint tests)
}

var (
	toWeiTests = []test5{
		{inputOne: `1`, inputTwo: []string{}, want: "1e+18"},
		{inputOne: `1`, inputTwo: []string{"bifer"}, want: "1e+18"},
		{inputOne: `1`, inputTwo: []string{"finney"}, want: "1e+15"},
		{inputOne: `1`, inputTwo: []string{"szabo"}, want: "1e+12"},
		{inputOne: `1`, inputTwo: []string{"shannon"}, want: "1000000000"},
		{inputOne: `1`, inputTwo: []string{"shannon", "bifer"}, wantErr: ErrParameter},
		{inputOne: `1`, inputTwo: []string{"shan"}, wantErr: ErrUintNoExist},
		{inputOne: `1asd`, inputTwo: []string{}, wantErr: ErrNumberString},
	}
	fromWeiTests = []test5{
		{inputOne: `1`, inputTwo: []string{}, want: "1e-18"},
		{inputOne: `1`, inputTwo: []string{"bifer"}, want: "1e-18"},
		{inputOne: `1`, inputTwo: []string{"finney"}, want: "1e-15"},
		{inputOne: `1`, inputTwo: []string{"szabo"}, want: "1e-12"},
		{inputOne: `1`, inputTwo: []string{"shannon"}, want: "1e-09"},
		{inputOne: `1`, inputTwo: []string{"shannon", "bifer"}, wantErr: ErrParameter},
		{inputOne: `1`, inputTwo: []string{"shan"}, wantErr: ErrUintNoExist},
		{inputOne: `1asd`, inputTwo: []string{}, wantErr: ErrNumberString},
	}
)

func TestToWei(t *testing.T){
	for _, test := range toWeiTests {
		res, err := ToWei(test.inputOne, test.inputTwo...)
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
		res, err := FromWei(test.inputOne, test.inputTwo...)
		if !checkError(t, test.inputOne, err, test.wantErr) {
			continue
		}
		if res != test.want.(string) {
			t.Errorf("input %s %v: value mismatch: got %s, want %s", test.inputOne, test.inputTwo, res, test.want)
			continue
		}
	}
}