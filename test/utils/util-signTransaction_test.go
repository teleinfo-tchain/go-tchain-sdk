package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
	"testing"
)

// 使用国密对交易进行签名转账
func TestCoreSignTransactionSm2(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	nonce, err := connection.Core.GetTransactionCount(resources.AddressSM2, block.LATEST)
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
	privateKey := resources.AddressPriKey
	var cryType uint = 0
	sender, err := utils.GetAddressFromPrivate(privateKey, cryType)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from := common.StringToAddress(sender)
	if from != common.StringToAddress(resources.AddressSM2){
		t.Errorf("Address and private key do not match")
		t.FailNow()
	}

	to := common.StringToAddress(resources.AddressTwo)
	tx := &utils.Txdata{
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
	}

	res, _ := utils.SignTransaction(tx, privateKey, int64(chainId))

	txIDRaw, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txIDRaw)
}

// 使用非国密对交易进行签名转账
func TestCoreSignTransactionNoSm2(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	nonce, err := connection.Core.GetTransactionCount(resources.CoinBase, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Errorf("Failed getChainId")
		t.FailNow()
	}

	privateKey := resources.CoinBasePriKey
	var cryType uint = 1
	sender, err := utils.GetAddressFromPrivate(privateKey, cryType)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from := common.StringToAddress(sender)
	if from != common.StringToAddress(resources.CoinBase){
		t.Errorf("Address and private key do not match")
		t.FailNow()
	}

	to := common.StringToAddress(resources.AddressSM2)
	tx := &utils.Txdata{
		AccountNonce: nonce.Uint64(),
		Price:     big.NewInt(35),
		GasLimit:  2000000,
		Sender:    &from,
		Recipient: &to,
		Amount:    big.NewInt(50000000000),
		Payload:   nil,
		V:         new(big.Int),
		R:         new(big.Int),
		S:         new(big.Int),
		T:         big.NewInt(0),
	}

	res, _ := utils.SignTransaction(tx, privateKey, int64(chainId))

	txIDRaw, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(txIDRaw)
}
