package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/crypto"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// Errors
var (
	ErrInvalidSha3    = &decError{"invalid input"}
	ErrInvalidAddress = &decError{"invalid Bif Address"}
	ErrNumberString = &decError{"invalid number string"}
	ErrNumberInput = &decError{"invalid number input"}
)
type decError struct{ msg string }

func (err decError) Error() string { return err.msg }

const Sha3Null = "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
// Sm的null是什么？？？
const Sm2Null = "?????"

func ByteToHex(byteArr []byte) string{
	return "0x"+hex.EncodeToString(byteArr)
}

func IsBN(){
}
//
func IsBigNumber() bool{
	return false
}

// 应该是和Sha3类似吧？只是加密的方式不同
func Sm2(str string) (string, error){
	var hexBytes []byte
	if hexutil.IsHexStrict(str){
		hexBytes = common.Hex2Bytes(str[2:])
	}else{
		hexBytes = []byte(str)
	}
	resStr := ByteToHex(crypto.Keccak256(crypto.SM2, hexBytes))
	if resStr == Sm2Null {
		return "", ErrInvalidSha3
	}else{
		return resStr, nil
	}
}

//func Sha3Big(str string) (string,error){
//	var hexBytes []byte
//	if hexutil.IsHexStrict(str){
//		hexBytes = common.Hex2Bytes(str[2:])
//	}else{
//		hexBytes = []byte(str)
//	}
//	resStr := ByteToHex(crypto.Keccak256(crypto.SECP256K1, hexBytes))
//	if resStr == Sha3Null {
//		return "", ErrInvalidSha3
//	}else{
//		return resStr, nil
//	}
//}

// 如果是big Number的话
func Sha3(str string) (string,error){
	var hexBytes []byte
	if hexutil.IsHexStrict(str){
		hexBytes = common.Hex2Bytes(str[2:])
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
	if hexutil.IsHexStrict(str){
		hexBytes = common.Hex2Bytes(str[2:])
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
	//fmt.Println(address)
	//fmt.Println(strings.ToLower(address))
	//fmt.Println(addressHash)
	for i := 0; i < 40; i++ {
		//the nth letter should be uppercase if the nth digit of casemap is 1
		// 关于校验和的判断依据为什么是这个？？
		char := string(address[i])
		hashChar:= string(addressHash[i])
		upperChar := strings.ToUpper(char)
		lowerChar := strings.ToLower(char)
		//replace(/^0x/i,'')
		parseHashChar, _ := strconv.ParseInt(hashChar, 16, 64)
		if parseHashChar>7 && upperChar != char || parseHashChar <= 7 && lowerChar != char {
			fmt.Println(i, char, hashChar,upperChar, lowerChar, parseHashChar)
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
	checkSumAddress := "0x"
	fmt.Println(addressHash)
	for i := 0; i < len(address); i++ {
		hashChar:= string(addressHash[i])
		parseHashChar, _ := strconv.ParseInt(hashChar, 16, 64)
		if parseHashChar >7 {
			checkSumAddress += strings.ToUpper(string(address[i]))
		}else{
			checkSumAddress += string(address[i])
		}
	}
	fmt.Println(checkSumAddress)
	return checkSumAddress, nil
}

// helper.go 也有一个但是是接收大整整型的，后面去合并
func Int64ToHex(number int64) string{
	n := new(big.Int)
	n = n.SetInt64(number)
	return fmt.Sprintf("0x%x", n)
}

func BigIntToHex(number *big.Int) string{
	return fmt.Sprintf("0x%x", number)
}

// 只是10进制和16进制的str
func StrToHex(str string) (string, error){
	n := new(big.Int)
	if hexutil.IsHexStrict(str){
		n, ok := n.SetString(str[2:], 16)
		if !ok{
			return "", ErrNumberString
		}
		return fmt.Sprintf("0x%x", n), nil
	}
	n, ok := n.SetString(str, 10)
	if !ok{
		return "", ErrNumberString
	}
	return fmt.Sprintf("0x%x", n), nil
}

func NumberToHex(input interface{}) (string, error){
	switch input.(type) {
	case string:
		return StrToHex(input.(string))
	case int64:
		return Int64ToHex(input.(int64)), nil
	case int:
		return Int64ToHex(int64(input.(int))), nil
	case *big.Int:
		return BigIntToHex(input.(*big.Int)), nil
	default:
		return "", ErrNumberInput
	}
}

//func ToHex() (string, error){
//
//}
