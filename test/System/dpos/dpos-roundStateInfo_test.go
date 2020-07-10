package dpos

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestRoundStateInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	Dpos := connection.System.NewDpos()
	roundStateInfo, err := Dpos.RoundStateInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("roundStateInfo is %#v \n", roundStateInfo)
}
