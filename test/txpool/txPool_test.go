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

package test

import (
	"github.com/tchain/go-tchain-sdk/test/resources"
	"github.com/tchain/go-tchain-sdk/txpool"
	"strconv"
	"testing"

	"github.com/tchain/go-tchain-sdk/providers"
)

func TestGetStatus(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	status, err := connection.GetStatus()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(status)

}

func TestTxPoolInspect(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	inspect, err := connection.Inspect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(inspect)
}

func TestTxPoolContent(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	transactions, err := connection.Content()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactions)
}