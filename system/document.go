package system

import (
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/dto"
	"strings"
)

const (
	DocContractAddr = "did:bid:00000000000000000000000a"
)

// const (
// 	bid = iota
// 	publicKeyPem
// 	publicKeyJwk
// 	publicKeyHex
// 	publicKeyBase64
// 	publicKeyBase58
// 	publicKeyMultibase
// 	ethereumAddress
// )
//
// var (
// 	ErrBidDDOExistsInvalid = errors.New("the DDO document for this bid has been initialized")
// 	ErrNotHaveAuthority    = errors.New("对不起，您没有权限做该操作")
// )

// did文档的AbiJson数据
const DocAbiJSON = `[
{"constant": false,"name":"init","inputs":[{"name":"bid_type","type":"uint64"}],"outputs":[],"type":"function"},
{"constant": false,"name":"setBidName","inputs":[{"name":"id","type":"string"},{"name":"bid_name","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addPublic","inputs":[{"name":"id","type":"string"},{"name":"public_type","type":"string"},{"name":"public_auth","type":"string"},{"name":"public_key","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delPublic","inputs":[{"name":"id","type":"string"},{"name":"public_key","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addAuth","inputs":[{"name":"id","type":"string"},{"name":"auth","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delAuth","inputs":[{"name":"id","type":"string"},{"name":"auth","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addService","inputs":[{"name":"id","type":"string"},{"name":"service_id","type":"string"},{"name":"service_type","type":"string"},{"name":"service_endpoint","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delService","inputs":[{"name":"id","type":"string"},{"name":"service_id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addProof","inputs":[{"name":"id","type":"string"},{"name":"proof_type","type":"string"},{"name":"proof_creator","type":"string"},{"name":"proof_sign","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delProof","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addExtra","inputs":[{"name":"id","type":"string"},{"name":"extra","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delExtra","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"enable","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"disable","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"bidEvent","type":"event"}
]`

// Doc - The Doc Module
type Doc struct {
	super *System
	abi   abi.ABI
}

// NewDoc - NewDoc初始化
func (sys *System) NewDoc() *Doc {
	parsedAbi, _ := abi.JSON(strings.NewReader(DocAbiJSON))

	doc := new(Doc)
	doc.abi = parsedAbi
	doc.super = sys

	return doc
}

