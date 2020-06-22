package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/common/math"
	"github.com/bif/bif-sdk-go/crypto"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type StringType uint8

const (
	Hex       StringType = 0 // 十六进制
	Decimal   StringType= 1 // 十进制
	OtherStr  StringType= 2 // 其他类型（当做字符串处理）
)

// Errors
var (
	ErrInvalidSha3    = errors.New("invalid input")
	ErrInvalidAddress = errors.New("invalid Bif Address")
	ErrNumberString = errors.New("invalid number string")
	ErrNumberInput = errors.New("invalid number input")
	ErrBigNegInt = errors.New("negative big int")
	ErrNegInt = errors.New("negative int")
)

const Sha3Null = "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"

//  ????
const Sm3Null = "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"

func ByteToHex(byteArr []byte) string{
	return "0x"+hex.EncodeToString(byteArr)
}

func Sm3(str string) (string, error){
	var hexBytes []byte
	if IsHexStrict(str){
		hexBytes = common.Hex2Bytes(str[2:])
	}else{
		hexBytes = []byte(str)
	}
	resStr := ByteToHex(crypto.Keccak256(crypto.SM2, hexBytes))
	if resStr == Sm3Null {
		return "", ErrInvalidSha3
	}else{
		return resStr, nil
	}
}

// 是否还处理bigNumber
func Sha3(str string) (string,error){
	var hexBytes []byte
	if IsHexStrict(str){
		hexBytes, _ = hex.DecodeString(str[2:])
	}else{
		hexBytes = []byte(str)
	}
	resStr := ByteToHex(crypto.Keccak256(crypto.SECP256K1, hexBytes))
	if resStr == Sha3Null {
		return "", ErrInvalidSha3
	}else{
		return resStr, nil
	}
}

func Sha3Raw(str string) string{
	var hexBytes []byte
	if IsHexStrict(str){
		hexBytes, _ = hex.DecodeString(str[2:])
	}else{
		hexBytes = []byte(str)
	}
	resStr := ByteToHex(crypto.Keccak256(1, hexBytes))
	return resStr
}

// address string with "0x"
func CheckAddressChecksum(address string) bool{
	address = address[2:]
	addressHash, _ := Sha3(strings.ToLower(address))
	addressHash = addressHash[2:]
	for i := 0; i < 40; i++ {
		// 校验和的判断依据是根据其生成的方式判断的
		char := string(address[i])
		hashChar:= string(addressHash[i])
		upperChar := strings.ToUpper(char)
		lowerChar := strings.ToLower(char)
		parseHashChar, _ := strconv.ParseUint(hashChar, 16, 64)
		if parseHashChar>7 && upperChar != char || parseHashChar <= 7 && lowerChar != char {
			return false
		}
	}
	return true
}

func IsAddress(address string) bool{
	res1, _ := regexp.Compile("^(0[x,X])?[A-F, a-f, 0-9]{40}$")
	res2, _ := regexp.Compile("^(0[x,X])?[a-f, 0-9]{40}$")
	res3, _ := regexp.Compile("^(0[x,X])?[A-F, 0-9]{40}$")
	if !res1.MatchString(address){
		return false
	}else if res2.MatchString(address) || res3.MatchString(address) {
		return true
	}else {
		return CheckAddressChecksum(address)
	}
}

func ToChecksumAddress(address string) (string, error){
	res, _ := regexp.Compile("^(0[x,X])?[A-F, a-f, 0-9]{40}$")
	if !res.MatchString(address){
		return "", ErrInvalidAddress
	}
	address = strings.ToLower(address)[2:]
	addressHash , _:= Sha3(address)
	addressHash = addressHash[2:]
	checkSumAddress := "0x"
	for i := 0; i < len(address); i++ {
		hashChar:= string(addressHash[i])
		parseHashChar, _ := strconv.ParseUint(hashChar, 16, 64)
		if parseHashChar >7 {
			checkSumAddress += strings.ToUpper(string(address[i]))
		}else{
			checkSumAddress += string(address[i])
		}
	}
	return checkSumAddress, nil
}

// 如果是负数则不转换，提示错误
func intToHex(int interface{}) (string,error){
	str := fmt.Sprintf("%x", int)
	// represent  char "-"
	if string(str[0]) == "-" {
		return "", ErrNegInt
	}
	return "0x"+str, nil
}


