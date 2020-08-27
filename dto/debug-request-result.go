package dto

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type DebugRequestResult struct {
	RequestResult
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"wallet"`
}

type DumpAccount struct {
	Balance  string            `json:"balance"`
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"codeHash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

func (pointer *DebugRequestResult) ToDumpBlock() (*Dump, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}
	fmt.Printf("%#v \n", result)
	for k,v := range result{
		fmt.Printf("%s ,%s, %#v \n", k,reflect.TypeOf(v), v)
	}
	dump := &Dump{}

	marshal, err := json.Marshal(result)

	err = json.Unmarshal(marshal, dump)

	return dump, err
}
