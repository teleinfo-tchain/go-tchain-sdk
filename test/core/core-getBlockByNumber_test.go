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
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestCoreGetBlockByNumber(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	for _, test := range []struct {
		transactionDetails bool
	}{
		{true},
		// {false},
	} {
		_, err := connection.Core.GetBlockNumber()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		block, err := connection.Core.GetBlockByNumber("pending", test.transactionDetails)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if block == nil {
			t.Error("Block returned is nil")
			t.FailNow()
		}

		if test.transactionDetails {
			fmt.Println("this is detail ")
			fmt.Println("Number:", block.(*dto.BlockDetails).Number)
			fmt.Println("Hash:", block.(*dto.BlockDetails).Hash)
			fmt.Println("ParentHash:", block.(*dto.BlockDetails).ParentHash)
			fmt.Println("LogsBloom:", block.(*dto.BlockDetails).LogsBloom)
			fmt.Println("StateRoot:", block.(*dto.BlockDetails).StateRoot)
			fmt.Println("Generator:", block.(*dto.BlockDetails).Generator)
			fmt.Println("Regulatory:", block.(*dto.BlockDetails).Regulatory)
			fmt.Println("ExtraData:", block.(*dto.BlockDetails).ExtraData)
			fmt.Println("Size:", block.(*dto.BlockDetails).Size)
			fmt.Println("Timestamp:", block.(*dto.BlockDetails).Timestamp)
			fmt.Println("TransactionsRoot:", block.(*dto.BlockDetails).TransactionsRoot)
			fmt.Println("ReceiptsRoot:", block.(*dto.BlockDetails).ReceiptsRoot)
			for _, val := range block.(*dto.BlockDetails).Transactions {
				fmt.Printf("transactions: %#v \n", val)
			}
		} else {
			fmt.Println("this is no detail ")
			fmt.Println("Number:", block.(*dto.BlockNoDetails).Number)
			fmt.Println("Hash:", block.(*dto.BlockNoDetails).Hash)
			fmt.Println("ParentHash:", block.(*dto.BlockNoDetails).ParentHash)
			fmt.Println("LogsBloom:", block.(*dto.BlockNoDetails).LogsBloom)
			fmt.Println("StateRoot:", block.(*dto.BlockNoDetails).StateRoot)
			fmt.Println("Generator:", block.(*dto.BlockNoDetails).Generator)
			fmt.Println("Regulatory:", block.(*dto.BlockNoDetails).Regulatory)
			fmt.Println("ExtraData:", block.(*dto.BlockNoDetails).ExtraData)
			fmt.Println("Size:", block.(*dto.BlockNoDetails).Size)
			fmt.Println("Timestamp:", block.(*dto.BlockNoDetails).Timestamp)
			fmt.Println("TransactionsRoot:", block.(*dto.BlockNoDetails).TransactionsRoot)
			fmt.Println("ReceiptsRoot:", block.(*dto.BlockNoDetails).ReceiptsRoot)
			for idx, val := range block.(*dto.BlockNoDetails).Transactions {
				fmt.Println("transactions:[0]", val[idx])
			}
		}
	}

}
