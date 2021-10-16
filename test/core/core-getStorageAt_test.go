package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"strconv"
	"testing"
)

func TestCoreGetStorageAt(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generator, _ := connection.Core.GetGenerator()

	value, err := connection.Core.GetStorageAt(generator, big.NewInt(0), block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(value)

}
