package common

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
