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
 * @file net.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package node

import (
	"github.com/bif/bifGo/dto"
	"github.com/bif/bifGo/providers"
	"github.com/bif/go-bif/p2p"
)

// Net - The Net Module
type Node struct {
	provider providers.ProviderInterface
}

// NewNet - Net Module constructor to set the default provider
func NewNode(provider providers.ProviderInterface) *Node {
	net := new(Node)
	net.provider = provider
	return net
}

func (node *Node) GetPeers() ([]*p2p.PeerInfo, error) {

	pointer := &dto.RequestResult{}

	err := node.provider.SendRequest(pointer, "admin_peers", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToPeerInfo()
}

func (node *Node) AddPeer(url string) (bool, error) {

	pointer := &dto.RequestResult{}

	params := make([]string, 1)
	params[0] = url

	err := node.provider.SendRequest(pointer, "admin_addPeer", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

func (node *Node) RemovePeer(url string) (bool, error) {

	pointer := &dto.RequestResult{}

	params := make([]string, 1)
	params[0] = url

	err := node.provider.SendRequest(pointer, "admin_removePeer", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}
