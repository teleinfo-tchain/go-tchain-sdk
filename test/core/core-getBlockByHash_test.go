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
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestCoreGetBlockByHash(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	const transactionDetails = false
	blockByNumber, err := connection.Core.GetBlockByNumber(block.LATEST, transactionDetails)

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
		(blockByNumber.(*dto.BlockNoDetails).Miner != blockByHash.(*dto.BlockNoDetails).Miner) ||
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
