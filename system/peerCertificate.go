package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
	"strings"
)

//监管节点向节点颁发证书，节点由普通节点变成可信节点。

const (
	PeerCertificateContractAddr            = "did:bid:00000000000000000000000d"
	Year                    = uint64(24) * 3600 * 365
	regulatoryAddressLength = 12 // 地址，除去did:bid:还有12字节
)

const PeerCertificateAbiJSON = `[
{"constant": false,"name":"registerCertificate","inputs":[{"name":"id","type":"string"},{"name":"publicKey","type":"string"},{"name":"nodeName","type":"string"},{"name":"messageSha3","type":"string"},{"name":"signature","type":"string"},{"name":"nodeType","type":"uint64"},{"name":"period","type":"uint64"},{"name":"ip","type":"string"},{"name":"port","type":"uint64"},{"name":"companyName","type":"string"},{"name":"companyCode","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"peerEvent","type":"event"}
]`

type peerCertificate struct {
	super *System
	abi abi.ABI
}

func(sys *System) NewPeerCertificate() *peerCertificate {
	parseAbi, _ := abi.JSON(strings.NewReader(PeerCertificateAbiJSON))

	peerCertificate := new(peerCertificate)
	peerCertificate.super = sys
	peerCertificate.abi = parseAbi
	return peerCertificate
}

//为节点颁发可信证书可信
func (peerCer *peerCertificate) RegisterCertificate(from common.Address, registerCertificate *dto.RegisterCertificateInfo) (string, error){
	// encoding
	// registerCertificate is a struct we need to use the components.
	registerCertificate.Id = registerCertificate.PublicKey
	var values []interface{}
	values = peerCer.super.StructToInterface(*registerCertificate,values)
	inputEncode, err := peerCer.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}

	transaction := peerCer.super.PrePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return peerCer.super.SendTransaction(transaction)
}

//吊销节点的可信证书可信
func (peerCer *peerCertificate) RevokedCertificate(from common.Address, publicKey string)(string, error){
	// encoding
	inputEncode, err := peerCer.abi.Pack("revokedCertificate", publicKey)
	if err != nil {
		return "", err
	}

	transaction := peerCer.super.PrePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return peerCer.super.SendTransaction(transaction)
}

// 证书有效期
func (peerCer *peerCertificate) GetPeriod(publicKey string)(uint64, error){
	params := make([]string, 1)
	params[0] = publicKey

	pointer := &dto.RequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_period", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToPeerPeriod()
}

// 证书是否可用
func (peerCer *peerCertificate) GetActive(publicKey string)(bool, error){
	params := make([]string, 1)
	params[0] = publicKey

	pointer := &dto.RequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_active", params)
	if err != nil {
		return false, err
	}

	return pointer.ToPeerActive()
}


// 证书的信息
func (peerCer *peerCertificate) GetPeerCertificate(publicKey string)(*dto.PeerCertificate, error){
	params := make([]string, 1)
	params[0] = publicKey

	pointer := &dto.RequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_peerCertificate", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToPeerCertificate()
}