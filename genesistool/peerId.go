package genesistool

import (
	libp2pcorecrypto "github.com/libp2p/go-libp2p-core/crypto"
	libp2pcorepeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/tchain/go-tchain-sdk/crypto"
	"github.com/tchain/go-tchain-sdk/crypto/config"
)

func PeerIdGen(nodePrivateHex string) (string, error) {
	p2pPrivateKey, err := crypto.HexToECDSA(nodePrivateHex, config.SECP256K1)
	if err != nil {
		return "", err
	}

	var peerId libp2pcorepeer.ID
	if peerId, err = libp2pcorepeer.IDFromPrivateKey((*libp2pcorecrypto.Secp256k1PrivateKey)(p2pPrivateKey)); err != nil {
		return "", err
	}

	return peerId.String(), nil

}
