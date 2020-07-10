package dpos

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

func TestGetValidators(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	blockNumber := big.NewInt(157532)
	Dpos := connection.System.NewDpos()
	validators, err := Dpos.GetValidators(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(validators)
}
