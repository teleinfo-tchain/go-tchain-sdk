/********************************************************************************
   This file is part of go-bif.
   go-bif is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-bif is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-bif.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package dto

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/tchain/go-tchain-sdk/utils/types"

	"encoding/json"
	"fmt"
	"math/big"
)

type RequestResult struct {
	ID      int         `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error,omitempty"`
	Data    string      `json:"data,omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (pointer *RequestResult) ToStringArray() ([]string, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var data0 []byte
	var err0 error
	if value, ok := pointer.Result.(string); ok {
		//come from websock
		data0 = []byte(value)
	} else {
		//come from http
		resultSingle := pointer.Result.([]interface{})

		if len(resultSingle) == 0 {
			return nil, EMPTYRESPONSE
		}
		data0, err0 = json.Marshal(resultSingle)
	}

	if err0 != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	var resultLi []interface{}

	err0 = json.Unmarshal(data0, &resultLi)

	stringArray := make([]string, len(resultLi))
	for i, v := range resultLi {
		var data []byte
		var err error
		if value, ok := v.(string); ok {
			//come from websock
			data = []byte(value)
		} else {
			//come from http
			result := v.(map[string]interface{})

			if len(result) == 0 {
				return nil, EMPTYRESPONSE
			}
			data, err = json.Marshal(result)
		}

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		var stringInfo string

		err = json.Unmarshal(data, &stringInfo)

		stringArray[i] = stringInfo
	}

	return stringArray, nil

}

func (pointer *RequestResult) ToComplexString() (types.ComplexString, error) {

	if err := pointer.checkResponse(); err != nil {
		return "", err
	}

	result := (pointer).Result.(interface{})

	return types.ComplexString(result.(string)), nil

}

func (pointer *RequestResult) ToString() (string, error) {

	if err := pointer.checkResponse(); err != nil {
		return "", err
	}

	result := (pointer).Result.(interface{})

	return result.(string), nil

}

func (pointer *RequestResult) ToInt64() (int64, error) {

	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	hex := result.(string)[2:]

	numericResult, err := strconv.ParseInt(hex, 16, 64)

	return numericResult, err

}

func (pointer *RequestResult) ToUint64() (uint64, error) {

	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})
	hex := result.(string)

	if strings.HasPrefix(hex, "0") {
		hex = hex[2:]
	} else {
		hex = hex[3 : len(hex)-1]
	}

	numericResult, err := strconv.ParseUint(hex, 16, 64)

	return numericResult, err

}

func (pointer *RequestResult) ToBigInt() (*big.Int, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	res := (pointer).Result.(interface{})

	//ret, success := big.NewInt(0).SetString(string(res.([]uint8))[2:], 16)
	res = strings.Replace(res.(string), "\"", "", -1)
	ret, success := big.NewInt(0).SetString(res.(string)[2:], 16)

	if !success {
		return nil, errors.New(fmt.Sprintf("Failed to convert %s to BigInt", res.(string)))
	}

	return ret, nil
}

func (pointer *RequestResult) ToComplexIntResponse() (types.ComplexIntResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return types.ComplexIntResponse(rune(0)), err
	}

	result := (pointer).Result.(interface{})

	var hex string

	switch v := result.(type) {
	// Testrpc returns a float64
	case float64:
		hex = strconv.FormatFloat(v, 'E', 16, 64)
		break
	default:
		hex = result.(string)
	}

	cleaned := strings.TrimPrefix(hex, "0x")

	return types.ComplexIntResponse(cleaned), nil

}

func (pointer *RequestResult) ToBoolean() (bool, error) {

	if err := pointer.checkResponse(); err != nil {
		return false, err
	}

	result := (pointer).Result.(interface{})

	return result.(bool), nil

}

func (pointer *RequestResult) ToTransactionReceipt() (*TransactionReceipt, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	//result := (pointer).Result.(map[string]interface{})

	//if len(result) == 0 {
	//	return nil, EMPTYRESPONSE
	//}

	transactionReceipt := &TransactionReceipt{}

	marshal, err := json.Marshal((pointer).Result)

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, transactionReceipt)

	return transactionReceipt, err

}

// To avoid a conversion of a nil interface
func (pointer *RequestResult) checkResponse() error {

	if pointer.Error != nil {
		return errors.New(pointer.Error.Message)
	}

	if pointer.Result == nil {
		return EMPTYRESPONSE
	}

	if value, ok := pointer.Result.(string); ok {
		if strings.Index(value, "0x") == 0 {
			return nil
		}

		if rvalue, err := base64.StdEncoding.DecodeString(value); err == nil {
			pointer.Result = string(rvalue)
		}
	}

	return nil

}
