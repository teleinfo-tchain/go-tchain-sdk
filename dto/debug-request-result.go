package dto

import (
	"encoding/json"
	"fmt"
	customerror "github.com/bif/bif-sdk-go/constants"
)

type DebugRequestResult struct {
	RequestResult
}

type Dump struct {

}

func (pointer *DebugRequestResult) ToDumpBlock() (*Dump, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	fmt.Printf("dump is %#v \n", result)
	dump := &Dump{}

	marshal, err := json.Marshal(result)

	err = json.Unmarshal(marshal, dump)

	return dump, err
}