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

// Various big integer limit values.
var (
	tt255     = math.BigPow(2, 255)
	MaxBig255   = new(big.Int).Sub(tt255, big.NewInt(1))
	MinBig255   = new(big.Int).Neg(tt255)
)

// Errors
var (
	ErrInvalidSha3    = errors.New("invalid input")
	ErrInvalidAddress = errors.New("invalid Bif Address")
	ErrNumberString = errors.New("invalid number string")
	ErrNumberInput = errors.New("invalid number input")
	ErrBigInt = errors.New("big int not in -2*255——2*255-1")
)

const Sha3Null = "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"

//  ????
const Sm3Null = "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"

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

// calculate the sha3 of the input (if input is invalid, it will return error)
// 是否还处理bigNumber?
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

//calculate the sha3 of the input(if input is invalid, it will return Sha3Null)
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

// address string with "0x", Checks the checksum of a given address
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

//Checks if a given string is a valid Bif address. It will also check the checksum, if the address has upper and lowercase letters.
func IsAddress(address string) bool{
	res1, _ := regexp.Compile("^(0[x,X])?[A-F, a-f0-9]{40}$")
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

//  convert an upper or lowercase Bif address to a checksum address
func ToChecksumAddress(address string) (string, error){
	res, _ := regexp.Compile("^(0[x,X])?[A-F, a-f0-9]{40}$")
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

// encode int64 to hex string
func EncodeInt64(i int64) string {
	prefix := ""
	if i< 0{
		i = -i
		prefix = "-"
	}
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return prefix+string(strconv.AppendInt(enc, i, 16))
}

// convert byte to hex string
func ByteToHex(byteArr []byte) string{
	return "0x"+hex.EncodeToString(byteArr)
}

// convert number string to hex string
func NumberStrToHex(str string)(string, error){
	n := new(big.Int)
	n, ok := n.SetString(str, 0)
	if !ok{
		return "", ErrNumberString
	}
	return hexutil.EncodeBig(n), nil
}

// convert string(include number string) to hex string
func StringToHex(str string) (string, error){
	n := new(big.Int)
	n, ok := n.SetString(str, 0)
	if !ok{
		return "0x"+hex.EncodeToString([]byte(str)), nil
	}
	return hexutil.EncodeBig(n), nil
}


// convert number string/(u)int(8,16,32,64)/big.Int to hex string
func NumberToHex(input interface{}) (string, error){
	switch input.(type) {
	case string:
		return NumberStrToHex(input.(string))
	case uint:
		return hexutil.EncodeUint64(uint64(input.(uint))), nil
	case uint8:
		return hexutil.EncodeUint64(uint64(input.(uint8))), nil
	case uint16:
		return hexutil.EncodeUint64(uint64(input.(uint16))), nil
	case uint32:
		return hexutil.EncodeUint64(uint64(input.(uint32))), nil
	case uint64:
		return hexutil.EncodeUint64(input.(uint64)), nil
	case int:
		return EncodeInt64(int64(input.(int))), nil
	case int8:
		return EncodeInt64(int64(input.(int))), nil
	case int16:
		return EncodeInt64(int64(input.(int))), nil
	case int32:
		return EncodeInt64(int64(input.(int))), nil
	case int64:
		return EncodeInt64(int64(input.(int))), nil
	case *big.Int:
		return hexutil.EncodeBig(input.(*big.Int)), nil
	default:
		return "", ErrNumberInput
	}
}

// convert string/(u)int(8,16,32,64)/big.Int to hex string
func ToHex(input interface{}) (string, error){
	switch input.(type) {
	case string:
		return StringToHex(input.(string))
	case uint:
		return hexutil.EncodeUint64(uint64(input.(uint))), nil
	case uint8:
		return hexutil.EncodeUint64(uint64(input.(uint8))), nil
	case uint16:
		return hexutil.EncodeUint64(uint64(input.(uint16))), nil
	case uint32:
		return hexutil.EncodeUint64(uint64(input.(uint32))), nil
	case uint64:
		return hexutil.EncodeUint64(input.(uint64)), nil
	case int:
		return EncodeInt64(int64(input.(int))), nil
	case int8:
		return EncodeInt64(int64(input.(int))), nil
	case int16:
		return EncodeInt64(int64(input.(int))), nil
	case int32:
		return EncodeInt64(int64(input.(int))), nil
	case int64:
		return EncodeInt64(int64(input.(int))), nil
	case *big.Int:
		return hexutil.EncodeBig(input.(*big.Int)), nil
	default:
		return "", ErrNumberInput
	}
}

// convert number string to big.Int
func numberStrToBN(str string) (*big.Int,error){
	n := new(big.Int)
	n, ok := n.SetString(str, 0)
	if !ok{
		return nil, ErrNumberString
	}
	return n, nil
}

// convert string/(u)int(8,16,32,64)/big.Int to big.Int
func ToBN(input interface{}) (*big.Int, error){
	switch input.(type) {
	case string:
		return numberStrToBN(input.(string))
	case *big.Int:
		return input.(*big.Int), nil
	case uint:
		return new(big.Int).SetUint64(uint64(input.(uint))), nil
	case uint8:
		return new(big.Int).SetUint64(uint64(input.(uint8))), nil
	case uint16:
		return new(big.Int).SetUint64(uint64(input.(uint16))), nil
	case uint32:
		return new(big.Int).SetUint64(uint64(input.(uint32))), nil
	case uint64:
		return new(big.Int).SetUint64(input.(uint64)), nil
	case int:
		return new(big.Int).SetInt64(int64(input.(int))), nil
	case int8:
		return new(big.Int).SetInt64(int64(input.(int))), nil
	case int16:
		return new(big.Int).SetInt64(int64(input.(int))), nil
	case int32:
		return new(big.Int).SetInt64(int64(input.(int))), nil
	case int64:
		return new(big.Int).SetInt64(int64(input.(int))), nil
	default:
		return nil, ErrNumberInput
	}
}

// judge if input is a big.Int
func IsBN(input interface{}) bool{
	_, ok := input.(*big.Int)
	if ok{
		return true
	}
	return false
}

// 不支持小数，不支持超过256位表示的int补码转换（即从-2*255到2*255-1）
func ToTwosComplement(input interface{}) (string, error){
	bigInt, err := ToBN(input)
	if err != nil{
		return "", ErrNumberInput
	}
	if bigInt.Cmp(MaxBig255)==1 || bigInt.Cmp(MinBig255)==-1{
		return "", ErrBigInt
	}
	if bigInt.Sign() == 1 {
		nStr := fmt.Sprintf("%064x", bigInt)
		// 如果超过64位的话, 截断
		return "0x"+nStr, nil
	}else if bigInt.Sign() == -1 {
		return "0x"+fmt.Sprintf("%x", math.U256(bigInt)), nil
	}else {
		return "0x"+PadLeft("0", 64), nil
	}
}