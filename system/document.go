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
	DocContractAddr = "did:bid:00000000000000000000000a"
)

//const (
//	bid = iota
//	publicKeyPem
//	publicKeyJwk
//	publicKeyHex
//	publicKeyBase64
//	publicKeyBase58
//	publicKeyMultibase
//	ethereumAddress
//)

//var (
//	ErrBidDDOExistsInvalid = errors.New("the DDO document for this bid has been initialized")
//	ErrNotHaveAuthority    = errors.New("对不起，您没有权限做该操作")
//)

// did文档的AbiJson数据
const DocAbiJSON = `[
{"constant": false,"name":"InitializationDDO","inputs":[{"name":"bidType","type":"uint64"}],"outputs":[],"type":"function"},
{"constant": false,"name":"SetBidName","inputs":[{"name":"bidName","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"AddPublicKey","inputs":[{"name":"type","type":"string"},{"name":"authority","type":"string"},{"name":"publicKey","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"DeletePublicKey","inputs":[{"name":"publicKey","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"AddProof","inputs":[{"name":"type","type":"string"},{"name":"proofID","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"DeleteProof","inputs":[{"name":"proofID","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"AddAttr","inputs":[{"name":"type","type":"string"},{"name":"value","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"DeleteAttr","inputs":[{"name":"type","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"Enable","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"Disable","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"bidEvent","type":"event"}]`

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
InitializationDDO: did文档初始化

Params:
	- from: [20]byte，交易发送方地址
	- bidType: uint64，0: 普通用户,1:智能合约以及设备，2: 企业或者组织，BID类型一经设置，永不能变

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？
*/
func (doc *Doc) InitializationDDO(from common.Address, bidType uint64) (string, error) {
	// encoding
	inputEncode, _ := doc.abi.Pack("InitializationDDO", bidType)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
SetBidName: 设置bid标识符昵称(bidName)

Params:
	- from: [20]byte，交易发送方地址
	- bidName: string

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？
*/
func (doc *Doc) SetBidName(from common.Address, bidName string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("SetBidName", bidName)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
GetDocument: 查询文档的信息

Params:
	- id: string，文档的bid或文档设置的名字

Returns:
	- *dto.Document
		Id              common.Address `json:"id"`              //bid
		Contexts        []byte         `json:"context"`
		Name            []byte         `json:"name"`            //bid标识符昵称
		Type            []byte         `json:"type"`            //bid的类型，包括0: 普通用户,1:智能合约以及设备，2: 企业或者组织，BID类型一经设置，永不能变
		PublicKeys      []byte         `json:"publicKeys"`      //用户用于身份认证的公钥信息
		Authentications []byte         `json:"authentications"` //用户身份认证列表信息
		Attributes      []byte         `json:"attributes"`      //用户填写的个人信息值
		IsEnable        []byte         `json:"is_enable"`       //该BID是否启用
		CreateTime      time.Time      `json:"createTime"`
		UpdateTime      time.Time      `json:"updateTime"`
	- error

Call permissions: Anyone
*/
func (doc *Doc) GetDocument(id string) (*dto.Document, error) {
	params := make([]interface{}, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := doc.super.provider.SendRequest(pointer, "document_document", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToDocument()
}

/*
AddPublicKey: 增加用户did身份认证

Params:
	- from: [20]byte，交易发送方地址
	- addType: string,bid的类型，包括0: 普通用户,1:智能合约以及设备，2: 企业或者组织？？？ 和初始化的是否类似？能否变更
	- authority: string,？？
	- publicKey: string,

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) AddPublicKey(from common.Address, addType string, authority string, publicKey string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("AddPublicKey", addType, authority, publicKey)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
DeletePublicKey: 删除用户did身份认证

Params:
	- from: [20]byte，交易发送方地址
	- publicKey: string, ??

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) DeletePublicKey(from common.Address, publicKey string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("DeletePublicKey", publicKey)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
AddProof: 增加证明

Params:
	- from: [20]byte，交易发送方地址
	- issuer: string, ??
	- proofID: string, ??

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) AddProof(from common.Address, issuer string, proofID string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("AddProof", issuer, proofID)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
DeleteProof: 删除证明

Params:
	- from: [20]byte，交易发送方地址
	- proofID: string, ??

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) DeleteProof(from common.Address, proofID string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("DeleteProof", proofID)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
AddAttribute: 添加用户的基本信息

Params:
	- from: [20]byte，交易发送方地址
	- attrType: string, ??
	- value: string, ??

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) AddAttribute(from common.Address, attrType string, value string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("AddAttr", attrType, value)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
DeleteAttribute: 删除用户的基本信息

Params:
	- from: [20]byte，交易发送方地址
	- addType: string, ??
	- value: string, ??

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) DeleteAttribute(from common.Address, addType string, value string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("DeleteAttr", addType, value)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
Enable: 使用户的Did身份可用

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) Enable(from common.Address) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("Enable")

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
Disable: 使用户的Did身份不可用

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

Call permissions: ？？
*/
func (doc *Doc) Disable(from common.Address) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("Disable")

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
IsEnable: 查询文档是否可用

Params:
	- id: string，文档的bid或文档设置的名字

Returns:
	- bool, true可用，false不可用
	- error

Call permissions: Anyone
*/
func (doc *Doc) IsEnable(id string) (bool, error) {
	params := make([]interface{}, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := doc.super.provider.SendRequest(pointer, "", params)
	if err != nil {
		return false, err
	}

	return pointer.ToDocIsEnable()
}
