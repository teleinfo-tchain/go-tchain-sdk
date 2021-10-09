package system

import (
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils"
	"strings"
)



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

  Call permissions: 只能初始自己的文档，且只能初始化一次
*/
func (doc *Doc) Init(signTxParams *SysTxParams, bidType uint64) (string, error) {
	if bidType != 0 && bidType != 1 && bidType != 2 {
		return "", errors.New("bidType should be 0, 1, or 2")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("init", bidType)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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

  Call permissions: 权限为`all`
*/
func (doc *Doc) SetBidName(signTxParams *SysTxParams, id string, bidName string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}

	if len(bidName) < 6 || len(bidName) > 20 || isBlankCharacter(bidName) {
		return "", errors.New("bidName's length must be 6 to 20 and can't be blank character")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("setBidName", id, bidName)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
	- did: string，did文档的地址或bidName

  Returns:
  	- dto.Document
		Id              string       `json:"id"` // bid
		Contexts        string       `json:"context"`
		Name            string       `json:"name"`           // bid标识符昵称
		Type            uint64       `json:"type"`           // bid的类型，包括0：普通用户,1:智能合约以及设备，2：企业或者组织，BID类型一经设置，永不能变
		PublicKeys      []*PublicKey `json:"publicKey"`      // 用户用于身份认证的公钥信息
		Authentications []string     `json:"authentication"` // 用户身份认证列表信息
		Services        []*Service   `json:"service"`        // 用户填写的服务端点信息
		Proof           *Proof       `json:"proof"`          // 用户填写的证明信息值
		Extra           string       `json:"extra"`          // 用户填写的备注
		IsEnable        bool         `json:"isEnable"`       // 该BID是否启用
		CreateTime      string       `json:"created"`
		UpdateTime      string       `json:"updated"`
	- error

  Call permissions: Anyone
  todo:文档初始化的时间采用的是0时区的，所以时间会差8小时，这个后期根据需要再修改
*/
func (doc *Doc) GetDocument(did string) (dto.Document, error) {
	var document dto.Document
	if len(did) == 0 || isBlankCharacter(did) {
		return document, errors.New("did can't be empty or blank character")
	}

	params := make([]string, 1)
	params[0] = did

	pointer := &dto.SystemRequestResult{}

	err := doc.super.provider.SendRequest(pointer, "document_document", params)
	if err != nil {
		return document, err
	}

	res, err := pointer.ToDocument()
	if err != nil{
		return document, err
	}

	return *res, nil
}

/*
  AddPublic:
   	EN -
	CN - 增加用户did身份认证
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- publicType: string, 公钥类型（随机字符串）
	- publicAuth: string, 公钥权限（通常用all, update, ban）
	- publicKey: string,  公钥（十六进制的字符串(130或带0x前缀的132)）

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: 权限为`all`
*/
func (doc *Doc) AddPublic(signTxParams *SysTxParams, id string, publicType string, publicAuth string, publicKey string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}
	if len(publicType) == 0 || isBlankCharacter(publicType) {
		return "", errors.New("publicType can't be empty or blank character")
	}
	if publicAuth != "all" && publicAuth != "update" && publicAuth != "ban" {
		return "", errors.New("publicAuth should be all 、update or ban")
	}

	if utils.Has0xPrefix(publicKey) {
		publicKey = publicKey[2:]
	}
	// 检测公钥合法性
	if !(utils.IsHex(publicKey) && len(publicKey) == 130) {
		return "", errors.New("publicKey is not a hexadecimal string or the length is less than 130(132 with prefix '0x'")
	}
	if !(publicKey[:2] == "01" || publicKey[:2] == "02" || publicKey[:2] == "04") {
		return "", errors.New("publicKey should be with prefix 01、 02 、04 or 0x01、 0x02、 0x04")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("addPublic", id, publicType, publicAuth, publicKey)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
	- publicKey: string,   公钥（十六进制的字符串(130或带0x前缀的132)）

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: 权限为`all`
*/
func (doc *Doc) DelPublic(signTxParams *SysTxParams, id string, publicKey string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}
	if utils.Has0xPrefix(publicKey) {
		publicKey = publicKey[2:]
	}

	// 检测公钥合法性
	if !(utils.IsHex(publicKey) && len(publicKey) == 130) {
		return "", errors.New("publicKey is not a hexadecimal string or the length is less than 130(132 with prefix '0x'")
	}
	if !(publicKey[:2] == "01" || publicKey[:2] == "02" || publicKey[:2] == "04") {
		return "", errors.New("publicKey should be with prefix 01、 02 、04 or 0x01、 0x02、 0x04")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("delPublic", id, publicKey)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- auth: string,  权限(随机字符串)

  Returns:
	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: 权限为all、update
*/
func (doc *Doc) AddAuth(signTxParams *SysTxParams, id string, auth string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}
	if len(auth) == 0 || isBlankCharacter(auth) {
		return "", errors.New("auth can't be empty or blank character")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("addAuth", id, auth)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- auth: string,  权限(随机字符串)

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: 权限为all、update
*/
func (doc *Doc) DelAuth(signTxParams *SysTxParams, id string, auth string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}
	if len(auth) == 0 || isBlankCharacter(auth) {
		return "", errors.New("auth can't be empty or blank character")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("delAuth", id, auth)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string             bid地址
	- serviceId: string		 bid地址
	- serviceType: string,  服务类型(随机字符串)
	- serviceEndpoint: string, 服务url(随机字符串)

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: 权限为all、update
*/
func (doc *Doc) AddService(signTxParams *SysTxParams, id string, serviceId string, serviceType string, serviceEndpoint string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}
	if !isValidHexAddress(serviceId) {
		return "", errors.New("serviceId is not valid bid")
	}
	if len(serviceType) == 0 || isBlankCharacter(serviceType) {
		return "", errors.New("serviceType can't be empty or blank character")
	}
	if len(serviceEndpoint) == 0 || isBlankCharacter(serviceEndpoint) {
		return "", errors.New("serviceEndpoint can't be empty or blank character")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("addService", id, serviceId, serviceType, serviceEndpoint)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string             bid地址
	- serviceId: string		 bid地址

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: 权限为all、update
*/
func (doc *Doc) DelService(signTxParams *SysTxParams, id string, serviceId string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}
	if !isValidHexAddress(serviceId) {
		return "", errors.New("serviceId is not valid bid")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("delService", id, serviceId)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
	- proofType: string, 证书类型（随机字符串）
	- proofCreator: string, 证书创建人（随机字符串）
	- proofSign: string, 证书的签名（随机字符串）

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: 权限为all、update
*/
func (doc *Doc) AddProof(signTxParams *SysTxParams, id string, proofType string, proofCreator string, proofSign string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}
	if len(proofType) == 0 || isBlankCharacter(proofType) {
		return "", errors.New("proofType can't be empty or blank character")
	}
	if len(proofCreator) == 0 || isBlankCharacter(proofCreator) {
		return "", errors.New("proofCreator can't be empty or blank character")
	}
	if len(proofSign) == 0 || isBlankCharacter(proofSign) {
		return "", errors.New("proofSign can't be empty or blank character")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("addProof", id, proofType, proofCreator, proofSign)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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

  Call permissions: 权限为all、update
*/
func (doc *Doc) DelProof(signTxParams *SysTxParams, id string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("delProof", id)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
	- extra: string, 随机字符串

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: publicAuth权限为all、update
*/
func (doc *Doc) AddExtra(signTxParams *SysTxParams, id string, extra string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}

	if len(extra) == 0 || isBlankCharacter(extra) {
		return "", errors.New("extra can't be empty or blank character")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("addExtra", id, extra)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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

  Call permissions: publicAuth权限为all、update
*/
func (doc *Doc) DelExtra(signTxParams *SysTxParams, id string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}

	// encoding
	inputEncode, err := doc.abi.Pack("delExtra", id)
	if err != nil {
		return "", err
	}

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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

  Call permissions: 自身调用或者拥有相关权限的其他地址
*/
func (doc *Doc) Enable(signTxParams *SysTxParams, id string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}

	// encoding
	inputEncode, _ := doc.abi.Pack("enable", id)

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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

  Call permissions: 自身调用或者拥有相关权限的其他地址
*/
func (doc *Doc) Disable(signTxParams *SysTxParams, id string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}

	// encoding
	inputEncode, _ := doc.abi.Pack("disable", id)

	signedTx, err := doc.super.prePareSignTransaction(signTxParams, inputEncode, DocumentContract)
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
	- did: string，did文档的地址或bidName

  Returns:
  	- bool, true可用，false不可用
	- error

  Call permissions: Anyone
*/
func (doc *Doc) IsEnable(did string) (bool, error) {
	if len(did) == 0 || isBlankCharacter(did) {
		return false, errors.New("did can't be empty or blank character")
	}

	params := make([]string, 1)
	params[0] = did

	pointer := &dto.SystemRequestResult{}

	err := doc.super.provider.SendRequest(pointer, "document_isEnable", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}
