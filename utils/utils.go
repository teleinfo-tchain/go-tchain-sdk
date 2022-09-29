package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tchain/go-tchain-sdk/utils/hexutil"
	"github.com/tchain/go-tchain-sdk/utils/math"
	"github.com/tchain/go-tgmsm/sm3"
	"golang.org/x/crypto/sha3"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// define some const
const bif string = "bif"
const Sha3Null = "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"

// 现在的测试不能确保正确，因为没有一个100%确定正确的做对照
const Sm3Null = "0x1ab21d8355cfa17f8e61194831e81a8f22bec8c728fefb747ed035eb5082aa2b"

// Various big integer limit values.
var (
	tt255     = math.BigPow(2, 255)
	MaxBig255 = new(big.Int).Sub(tt255, big.NewInt(1))
	MinBig255 = new(big.Int).Neg(tt255)
)

// Errors
var (
	ErrInvalidSha3 = errors.New("invalid input, input is null")
	ErrInvalidSm3  = errors.New("invalid input, input is null")
	ErrBigInt      = errors.New("big int not in -2*255——2*255-1")
)

// bifUint
var bifUint = map[string]string{
	"nobif":    "0",
	"wei":      "1",
	"kwei":     "1000",
	"Kwei":     "1000",
	"mwei":     "1000000",
	"Mwei":     "1000000",
	"gwei":     "1000000000",
	"Gwei":     "1000000000",
	"microbif": "1000000000000",
	"micro":    "1000000000000",
	"millibif": "1000000000000000",
	"milli":    "1000000000000000",
	"bif":      "1000000000000000000",
}

/*
  ToWei:
   	EN - Converts any bif value value into wei
 	CN - 将任何bif值转换为wei
  Params:
  	- balance,string 要转换的bif金额
  	- uint（可选）,  string  转换的单位，默认为bifer

  Returns:
   	- *big.Int，转换为以wei为单位的bif余额
   	- error

  Call permissions: Anyone
*/
func ToWei(balance string, uint ...string) (*big.Int, error) {
	var number string
	if len(uint) >= 1 {
		value, ok := bifUint[strings.ToLower(uint[0])]
		if !ok {
			return nil, errors.New("uint not exist")
		}
		number = value
	} else {
		number = bifUint[bif]
	}
	value, _ := new(big.Int).SetString(number, 10)

	// 判断输入的正负
	var isNeg = false
	if balance[0] == '-' {
		isNeg = true
		balance = balance[1:]
	}

	balSplit := strings.Split(balance, ".")
	if len(balSplit) > 2 {
		return nil, errors.New(fmt.Sprintf("while converting number %s to wei,  too many decimal points not met", balance))
	}

	whole, ok := new(big.Int).SetString(balSplit[0], 10)
	if !ok {
		return nil, errors.New("trans fail")
	}

	if len(balSplit) > 1 {
		if len(balSplit[1]) > len(number)-1 {
			return nil, errors.New(fmt.Sprintf("while converting number %s to wei,  too many decimal points not met", balance))
		}
		for len(balSplit[1]) < len(number)-1 {
			balSplit[1] += "0"
		}
		fraction, ok := new(big.Int).SetString(balSplit[1], 10)
		if !ok {
			return nil, errors.New("trans fail")
		}
		whole.Add(whole.Mul(whole, value), fraction)
		if isNeg {
			return whole.Neg(whole), nil
		}
		return whole, nil
	} else {
		if isNeg {
			whole.Mul(whole, value)
			return whole.Neg(whole), nil
		}
		return whole.Mul(whole, value), nil
	}

	// value, err := decimal.NewFromString(number)
	// if err != nil {
	// 	return "", errors.New("trans fail")
	// }
	//
	// bal, err := decimal.NewFromString(balance)
	// if err != nil {
	// 	return "", errors.New("trans fail")
	// }
	//
	// return bal.Mul(value).String(), nil
}

