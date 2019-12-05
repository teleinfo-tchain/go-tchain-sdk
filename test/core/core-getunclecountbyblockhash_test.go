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
 * @file core-getunclecountbyblockhash_test.go
 * @authors:
 * 		Sigma Prime <sigmaprime.io>
 * @date 2018
 */

package test

import (
	"github.com/bif/bifGo"
	"github.com/bif/bifGo/providers"
	"testing"
)

func TestGetUncleCountByBlockHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	blockNumber, err := connection.Core.GetBlockNumber()

	blockByNumber, err := connection.Core.GetBlockByNumber(blockNumber, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	uncleByHash, err := connection.Core.GetUncleCountByBlockHash(blockByNumber.Hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(uncleByHash.Int64())

	if uncleByHash.Int64() != 0 {
		t.Errorf("Returned uncle for block with no uncle.")
		t.FailNow()
	}

	_, err = connection.Core.GetUncleCountByBlockHash("0x1234")

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}

	_, err = connection.Core.GetUncleCountByBlockHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa0")

	if err == nil {
		t.Errorf("Invalid hash not rejected")
		t.FailNow()
	}
}
