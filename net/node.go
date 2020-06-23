package net

import (
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
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

func (node *Node) GetPeers() ([]*common.PeerInfo, error) {

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
