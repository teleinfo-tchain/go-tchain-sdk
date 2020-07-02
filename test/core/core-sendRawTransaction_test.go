package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/core"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

//test/core/core-sendRawTransaction_test.go
func TestCoreSendRawTransaction(t *testing.T){
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	nonce, err := connection.Core.GetTransactionCount(resources.CoinBase, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	privKey := resources.CoinBasePriKey
	from := common.StringToAddress(resources.CoinBase)
	to := common.StringToAddress(resources.AddressTwo)
	tx := &core.Txdata{
		AccountNonce: nonce.Uint64(),
		//AccountNonce: 10,
		Price: big.NewInt(25),
		GasLimit: 2000000,
		Sender: &from,
		Recipient: &to,
		Amount: big.NewInt(50000000000),
		Payload: nil,
		V: new(big.Int),
		R: new(big.Int),
		S: new(big.Int),
		T: big.NewInt(0),
	}
	res, _ := core.SignTransaction(tx, privKey, 666)

	txIDRaw, err := connection.Core.SendRawTransaction(common.ToHex(res.Raw))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txIDRaw)

}
