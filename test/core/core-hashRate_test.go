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
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"

	web3 "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestCoreHashRate(t *testing.T) {

	var connection = web3.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	rate, err := connection.Core.GetHashRate()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if rate.Int64() < int64(0) {
		t.Errorf("Less than 0 hash rate")
		t.FailNow()
	}

	t.Log(rate)
}
