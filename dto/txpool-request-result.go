package dto

import (
	"encoding/json"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
)

type TxPoolRequestResult struct {
	RequestResult
}

type RPCTransaction struct {
	BlockHash        utils.Hash     `json:"blockHash"`
	BlockNumber      *hexutil.Big   `json:"blockNumber"`
	From             string         `json:"from"`
	Gas              hexutil.Uint64 `json:"gas"`
	GasPrice         *hexutil.Big   `json:"gasPrice"`
	Hash             utils.Hash     `json:"hash"`
	Input            hexutil.Bytes  `json:"input"`
	Nonce            hexutil.Uint64 `json:"nonce"`
	To               string         `json:"to"`
	TransactionIndex hexutil.Uint   `json:"transactionIndex"`
	Value            *hexutil.Big   `json:"value"`
	V                *hexutil.Big   `json:"v"`
	R                *hexutil.Big   `json:"r"`
	S                *hexutil.Big   `json:"s"`
}

func (pointer *TxPoolRequestResult) ToTxPoolStatus() (map[string]hexutil.Uint, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}


	var data []byte
	var err error
	var status map[string]hexutil.Uint
	if value, ok := pointer.Result.(string); ok {
		//come from websock
		data = []byte(value)
		//status = make(map[string]hexutil.Uint, len(value))
	} else {
		//come from http
		result := (pointer).Result.(map[string]interface{})
		//status = make(map[string]hexutil.Uint, len(result))
		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}
		data, err = json.Marshal(result)
	}

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(data, &status)

	return status, err
}

func (pointer *TxPoolRequestResult) ToTxPoolInspect() (map[string]map[string]map[string]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var data []byte
	var err error
	var inspect map[string]map[string]map[string]string
	if value, ok := pointer.Result.(string); ok {
		//come from websock
		data = []byte(value)
	} else {
		//come from http
		result := (pointer).Result.(map[string]interface{})
		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}
		data, err = json.Marshal(result)
	}

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(data, &inspect)

	return inspect, err
}

func (pointer *TxPoolRequestResult) ToTxPoolContent() (map[string]map[string]map[string]*RPCTransaction, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var data []byte
	var err error
	var transactions map[string]map[string]map[string]*RPCTransaction

	if value, ok := pointer.Result.(string); ok {
		//come from websock
		data = []byte(value)
	} else {
		//come from http
		result := (pointer).Result.(map[string]interface{})
		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}
		data, err = json.Marshal(result)
	}

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(data, &transactions)

	return transactions, err
}
