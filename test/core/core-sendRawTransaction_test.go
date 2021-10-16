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
	"strconv"
	"testing"
)

// 测试发送RawTransaction
func TestCoreSendRawTransaction(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	nonce, err := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	privKey := "bd0bc93d5bfd6f329f3b1fb61109fdb81290dfb9a78de491aa84276af4a713a2"

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var recipient utils.Address
	recipient = utils.StringToAddress(resources.RegisterAllianceOne)

	tx := &account.SignTxParams{
		Recipient:    &recipient,
		AccountNonce: nonce.Uint64(),
		GasPrice:     big.NewInt(2000000),
		GasLimit:     uint64(41000),
		Amount:       big.NewInt(50000000000),
		Payload:      nil,
		ChainId:      chainId,
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

