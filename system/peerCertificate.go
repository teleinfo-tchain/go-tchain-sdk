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
	// regulatoryAddressLength     = 12 // 地址，除去did:bid:还有12字节
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

Params:
	- from: [20]byte，交易发送方地址
	- registerCertificate:  *dto.RegisterCertificateInfo，包含可信证书的信息
		Id          string //节点证书的bid,，必须和public_key相同
		PublicKey   string //53个字符的公钥
		NodeName    string //节点名称，不含敏感词的字符串
		MessageSha3 string //消息sha3后的16进制字符串
		Signature   string //对上一个字段消息的签名，16进制字符串
		NodeType    uint64 //节点类型，0企业，1个人
		Period      uint64 //证书有效期，以年为单位的整型
		IP          string // ip
		Port        uint64 // port
		CompanyName string //公司名（如果是个人，则是个人姓名）
		CompanyCode string //公司代码

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

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

Params:
	- from: [20]byte，交易发送方地址
	- id: string，节点证书的bid

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: 只有监管节点地址可以调用
*/
func (peerCer *PeerCertificate) RevokedCertificate(from common.Address, id string) (string, error) {
	// encoding
	inputEncode, err := peerCer.abi.Pack("revokedCertificate", id)
	if err != nil {
		return "", err
	}

	transaction := peerCer.super.prePareTransaction(from, PeerCertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return peerCer.super.sendTransaction(transaction)
}

/*
GetPeriod: 查看证书有效期

Params:
	- id: string，节点证书的bid

Returns:
	- uint64，返回证书有效期，如果证书被吊销，则有效期是0
	- error

Call permissions: Anyone
*/
func (peerCer *PeerCertificate) GetPeriod(id string) (uint64, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_period", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
GetActive: 查看证书是否有效

Params:
	- id: string，节点证书的bid

Returns:
	- bool，true可用，false不可用
	- error

Call permissions: Anyone
*/
func (peerCer *PeerCertificate) GetActive(id string) (bool, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_active", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
GetPeerCertificate: 查看证书的信息

Params:
	- id: string，节点证书的bid

Returns:
	- *dto.PeerCertificate
		Id          string   `json:"id"`          //唯一索引
		Issuer      string   `json:"issuer"`      //颁发者地址
		PublicKey   string   `json:"publicKey"`   //节点公钥
		NodeName    string   `json:"nodeName"`    //节点名称
		Signature   string   `json:"signature"`   //节点签名内容
		NodeType    uint64   `json:"nodeType"`    //节点类型0企业，1个人
		CompanyName string   `json:"companyName"` //公司名称
		CompanyCode string   `json:"companyCode"` //公司信用代码
		IssuedTime  *big.Int `json:"issuedTime"`  //颁发时间
		Period      uint64   `json:"period"`      //有效期
		IsEnable    bool     `json:"isEnable"`    //true 凭证有效，false 凭证已撤销
	- error

Call permissions: Anyone
*/
func (peerCer *PeerCertificate) GetPeerCertificate(id string) (*dto.PeerCertificate, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_peerCertificate", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToPeerCertificate()
}
