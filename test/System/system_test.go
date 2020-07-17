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
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	log, err := connection.System.SystemLogDecode("0x153e77acfbf7b29dea4f4739c89df860e955a00b46021de7939672b38b5c1430")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(log.Method)
	t.Log(log.Status)
	t.Log(log.Result)
}
