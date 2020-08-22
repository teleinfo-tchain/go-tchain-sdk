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

// 测试发送RawTransaction
func TestCoreSendRawTransaction(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddr, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	privKey := resources.TestAddrPri

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tx := &account.SignTxParams{
		Version: 1,
		To:       resources.NewAddrE,
		Nonce:    nonce.Uint64(),
		Gas:      2000000,
		GasPrice: big.NewInt(30),
		Value:    big.NewInt(50000000000),
		Data:     nil,
		ChainId:  chainId,
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

