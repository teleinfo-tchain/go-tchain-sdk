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

	t.Log(len(peers))
}

// 如何添加和移除需要确定
func TestAddPeer(t *testing.T) {

	url := "/ip4/169.254.248.29/tcp/44051/p2p/16Uiu2HAm4TSmKV3QAVzbd1V8mMpgqA3xvTEiLF71WbdojVvf3vt1"

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	r, err := connection.AddPeer(url)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(r)
}

func TestRemovePeer(t *testing.T) {

	url := "/ip4/169.254.248.29/tcp/44051/p2p/16Uiu2HAm4TSmKV3QAVzbd1V8mMpgqA3xvTEiLF71WbdojVvf3vt1"

	var connection = net.NewNet(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	r, err := connection.RemovePeer(url)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(r)
}
