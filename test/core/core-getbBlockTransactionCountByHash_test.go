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
	"math/big"
	"strconv"
	"testing"
	"time"
)

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

	time.Sleep(time.Second*6)

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
