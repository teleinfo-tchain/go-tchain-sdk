package utils

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/utils/math"
	"github.com/teleinfo-bif/bit-gmsm/sm3"
	"golang.org/x/crypto/sha3"
	"math/big"
	"regexp"
	"strings"
)

// define some const
const bifer string = "bifer"
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
	ErrInvalidSha3  = errors.New("invalid input, input is null")
	ErrInvalidSm3   = errors.New("invalid input, input is null")
	ErrBigInt       = errors.New("big int not in -2*255——2*255-1")
	ErrUintNoExist  = errors.New("uint not exist")
)

// Utils - The Utils Module
type Utils struct {
	BiferUint map[string]string
}

// NewUtils - Utils Module constructor to set the default provider
func NewUtils() *Utils {
	utils := new(Utils)
	utils.BiferUint = biferUint
	return utils
}

// bifer uints
var biferUint = map[string]string{
	"nobifer":    "0",
	"wei":        "1",
	"kwei":       "1000",
	"babbage":    "1000",
	"femtobifer": "1000",
	"mwei":       "1000000",
	"lovelace":   "1000000",
	"picobifer":  "1000000",
	"gwei":       "1000000000",
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

/*
  ToWei:
   	EN - Converts any bif value value into wei
 	CN - 将任何bif值转换为wei
  Params:
  	- balance,string 要转换的bif金额
  	- uint（可选）,  string  转换的单位，默认为bifer

  Returns:
   	- string，转换为以wei为单位的bif余额
   	- error

  Call permissions: Anyone

  Debug: 转换时，将浮点数转换为int，（同时判断，是否符合bif的单位转化，这点需要添加）
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
		number = biferUint[bifer]
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
		if len(balSplit[1]) > len(number)-1{
			return nil, errors.New(fmt.Sprintf("while converting number %s to wei,  too many decimal points not met", balance))
		}
		for len(balSplit[1]) < len(number)-1{
			balSplit[1] += "0"
		}
		fraction, ok := new(big.Int).SetString(balSplit[1], 10)
		if !ok {
			return nil, errors.New("trans fail")
		}
		whole.Add(whole.Mul(whole, value), fraction)
		if isNeg{
			return whole.Neg(whole), nil
		}
		return whole, nil
	} else {
		if isNeg{
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
  Debug: 精度问题
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
		number = biferUint[bifer]
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
  CheckBidChecksum:
   	EN - bid string with "0x", Checks the checksum of a given Bid
	CN - bid字符串(必须带有0x前缀，否则一直返回false)，检查给定bid的校验和
  Params:
  	- bid string

  Returns:
  	- bool,如果地址的校验和有效则为true，否则为false

  Call permissions: Anyone
*/
func (util *Utils) CheckBidChecksum(bid string) bool {
	unCheckBid := HexToAddress(bid).Hex()
	return bid == unCheckBid
}

/*
  IsBid:
   	EN - Checks if a given string is a valid Bif bid. It will also check the checksum, if the bid has upper and lowercase letters.
	CN - 检查给定的字符串是否为有效的Bid。 如果Bid包含大小写字母，它还将检查校验和。(如果全是hex表示，则必须带有0x前缀)
  Params:
  	- bid string

  Returns:
  	- bool,如果为有效的bid则返回true，否则为false

  Call permissions: Anyone
  Debug 这里地址判单没有使用common中提供的，因为其没有对其进行校验和检查
*/
func (util *Utils) IsBid(bid string) bool {
	if Has0xPrefix(bid) {
		bid = bid[2:]
		return longBidCheck(bid)
	} else if HasDidBidPrefix(bid) {
		bid = bid[8:]
		return shortBidCheck(bid)
	}
	return false
}

func longBidCheck(bid string) bool {
	res1, _ := regexp.Compile("^[A-F, a-f0-9]{40}$")
	res2, _ := regexp.Compile("^[a-f, 0-9]{40}$")
	res3, _ := regexp.Compile("^[A-F, 0-9]{40}$")
	if !res1.MatchString(bid) {
		return false
	} else if res2.MatchString(bid) || res3.MatchString(bid) {
		return true
	} else {
		return bid == HexToAddress(bid).Hex()
	}
}

func shortBidCheck(bid string) bool {
	res1, _ := regexp.Compile("^[A-F, a-f0-9]{24}$")
	res2, _ := regexp.Compile("^[a-f, 0-9]{24}$")
	res3, _ := regexp.Compile("^[A-F, 0-9]{24}$")
	if !res1.MatchString(bid) {
		return false
	} else if res2.MatchString(bid) || res3.MatchString(bid) {
		return true
	} else {
		bid = "0x6469643A6269643A" + bid
		return bid == HexToAddress(bid).Hex()
	}
}

/*
  ToChecksumBid:
   	EN - convert bid to a checksum bid
	CN - 将bid转换为校验和bid
  Params:
  	- bid string

  Returns:
  	- string 校验和bid

  Call permissions: Anyone
  Debug: 是否还判断是否前缀为did:bid:？？？
*/
func (util *Utils) ToChecksumBid(bid string) string {
	return HexToAddress(bid).Hex()
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

func (util *Utils) ByteCodeDeploy(abi string, byteCode string, args ...interface{}) (string, error) {
	functions, err := getFunctions(abi)
	if err != nil {
		return "", err
	}

	funcTypeArr, ok := functions["constructor"]
	if !ok {
		return "", errors.New("not exist constructor")
	}

	if len(args) != len(funcTypeArr) {
		return "", errors.New("the number of function parameters does not match")
	}

	for index := 0; index < len(funcTypeArr); index++ {
		tmpBytes, err := GetHexValue(funcTypeArr[index], args[index])

		if err != nil {
			return "", err
		}
		// fmt.Println("tmpBytes ", tmpBytes)
		byteCode += tmpBytes
	}

	// fmt.Println("byteCode is ", byteCode)
	return byteCode, nil

}

func (util *Utils) ByteCodeInteract(abi string, functionName string, args ...interface{}) (string, error) {
	functions, err := getFunctions(abi)
	if err != nil {
		return "", err
	}

	funcTypeArr, ok := functions[functionName]
	if !ok {
		return "", errors.New("not exist functionName")
	}

	if len(args) != len(funcTypeArr) {
		return "", errors.New("the number of function parameters does not match")
	}

	fullFunction := fmt.Sprintf("%s(%s)", functionName, strings.Join(funcTypeArr, ","))
	sha3Function, err := NewUtils().Sha3(fullFunction)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for index := 0; index < len(funcTypeArr); index++ {
		currentData, err := GetHexValue(funcTypeArr[index], args[index])

		if err != nil {
			return "", err
		}
		builder.WriteString(currentData)
	}

	byteCode := fmt.Sprintf("%s%s", sha3Function[0:10], builder.String())
	return byteCode, nil

}

func getFunctions(abi string) (map[string][]string, error) {
	var mockInterface interface{}
	err := json.Unmarshal([]byte(abi), &mockInterface)

	if err != nil {
		return nil, err
	}

	jsonInterface := mockInterface.([]interface{})
	functions := make(map[string][]string)
	for index := 0; index < len(jsonInterface); index++ {
		function := jsonInterface[index].(map[string]interface{})

		if function["type"] == "constructor" || function["type"] == "fallback" {
			function["name"] = function["type"]
		}

		functionName := function["name"].(string)
		functions[functionName] = make([]string, 0)

		if function["inputs"] == nil {
			continue
		}

		inputs := function["inputs"].([]interface{})
		for paramIndex := 0; paramIndex < len(inputs); paramIndex++ {
			params := inputs[paramIndex].(map[string]interface{})
			functions[functionName] = append(functions[functionName], params["type"].(string))
		}

	}
	return functions, nil
}

// 编码方式需要修改，这种不可取
func GetHexValue(inputType string, value interface{}) (string, error) {

	var builder strings.Builder
	if strings.HasPrefix(inputType, "int") ||
		strings.HasPrefix(inputType, "uint") ||
		strings.HasPrefix(inputType, "fixed") ||
		strings.HasPrefix(inputType, "ufixed") {

		bigVal := value.(*big.Int)

		// Checking that the string actually is the correct inputType
		if strings.Contains(inputType, "128") {
			// 128 bit
			if bigVal.BitLen() > 128 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		} else if strings.Contains(inputType, "256") {
			// 256 bit
			if bigVal.BitLen() > 256 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		}

		builder.WriteString(fmt.Sprintf("%064s", fmt.Sprintf("%x", bigVal)))
	}

	if strings.Compare("address", inputType) == 0 {
		builder.WriteString(fmt.Sprintf("%064d", len(value.(string)[:])))
		builder.WriteString(fmt.Sprintf("%064d", len(value.(string)[:])))
		builder.WriteString(fmt.Sprintf("%064s", value.(string)[:]))
	}

	if strings.Compare("string", inputType) == 0 {
		builder.WriteString(fmt.Sprintf("%064s", fmt.Sprintf("%x", 32)))
		builder.WriteString(fmt.Sprintf("%064s", fmt.Sprintf("%x", len(value.(string)))))
		builder.WriteString(fmt.Sprintf("%064s", fmt.Sprintf("%x", value.(string))))
	}

	return builder.String(), nil

}
