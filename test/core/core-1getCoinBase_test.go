package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestGetCoinBase(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	coinBase, err := connection.Core.GetCoinBase()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(coinBase)

	_, err = connection.Core.GetBalance(coinBase, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}

