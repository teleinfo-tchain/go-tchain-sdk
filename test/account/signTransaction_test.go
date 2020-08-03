package account

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


func TestCoreSignTransaction(t *testing.T){
	for _, test := range []struct{
		id string
		address string
		privateKey string
		to string
		cryType uint

	}{
		{"国密签名转账", resources.AddressSM2, resources.AddressPriKey, resources.CoinBase,0},
		{"非国密签名转账", resources.CoinBase, resources.CoinBasePriKey, resources.AddressSM2,1},
	}{
		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
		nonce, err := connection.Core.GetTransactionCount(test.address, block.LATEST)
		if err != nil {
			t.Errorf("Failed connection")
			t.FailNow()
		}

		chainId, err := connection.Core.GetChainId()
		if err != nil {
			t.Errorf("Failed getChainId")
			t.FailNow()
		}

		// fmt.Println("nonce is ",nonce.Uint64())
		sender, err := account.GetAddressFromPrivate(test.privateKey, test.cryType)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		from := utils.StringToAddress(sender)
		if from != utils.StringToAddress(test.address) {
			t.Errorf("Address and private key do not match")
			t.FailNow()
		}

		to := utils.StringToAddress(test.to)
		tx := &account.TxData{
			AccountNonce: nonce.Uint64(),
			Price:        big.NewInt(30),
			GasLimit:     2000000,
			Sender:       &from,
			Recipient:    &to,
			Amount:       big.NewInt(50000000000),
			Payload:      nil,
			V:            new(big.Int),
			R:            new(big.Int),
			S:            new(big.Int),
			T:            big.NewInt(0),
			// NT:           new(big.Int),
			// NV:           new(big.Int),
			// NR:           new(big.Int),
			// NS:           new(big.Int),
		}

		acc := account.NewAccount()
		res, _ := acc.SignTransaction(tx, test.privateKey, big.NewInt(0).SetUint64(chainId))

		transactionHash, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))

		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("id: %s, transactionHash is %s \n", test.id, transactionHash)
	}
}
