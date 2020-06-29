package dpos

import (
	"fmt"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)


func TestGetValidators(t *testing.T){
	var connection = system.NewDpos(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	validators, err := connection.GetValidators(block.LATEST)

	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(validators)
}

