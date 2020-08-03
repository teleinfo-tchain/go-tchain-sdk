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

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	status := make(map[string]hexutil.Uint, len(result))

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, &status)

	return status, err
}

func (pointer *TxPoolRequestResult) ToTxPoolInspect() (map[string]map[string]map[string]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	inspect := make(map[string]map[string]map[string]string, len(result))

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, &inspect)

	return inspect, err
}

func (pointer *TxPoolRequestResult) ToTxPoolContent() (map[string]map[string]map[string]*RPCTransaction, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	transactions := make(map[string]map[string]map[string]*RPCTransaction, len(result))

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, &transactions)

	return transactions, err
}
