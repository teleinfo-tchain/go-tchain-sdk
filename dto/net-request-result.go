package dto

import (
	"encoding/json"
)

type NetRequestResult struct {
	RequestResult
}

type PeerInfo struct {
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
}

type NodeInfo struct {
	ID    string `json:"id"`    // Unique node identifier (also the encryption key)
	Name  string `json:"name"`  // Name of the node, including client type, version, OS, custom data
	Enode string `json:"enode"` // Enode URL for adding this peer from remote peers
	IP    string `json:"ip"`    // IP address of the node
	Ports struct {
		Discovery int `json:"discovery"` // UDP listening port for discovery protocol
		Listener  int `json:"listener"`  // TCP listening port for RLPx
	} `json:"ports"`
	ListenAddr string                 `json:"listenAddr"`
	Protocols  map[string]interface{} `json:"protocols"`
}

func (pointer *NetRequestResult) ToPeerInfo() ([]*PeerInfo, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	resultLi := (pointer).Result.([]interface{})

	peerInfoLi := make([]*PeerInfo, len(resultLi))

	for i, v := range resultLi {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}

		info := &PeerInfo{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, info)

		peerInfoLi[i] = info

	}

	return peerInfoLi, nil

}

func (pointer *NetRequestResult) ToNodeInfo() (*NodeInfo, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	nodeInfo := &NodeInfo{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, nodeInfo)

	return nodeInfo, err
}
