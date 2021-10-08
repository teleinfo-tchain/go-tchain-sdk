package test

import (
	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"strconv"
	"testing"
)

func TestGetNodeInfo(t *testing.T) {

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	nodeInfo, err := connection.GetNodeInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n ", nodeInfo)
}
