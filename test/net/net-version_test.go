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
	"errors"
	"github.com/bif/bif-sdk-go/test/resources"
	"sort"
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestNetVersion(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	//Possible options
	po := []string{"1", "2", "3", "4", "42"}

	version, err := connection.Net.GetVersion()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(version)

	if found := sort.SearchStrings(po, version); found < len(po) && po[found] != version {
		t.Error(errors.New("invalid network"))
		t.Fail()
	}

}
