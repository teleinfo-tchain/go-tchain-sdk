package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"testing"
)

func TestCoreSendRawTransaction(t *testing.T){
	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:44002", 10, false))
	tx := "0xf87b010a831e8480946469643a6269643a73f6a70d05af2141dd4ad995946469643a6269643a73890cf407f6c883e9a427358307a120808082039da0fd5f07051af57b7e30ddafb354979f03a0ad9b0213d59c5441fd1d6daadfe321a0ea1d20b57d381d7fec266c184caf4333090c7fd6a878918b1f259aafef3ccdd4"
	txID, err := connection.Core.SendRawTransaction(tx)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txID)
}
