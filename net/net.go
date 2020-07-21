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


package net

import (
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"math/big"
)

// Net - The Net Module
// Net - 网络模块
type Net struct {
	provider providers.ProviderInterface
}

// NewNet - Net Module constructor to set the default provider
// NewNet - Net Module 构造函数来初始化
func NewNet(provider providers.ProviderInterface) *Net {
	net := new(Net)
	net.provider = provider
	return net
}

 /*
  IsListening:
 	EN - Returns true if client is actively listening for network connections.
  	CN - 判断客户端是否正常监听网络连接
  Params:
	- None

  Returns:
	- bool, true是正常，false是不正常
	- error
  Call permissions: Anyone
  */
func (net *Net) IsListening() (bool, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_listening", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}

 /*
  GetPeerCount:
  	EN - Returns number of peers currently connected to the client.
   	CN - 返回当前连接到客户端的对等节点的数量。
  Params:
  	- None

  Returns:
  	- *big.Int, 连接的对等节点的数量
	- error

  Call permissions: Anyone
  */
func (net *Net) GetPeerCount() (*big.Int, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_peerCount", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()

}

 /*
  GetVersion:
   	EN - Returns the current network id.
	CN - 返回当前的网络ID
  Params:
  	- None

  Returns:
  	- string,当前的网络ID
 	- error

  Call permissions: Anyone
  */
func (net *Net) GetVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_version", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

 /*
  GetPeers:
   	EN -
 	CN -
  Params:
  	-
  	-

  Returns:
  	-
 	- err

  Call permissions: Anyone
  */
func (net *Net) GetPeers() ([]*common.PeerInfo, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "admin_peers", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToPeerInfo()
}

func (net *Net) AddPeer(url string) (bool, error) {

	pointer := &dto.RequestResult{}

	params := make([]string, 1)
	params[0] = url

	err := net.provider.SendRequest(pointer, "admin_addPeer", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

func (net *Net) RemovePeer(url string) (bool, error) {

	pointer := &dto.RequestResult{}

	params := make([]string, 1)
	params[0] = url

	err := net.provider.SendRequest(pointer, "admin_removePeer", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}
