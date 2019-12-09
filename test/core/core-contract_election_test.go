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

func TestCoreElectionContract(t *testing.T) {
	const (
		ContractAddr = "did:bid:000000000000000000000009"
		AbiJSON      = `[
{"name":"registerWitness","inputs":[{"name":"nodeUrl","type":"string"},{"name":"website","type":"string"},{"name":"name","type":"string"}],"outputs":[],"type":"function"}, 
{"name":"unregisterWitness","inputs":[],"outputs":[],"type":"function"},
{"name":"voteWitnesses","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"name":"cancelVote","inputs":[],"outputs":[],"type":"function"},
{"name":"startProxy","inputs":[],"outputs":[],"type":"function"},
{"name":"stopProxy","inputs":[],"outputs":[],"type":"function"},
{"name":"cancelProxy","inputs":[],"outputs":[],"type":"function"},
{"name":"setProxy","inputs":[{"name":"proxy","type":"string"}],"outputs":[],"type":"function"},
{"name":"stake","inputs":[{"name":"stakeCount","type":"uint256"}],"outputs":[],"type":"function"},
{"name":"unStake","inputs":[],"outputs":[],"type":"function"},
{"name":"extractOwnBounty","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"electEvent","type":"event"},
{"name":"issueAddtitionalBounty","inputs":[],"outputs":[],"type":"function"},
{"name":"queryCandidates","inputs":[{"name":"candiaddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"votecount","type":"uint64"},{"name":"active","type":"bool"},{"name":"url","type":"string"},{"name":"totalbounty","type":"uint64"},{"name":"extractedbounty","type":"uint64"},{"name":"lastExtractTime","type":"uint64"},{"name":"website","type":"string"},{"name":"name","type":"string"}],"type":"function"},
{"name":"queryAllCandidates","inputs":[],"outputs":[{"name":"listnum","type":"uint64"},{"name":"candinfolist","type":"string"}],"type":"function"},
{"name":"queryVoter","inputs":[{"name":"voteraddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"isproxy","type":"bool"},{"name":"proxyvotecount","type":"uint64"},{"name":"proxy","type":"string"},{"name":"lastvotecount","type":"uint64"},{"name":"timestamp","type":"uint64"},{"name":"votecandidateslist","type":"string"}],"type":"function"},
{"name":"queryStake","inputs":[{"name":"voteraddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"stakecount","type":"uint64"},{"name":"timestamp","type":"uint64"}],"type":"function"},
{"name":"queryVoterList","inputs":[{"name":"candiaddress","type":"string"}],"outputs":[{"name":"candiaddress","type":"string"},{"name":"voternum","type":"uint64"},{"name":"voterlist","type":"string"}],"type":"function"}
]`
	)

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))
	contract, err := connection.Core.NewContract(AbiJSON)

	transaction := new(dto.TransactionParameters)
	coinbase, err := connection.Core.GetCoinbase()
	transaction.From = coinbase
	transaction.Gas = big.NewInt(4000000)
	transaction.To = ContractAddr

	// ============= 查询 ===============
	result, err := contract.Call(transaction, "queryAllCandidates")
	if err == nil && result != nil {
		fmt.Println(result)
	}

	// =========== 修改 ==========
	hash, err := contract.Send(transaction, "startProxy")

	var receipt *dto.TransactionReceipt

	if receipt == nil {
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err == nil && receipt != nil {
		fmt.Println(receipt)
	}
}
