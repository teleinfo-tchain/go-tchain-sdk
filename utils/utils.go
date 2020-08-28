package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/utils/math"
	"github.com/teleinfo-bif/bit-gmsm/sm3"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
	"sync"
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
	ErrUintNoExist = errors.New("uint not exist")
)

// Utils - The Utils Module
type Utils struct {
	BiferUint map[string]string
}

var utils *Utils
var once sync.Once

// NewUtils - Utils Module constructor to set the default provider
func NewUtils() *Utils {
	once.Do(func() {
		utils = new(Utils)
		utils.BiferUint = biferUint

	})
	return utils
}

// bifer uints
var biferUint = map[string]string{
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
func (util *Utils) ToWei(balance string, uint ...string) (*big.Int, error) {
	var number string
	if len(uint) >= 1 {
		value, ok := biferUint[strings.ToLower(uint[0])]
		if !ok {
			return nil, errors.New("uint not exist")
		}
		number = value
	} else {
		number = biferUint[bif]
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
func (util *Utils) FromWei(balance *big.Int, uint ...string) (string, error) {
	var number string
	if len(uint) >= 1 {
		value, ok := biferUint[strings.ToLower(uint[0])]
		if !ok {
			return "", errors.New("uint not exist")
		}
		number = value
	} else {
		number = biferUint[bif]
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

/*
  Sm3:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: Anyone
*/
func (util *Utils) Sm3(str string) (string, error) {
	var hexBytes []byte
	if util.IsHexStrict(str) {
		hexBytes = Hex2Bytes(str[2:])
	} else {
		hexBytes = []byte(str)
	}
	resStr := util.ByteToHex(keccak256(0, hexBytes))
	if resStr == Sm3Null {
		return "", ErrInvalidSm3
	} else {
		return resStr, nil
	}
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
func (util *Utils) Sha3(input string) (string, error) {
	var hexBytes []byte
	if util.IsHexStrict(input) {
		hexBytes, _ = hex.DecodeString(input[2:])
	} else {
		hexBytes = []byte(input)
	}
	resStr := util.ByteToHex(keccak256(1, hexBytes))
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
func (util *Utils) Sha3Raw(str string) string {
	var hexBytes []byte
	if util.IsHexStrict(str) {
		hexBytes, _ = hex.DecodeString(str[2:])
	} else {
		hexBytes = []byte(str)
	}
	resStr := util.ByteToHex(keccak256(1, hexBytes))
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
func (util *Utils) ByteToHex(byteArr []byte) string {
	return "0x" + hex.EncodeToString(byteArr)
}

/*
  LeftPadBytes:
   	EN - zero-pads slice to the left up to length l.
	CN - 在字节数组的左侧填充0至长度l
  Params:
  	- slice []byte 要填充的字节数组
  	- l int 填充的长度（如果长度小于字节数组，则返回原字节数组）

  Returns:
  	- []byte 填充的字节数组

  Call permissions: Anyone
*/
func (util *Utils) LeftPadBytes(slice []byte, l int) []byte {
	return LeftPadBytes(slice, l)
}

/*
  RightPadBytes:
   	EN - zero-pads slice to the right up to length l.
	CN - 在字节数组的右侧填充0至长度l
  Params:
  	- slice []byte 要填充的字节数组
  	- l int 填充的长度（如果长度小于字节数组，则返回原字节数组）

  Returns:
  	- []byte 填充的字节数组

  Call permissions: Anyone
*/
func (util *Utils) RightPadBytes(slice []byte, l int) []byte {
	return RightPadBytes(slice, l)
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
func (util *Utils) ToTwosComplement(input *big.Int) (string, error) {
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
		// return "0x" + util.PadLeft("0", 64), nil
		return fmt.Sprintf("%064x", 0), nil
	}
}
