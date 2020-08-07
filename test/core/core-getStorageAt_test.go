package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

func TestCoreGetStorageAt(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	coinBase, _ := connection.Core.GetCoinBase()

	value, err := connection.Core.GetStorageAt(coinBase, big.NewInt(0), block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(value)

}
