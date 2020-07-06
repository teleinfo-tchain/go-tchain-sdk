/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

/**
 * @file web3-clientVersion.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/txpool"
	"testing"

	"github.com/bif/bif-sdk-go/providers"
)

func TestTxpoolStatus(t *testing.T) {

	var connection = txpool.NewTxpool(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	 status, err := connection.Status()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("pending is %d\n", status["pending"])
	fmt.Printf("queued is %d\n", status["queued"])

}

func TestTxpoolInspect(t *testing.T) {

	var connection = txpool.NewTxpool(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	_, err := connection.Inspect()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}