/*
  init:
   	EN -
	CN - did文档初始化
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- bidType: uint64，0: 普通用户,1:智能合约以及设备，2: 企业或者组织，BID类型一经设置，永不能变

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
*/
func (doc *Doc) Init(signTxParams *SysTxParams, bidType uint64) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("init", bidType)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  SetBidName:
   	EN -
	CN - 设置昵称
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- bidName: string 昵称(字符串长度6~20)

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
*/
func (doc *Doc) SetBidName(signTxParams *SysTxParams, id string, bidName string) (string, error) {
	if len(bidName) < 6 || len(bidName) > 20 {
		return "", errors.New("nickname's length must 6 to 20")
	}
	// encoding
	inputEncode, err := doc.abi.Pack("setBidName", id, bidName)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  GetDocument:
   	EN -
	CN - 查询文档的信息
  Params:
  	- isDidAddress: bool，如果是true则第二个参数传did的地址；否则传入bidName
	- did: string，did文档的地址或bidName

  Returns:
  	- *dto.Document
		Id              utils.Address `json:"id"` // bid
		Contexts        []byte        `json:"context"`
		Name            []byte        `json:"name"`            // bid标识符昵称
		Type            []byte        `json:"type"`            // bid的类型，包括0: 普通用户,1:智能合约以及设备，2: 企业或者组织，BID类型一经设置，永不能变
		PublicKeys      []byte        `json:"publicKeys"`      // 用户用于身份认证的公钥信息
		Authentications []byte        `json:"authentications"` // 用户身份认证列表信息
		Attributes      []byte        `json:"attributes"`      // 用户填写的个人信息值
		IsEnable        []byte        `json:"is_enable"`       // 该BID是否启用
		CreateTime      time.Time     `json:"createTime"`
		UpdateTime      time.Time     `json:"updateTime"`
	- error

  Call permissions: ？？
*/
func (doc *Doc) GetDocument(isDidAddress bool, did string) (*dto.Document, error) {
	params := make([]string, 2)
	if isDidAddress {
		params[0] = "0"
	} else {
		params[0] = "1"
	}
	params[1] = did

	pointer := &dto.SystemRequestResult{}

	err := doc.super.provider.SendRequest(pointer, "document_document", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToDocument()
}

/*
  AddPublic:
   	EN -
	CN - 增加用户did身份认证
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- publicType: string, 公钥类型（secp256k1、SM2）
	- publicAuth: string, 公钥权限（all、update、ban）
	- publicKey: string,  公钥（十六进制的字符串）

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ？？
*/
func (doc *Doc) AddPublic(signTxParams *SysTxParams, id string, publicType string, publicAuth string, publicKey string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("addPublic", id, publicType, publicAuth, publicKey)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  DelPublic:
   	EN -
	CN - 删除用户did身份认证
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- publicKey: string,   公钥

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
*/
func (doc *Doc) DelPublic(signTxParams *SysTxParams, id string, publicKey string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("delPublic", id, publicKey)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  AddAuth:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: ??
*/
func (doc *Doc) AddAuth(signTxParams *SysTxParams, id string, auth string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("addAuth", id, auth)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  DelAuth:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: ??
*/
func (doc *Doc) DelAuth(signTxParams *SysTxParams, id string, auth string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("delAuth", id, auth)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  AddService:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: Anyone
*/
func (doc *Doc) AddService(signTxParams *SysTxParams, id string, serviceId string, serviceType string, serviceEndpoint string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("addService", id, serviceId, serviceType, serviceEndpoint)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  DelService:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: Anyone
*/
func (doc *Doc) DelService(signTxParams *SysTxParams, id string, serviceId string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("delService", id, serviceId)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}
	return doc.super.sendRawTransaction(signedTx)
}

/*
  AddProof:
   	EN -
	CN - 增加证明
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- proofType: string, ??
	- proofCreator: string, ??
	- proofSign: string, ??

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
*/
func (doc *Doc) AddProof(signTxParams *SysTxParams, id string, proofType string, proofCreator string, proofSign string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("addProof", id, proofType, proofCreator, proofSign)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  DelProof:
   	EN -
	CN - 删除证明
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ？？
*/
func (doc *Doc) DelProof(signTxParams *SysTxParams, id string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("delProof", id)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  AddExtra:
   	EN -
	CN - 添加用户的基本信息
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string, bid
	- extra: string, ??

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
*/
func (doc *Doc) AddExtra(signTxParams *SysTxParams, id string, extra string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("addExtra", id, extra)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  DelExtra:
   	EN -
	CN - 删除用户的基本信息
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string, bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: Anyone
*/
func (doc *Doc) DelExtra(signTxParams *SysTxParams, id string) (string, error) {
	// encoding
	inputEncode, err := doc.abi.Pack("delExtra", id)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  Enable:
   	EN -
	CN - 使用户的Did身份可用
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
  	- id string bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
*/
func (doc *Doc) Enable(signTxParams *SysTxParams, id string) (string, error) {
	// encoding
	inputEncode, _ := doc.abi.Pack("enable", id)

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  Disable:
   	EN -
	CN - 使用户的Did身份不可用
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
  	- id string bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
*/
func (doc *Doc) Disable(signTxParams *SysTxParams, id string) (string, error) {
	// encoding
	inputEncode, _ := doc.abi.Pack("disable", id)

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocContractAddr)
	if err != nil {
		return "", err
	}

	return doc.super.sendRawTransaction(signedTx)
}

/*
  IsEnable:
   	EN -
	CN - 查询文档是否可用
  Params:
  	- isDidAddress: bool，如果是true则第二个参数传did的地址；否则传入bidName
	- did: string，did文档的地址或bidName

  Returns:
  	- bool, true可用，false不可用
	- error

  Call permissions: Anyone
*/
func (doc *Doc) IsEnable(isDidAddress bool, did string) (bool, error) {
	params := make([]string, 2)
	if isDidAddress {
		params[0] = "0"
	} else {
		params[0] = "1"
	}
	params[1] = did

	pointer := &dto.SystemRequestResult{}

	err := doc.super.provider.SendRequest(pointer, "document_isEnable", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}
