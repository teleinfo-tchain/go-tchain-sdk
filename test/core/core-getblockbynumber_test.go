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
	"fmt"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestCoreGetBlockByNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	const transactionDetails = true

	blockNumber, err := connection.Core.GetBlockNumber()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	block, err := connection.Core.GetBlockByNumber(hexutil.EncodeBig(blockNumber), transactionDetails)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if block == nil {
		t.Error("Block returned is nil")
		t.FailNow()
	}

	if transactionDetails {
		fmt.Println("gasLimit:", block.(*dto.BlockDetails).GasLimit)
		fmt.Println("extraData:", block.(*dto.BlockDetails).ExtraData)
		fmt.Println("LogsBloom:", block.(*dto.BlockDetails).LogsBloom)
		fmt.Println("MixHash:", block.(*dto.BlockDetails).MixHash)
		fmt.Println("ReceiptsRoot:", block.(*dto.BlockDetails).ReceiptsRoot)
		fmt.Println("StateRoot:", block.(*dto.BlockDetails).StateRoot)
		for _, val := range block.(*dto.BlockDetails).Transactions {
			fmt.Println("transactions:", val.Hash)
		}
		fmt.Println("transactionLen:", len(block.(*dto.BlockDetails).Transactions))
		fmt.Println("TransactionsRoot:", block.(*dto.BlockDetails).TransactionsRoot)
	} else {
		fmt.Println("gasLimit:", block.(*dto.BlockNoDetails).GasLimit)
		fmt.Println("extraData:", block.(*dto.BlockNoDetails).ExtraData)
		fmt.Println("LogsBloom:", block.(*dto.BlockNoDetails).LogsBloom)
		fmt.Println("MixHash:", block.(*dto.BlockNoDetails).MixHash)
		fmt.Println("ReceiptsRoot:", block.(*dto.BlockNoDetails).ReceiptsRoot)
		fmt.Println("StateRoot:", block.(*dto.BlockNoDetails).StateRoot)
		for idx, val := range block.(*dto.BlockNoDetails).Transactions {
			fmt.Println("transactions:[0]", val[idx])
		}
		fmt.Println("transactionLen:", len(block.(*dto.BlockNoDetails).Transactions))
		fmt.Println("TransactionsRoot:", block.(*dto.BlockNoDetails).TransactionsRoot)
	}
}

