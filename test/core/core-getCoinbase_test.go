package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestGetCoinbase(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	coinbase, err := connection.Core.GetCoinbase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(coinbase)

	_, err = connection.Core.GetBalance(coinbase, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}

