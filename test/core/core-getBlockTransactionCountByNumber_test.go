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
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"strconv"
	"testing"
)

func TestGetBlockTransactionCountByNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))


	// chainId, _ := connection.Core.GetChainId()
	//
	// nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)
	//
	// var sender utils.Address
	// sender = utils.StringToAddress(resources.Addr1)
	// var recipient utils.Address
	// recipient = utils.StringToAddress(resources.Addr2)
	//
	// tx := &account.SignTxParams{
	// 	ChainId:   chainId,
	// 	Nonce:     nonce,
	// 	GasPrice:  big.NewInt(200),
	// 	GasLimit:  200000,
	// 	Sender:    &sender,
	// 	Recipient: &recipient,
	// 	Amount:    big.NewInt(10000),
	// }
	//
	// res, err := account.SignTransaction(tx, resources.Addr1Pri, false)
	//
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	//
	// txHah, err := connection.Core.SendRawTransaction(res.Raw.String())
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	//
	// time.Sleep(time.Second*6)

	txRes, _ := connection.Core.GetTransactionByHash("0x3652b5e090719d7a85a2efac75626b5553eaba7cbe0bf24e9df05ac768939421")

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

	txCount, err = connection.Core.GetBlockTransactionCountByNumber(block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if txCount != 1 {
		t.Error("invalid block transaction count")
		t.FailNow()
	}
}
