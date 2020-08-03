package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"testing"
)

// 测试发送RawTransaction
func TestCoreSendRawTransaction(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	nonce, err := connection.Core.GetTransactionCount(resources.CoinBase, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	privKey := resources.CoinBasePriKey
	from := utils.StringToAddress(resources.CoinBase)
	to := utils.StringToAddress(resources.AddressTwo)
	tx := &account.TxData{
		AccountNonce: nonce.Uint64(),
		Price:        big.NewInt(25),
		GasLimit:     2000000,
		Sender:       &from,
		Recipient:    &to,
		Amount:       big.NewInt(50000000000),
		Payload:      nil,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
		T:            new(big.Int),
		// NT:           new(big.Int),
		// NV:           new(big.Int),
		// NR:           new(big.Int),
		// NS:           new(big.Int),
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	acc := account.NewAccount()
	res, _ := acc.SignTransaction(tx, privKey, big.NewInt(0).SetUint64(chainId))
	txIDRaw, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txIDRaw)

}

