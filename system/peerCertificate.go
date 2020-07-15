package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"strings"
)

const (
	PeerCertificateContractAddr            = "did:bid:00000000000000000000000d"
	//Year                    = int64(24) * 3600 * 365
	//regulatoryAddressLength = 12 // 地址，除去did:bid:还有12字节
)

const PeerCertificateAbiJSON = `[
{"constant": false,"name":"registerCertificate","inputs":[{"name":"publicKey","type":"string"},{"name":"period","type":"uint64"},{"name":"ip","type":"string"},{"name":"port","type":"uint64"},{"name":"company_name","type":"string"},{"name":"company_code","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"publicKey","type":"string"}],"outputs":[],"type":"function"},
{"constant": true,"name":"queryPeriod","inputs":[{"name":"publicKey", "type":"string"}],"outputs":[{"name":"period","type":"uint64"}],"type":"function"},
{"constant": true,"name":"queryActive","inputs":[{"name":"publicKey", "type":"string"}],"outputs":[{"name":"isEnable","type":"bool"}],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"cerdEvent","type":"event"}
]`

type peerCertificate struct {
	super *System
	abi abi.ABI
}

type RegisterCertificateInfo struct {
	PublicKey string
	Period uint64
	IP string
	Port uint64
	CompanyName string
	CompanyCode string
}

func(sys *System) NewPeerCertificate() (*peerCertificate, error){
	parseAbi, err := abi.JSON(strings.NewReader(PeerCertificateAbiJSON))

	if err != nil{
		return nil, err
	}

	peerCertificate := new(peerCertificate)
	peerCertificate.super = sys
	peerCertificate.abi = parseAbi
	return peerCertificate, nil
}

//"registerCertificate","inputs":[{"name":"publicKey","type":"string"},{"name":"period","type":"uint64"},{"name":"ip","type":"string"},{"name":"port","type":"uint64"},{"name":"company_name","type":"string"},{"name":"company_code","type":"string"}],"outputs":[]
func (peerCer *peerCertificate) RegisterCertificate(from common.Address, registerCertificate *RegisterCertificateInfo) (string, error){
	// encoding
	// registerCertificate is a struct we need to use the components.
	var values []interface{}
	values = peerCer.super.StructToInterface(*registerCertificate,values)
	inputEncode, err := peerCer.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}

	transaction := peerCer.super.PrePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return peerCer.super.SendTransaction(transaction)
}

//"revokedCertificate","inputs":[{"name":"publicKey","type":"string"}],"outputs":[]
func (peerCer *peerCertificate) RevokedCertificate(from common.Address, publicKey string)(string, error){
	// encoding
	inputEncode, err := peerCer.abi.Pack("revokedCertificate", publicKey)
	if err != nil {
		return "", err
	}

	transaction := peerCer.super.PrePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return peerCer.super.SendTransaction(transaction)
}

//"queryPeriod","inputs":[{"name":"publicKey", "type":"string"}],"outputs":[{"name":"period","type":"uint64"}]
func (peerCer *peerCertificate) GetPeriod(from common.Address, publicKey string)(uint64, error){
	// encoding
	inputEncode, err := peerCer.abi.Pack("queryPeriod", publicKey)
	if err != nil {
		return 0, err
	}

	transaction := peerCer.super.PrePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := peerCer.super.Call(transaction)
	if err != nil {
		return 0, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))

	var period uint64
	err = peerCer.abi.Methods["queryPeriod"].Outputs.Unpack(&period, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return 0, err
	}

	if period == 0 {
		return 0, ErrCertificateNotExist
	}
	return period, err
}

//"queryActive","inputs":[{"name":"publicKey", "type":"string"}],"outputs":[{"name":"isEnable","type":"bool"}]
func (peerCer *peerCertificate) GetActive(from common.Address, publicKey string)(bool, error){
	// encoding
	inputEncode, err := peerCer.abi.Pack("queryActive", publicKey)
	if err != nil {
		return false, err
	}

	transaction := peerCer.super.PrePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := peerCer.super.Call(transaction)
	if err != nil {
		return false, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var isEnable bool
	err = peerCer.abi.Methods["queryActive"].Outputs.Unpack(&isEnable, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return false, err
	}

	return isEnable, err
}
