package utils

import (
	"errors"
	"math/big"
	"strings"
)

const bifer string = "bifer"

// Errors
var (
	ErrUintNoExist   = errors.New("uint not exist")
	ErrParameter     = errors.New("error parameter number")
)

var BiferUint = map[string]string {
	"nobifer": "0",
	"wei":        "1",
	"kwei":       "1000",
	"Kwei":       "1000",
	"babbage":    "1000",
	"femtobifer": "1000",
	"mwei":       "1000000",
	"Mwei":       "1000000",
	"lovelace":   "1000000",
	"picobifer":  "1000000",
	"gwei":       "1000000000",
	"Gwei":       "1000000000",
	"shannon":    "1000000000",
	"nanobifer":  "1000000000",
	"nano":       "1000000000",
	"szabo":      "1000000000000",
	"microbifer": "1000000000000",
	"micro":      "1000000000000",
	"finney":     "1000000000000000",
	"millibifer": "1000000000000000",
	"milli":      "1000000000000000",
	"bifer":      "1000000000000000000",
	"kbifer":     "1000000000000000000000",
	"grand":      "1000000000000000000000",
	"mbifer":     "1000000000000000000000000",
	"gbifer":     "1000000000000000000000000000",
	"tbifer":     "1000000000000000000000000000000",
}

func getUintValue(uint string) (string, error){
	uint = strings.ToLower(uint)
	if BiferUint[uint]==""{
		return "", ErrUintNoExist
	}
	return BiferUint[uint],nil
}

// 是否只是接收字符串？？另外如果数太大/太小，go自动转换为科学计数法，是否合适
func ToWei(number string, values ...string) (string, error){
	if len(values)>1{
		return "", ErrParameter
	}
	uints := bifer
	for _, value := range values {
		uints = value
	}
	uintValue, err := getUintValue(uints)
	if err != nil{
		return "", err
	}
	multiplier0ne := new(big.Float)
	m, ok := multiplier0ne.SetString(number)
	if !ok {
		return "", ErrNumberString
	}
	multiplierTwo := new(big.Float)
	n, _ := multiplierTwo.SetString(uintValue)
	resValue := new(big.Float)
	return resValue.Mul(m,n).String(), nil
}

func FromWei(number string, values ...string) (string, error){
	if len(values)>1{
		return "", ErrParameter
	}
	uints := bifer
	for _, value := range values {
		uints = value
	}
	uintValue, err := getUintValue(uints)
	if err != nil{
		return "", err
	}
	// 被除数
	dividend := new(big.Float)
	m, ok := dividend.SetString(number)
	if !ok {
		return "", ErrNumberString
	}
	// 除数
	divisor := new(big.Float)
	n, _ := divisor.SetString(uintValue)
	resValue := new(big.Float)
	return resValue.Quo(m,n).String(), nil
}
