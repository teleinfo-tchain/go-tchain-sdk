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
	"fmt"
	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/test/resources"
	"sort"
	"strconv"
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func TestNetPeerCount(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	peers, err := connection.Net.GetPeerCount()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(peers.Uint64())

}

func TestNetListening(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	listening, err := connection.Net.IsListening()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !listening {
		t.Error("not listening")
		t.Fail()
	}
}

func TestGetDataDir(t *testing.T) {

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	dataDir, err := connection.GetDataDir()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(dataDir)
}

func TestGetNodeInfo(t *testing.T) {

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	nodeInfo, err := connection.GetNodeInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n ", nodeInfo)
}

func TestNetVersion(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	// Possible options
	po := []string{"1", "2", "3", "4", "42", "333"}

	version, err := connection.Net.GetVersion()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Printf("version is %s \n", version)
	sort.Strings(po)
	if found := sort.SearchStrings(po, version); found < len(po) && po[found] != version {
		t.Error("invalid network")
		t.Fail()
	}

}
