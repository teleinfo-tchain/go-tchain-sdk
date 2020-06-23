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
 * @file bif-sha3_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package test

import (
	"errors"
	"strings"
	"testing"

	bif "github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestUtilsSha3(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("172.20.3.21:44032", 10, false))

	sha3String, err := connection.Utils.Sha3("test")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(sha3String)

	result := strings.Compare("0x9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658", sha3String)

	if result != 0 {
		t.Error(errors.New("Sha3 string not equal calculed hash"))
		t.Fail()
	}

}
