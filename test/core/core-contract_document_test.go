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

/**
 * @file contract.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2018
 */

package test

import (
	"fmt"
	bif "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"math/big"
	"testing"
)

func TestCoreDocumentContract(t *testing.T) {
	const (
		ContractAddr = "did:bid:00000000000000000000000a"
		AbiJSON      = `[
{"name":"InitializationDDO","inputs":[{"name":"bidType","type":"uint64"}],"outputs":[],"type":"function"},
{"name":"SetBidName","inputs":[{"name":"bidName","type":"string"}],"outputs":[],"type":"function"},
{"name":"FindDDOByType","inputs":[{"name":"key","type":"uint64"},{"name":"value","type":"string"}],"outputs":[{"name":"result","type":"string"}],"type":"function"},
{"name":"AddPublicKey","inputs":[{"name":"type","type":"string"},{"name":"authority","type":"string"},{"name":"publicKey","type":"string"}],"outputs":[],"type":"function"},
{"name":"DeletePublicKey","inputs":[{"name":"publicKey","type":"string"}],"outputs":[],"type":"function"},
{"name":"AddProof","inputs":[{"name":"type","type":"string"},{"name":"proofID","type":"string"}],"outputs":[],"type":"function"},
{"name":"DeleteProof","inputs":[{"name":"proofID","type":"string"}],"outputs":[],"type":"function"},
{"name":"AddAttr","inputs":[{"name":"type","type":"string"},{"name":"value","type":"string"}],"outputs":[],"type":"function"},
{"name":"DeleteAttr","inputs":[{"name":"type","type":"string"}],"outputs":[],"type":"function"},
{"name":"Enable","inputs":[],"outputs":[],"type":"function"},
{"name":"Disable","inputs":[],"outputs":[],"type":"function"},
{"name":"IsEnable","inputs":[{"name":"key","type":"uint64"},{"name":"value","type":"string"}],"outputs":[{"name":"result","type":"bool"}],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"bidEvent","type":"event"}
]`
	)

	var connection = bif.NewBif(providers.NewHTTPProvider("172.20.3.21:44032", 10, false))
	contract, err := connection.Core.NewContract(AbiJSON)

	transaction := new(dto.TransactionParameters)
	coinbase, err := connection.Core.GetCoinbase()
	transaction.From = coinbase
	transaction.Gas = big.NewInt(4000000)
	transaction.To = ContractAddr

	// ============= 查询 ===============
	result, err := contract.Call(transaction, "FindDDOByType", uint64(1), coinbase)
	if err == nil && result != nil {
		fmt.Println(result)
	}

	// =========== 修改 ==========
	hash, err := contract.Send(transaction, "InitializationDDO", coinbase)

	var receipt *dto.TransactionReceipt

	if receipt == nil {
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err == nil && receipt != nil {
		fmt.Println(receipt)
	}
}
