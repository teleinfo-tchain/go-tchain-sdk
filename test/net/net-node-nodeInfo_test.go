package test

import (
	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestGetNodeInfo(t *testing.T) {

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP51+":"+resources.Port, 10, false))

	nodeInfo, err := connection.GetNodeInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(nodeInfo)
}