/*
  FromWei:
   	EN - Converts any wei value into a bif value
 	CN - 将任何以wei为单位的数值转换为其他单位的数值
  Params:
  	- balance, *big.Int 要转换的bif值
  	- uint（可选）,  string  转换的单位，默认为bifer

  Returns:
   	- string，转换为以wei为单位的bif余额
   	- error

  Call permissions: Anyone
*/
func FromWei(balance *big.Int, uint ...string) (string, error) {
	var number string
	if len(uint) >= 1 {
		value, ok := bifUint[strings.ToLower(uint[0])]
		if !ok {
			return "", errors.New("uint not exist")
		}
		number = value
	} else {
		number = bifUint[bif]
	}
	value, _ := new(big.Int).SetString(number, 10)
	// 判断输入的正负
	sign := balance.Sign()
	if sign == -1 {
		balance = balance.Neg(balance)
	}
	fraction := new(big.Int)
	whole := new(big.Int)
	fraction.Mod(balance, value)
	whole.Div(balance, value)
	if sign == -1 {
		whole = whole.Neg(whole)
	}

	// 根据模是否为0判断是否还需拼接
	if fraction.Cmp(big.NewInt(0)) == 0 {
		return whole.String(), nil
	} else {
		res := whole.String() + "." + string(bytes.Repeat([]byte{48}, len(number)-len(fraction.String())-1)) + fraction.String()
		return strings.TrimRight(res, "0"), nil
	}

}

func Sm3(str string) (string, error) {
	var hexBytes []byte
	if IsHexStrict(str) {
		hexBytes = Hex2Bytes(str[2:])
	} else {
		hexBytes = []byte(str)
	}
	resStr := ByteToHex(keccak256(0, hexBytes))
	if resStr == Sm3Null {
		return "", ErrInvalidSm3
	} else {
		return resStr, nil
	}
}

/*
  Sha3:
   	EN - calculate the sha3 of the input (if input is invalid, it will return error)
	CN - 计算输入的sha3（如果输入为空，则将返回错误）
  Params:
  	- input string

  Returns:
  	- string
	- error

  Call permissions: Anyone
*/
func Sha3(input string) (string, error) {
	var hexBytes []byte
	if IsHexStrict(input) {
		hexBytes, _ = hex.DecodeString(input[2:])
	} else {
		hexBytes = []byte(input)
	}
	resStr := ByteToHex(keccak256(1, hexBytes))
	if resStr == Sha3Null {
		return "", ErrInvalidSha3
	} else {
		return resStr, nil
	}
}

/*
  Sha3Raw:
   	EN - calculate the sha3 of the input(if input is invalid, it will return Sha3Null)
	CN - 计算输入的sha3（如果输入为空，则将返回Sha3Null）
  Params:
  	- input string

  Returns:
  	- string

  Call permissions: Anyone
*/
func Sha3Raw(str string) string {
	var hexBytes []byte
	if IsHexStrict(str) {
		hexBytes, _ = hex.DecodeString(str[2:])
	} else {
		hexBytes = []byte(str)
	}
	resStr := ByteToHex(keccak256(1, hexBytes))
	return resStr
}

/*
  ByteToHex:
   	EN - Convert bytes to hexadecimal string
	CN - 将字节数组转换为十六进制字符串
  Params:
  	- byteArr []byte 字节数组

  Returns:
  	- string hex 字符串
	- error

  Call permissions: Anyone
*/
func ByteToHex(byteArr []byte) string {
	return "0x" + hex.EncodeToString(byteArr)
}

/*
  ToTwosComplement:
   	EN -
	CN - 不支持小数，不支持超过256位表示的int补码转换（即从-2*255到2*255-1）
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: Anyone
*/
func ToTwosComplement(input *big.Int) (string, error) {
	if input.Cmp(MaxBig255) == 1 || input.Cmp(MinBig255) == -1 {
		return "", ErrBigInt
	}
	if input.Sign() == 1 {
		nStr := fmt.Sprintf("%064x", input)
		// 如果超过64位的话, 截断
		return "0x" + nStr, nil
	} else if input.Sign() == -1 {
		return "0x" + fmt.Sprintf("%x", math.U256(input)), nil
	} else {
		// return "0x" + PadLeft("0", 64), nil
		return fmt.Sprintf("%064x", 0), nil
	}
}

