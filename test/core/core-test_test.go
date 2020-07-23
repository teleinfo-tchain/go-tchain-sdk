package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func Test1(t *testing.T){
	hash := "0x216a911507110236292ac35c96951ed53ca21495a95c08ded4523a242acd67c9"
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	res,err := connection.Core.GetTransactionReceipt(hash)
	fmt.Println(res,err)
	// fmt.Println(res.BlockNumber)
}

