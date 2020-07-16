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

type doc struct {
	super *System
	abi abi.ABI
}

func (sys *System) NewDoc() *doc {
	parsedAbi, _ := abi.JSON(strings.NewReader(DocAbiJSON))

	doc := new(doc)
	doc.abi = parsedAbi
	doc.super = sys

	return doc
}

//"InitializationDDO","inputs":[{"name":"bidType","type":"uint64"}],"outputs":[]
func (doc *doc) InitializationDDO(from common.Address, bidType uint64)(string, error){
	// encoding
	inputEncode, _ := doc.abi.Pack("InitializationDDO", bidType)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"SetBidName","inputs":[{"name":"bidName","type":"string"}],"outputs":[]
func (doc *doc) SetBidName(from common.Address, bidName string) (string, error) {
	//encoding
	inputEncode, _ := doc.abi.Pack("SetBidName", bidName)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//did文档内容
func (doc * doc) GetDocument(key uint64, value string)(*dto.Document, error){
	params := make([]interface{}, 2)
	params[0] = key
	params[1] = value

	pointer := &dto.RequestResult{}

	err := doc.super.provider.SendRequest(pointer, "document_document", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToDocument()
}

//"AddPublicKey","inputs":[{"name":"type","type":"string"},{"name":"authority","type":"string"},{"name":"publicKey","type":"string"}],"outputs":[]
func(doc *doc) AddPublicKey(from common.Address, addType string, authority string, publicKey string)(string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("AddPublicKey", addType, authority, publicKey)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"DeletePublicKey","inputs":[{"name":"publicKey","type":"string"}],"outputs":[]
func(doc *doc) DeletePublicKey(from common.Address, publicKey string) (string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("DeletePublicKey", publicKey)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"AddProof","inputs":[{"name":"type","type":"string"},{"name":"proofID","type":"string"}],"outputs":[]
func(doc *doc) AddProof(from common.Address, addType string, proofID string) (string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("AddProof", addType, proofID)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"DeleteProof","inputs":[{"name":"proofID","type":"string"}],"outputs":[]
func(doc *doc) DeleteProof(from common.Address, proofID string) (string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("DeleteProof", proofID)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"AddAttr","inputs":[{"name":"type","type":"string"},{"name":"value","type":"string"}],"outputs":[]
func(doc *doc) AddAttribute(from common.Address, addType string, value string) (string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("AddAttr", addType, value)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"DeleteAttr","inputs":[{"name":"type","type":"string"}],"outputs":[]
func(doc *doc) DeleteAttribute(from common.Address, addType string, value string) (string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("DeleteAttr", addType, value)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"Enable","inputs":[],"outputs":[]
func(doc *doc) Enable(from common.Address) (string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("Enable")

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//"Disable","inputs":[],"outputs":[]
func(doc *doc) Disable(from common.Address) (string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("Disable")

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return doc.super.SendTransaction(transaction)
}

//did文档是否可用
func(doc *doc) IsEnable(key uint64, value string) (bool, error){
	params := make([]interface{}, 2)
	params[0] = key
	params[1] = value

	pointer := &dto.RequestResult{}

	err := doc.super.provider.SendRequest(pointer, "", params)
	if err != nil{
		return false, err
	}

	return pointer.ToDocIsEnable()
}