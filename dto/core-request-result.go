package dto

import (
	"encoding/json"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
)

type CoreRequestResult struct {
	RequestResult
}

type SyncingResponse struct {
	StartingBlock *big.Int `json:"startingBlock"`
	CurrentBlock  *big.Int `json:"currentBlock"`
	HighestBlock  *big.Int `json:"highestBlock"`
}

type AccountResult struct {
	Address      utils.Address   `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     utils.Hash      `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  utils.Hash      `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

func (pointer *CoreRequestResult) ToSyncingResponse() (*SyncingResponse, error) {

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
		return nil, UNPARSEABLEINTERFACE
	}

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	syncingResponse := &SyncingResponse{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, syncingResponse)

	return syncingResponse, err

}

func (pointer *CoreRequestResult) ToTransactionResponse() (*TransactionResponse, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	var data []byte
	var err error
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

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}
	}

	transactionResponse := &TransactionResponse{}

	err = json.Unmarshal(data, transactionResponse)

	return transactionResponse, err

}

func (pointer *CoreRequestResult) ToSignTransactionResponse() (*SignTransactionResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var data []byte
	var err error
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
		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}
	}

	signTransactionResponse := &SignTransactionResponse{}

	err = json.Unmarshal(data, signTransactionResponse)

	return signTransactionResponse, err
}

func (pointer *CoreRequestResult) ToBlock(transactionDetails bool) (interface{}, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var data []byte
	var err error
	if value, ok := pointer.Result.([]byte); ok {
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

	var block interface{}
	if transactionDetails {
		block = &BlockDetails{}
		err = json.Unmarshal(data, block)
	} else {
		block = &BlockNoDetails{}
		err = json.Unmarshal(data, block)
	}
	return block, err
}

func (pointer *CoreRequestResult) ToPendingTransactions() ([]*TransactionResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	results := (pointer).Result.([]interface{})

	pendingTransactions := make([]*TransactionResponse, len(results))

	for i, v := range results {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}

		info := &TransactionResponse{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, info)
		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		pendingTransactions[i] = info

	}

	return pendingTransactions, nil

}

func (pointer *CoreRequestResult) ToProof() (*AccountResult, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var data []byte
	var err error
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

	accountResult := &AccountResult{}

	err = json.Unmarshal(data, accountResult)

	return accountResult, err

}
