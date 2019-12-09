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
 * @file ipc-provider_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */
package test

import (
	"testing"

	bif "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func Test_IPCProvider(t *testing.T) {

	var ethClient = bif.NewBif(providers.NewIPCProvider("/tmp/geth.ipc"))

	var _, error = ethClient.ClientVersion()

	if error != nil {
		t.Error(error)
		t.Fail()
	}

}
