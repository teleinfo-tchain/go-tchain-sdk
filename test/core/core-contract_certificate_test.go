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
	bif "github.com/bif/bifGo"
	"github.com/bif/bifGo/dto"
	"github.com/bif/bifGo/providers"
	"math/big"
	"testing"
)

func TestCoreCertificateContract(t *testing.T) {
	const (
		ContractAddr = "did:bid:00000000000000000000000b"
		AbiJSON      = `[
{"name":"registerCertificate","inputs":[{"name":"Id","type":"string"},{"name":"Context","type":"string"},{"name":"Subject","type":"string"},{"name":"Period","type":"uint64"},{"name":"IssuerAlgorithm","type":"string"},{"name":"IssuerSignature","type":"string"},{"name":"SubjectPublicKey","type":"string"},{"name":"SubjectAlgorithm","type":"string"},{"name":"SubjectSignature","type":"string"}],"outputs":[],"type":"function"},
{"name":"revocedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"name":"queryPeriod","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"period","type":"uint64"}],"type":"function"},
{"name":"queryActive","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"isEnable","type":"bool"}],"type":"function"},
{"name":"queryIssuer","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Issuer","type":"string"}],"type":"function"},
{"name":"queryIssuerSignature","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Id","type":"string"},{"name":"PublicKey","type":"string"},{"name":"Algorithm","type":"string"},{"name":"Signature","type":"string"}],"type":"function"},
{"name":"querySubjectSignature","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Id","type":"string"},{"name":"PublicKey","type":"string"},{"name":"Algorithm","type":"string"},{"name":"Signature","type":"string"}],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"cerdEvent","type":"event"}
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
	result, err := contract.Call(transaction, "queryPeriod", coinbase)
	if err == nil && result != nil {
		fmt.Println(result)
	}

	// =========== 修改 ==========
	hash, err := contract.Send(transaction, "revocedCertificate", coinbase)

	var receipt *dto.TransactionReceipt

	if receipt == nil {
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err == nil && receipt != nil {
		fmt.Println(receipt)
	}
}
