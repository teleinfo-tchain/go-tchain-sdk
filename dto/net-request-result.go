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

	var data0 []byte
	var err0 error
	if value, ok := pointer.Result.(string); ok {
		//come from websock
		data0 = []byte(value)
	} else {
		//come from http
		resultLi := pointer.Result.([]interface{})

		if len(resultLi) == 0 {
			return nil, EMPTYRESPONSE
		}
		data0, err0 = json.Marshal(resultLi)
	}

	if err0 != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	var resultLi []interface{}

	err0 = json.Unmarshal(data0, &resultLi)

	peerInfoLi := make([]*PeerInfo, len(resultLi))

	for i, v := range resultLi {

		var data []byte
		var err error
		if value, ok := v.(string); ok {
			//come from websock
			data = []byte(value)
		} else {
			//come from http
			result := v.(map[string]interface{})

			if len(result) == 0 {
				return nil, EMPTYRESPONSE
			}
			data, err = json.Marshal(result)
		}

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		nodeInfo := &PeerInfo{}

		err = json.Unmarshal(data, nodeInfo)

		peerInfoLi[i] = nodeInfo
	}
	return peerInfoLi, nil
}

func (pointer *NetRequestResult) ToNodeInfo() (*NodeInfo, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	var data []byte
	var err error
	if value, ok := pointer.Result.(string); ok {
		//come from websock
		data = []byte(value)
	} else {
		//come from http
		result := (pointer).Result.(map[string]interface{})

		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}
		data, err = json.Marshal(result)
	}

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	nodeInfo := &NodeInfo{}

	err = json.Unmarshal(data, nodeInfo)

	return nodeInfo, err
}
