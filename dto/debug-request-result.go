package dto

import (
	"encoding/json"
	customerror "github.com/bif/bif-sdk-go/constants"
)

type DebugRequestResult struct {
	RequestResult
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

type DumpAccount struct {
	Balance  string            `json:"balance"`
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"code_hash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

func (pointer *DebugRequestResult) ToDumpBlock() (*Dump, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	dump := &Dump{}

	marshal, err := json.Marshal(result)

	err = json.Unmarshal(marshal, dump)

	return dump, err
}
