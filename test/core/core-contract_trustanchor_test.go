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
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

func TestCoreTrustanchorContract(t *testing.T) {
	const (
		ContractAddr = "did:bid:00000000000000000000000c"
		AbiJSON      = `[
{"name":"registerTrustAnchor","inputs":[{"name":"anchor","type":"string"},{"name":"anchortype","type":"uint64"},{"name":"anchorname","type":"string"},{"name":"company","type":"string"},{"name":"companyurl","type":"string"},{"name":"website","type":"string"},{"name":"documenturl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"name":"unRegisterTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[],"type":"function"},
{"name":"isTrustAnchor","inputs":[{"name":"address","type":"string"}],"outputs":[{"name":"trustanchor","type":"bool"}],"type":"function"},
{"name":"updateBaseAnchorInfo","inputs":[{"name":"anchor","type":"string"},{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"name":"updateExependAnchorInfo","inputs":[{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"name":"extractOwnBounty","inputs":[{"name":"anchor","type":"string"}],"outputs":[],"type":"function"},
{"name":"queryTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[{"name":"id","type":"string"},{"name":"name","type":"string"},{"name":"company","type":"string"},{"name":"CompanyUrl","type":"string"},{"name":"website","type":"string"},{"name":"ServerUrl","type":"string"},{"name":"DocumentUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"},{"name":"TrustAnchorType","type":"uint64"},{"name":"status","type":"uint64"},{"name":"active","type":"bool"},{"name":"totalbounty","type":"uint64"},{"name":"extractedBounty","type":"uint64"},{"name":"lastExtracTime","type":"uint64"},{"name":"votecount","type":"uint64"},{"name":"stake","type":"uint64"},{"name":"createDate","type":"uint64"},{"name":"certificateAcount","type":"uint64"}],"type":"function"},
{"name":"queryTrustAnchorStatus","inputs":[{"name":"anchor","type":"string"}],"outputs":[{"name":"anchorstatus","type":"uint8"}],"type":"function"},
{"name":"queryBaseTrustAnchorList","inputs":[],"outputs":[{"name":"baseList","type":"string"}],"type":"function"},
{"name":"queryBaseTrustAnchorNum","inputs":[],"outputs":[{"name":"baseListNum","type":"uint64"}],"type":"function"},
{"name":"queryExpendTrustAnchorList","inputs":[],"outputs":[{"name":"expendList","type":"string"}],"type":"function"},
{"name":"queryExpendTrustAnchorNum","inputs":[],"outputs":[{"name":"expendListNum","type":"uint64"}],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"trustAnchorEvent","type":"event"},
{"name":"voteElect","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"name":"cancelVote","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"name":"queryVoter","inputs":[{"name":"voterAddress","type":"string"}],"outputs":[{"name":"voterInfo","type":"string"}],"type":"function"},
{"name":"checkSenderAddress","inputs":[{"name":"address","type":"string"}],"outputs":[{"name":"supernode","type":"bool"}],"type":"function"}
]`
	)

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	contract, err := connection.Core.NewContract(AbiJSON)

	transaction := new(dto.TransactionParameters)
	coinbase, err := connection.Core.GetCoinbase()
	transaction.From = coinbase
	transaction.Gas = big.NewInt(4000000)
	transaction.To = ContractAddr

	// ============= 查询 ===============
	result, err := contract.Call(transaction, "checkSenderAddress", coinbase)
	if err == nil && result != nil {
		fmt.Println(result)
	}

	// =========== 修改 ==========
	hash, err := contract.Send(transaction, "unRegisterTrustAnchor", coinbase)

	var receipt *dto.TransactionReceipt

	if receipt == nil {
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err == nil && receipt != nil {
		fmt.Println(receipt)
	}
}
