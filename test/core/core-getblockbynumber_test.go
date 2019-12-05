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
 * @file core-getBlockByNumber_test.go
 * @authors:
 *   Jérôme Laurens <jeromelaurens@gmail.com>
 * @date 2017
 */

package test

import (
	"testing"

	"github.com/bif/bifGo"
	"github.com/bif/bifGo/providers"
)

func TestCoreGetBlockByNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	blockNumber, err := connection.Core.GetBlockNumber()

	block, err := connection.Core.GetBlockByNumber(blockNumber, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if block == nil {
		t.Error("Block returned is nil")
		t.FailNow()
	}

	if block.Number.Int64() == 0 {
		t.Error("Block not found")
		t.FailNow()
	}

}