/*
  描述:
   	EN - judge if input is hex strict string with prefix 0[x|X]
	CN - 判断输入的字符串是否为十六进制字符串(必须带前缀0[x|X])
  Params:
  	- input string, 输入的字符串

  Returns:
  	- bool， true表示是hex string，false表示不是hex string

  Call permissions: Anyone
*/
func IsHexStrict(input string) bool {
	if Has0xPrefix(input) {
		input = input[2:]
	} else {
		return false
	}
	return IsHex(input)
}

/*
  RandomHex:
   	EN - generate cryptographically strong pseudo-random HEX strings from a given byte size
	CN - 从给定的字节大小生成具有加密强度的伪随机HEX字符串
  Params:
  	- size int 给定字节的大小

  Returns:
  	- string， HEX字符串

  Call permissions: Anyone
*/
func RandomHex(size int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size*2; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

/*
  描述:
   	EN -  convert hex string to byte
	CN - 将十六进制字符串转换为字节
  Params:
  	- input string 输入的hex string，必须带有前缀0[x|X],否则报错

  Returns:
  	- []byte
	- error

  Call permissions: Anyone
*/
func HexToBytes(input string) ([]byte, error) {
	return hexutil.Decode(input)
}

/*
  HexToUint64:
   	EN - convert hex string to uint64
	CN - 将十六进制字符串转换为uint64
  Params:
  	- input string，输入的hex 字符串

  Returns:
  	- uint64， 转换得到的值
	- error

  Call permissions: Anyone
*/
func HexToUint64(input string) (uint64, error) {
	return hexutil.DecodeUint64(input)
}

/*
  描述:
   	EN - convert hex string to *big.Int
	CN - 将十六进制字符串转换为*big.Int
  Params:
	- input string，输入的hex 字符串
  	-

  Returns:
  	- *big.Int， 转换得到的值
	- error

  Call permissions: Anyone
*/
func HexToBigInt(input string) (*big.Int, error) {
	return hexutil.DecodeBig(input)
}

/*
  描述:
   	EN - convert hex string to utf8 string
	CN - 将十六进制字符串转换为utf8字符串
  Params:
  	- input string hex字符串

  Returns:
  	- string UTF-8字符串
	- error

  Call permissions: Anyone
*/
func HexToUtf8(input string) (string, error) {
	res, err := hexutil.Decode(input)
	if err != nil {
		return "", err
	} else {
		return string(res), nil
	}
}

/*
  HexToAscii:
   	EN - convert hex string to ascii string
	CN - 将十六进制字符串转换为ASCII字符串
  Params:
  	- input string hex字符串

  Returns:
  	- string ASCII字符串
	- error

  Call permissions: Anyone
*/
func HexToAscii(input string) (string, error) {
	res, err := hexutil.Decode(input)
	if err != nil {
		return "", err
	} else {
		return string(res), nil
	}
}

/*
  AsciiToHex:
   	EN - convert ASCII string to hex string
	CN - 将ASCII字符串转换为十六进制字符串
  Params:
  	- input string,给定的ASCII字符串

  Returns:
  	- string，hex字符串

  Call permissions: Anyone
*/
func AsciiToHex(input string) string {
	runeInput := []rune(input)
	hexStr := "0x"
	for _, v := range runeInput {
		hexStr += strconv.FormatInt(int64(v), 16)
	}
	return hexStr
}

/*
  Utf8ToHex:
   	EN - convert utf8 string to hex string
	CN - 将UTF-8字符串转换为十六进制字符串
  Params:
  	- input string,给定的UTF-8字符串

  Returns:
  	- string，hex字符串

  Call permissions: Anyone
*/
func Utf8ToHex(input string) string {
	return "0x" + hex.EncodeToString([]byte(input))
}

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(cryptoType uint8, data ...[]byte) []byte {
	switch cryptoType {
	case 0:
		return keccak256Sm2(data...)
	case 1:
		return keccak256Btc(data...)
	default:
		return keccak256Sm2(data...)
	}
}

func keccak256Sm2(data ...[]byte) []byte {
	d := sm3.New()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// keccak256Btc calculates and returns the Keccak256 hash of the input data.
func keccak256Btc(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
