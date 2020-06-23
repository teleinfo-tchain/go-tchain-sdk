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
 * @file bif-blocknumber_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package test

import (
	"fmt"
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestCoreBlockNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("172.20.3.21:44032", 10, false))

	blockNumber, err := connection.Core.GetBlockNumber()

	fmt.Println(blockNumber)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if blockNumber.Int64() < 0 {
		t.Errorf("Invalid Block Number")
		t.Fail()
	}
}
