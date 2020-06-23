package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"testing"
)

func TestCoreSendRawTransaction(t *testing.T){
	var connection = bif.NewBif(providers.NewHTTPProvider("172.20.3.21:44032", 10, false))
	tx := "0xf87b010a831e8480946469643a6269643a73f6a70d05af2141dd4ad995946469643a6269643a73890cf407f6c883e9a427358307a120808082039ea05572544e70662d7bd0b8be692db8c455ebd3ca42aef28ed5e09e43061c2d4a82a01a9a8e218a3ddfc1ed081018fa9cb78c84cffb8101265141987f4e5a0a215f48"
	txID, err := connection.Core.SendRawTransaction(tx)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txID)
}
