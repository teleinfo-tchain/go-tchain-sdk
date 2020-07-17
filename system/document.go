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

 Params：bidType类型，包括0：普通用户,1:智能合约以及设备，2：企业或者组织，BID类型一经设置，永不能变

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

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

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？
*/
func (doc *Doc) SetBidName(from common.Address, bidName string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("SetBidName", bidName)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
 GetDocument: 查询did文档内容

 Returns： *dto.Document
*/
func (doc *Doc) GetDocument(key uint64, value string) (*dto.Document, error) {
	params := make([]interface{}, 2)
	params[0] = key
	params[1] = value

	pointer := &dto.RequestResult{}

	err := doc.super.provider.SendRequest(pointer, "document_document", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToDocument()
}

/*
AddPublicKey: 增加用户did身份认证

	Params：
	- addType： bid的类型，包括0：普通用户,1:智能合约以及设备，2：企业或者组织？？？ 和初始化的是否类似？能否变更
	- authority：？？
	- publicKey：

Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

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

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

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

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

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

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

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

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？
*/
func (doc *Doc) AddAttribute(from common.Address, AttrType string, value string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("AddAttr", AttrType, value)

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
 DeleteAttribute: 删除用户的基本信息

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

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

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

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

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？
*/
func (doc *Doc) Disable(from common.Address) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("Disable")

	transaction := doc.super.prePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.sendTransaction(transaction)
}

/*
 IsEnable: 查询did文档是否可用

 Returns： bool, true可用，false不可用
*/
func (doc *Doc) IsEnable(key uint64, value string) (bool, error) {
	params := make([]interface{}, 2)
	params[0] = key
	params[1] = value

	pointer := &dto.RequestResult{}

	err := doc.super.provider.SendRequest(pointer, "", params)
	if err != nil {
		return false, err
	}

	return pointer.ToDocIsEnable()
}
