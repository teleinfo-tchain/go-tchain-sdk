package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
	"strings"
)

const (
	PeerCertificateContractAddr = "did:bid:00000000000000000000000d"
	Year                        = uint64(24) * 3600 * 365
	regulatoryAddressLength     = 12 // 地址，除去did:bid:还有12字节
)

// 节点可信的AbiJson数据
const PeerCertificateAbiJSON = `[
{"constant": false,"name":"registerCertificate","inputs":[{"name":"id","type":"string"},{"name":"publicKey","type":"string"},{"name":"nodeName","type":"string"},{"name":"messageSha3","type":"string"},{"name":"signature","type":"string"},{"name":"nodeType","type":"uint64"},{"name":"period","type":"uint64"},{"name":"ip","type":"string"},{"name":"port","type":"uint64"},{"name":"companyName","type":"string"},{"name":"companyCode","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"peerEvent","type":"event"}
]`

// PeerCertificate - The PeerCertificate Module
type PeerCertificate struct {
	super *System
	abi   abi.ABI
}

// NewPeerCertificate - 初始化PeerCertificate
func (sys *System) NewPeerCertificate() *PeerCertificate {
	parseAbi, _ := abi.JSON(strings.NewReader(PeerCertificateAbiJSON))

	peerCertificate := new(PeerCertificate)
	peerCertificate.super = sys
	peerCertificate.abi = parseAbi
	return peerCertificate
}

/*
 RegisterCertificate: 为节点颁发可信证书可信

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 只有监管节点地址可以调用
*/
func (peerCer *PeerCertificate) RegisterCertificate(from common.Address, registerCertificate *dto.RegisterCertificateInfo) (string, error) {
	// encoding
	// registerCertificate is a struct we need to use the components.
	registerCertificate.Id = registerCertificate.PublicKey
	var values []interface{}
	values = peerCer.super.structToInterface(*registerCertificate, values)
	inputEncode, err := peerCer.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}

	transaction := peerCer.super.prePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return peerCer.super.sendTransaction(transaction)
}

/*
 RevokedCertificate: 吊销节点的可信证书可信

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 只有监管节点地址可以调用
*/
func (peerCer *PeerCertificate) RevokedCertificate(from common.Address, publicKey string) (string, error) {
	// encoding
	inputEncode, err := peerCer.abi.Pack("revokedCertificate", publicKey)
	if err != nil {
		return "", err
	}

	transaction := peerCer.super.prePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return peerCer.super.sendTransaction(transaction)
}

/*
 GetPeriod: 查看证书有效期

 Returns： 返回证书有效期，如果证书被吊销，则有效期是0
*/
func (peerCer *PeerCertificate) GetPeriod(publicKey string) (uint64, error) {
	params := make([]string, 1)
	params[0] = publicKey

	pointer := &dto.RequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_period", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToPeerPeriod()
}

/*
 GetActive: 查看证书是否有效

 Returns： bool，true可用，false不可用
*/
func (peerCer *PeerCertificate) GetActive(publicKey string) (bool, error) {
	params := make([]string, 1)
	params[0] = publicKey

	pointer := &dto.RequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_active", params)
	if err != nil {
		return false, err
	}

	return pointer.ToPeerActive()
}

/*
 GetPeerCertificate: 查看证书的信息

 Returns： *dto.PeerCertificate
*/
func (peerCer *PeerCertificate) GetPeerCertificate(publicKey string) (*dto.PeerCertificate, error) {
	params := make([]string, 1)
	params[0] = publicKey

	pointer := &dto.RequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_peerCertificate", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToPeerCertificate()
}
