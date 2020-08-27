package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

// 测试系统合约的执行结果
func TestSystemLogDecode(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 0xe1fa5649868d8dba74ebb693e69edae57c446a8f00978653e53fc6da4e52a0c3
	log, err := connection.System.SystemLogDecode("0x9aaf3a93ff2d254f8ca64c578039b4ee64cdd43dc7f16b11f810ed0387645bd0")

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	t.Log("method is ", log.Method)
	t.Log("status is ", log.Status)
	t.Log("result is ", log.Result)
}
