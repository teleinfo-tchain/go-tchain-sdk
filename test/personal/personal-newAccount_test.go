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
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestPersonalNewAccount(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	address, err := connection.Personal.NewAccount("testSm2", true)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if address[8:10] != "73"{
		t.Error("generate account error")
		t.FailNow()
	}
	t.Log(address)

	address, err= connection.Personal.NewAccount("notSm2", false)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if address[8:10] == "73"{
		t.Error("generate account error")
		t.FailNow()
	}
	t.Log(address)
}
