package utils

import (
	"encoding/hex"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"math/big"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func IsHex(input string) bool {
	r, _ := regexp.Compile("^(0[x,X])?[A-F, a-f, 0-9]+$")
	return r.MatchString(input)
}

func IsHexStrict(input string) bool {
	r, _ := regexp.Compile("^(0[x,X])[A-F, a-f, 0-9]+$")
	return r.MatchString(input)
}

// generate cryptographically strong pseudo-random HEX strings from a given byte size
func  RandomHex(size int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size*2; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return "0x"+string(result)
}


// hexString to byte
func HexToBytes(input string) ([]byte, error){
	return hexutil.Decode(input)
}

func HexToUtf8(input string) (string, error){
	res, err := hexutil.Decode(input)
	if err != nil{
		return "", err
	}else {
		return string(res), nil
	}
}

func HexToAscii(input string) (string, error){
	res, err := hexutil.Decode(input)
	if err != nil{
		return "", err
	}else {
		return string(res), nil
	}
}


// hexutil.DecodeUint64和hexutil.DecodeBig
func HexToNumberString(input string) (string, error){
	if len(input) == 0 {
		return "", hexutil.ErrEmptyString
	}
	if !has0xPrefix(input) {
		return "", hexutil.ErrMissingPrefix
	}
	var bigint *big.Int
	var ok bool
	bigint, ok = new(big.Int).SetString(input[2:], 16)
	if ok && bigint.BitLen() > 256 {
		return "", hexutil.ErrBig256Range
	}
	return bigint.String(), nil
}

// 将HexToNumber 变成两个， HexToUint64Number， HexToBigNumber

// hexStr 应该是将hex转换为uint64
func HexToUint64Number(input string) (uint64, error){
	return hexutil.DecodeUint64(input)
}

// hexStr 应该是将hex转换为uint64
func HexToBigNumber(input string) (*big.Int, error){
	return hexutil.DecodeBig(input)
}

func AsciiToHex(input string) string{
	return "0x"+hex.EncodeToString([]byte(input))
}

func Utf8ToHex(input string) string{
	return "0x"+hex.EncodeToString([]byte(input))
}

// Adds a padding on the left of a string, Useful for adding paddings to HEX strings.
// 是否还需考虑接收hex数值
func PadLeft(str string, characterAmount int, signs ...string) string {
	sign := "0"
	if len(signs)>=1{
		sign = signs[0]
	}
	if has0xPrefix(str){
		count := characterAmount+2-len([]rune(str))
		if count>0{
			return "0x"+strings.Repeat(sign, count)[:count] +str[2:]
		}
		return "0x"+ str[2:]
	}else {
		count := characterAmount-len([]rune(str))
		if count>0{
			return strings.Repeat(sign, count)[:count] +str
		}
		return str
	}
}

// Adds a padding on the right of a string, Useful for adding paddings to HEX strings.
func PadRight(str string, characterAmount int, signs ...string) string {
	sign := "0"
	if len(signs)>=1{
		sign = signs[0]
	}
	if has0xPrefix(str){
		count := characterAmount+2-len([]rune(str))
		if count>0{
			return str+ strings.Repeat(sign, count)[:count]
		}
		return "0x"+ str[2:]
	}else {
		count := characterAmount-len([]rune(str))
		if count>0{
			return str+ strings.Repeat(sign, count)[:count]
		}
		return str
	}
}
