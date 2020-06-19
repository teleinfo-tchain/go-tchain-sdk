package utils

import (
	"encoding/hex"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// web3.js对于hex string奇偶没有判定，是否还需要对数值类型的判断？
func IsHex(input string) bool {
	r, _ := regexp.Compile("^(0[x,X])?[A-F, a-f, 0-9]+$")
	return r.MatchString(input)
}

func IsHexStrict(input string) bool {
	r, _ := regexp.Compile("^(0[x,X])[A-F, a-f, 0-9]+$")
	return r.MatchString(input)
}

// 给定的字节大小生成伪随机HEX字符串
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


// 是否需要接收hex数值？？
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

// 是否接收hex数值？？
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

// hex 应该是将hex转换为十进制数，不支持big.int
func HexToNumber(input string) (interface{}, error){
	if !has0xPrefix(input){
		return "", hexutil.ErrMissingPrefix
	}
	input = input[2:]
	res, error :=strconv.ParseInt(input, 16, 64)
	if error != nil{
		return "", hexutil.ErrUint64Range
	}else {
		return res, nil
	}
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
