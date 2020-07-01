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
 * @file net-getpeercount_test.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"

	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/providers"
)

func TestNodePeers(t *testing.T) {

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	peers, err := connection.GetPeers()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for i := range peers {
		fmt.Printf("%+v\n", peers[i])
	}

	fmt.Println(len(peers))
}

func TestAddPeer(t *testing.T) {

	url := "enode://fc83863fcd8bce46b4a722894c4c70bc5ca74cd12fa32974a84163d002a7c12b7ad28479e414c25ec355857221af4ae8f5d28a0fa8d752308a7ecf4e8eb720bd@192.168.150.20:51260"

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	r, err := connection.AddPeer(url)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(r)
}

func TestRemovePeer(t *testing.T) {

	url := "enode://fc83863fcd8bce46b4a722894c4c70bc5ca74cd12fa32974a84163d002a7c12b7ad28479e414c25ec355857221af4ae8f5d28a0fa8d752308a7ecf4e8eb720bd@192.168.150.20:51260"

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	r, err := connection.RemovePeer(url)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(r)
}
