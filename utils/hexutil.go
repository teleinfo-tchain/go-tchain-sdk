package utils

import (
	"encoding/hex"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"math/big"
	"math/rand"
	"regexp"
	"time"
)

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// judge if input is hex string with/without prefix 0[x,X]
func (util *Utils) IsHex(input string) bool {
	r, _ := regexp.Compile("^(0[x,X])?[A-F, a-f, 0-9]+$")
	return r.MatchString(input)
}

// judge if input is hex strict string with prefix 0[x,X]
func (util *Utils) IsHexStrict(input string) bool {
	r, _ := regexp.Compile("^(0[x,X])[A-F, a-f, 0-9]+$")
	return r.MatchString(input)
}

// generate cryptographically strong pseudo-random HEX strings from a given byte size
func (util *Utils) RandomHex(size int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size*2; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return "0x" + string(result)
}

// convert hex string to byte
func (util *Utils) HexToBytes(input string) ([]byte, error) {
	return hexutil.Decode(input)
}

// convert hex string to utf8 string
func (util *Utils) HexToUtf8(input string) (string, error) {
	res, err := hexutil.Decode(input)
	if err != nil {
		return "", err
	} else {
		return string(res), nil
	}
}

// convert hex string to ascii string
func (util *Utils) HexToAscii(input string) (string, error) {
	res, err := hexutil.Decode(input)
	if err != nil {
		return "", err
	} else {
		return string(res), nil
	}
}

// convert hex string to number string
func (util *Utils) HexToNumberString(input string) (string, error) {
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
// convert hex string to uint64
func (util *Utils) HexToUint64Number(input string) (uint64, error) {
	return hexutil.DecodeUint64(input)
}

// convert hex string to big.Int
func (util *Utils) HexToBigNumber(input string) (*big.Int, error) {
	return hexutil.DecodeBig(input)
}

//convert ascii string to hex string
func (util *Utils) AsciiToHex(input string) string {
	return "0x" + hex.EncodeToString([]byte(input))
}

//convert utf8 string to hex string
func (util *Utils) Utf8ToHex(input string) string {
	return "0x" + hex.EncodeToString([]byte(input))
}
