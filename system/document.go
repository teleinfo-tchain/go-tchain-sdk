package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"strings"
	"time"
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
{"constant": true,"name":"FindDDOByType","inputs":[{"name":"key","type":"uint64"},{"name":"value","type":"string"}],"outputs":[{"name":"result","type":"string"}],"type":"function"},
{"constant": false,"name":"AddPublicKey","inputs":[{"name":"type","type":"string"},{"name":"authority","type":"string"},{"name":"publicKey","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"DeletePublicKey","inputs":[{"name":"publicKey","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"AddProof","inputs":[{"name":"type","type":"string"},{"name":"proofID","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"DeleteProof","inputs":[{"name":"proofID","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"AddAttr","inputs":[{"name":"type","type":"string"},{"name":"value","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"DeleteAttr","inputs":[{"name":"type","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"Enable","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"Disable","inputs":[],"outputs":[],"type":"function"},
{"constant": true,"name":"IsEnable","inputs":[{"name":"key","type":"uint64"},{"name":"value","type":"string"}],"outputs":[{"name":"result","type":"bool"}],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"bidEvent","type":"event"}]`

type doc struct {
	super *System
	abi abi.ABI
}

type Document struct {
	Id              common.Address `json:"id"` //bid
	Contexts        []byte         `json:"context"`
	Name            []byte         `json:"name"`            //bid标识符昵称
	Type            []byte         `json:"type"`            // bid的类型，包括0：普通用户,1:智能合约以及设备，2：企业或者组织，BID类型一经设置，永不能变
	PublicKeys      []byte         `json:"publicKeys"`      //用户用于身份认证的公钥信息
	Authentications []byte         `json:"authentications"` //用户身份认证列表信息
	Attributes      []byte         `json:"attributes"`      //用户填写的个人信息值
	IsEnable        []byte         `json:"is_enable"`       //该BID是否启用
	CreateTime      time.Time      `json:"createTime"`
	UpdateTime      time.Time      `json:"updateTime"`
}

type PublicKey struct {
	Id         common.Address `json:"id"`
	KeyId      []byte         `json:"key_id"`
	Type       []byte         `json:"type"`
	Controller []byte         `json:"controller"`
	Authority  []byte         `json:"authority"` //公钥权限
	PublicKey  []byte         `json:"publicKey"`
}

type Authentication struct {
	Id        common.Address `json:"id"`
	ProofId   []byte         `json:"proofId"`
	Issuer    common.Address `json:"type"`
	PublicKey []byte         `json:"public_key"`
}

type Attribute struct {
	Id       common.Address `json:"id"`
	AttrType []byte         `json:"attr_type"`
	Value    []byte         `json:"value"`
}

func (sys *System) NewDoc()(*doc, error){
	parsedAbi, err := abi.JSON(strings.NewReader(DocAbiJSON))
	if err != nil {
		return nil, err
	}

	doc := new(doc)
	doc.abi = parsedAbi
	doc.super = sys

	return doc, nil
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

//"FindDDOByType","inputs":[{"name":"key","type":"uint64"},{"name":"value","type":"string"}],"outputs":[{"name":"result","type":"string"}]
func (doc *doc) FindDDOByType(from common.Address, key uint64, value string)(string, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("FindDDOByType", key, value)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	requestResult, err := doc.super.Call(transaction)

	if err != nil {
		return "", err
	}

	//fmt.Println("result is ", requestResult.Result.(string))
	var result string
	err = doc.abi.Methods["FindDDOByType"].Outputs.Unpack(&result, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return "", err
	}

	return result, err
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

//"IsEnable","inputs":[{"name":"key","type":"uint64"},{"name":"value","type":"string"}],"outputs":[{"name":"result","type":"bool"}]
func(doc *doc) IsEnable(from common.Address, key uint64, value string) (bool, error){
	//encoding
	inputEncode, _ := doc.abi.Pack("IsEnable", key, value)

	transaction := doc.super.PrePareTransaction(from, DocContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	requestResult, err := doc.super.Call(transaction)

	if err != nil {
		return false, err
	}

	//fmt.Println("result is ", requestResult.Result.(string))
	var result bool
	err = doc.abi.Methods["IsEnable"].Outputs.Unpack(&result, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return false, err
	}

	return result, err
}