// 判断字符串为10进制还是16进制(只是正数)，不支持其他进制判断
func judgeStrNumber(str string) StringType {
	if IsHexStrict(str){
		return Hex
	}
	r, _ := regexp.Compile("^[0-9]+$")
	if r.MatchString(str){
		return Decimal
	}
	return  OtherStr
}

// decimal str to hex str
func DecimalStrToHex(str string) (string, error){
	n := new(big.Int)
	n, ok := n.SetString(str, 10)
	if !ok{
		return "", ErrNumberString
	}
	return fmt.Sprintf("0x%x", n), nil
}

func NumberStrToHex(str string)(string, error){
	strType := judgeStrNumber(str)
	if strType == Decimal{
		return DecimalStrToHex(str)
	} else if strType == Hex{
		return "0x"+str[2:], nil
	}else {
		return "", ErrNumberString
	}
}

// 只对十进制和十六进制的字符串进行hex处理
func StrToHex(str string) (string, error){
	strType := judgeStrNumber(str)
	if strType == Decimal{
		return DecimalStrToHex(str)
	} else if strType == Hex{
		return "0x"+str[2:], nil //hex ,hexutil
	}else {
		return "0x"+hex.EncodeToString([]byte(str)), nil
	}
}

func NumberToHex(input interface{}) (string, error){
	switch input.(type) {
	case string:
		return NumberStrToHex(input.(string))
	case uint, int, uint8,uint16,uint32,uint64, int8, int16, int32,int64:
		return intToHex(input)
	case *big.Int:
		return hexutil.EncodeBig(input.(*big.Int)), nil
	default:
		return "", ErrNumberInput
	}
}


func ToHex(input interface{}) (string, error){
	switch input.(type) {
	case string:
		return StrToHex(input.(string))
	case uint, int, uint8,uint16,uint32,uint64, int8, int16, int32,int64:
		return intToHex(input)
	case *big.Int:
		return hexutil.EncodeBig(input.(*big.Int)), nil
	default:
		return "", ErrNumberInput
	}
}

// 判断字符串为10进制（包含正负数）还是16进制(只是正数)，不支持其他进制判断
func judgeStrNumberNeg(str string) StringType {
	if IsHexStrict(str){
		return Hex
	}
	r, _ := regexp.Compile("^(-)?[0-9]+$")
	if r.MatchString(str){
		return Decimal
	}
	return  OtherStr
}

// 对数值型字符串10进制和16进制(包含正负数)进行处理
func numberStrToBN(str string) (*big.Int,error){
	strType := judgeStrNumberNeg(str)
	n := new(big.Int)
	if strType == Decimal{
		n, ok := n.SetString(str, 10)
		if ok{
			return n,nil
		}
		return nil, ErrNumberString
	}else if strType == Hex{
		n, ok := n.SetString(str[2:], 16)
		if ok{
			return n,nil
		}
		return nil, ErrNumberString
	}else {
		return nil, ErrNumberString
	}

}

func numberToBN(number interface{}) (*big.Int, error){
	n := new(big.Int)
	n, ok := n.SetString(fmt.Sprintf("%d", number), 10)
	if ok{
		return n,nil
	}
	return nil, ErrNumberString
}

// 转换为大整型，对数值型字符串10进制和16进制(包含正负数)进行处理；  对数值进行大整型转换
func ToBN(input interface{}) (*big.Int, error){
	switch input.(type) {
	case string:
		return numberStrToBN(input.(string))
	case *big.Int:
		return input.(*big.Int), nil
	case uint, int, uint8,uint16,uint32,uint64, int8, int16, int32,int64:
		return numberToBN(input)
	default:
		return nil, ErrNumberInput
	}
}

// judge if is a big.Int
func IsBN(input interface{}) bool{
	_, ok := input.(*big.Int)
	if ok{
		return true
	}
	return false
}

// 转换为256位的hex字符串 不支持小数
func ToTwosComplement(input interface{}) (string, error){
	bigInt, err := ToBN(input)
	if err != nil{
		return "", ErrNumberInput
	}
	if bigInt.Sign() == 1 {
		nStr := fmt.Sprintf("%064x", bigInt)
		// 如果超过64位的话，是否截断？？？
		return "0x"+nStr, nil
	}else if bigInt.Sign() == -1 {
		return "0x"+fmt.Sprintf("%x", math.U256(bigInt)), nil
	}else {
		return "0x"+PadLeft("0", 64), nil
	}
}