package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

// 测试系统合约的执行结果
func TestSystemLogDecode(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP55+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 0x25b5c30ecb5a089e02c419b71a4d5b9aa50e683f6b7dce6da33b380ba5a10ba6
	log, err := connection.System.SystemLogDecode("0x2def8334a44574c7a0b492ea68201c50ac0d16bd5054fc863a4f0d0de64b791a")

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	t.Log("method is ", log.Method)
	t.Log("status is ", log.Status)
	t.Log("result is ", log.Result)
}
