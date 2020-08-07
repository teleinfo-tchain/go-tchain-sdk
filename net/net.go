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
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"math/big"
)

// Net - The Net Module
// Net - ����ģ��
type Net struct {
	provider providers.ProviderInterface
}

// NewNet - Net Module constructor to set the default provider
// NewNet - Net Module ���캯������ʼ��
func NewNet(provider providers.ProviderInterface) *Net {
	net := new(Net)
	net.provider = provider
	return net
}

/*
  IsListening:
 	EN - Returns true if client is actively listening for network connections.
  	CN - �жϿͻ����Ƿ�����������������
  Params:
	- None

  Returns:
	- bool, true��������false�ǲ�����
	- error
  Call permissions: Anyone
*/
func (net *Net) IsListening() (bool, error) {

	pointer := &dto.NetRequestResult{}

	err := net.provider.SendRequest(pointer, "net_listening", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}

/*
  GetPeerCount:
  	EN - Returns number of peers currently connected to the client.
   	CN - ���ص�ǰ���ӵ��ͻ��˵ĶԵȽڵ��������
  Params:
  	- None

  Returns:
  	- *big.Int, ���ӵĶԵȽڵ������
	- error

  Call permissions: Anyone
*/
func (net *Net) GetPeerCount() (*big.Int, error) {

	pointer := &dto.NetRequestResult{}

	err := net.provider.SendRequest(pointer, "net_peerCount", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()

}

/*
  GetVersion:
   	EN - Returns the current network id.
	CN - ���ص�ǰ������ID
  Params:
  	- None

  Returns:
  	- string,��ǰ������ID
 	- error

  Call permissions: Anyone
*/
func (net *Net) GetVersion() (string, error) {

	pointer := &dto.NetRequestResult{}

	err := net.provider.SendRequest(pointer, "net_version", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  GetPeers:
   	EN - Get all peer nodes info
 	CN - ��ȡ�������ӵĶԵȽڵ����ϸ��Ϣ
  Params:
  	- None

  Returns:
  	- []*common.PeerInfo, �ڵ���Ϣ�б�
 		Enode   string   `json:"enode"` // Node URL
		ID      string   `json:"id"`    // Unique node identifier
		Name    string   `json:"name"`  // Name of the node, including client type, version, OS, custom data
		Caps    []string `json:"caps"`  // Protocols advertised by this peer
		Network struct {
			LocalAddress  string `json:"localAddress"`  // Local endpoint of the TCP data connection
			RemoteAddress string `json:"remoteAddress"` // Remote endpoint of the TCP data connection
			Inbound       bool   `json:"inbound"`
			Trusted       bool   `json:"trusted"`
			Static        bool   `json:"static"`
		} `json:"network"`
		Protocols map[string]interface{} `json:"protocols"` // Sub-protocol specific metadata fields
 	- err

  Call permissions: Anyone
*/
func (net *Net) GetPeers() ([]*dto.PeerInfo, error) {

	pointer := &dto.NetRequestResult{}

	err := net.provider.SendRequest(pointer, "admin_peers", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToPeerInfo()
}

/*
  AddPeer:
   	EN - Add peer
	CN - �����µĽڵ�
  Params:
  	- url, string, �ڵ��url

  Returns:
  	- bool��true���ӳɹ���false����ʧ��
	- error

  Call permissions: Anyone
*/
func (net *Net) AddPeer(url string) (bool, error) {

	pointer := &dto.NetRequestResult{}

	params := make([]string, 1)
	params[0] = url

	err := net.provider.SendRequest(pointer, "admin_addPeer", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
  RemovePeer:
   	EN - Remove peer
	CN - �Ƴ����ӵĽڵ�
  Params:
  	- url, string, �ڵ��url

  Returns:
   	- bool��true�Ƴ��ɹ���false�Ƴ�ʧ��
 	- error

  Call permissions: Anyone
*/
func (net *Net) RemovePeer(url string) (bool, error) {

	pointer := &dto.NetRequestResult{}

	params := make([]string, 1)
	params[0] = url

	err := net.provider.SendRequest(pointer, "admin_removePeer", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

 /*
  GetNodeInfo:
   	EN - The host node info
 	CN - �����ڵ����Ϣ
  Params:
  	- None

  Returns:
  	- *dto.NodeInfo
 		ID    string `json:"id"`    // Unique node identifier (also the encryption key)
		Name  string `json:"name"`  // Name of the node, including client type, version, OS, custom data
		Enode string `json:"enode"` // Enode URL for adding this peer from remote peers
		IP    string `json:"ip"` // IP address of the node
		Ports struct {
			Discovery int `json:"discovery"` // UDP listening port for discovery protocol
			Listener  int `json:"listener"`  // TCP listening port for RLPx
		} `json:"ports"`
		ListenAddr string                 `json:"listenAddr"`
		Protocols  map[string]interface{} `json:"protocols"`
 	- error

  Call permissions: Anyone
  */
func (net *Net) GetNodeInfo() (*dto.NodeInfo, error) {

	pointer := &dto.NetRequestResult{}

	err := net.provider.SendRequest(pointer, "admin_nodeInfo", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToNodeInfo()
}

 /*
  GetDataDir:
   	EN - retrieves the current data directory the node is using
 	CN - �����ڵ�����ʹ�õĵ�ǰ����Ŀ¼
  Params:
  	- None

  Returns:
  	- string
 	- error

  Call permissions: Anyone
  */
func (net *Net) GetDataDir() (string, error) {

	pointer := &dto.NetRequestResult{}

	err := net.provider.SendRequest(pointer, "admin_datadir", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}