package utils

import (
	"encoding/hex"
	"encoding/json"
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

// define some const
const bifer string = "bifer"
const Sha3Null = "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
// 现在的测试不能确保正确，因为没有一个100%确定正确的做对照
const Sm3Null = "0x1ab21d8355cfa17f8e61194831e81a8f22bec8c728fefb747ed035eb5082aa2b"

// Various big integer limit values.
var (
	tt255     = math.BigPow(2, 255)
	MaxBig255   = new(big.Int).Sub(tt255, big.NewInt(1))
	MinBig255   = new(big.Int).Neg(tt255)
)

// Errors
var (
	ErrInvalidSha3    = errors.New("invalid input, input is null")
	ErrInvalidSm3    = errors.New("invalid input, input is null")
	ErrInvalidBid = errors.New("invalid Bif Bid")
	ErrNumberString = errors.New("invalid number string")
	ErrNumberInput = errors.New("invalid number input")
	ErrBigInt = errors.New("big int not in -2*255——2*255-1")
	ErrUintNoExist   = errors.New("uint not exist")
	ErrParameter     = errors.New("error parameter number")
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
var biferUint = map[string]string {
	"nobifer":    "0",
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

// get uint value
func getUintValue(uint string) (string, error){
	uint = strings.ToLower(uint)
	if biferUint[uint]==""{
		return "", ErrUintNoExist
	}
	return biferUint[uint],nil
}

// Converts any bifer value value into wei
func (util *Utils) ToWei(number string, values ...string) (string, error){
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

// Converts any wei value into a bifer value.
func (util *Utils) FromWei(number string, values ...string) (string, error){
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
	// Dividend
	dividend := new(big.Float)
	m, ok := dividend.SetString(number)
	if !ok {
		return "", ErrNumberString
	}
	// divisor
	divisor := new(big.Float)
	n, _ := divisor.SetString(uintValue)
	resValue := new(big.Float)
	return resValue.Quo(m,n).String(), nil
}

func (util *Utils) Sm3(str string) (string, error){
	var hexBytes []byte
	if util.IsHexStrict(str){
		hexBytes = common.Hex2Bytes(str[2:])
	}else{
		hexBytes = []byte(str)
	}
	resStr := util.ByteToHex(crypto.Keccak256(crypto.SM2, hexBytes))
	if resStr == Sm3Null {
		return "", ErrInvalidSm3
	}else{
		return resStr, nil
	}
}

// calculate the sha3 of the input (if input is invalid, it will return error)
func (util *Utils) Sha3(str string) (string,error){
	var hexBytes []byte
	if util.IsHexStrict(str){
		hexBytes, _ = hex.DecodeString(str[2:])
	}else{
		hexBytes = []byte(str)
	}
	resStr := util.ByteToHex(crypto.Keccak256(crypto.SECP256K1, hexBytes))
	if resStr == Sha3Null {
		return "", ErrInvalidSha3
	}else{
		return resStr, nil
	}
}

//calculate the sha3 of the input(if input is invalid, it will return Sha3Null)
func (util *Utils) Sha3Raw(str string) string{
	var hexBytes []byte
	if util.IsHexStrict(str){
		hexBytes, _ = hex.DecodeString(str[2:])
	}else{
		hexBytes = []byte(str)
	}
	resStr := util.ByteToHex(crypto.Keccak256(1, hexBytes))
	return resStr
}

// bid string with "0x", Checks the checksum of a given Bid
func (util *Utils) CheckBidChecksum(bid string) bool{
	bid = bid[2:]
	bidHash, _ := util.Sha3(strings.ToLower(bid))
	bidHash = bidHash[2:]
	for i := 0; i < 40; i++ {
		// 校验和的判断依据是根据其生成的方式判断的
		char := string(bid[i])
		hashChar:= string(bidHash[i])
		upperChar := strings.ToUpper(char)
		lowerChar := strings.ToLower(char)
		parseHashChar, _ := strconv.ParseUint(hashChar, 16, 64)
		if parseHashChar>7 && upperChar != char || parseHashChar <= 7 && lowerChar != char {
			return false
		}
	}
	return true
}

//Checks if a given string is a valid Bif bid. It will also check the checksum, if the bid has upper and lowercase letters.
func (util *Utils) IsBid(bid string) bool{
	res1, _ := regexp.Compile("^(0[x,X])?[A-F, a-f0-9]{40}$")
	res2, _ := regexp.Compile("^(0[x,X])?[a-f, 0-9]{40}$")
	res3, _ := regexp.Compile("^(0[x,X])?[A-F, 0-9]{40}$")
	if !res1.MatchString(bid){
		return false
	}else if res2.MatchString(bid) || res3.MatchString(bid) {
		return true
	}else {
		return util.CheckBidChecksum(bid)
	}
}

//  convert an upper or lowercase Bif bid to a checksum bid
func (util *Utils) ToChecksumBid(bid string) (string, error){
	res, _ := regexp.Compile("^(0[x,X])?[A-F, a-f0-9]{40}$")
	if !res.MatchString(bid){
		return "", ErrInvalidBid
	}
	bid = strings.ToLower(bid)[2:]
	bidHash , _:= util.Sha3(bid)
	bidHash = bidHash[2:]
	checkSumBid := "0x"
	for i := 0; i < len(bid); i++ {
		hashChar:= string(bidHash[i])
		parseHashChar, _ := strconv.ParseUint(hashChar, 16, 64)
		if parseHashChar >7 {
			checkSumBid += strings.ToUpper(string(bid[i]))
		}else{
			checkSumBid += string(bid[i])
		}
	}
	return checkSumBid, nil
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
func (util *Utils) ByteToHex(byteArr []byte) string{
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
func (util *Utils) NumberToHex(input interface{}) (string, error){
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
func (util *Utils) ToHex(input interface{}) (string, error){
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
func (util *Utils) ToBN(input interface{}) (*big.Int, error){
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
func (util *Utils) IsBN(input interface{}) bool{
	_, ok := input.(*big.Int)
	if ok{
		return true
	}
	return false
}

// Adds a padding on the left of a string, Useful for adding paddings to HEX strings.
func (util *Utils) PadLeft(str string, characterAmount int, signs ...string) string {
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
func (util *Utils) PadRight(str string, characterAmount int, signs ...string) string {
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

// 不支持小数，不支持超过256位表示的int补码转换（即从-2*255到2*255-1）
func (util *Utils) ToTwosComplement(input interface{}) (string, error){
	bigInt, err := util.ToBN(input)
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
		return "0x"+ util.PadLeft("0", 64), nil
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

	if len(args) != len(funcTypeArr){
		return "", errors.New("the number of function parameters does not match")
	}

	for index := 0; index < len(funcTypeArr); index++ {
		tmpBytes, err := getHexValue(funcTypeArr[index], args[index])

		if err != nil {
			return "", err
		}
		//fmt.Println("tmpBytes ", tmpBytes)
		byteCode += tmpBytes
	}

	//fmt.Println("byteCode is ", byteCode)
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

	if len(args) != len(funcTypeArr){
		return "", errors.New("the number of function parameters does not match")
	}

	fullFunction := fmt.Sprintf("%s(%s)", functionName, strings.Join(funcTypeArr, ","))
	sha3Function, err := NewUtils().Sha3(fullFunction)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for index := 0; index < len(funcTypeArr); index++ {
		currentData, err := getHexValue(funcTypeArr[index], args[index])

		if err != nil {
			return "", err
		}
		builder.WriteString(currentData)
	}

	byteCode := fmt.Sprintf("%s%s", sha3Function[0:10], builder.String())
	return byteCode, nil

}


func getFunctions(abi string) (map[string][]string,error){
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

func getHexValue(inputType string, value interface{}) (string, error) {

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
		builder.WriteString(fmt.Sprintf("%064s", fmt.Sprintf("%x", 32)))
		builder.WriteString(fmt.Sprintf("%064s", fmt.Sprintf("%x", value.(string))))
	}

	return builder.String(), nil

}
