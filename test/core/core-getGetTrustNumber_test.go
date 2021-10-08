package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"strconv"
	"testing"
)

func TestCoreGetTrustNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	address, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	trustNumber, err := connection.Core.GetTrustNumber(address)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("trustNumber is ", trustNumber)
}
