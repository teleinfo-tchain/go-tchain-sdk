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

/**
 * @file request-result.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package dto

import (
	"errors"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"strconv"
	"strings"

	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/constants"

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

func (pointer *RequestResult) ToPeerInfo() ([]*common.PeerInfo, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result_list := (pointer).Result.([]interface{})

	new := make([]*common.PeerInfo, len(result_list))

	for i, v := range result_list {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, customerror.EMPTYRESPONSE
		}

		info := &common.PeerInfo{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, customerror.UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal([]byte(marshal), info)

		new[i] = info

	}

	return new, nil

}

func (pointer *RequestResult) ToStringArray() ([]string, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	new := make([]string, len(result))
	for i, v := range result {
		new[i] = v.(string)
	}

	return new, nil

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

func (pointer *RequestResult) ToInt() (int64, error) {

	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	hex := result.(string)

	numericResult, err := strconv.ParseInt(hex, 16, 64)

	return numericResult, err

}

func (pointer *RequestResult) ToBigInt() (*big.Int, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	res := (pointer).Result.(interface{})

	ret, success := big.NewInt(0).SetString(res.(string)[2:], 16)

	if !success {
		return nil, errors.New(fmt.Sprintf("Failed to convert %s to BigInt", res.(string)))
	}

	return ret, nil
}

func (pointer *RequestResult) ToComplexIntResponse() (types.ComplexIntResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return types.ComplexIntResponse(0), err
	}

	result := (pointer).Result.(interface{})

	var hex string

	switch v := result.(type) {
	//Testrpc returns a float64
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

func (pointer *RequestResult) ToSignTransactionResponse() (*SignTransactionResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	signTransactionResponse := &SignTransactionResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), signTransactionResponse)

	return signTransactionResponse, err
}

func (pointer *RequestResult) ToTransactionResponse() (*TransactionResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	transactionResponse := &TransactionResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), transactionResponse)

	return transactionResponse, err

}

func (pointer *RequestResult) ToTransactionReceipt() (*TransactionReceipt, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	transactionReceipt := &TransactionReceipt{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), transactionReceipt)

	return transactionReceipt, err

}

func (pointer *RequestResult) ToTxpoolInspect() (map[string]map[string]map[string]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	inspect := make(map[string]map[string]map[string]string, len(result))

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), &inspect)

	return inspect, err
}

func (pointer *RequestResult) ToTxpoolStatus() (map[string]hexutil.Uint, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	status := make(map[string]hexutil.Uint, len(result))

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), &status)

	return status, err
}

func (pointer *RequestResult) ToBlock(transactionDetails bool) (interface{}, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	var block interface{}
	if transactionDetails {
		block = &BlockDetails{}
		err = json.Unmarshal([]byte(marshal), block)
	} else {
		block = &BlockNoDetails{}
		err = json.Unmarshal([]byte(marshal), block)
	}
	return block, err

}

func (pointer *RequestResult) ToSyncingResponse() (*SyncingResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var result map[string]interface{}

	switch (pointer).Result.(type) {
	case bool:
		return &SyncingResponse{}, nil
	case map[string]interface{}:
		result = (pointer).Result.(map[string]interface{})
	default:
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	syncingResponse := &SyncingResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	json.Unmarshal([]byte(marshal), syncingResponse)

	return syncingResponse, nil

}

func (pointer *RequestResult) ToCandidatesResponse() ([]*CandidateResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	resultList := (pointer).Result.([]interface{})

	new := make([]*CandidateResponse, len(resultList))

	for i, v := range resultList {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, customerror.EMPTYRESPONSE
		}

		var candidate CandidateResponse

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, customerror.UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal([]byte(marshal), &candidate)

		new[i] = &candidate
	}

	return new, nil

}

func (pointer *RequestResult) ToVoterResponse() (*VoterResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	data := &VoterResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), data)

	return data, err

}

func (pointer *RequestResult) ToStakeResponse() (*StakeResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	data := &StakeResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), data)

	return data, err

}

func (pointer *RequestResult) ToPeerCertificateResponse() (*PeerCertificate, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	data := &PeerCertificate{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), data)

	return data, err

}

// To avoid a conversion of a nil interface
func (pointer *RequestResult) checkResponse() error {

	if pointer.Error != nil {
		return errors.New(pointer.Error.Message)
	}

	if pointer.Result == nil {
		return customerror.EMPTYRESPONSE
	}

	return nil

}


func (pointer *RequestResult) ToValidators() ([]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	validators := make([]string, len(result))
	for i, v := range result {
		validators[i] = v.(string)
	}

	return validators, nil
}

func (pointer *RequestResult) ToValidatorsAtHash()([]string, error){
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	validators := make([]string, len(result))
	for i, v := range result {
		validators[i] = v.(string)
	}

	return validators, nil
}

func(pointer *RequestResult) ToRoundStateInfo() (*RoundStateInfo, error){
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	roundStateInfo := &RoundStateInfo{}

	marshal, err := json.Marshal(result)

	err = json.Unmarshal([]byte(marshal), roundStateInfo)

	return roundStateInfo, err
}


func(pointer *RequestResult) ToRoundChangeSetInfo() (*RoundChangeSetInfo, error){
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	roundChangeSetInfo := &RoundChangeSetInfo{}

	marshal, err := json.Marshal(result)

	err = json.Unmarshal([]byte(marshal), roundChangeSetInfo)

	return roundChangeSetInfo, err
}


func(pointer *RequestResult) ToBacklogs() (map[string][]*Message, error){
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	backlogs := make(map[string][]*Message, len(result))

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal([]byte(marshal), &backlogs)

	return backlogs, err
}