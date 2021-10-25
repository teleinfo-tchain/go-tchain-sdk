package test

import (
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func signAndSingTx(con *bif.Bif, receipt string, sender string, signPriKey string, isSM2 bool) (string, error) {
	nonce, err := con.Core.GetTransactionCount(sender, block.LATEST)
	if err != nil {
		return "", errors.New("failed connection")
	}

	chainId, err := con.Core.GetChainId()
	if err != nil {
		return "", errors.New("failed getChainId")
	}

	var recipient utils.Address
	recipient = utils.StringToAddress(receipt)
	tx := &account.SignTxParams{
		Recipient: &recipient,
		Nonce:     nonce,
		GasPrice:  big.NewInt(0),
		GasLimit:  21000,
		Amount:    big.NewInt(500000000),
		Payload:   nil,
		ChainId:   chainId,
	}
	res, err := account.SignTransaction(tx, signPriKey, isSM2)

	if err != nil {
		return "", err
	}

	txHash, err := con.Core.SendRawTransaction(res.Raw.String())
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func TestGetRawTransactionByBlockHashAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipient := "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"

	txID, err := signAndSingTx(connection, recipient, sender, priKey, false)

	t.Log("Tx Submitted: ", txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  wait for a block
	time.Sleep(time.Second * 8)

	txFromHash, err := connection.Core.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// // if it fails, it may be that the time is too short and the transaction has not been executed
	tx, err := connection.Core.GetRawTransactionByBlockHashAndIndex(txFromHash.BlockHash, txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx)
	// test removing the 0x

	tx, err = connection.Core.GetRawTransactionByBlockHashAndIndex(txFromHash.BlockHash[2:], txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx)
}

func TestGetRawTransactionByBlockNumberAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipient := "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"

	txID, err := signAndSingTx(connection, recipient, sender, priKey, false)

	t.Log("Tx Submitted: ", txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  wait for a block
	time.Sleep(time.Second * 8)
	// time.Sleep(time.Second)

	txFromHash, err := connection.Core.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// if it fails, it may be that the time is too short and the transaction has not been executed
	tx, err := connection.Core.GetRawTransactionByBlockNumberAndIndex(hexutil.EncodeBig(txFromHash.BlockNumber), txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx)
}

func TestGetRawTransactionByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipient := "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"

	txID, err := signAndSingTx(connection, recipient, sender, priKey, false)
	tx, err := connection.Core.GetRawTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx)
}

func TestGetTransactionByBlockHashAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipient := "did:bid:qwer:zftAgNtnQzLMGJHKPMdn9quPvuikNWUZ"

	txID, err := signAndSingTx(connection, recipient, sender, priKey, false)

	t.Log("Tx Submitted: ", txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var txFromHash *dto.TransactionResponse
	for {
		time.Sleep(time.Second * 3)
		txFromHash, err = connection.Core.GetTransactionByHash(txID)
		if txFromHash != nil && fmt.Sprintf("%d", txFromHash.BlockNumber) != "0" {
			break
		}
	}

	tx, err := connection.Core.GetTransactionByBlockHashAndIndex(txFromHash.BlockHash, txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if tx.Sender != sender || tx.Recipient != recipient || tx.Hash != txID {
		t.Errorf("Incorrect transaction from hash and index")
		t.FailNow()
	}

	// test removing the 0x
	tx, err = connection.Core.GetTransactionByBlockHashAndIndex(txFromHash.BlockHash[2:], txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if tx.Sender != sender || tx.Recipient != recipient || tx.Hash != txID {
		t.Errorf("Incorrect transaction from hash and index")
		t.FailNow()
	}
}

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipient := "did:bid:qwer:zftAgNtnQzLMGJHKPMdn9quPvuikNWUZ"

	txID, err := signAndSingTx(connection, recipient, sender, priKey, false)

	t.Log("Tx Submitted: ", txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  如果交易没有执行，则用已有的交易hash测试
	// txID := "0x1abaf67025a43fd5d3ecc19d3b67003ee1caf863aa642fca5248c153ca7ea5fc"
	var txFromHash *dto.TransactionResponse
	for {
		time.Sleep(time.Second * 2)
		txFromHash, err = connection.Core.GetTransactionByHash(txID)
		if txFromHash != nil && fmt.Sprintf("%d", txFromHash.BlockNumber) != "0" {
			break
		}
	}
	tx, err := connection.Core.GetTransactionByBlockNumberAndIndex(hexutil.EncodeBig(txFromHash.BlockNumber), txFromHash.TransactionIndex)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if tx.Sender != sender || tx.Recipient != recipient || tx.Hash != txID {
		t.Errorf("Incorrect transaction from hash and index")
		t.FailNow()
	}
}

func TestGetTransactionByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipient := "did:bid:qwer:zftAgNtnQzLMGJHKPMdn9quPvuikNWUZ"

	txID, err := signAndSingTx(connection, recipient, sender, priKey, false)

	t.Log("txID:", txID)

	if err != nil {
		t.Errorf("Failed SendTransaction")
		t.Error(err)
		t.FailNow()
	}

	// Wait for a block
	time.Sleep(time.Second)

	tx, err := connection.Core.GetTransactionByHash(txID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(tx)
}

func TestCoreGetTransactionCount(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generator, _ := connection.Core.GetGenerator()

	count, err := connection.Core.GetTransactionCount(generator, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	countTwo, err := connection.Core.GetTransactionCount(generator, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// count should not change
	if count.Cmp(countTwo) != 0 {
		t.Errorf("Count incorrect, changed between calls")
		t.FailNow()
	}
}

func TestCoreGetTransactionReceipt(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipient := "did:bid:qwer:zftAgNtnQzLMGJHKPMdn9quPvuikNWUZ"

	txID, err := signAndSingTx(connection, recipient, sender, priKey, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var receipt *dto.TransactionReceipt
	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txID)
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(receipt.ContractAddress) == 0 {
		t.Log("No contract address")
	}

	if len(receipt.TransactionHash) == 0 {
		t.Error("No transaction hash")
		t.FailNow()
	}

	if receipt.TransactionIndex == nil {
		t.Error("No transaction index")
		t.FailNow()
	}

	if len(receipt.BlockHash) == 0 {
		t.Error("No block hash")
		t.FailNow()
	}

	if receipt.BlockNumber == nil || receipt.BlockNumber.Cmp(big.NewInt(0)) == 0 {
		t.Error("No block number")
		t.FailNow()
	}

	if receipt.Logs == nil || len(receipt.Logs) == 0 {
		t.Log("No logs")
	}

	if !receipt.Status {
		t.Error("False status")
		t.FailNow()
	}
}

// 测试发送RawTransaction
func TestCoreSendRawTransaction(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"

	nonce, err := connection.Core.GetTransactionCount(sender, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipientStr := "did:bid:qwer:zftAgNtnQzLMGJHKPMdn9quPvuikNWUZ"

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var recipient utils.Address
	recipient = utils.StringToAddress(recipientStr)

	tx := &account.SignTxParams{
		Recipient: &recipient,
		Nonce:     nonce,
		GasPrice:  big.NewInt(2000000),
		GasLimit:  uint64(41000),
		Amount:    big.NewInt(50000000000),
		Payload:   nil,
		ChainId:   chainId,
	}

	res, err := account.SignTransaction(tx, priKey, false)
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

// todo: 测试失败
func TestCoreSignTransaction(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	sender := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"

	nonce, err := connection.Core.GetTransactionCount(sender, block.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.ChainId = chainId
	transaction.Sender = sender
	transaction.Recipient = "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"
	transaction.AccountNonce = nonce.Uint64()
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(5), big.NewInt(1e17))
	transaction.GasLimit = uint64(50000)
	transaction.GasPrice = big.NewInt(1)
	transaction.Payload = "Sign Transfer bif test"

	txID, err := connection.Core.SignTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(txID.Raw)
	// t.Logf("%#v \n", txID.Transaction)

	toAddress := utils.StringToAddress(transaction.Recipient)
	if txID.Transaction.To != hexutil.Encode(toAddress[:]) {
		t.Errorf(fmt.Sprintf("Expected %s | Got: %s", transaction.Recipient, txID.Transaction.To))
		t.FailNow()
	}

	if txID.Transaction.Value.Cmp(transaction.Amount) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.Amount.Uint64(), txID.Transaction.Value.Uint64()))
		t.FailNow()
	}

	// if txID.Transaction.GasLimit.Cmp(transaction.GasLimit) != 0 {
	//	t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.GasLimit.Uint64(), txID.Transaction.GasLimit.Uint64()))
	//	t.FailNow()
	// }
	if txID.Transaction.GasPrice.Cmp(transaction.GasPrice) != 0 {
		t.Errorf(fmt.Sprintf("Expected %d | Got: %d", transaction.GasPrice.Uint64(), txID.Transaction.GasPrice.Uint64()))
		t.FailNow()
	}
}
