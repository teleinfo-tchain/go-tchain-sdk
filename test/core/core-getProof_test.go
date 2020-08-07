package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestCoreGetProof(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	coinBase, _ := connection.Core.GetCoinBase()

	proof, err := connection.Core.GetProof(coinBase, []string{"0", "1"}, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", proof)

}
