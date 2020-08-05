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
	log, err := connection.System.SystemLogDecode("0x5f4131b6ad51c3444dd1431b1cf263266f60e6a41282a117e23f1a8585d01deb")

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	t.Log("method is ", log.Method)
	t.Log("status is ", log.Status)
	t.Log("result is ", log.Result)
}
