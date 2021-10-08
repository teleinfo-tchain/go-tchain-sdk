package test

import (
	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"strconv"
	"testing"
)

func TestGetDataDir(t *testing.T) {

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	dataDir, err := connection.GetDataDir()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(dataDir)
}