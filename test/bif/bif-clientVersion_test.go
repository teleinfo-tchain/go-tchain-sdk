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
	"testing"

	b "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestBifClientVersion(t *testing.T) {

	var connection = b.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	client, err := connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(client)
}