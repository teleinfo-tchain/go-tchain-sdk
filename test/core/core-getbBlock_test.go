/********************************************************************************
   This file is part of go-bif.
   go-bif is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-bif is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-bif.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package test

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	block2 "github.com/bif/bif-sdk-go/core/block"
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

func TestCoreBlockNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	blockNumber, err := connection.Core.GetBlockNumber()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if blockNumber.Int64() < 0 {
		t.Errorf("Invalid Block Number")
		t.Fail()
	}

	t.Log(blockNumber)
}

func TestCoreGetBlockByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	const transactionDetails = false
	blockByNumber, err := connection.Core.GetBlockByNumber(block2.LATEST, transactionDetails)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	blockByHash, err := connection.Core.GetBlockByHash(blockByNumber.(*dto.BlockNoDetails).Hash, transactionDetails)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Ensure it's the same block
	if (blockByNumber.(*dto.BlockNoDetails).Number.Cmp(blockByHash.(*dto.BlockNoDetails).Number)) != 0 ||
		(blockByNumber.(*dto.BlockNoDetails).Hash != blockByHash.(*dto.BlockNoDetails).Hash) {
		t.Errorf("Not same block returned")
		t.FailNow()
	}

	t.Log(blockByHash.(*dto.BlockNoDetails).Hash, blockByNumber.(*dto.BlockNoDetails).Hash)

	_, err = connection.Core.GetBlockByHash("0x1234", transactionDetails)

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}

	_, err = connection.Core.GetBlockByHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa0", false)

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}

	_, err = connection.Core.GetBlockByHash("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa0", false)

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}

	blockByHash, err = connection.Core.GetBlockByHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false)

	if err == nil {
		t.Errorf("Found a block with incorrect hash?")
		t.FailNow()
	}
}

func TestCoreGetBlockByNumber(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	for _, test := range []struct {
		transactionDetails bool
	}{
		{true},
		{false},
	} {
		_, err := connection.Core.GetBlockNumber()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		block, err := connection.Core.GetBlockByNumber(block2.LATEST, test.transactionDetails)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if block == nil {
			t.Error("Block returned is nil")
			t.FailNow()
		}

		if test.transactionDetails {
			t.Log("this is detail ")
			t.Log("Number:", block.(*dto.BlockDetails).Number)
			t.Log("Hash:", block.(*dto.BlockDetails).Hash)
			t.Log("ParentHash:", block.(*dto.BlockDetails).ParentHash)
			t.Log("LogsBloom:", block.(*dto.BlockDetails).LogsBloom)
			t.Log("StateRoot:", block.(*dto.BlockDetails).StateRoot)
			t.Log("Generator:", block.(*dto.BlockDetails).Generator)
			t.Log("Regulatory:", block.(*dto.BlockDetails).Regulatory)
			t.Log("ExtraData:", block.(*dto.BlockDetails).ExtraData)
			t.Log("Size:", block.(*dto.BlockDetails).Size)
			t.Log("Timestamp:", block.(*dto.BlockDetails).Timestamp)
			t.Log("TransactionsRoot:", block.(*dto.BlockDetails).TransactionsRoot)
			t.Log("ReceiptsRoot:", block.(*dto.BlockDetails).ReceiptsRoot)
			for _, val := range block.(*dto.BlockDetails).Transactions {
				t.Logf("transactions: %#v \n", val)
			}
		} else {
			t.Log("this is no detail ")
			t.Log("Number:", block.(*dto.BlockNoDetails).Number)
			t.Log("Hash:", block.(*dto.BlockNoDetails).Hash)
			t.Log("ParentHash:", block.(*dto.BlockNoDetails).ParentHash)
			t.Log("LogsBloom:", block.(*dto.BlockNoDetails).LogsBloom)
			t.Log("StateRoot:", block.(*dto.BlockNoDetails).StateRoot)
			t.Log("Generator:", block.(*dto.BlockNoDetails).Generator)
			t.Log("Regulatory:", block.(*dto.BlockNoDetails).Regulatory)
			t.Log("ExtraData:", block.(*dto.BlockNoDetails).ExtraData)
			t.Log("Size:", block.(*dto.BlockNoDetails).Size)
			t.Log("Timestamp:", block.(*dto.BlockNoDetails).Timestamp)
			t.Log("TransactionsRoot:", block.(*dto.BlockNoDetails).TransactionsRoot)
			t.Log("ReceiptsRoot:", block.(*dto.BlockNoDetails).ReceiptsRoot)
			for idx, val := range block.(*dto.BlockNoDetails).Transactions {
				t.Log("transactions:[0]", val[idx])
			}
		}
	}

}

func TestGetBlockTransactionCountByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	block, err := connection.Core.GetBlockByNumber(block2.LATEST, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	txCount, err := connection.Core.GetBlockTransactionCountByHash(block.(*dto.BlockNoDetails).Hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txCount:", txCount)

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block2.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(resources.Addr1)
	var recipient utils.Address
	recipient = utils.StringToAddress(resources.Addr2)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  200000,
		Sender:    &sender,
		Recipient: &recipient,
		Amount:    big.NewInt(10000),
	}

	res, err := account.SignTransaction(tx, resources.Addr1Pri, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txHah, err := connection.Core.SendRawTransaction(res.Raw.String())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("txHash is %s \n", txHah)

	time.Sleep(time.Second*8)

	txCount, err = connection.Core.GetBlockTransactionCountByHash(block.(*dto.BlockNoDetails).Hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//  注意如果有其他用户发送交易，可能会大于1，该测试可能会失败
	if txCount != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}
	t.Log("txCount:", txCount)
}

func TestGetBlockTransactionCountByNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))


	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block2.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(resources.Addr1)
	var recipient utils.Address
	recipient = utils.StringToAddress(resources.Addr2)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  200000,
		Sender:    &sender,
		Recipient: &recipient,
		Amount:    big.NewInt(10000),
	}

	res, err := account.SignTransaction(tx, resources.Addr1Pri, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txHah, err := connection.Core.SendRawTransaction(res.Raw.String())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(time.Second*8)

	txRes, _ := connection.Core.GetTransactionByHash(txHah)

	t.Logf("txRes blockNumber %+v \n", txRes)

	blockNumber := hexutil.EncodeBig(txRes.BlockNumber)

	txCount, err := connection.Core.GetBlockTransactionCountByNumber(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}

	txCount, err = connection.Core.GetBlockTransactionCountByNumber(block2.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}
}