package utils

import (
	"encoding/hex"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

/*
  IsHex:
   	EN - judge if input is hex string with/without prefix 0[x|X]
	CN - 判断输入的字符串是否为十六进制字符串(可以带/不带前缀0[x|X])
  Params:
  	- input string, 输入的字符串

  Returns:
  	- bool， true表示是hex string，false表示不是hex string

  Call permissions: Anyone
*/
func (util *Utils) IsHex(input string) bool {
	if Has0xPrefix(input) {
		input = input[2:]
	}
	return IsHex(input)
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
func (util *Utils) IsHexStrict(input string) bool {
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
func (util *Utils) RandomHex(size int) string {
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
func (util *Utils) HexToBytes(input string) ([]byte, error) {
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
func (util *Utils) HexToUint64(input string) (uint64, error) {
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
func (util *Utils) HexToBigInt(input string) (*big.Int, error) {
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
func (util *Utils) HexToUtf8(input string) (string, error) {
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
func (util *Utils) HexToAscii(input string) (string, error) {
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
func (util *Utils) AsciiToHex(input string) string {
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
func (util *Utils) Utf8ToHex(input string) string {
	return "0x" + hex.EncodeToString([]byte(input))
}
