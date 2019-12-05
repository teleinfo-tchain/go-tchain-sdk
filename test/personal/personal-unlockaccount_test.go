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
 * @file personal-unlockaccount_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */
package test

import (
	"errors"
	"testing"

	"github.com/bif/bifGo"
	"github.com/bif/bifGo/providers"
)

func TestPersonalUnlockAccount(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider("192.168.104.35:33333", 10, false))

	accounts, err := connection.Personal.ListAccounts()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	result, err := connection.Personal.UnlockAccount(accounts[0], "", 100)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !result {
		t.Error(errors.New("Can't unlock account"))
		t.FailNow()
	}

}
