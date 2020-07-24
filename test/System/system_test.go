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

	// 0xc8ca132b2483b5f95644bc1472c55c607fa93873ee5072d8882f8c07f069200f
	log, err := connection.System.SystemLogDecode("0x210873a36d45f76a88504b063c7a9342b9d75d1702e5945526bec9a78c96e726")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(log.Method)
	t.Log(log.Status)
	t.Log(log.Result)
}
