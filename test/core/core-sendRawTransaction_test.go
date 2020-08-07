package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"testing"
)

// ≤‚ ‘∑¢ÀÕRawTransaction
func TestCoreSendRawTransaction(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	nonce, err := connection.Core.GetTransactionCount(resources.CoinBase, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	privKey := resources.CoinBasePriKey

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tx := &account.SignTxParams{
		To:       resources.AddressTwo,
		Nonce:    nonce.Uint64(),
		Gas:      2000000,
		GasPrice: big.NewInt(30),
		Value:    big.NewInt(50000000000),
		Data:     nil,
		ChainId:  big.NewInt(0).SetUint64(chainId),
	}

	res, err := connection.Account.SignTransaction(tx, privKey, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txIDRaw, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txIDRaw)

}

