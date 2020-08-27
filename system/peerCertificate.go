package system

import (
	"bytes"
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils"
	"strings"
)

const (
	PeerCertificateContractAddr = "did:bid:ZFTc4jj4bH95pomT8JkGXFTxLnjbbj"
	Year                        = uint64(24) * 3600 * 365
	// regulatoryAddressLength     = 12 // 地址，除去did:bid:还有12字节
)

// 节点可信的AbiJson数据
// registerCertificate增加了apply
const PeerCertificateAbiJSON = `[
{"constant": false,"name":"registerCertificate","inputs":[{"name":"id","type":"string"},{"name":"apply","type":"string"},{"name":"publicKey","type":"string"},{"name":"nodeName","type":"string"},{"name":"messageSha3","type":"string"},{"name":"signature","type":"string"},{"name":"nodeType","type":"uint64"},{"name":"period","type":"uint64"},{"name":"ip","type":"string"},{"name":"port","type":"uint64"},{"name":"companyName","type":"string"},{"name":"companyCode","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"peerEvent","type":"event"}
]`

// PeerCertificate - The PeerCertificate Module
type PeerCertificate struct {
	super *System
	abi   abi.ABI
}

func (peerCer *PeerCertificate) peerRegisterPreCheck(registerCertificate *dto.RegisterCertificateInfo) (bool, error) {
	// 查验参数是否输入合法
	if !isValidHexAddress(registerCertificate.Id) {
		return false, errors.New("registerCertificate id is not valid bid")
	}
	if !isValidHexAddress(registerCertificate.Apply) {
		return false, errors.New("registerCertificate Apply is not valid bid")
	}
	if len(registerCertificate.PublicKey) != 53  || isBlankCharacter(registerCertificate.PublicKey) {
		return false, errors.New("registerCertificate publicKey's len should be 53")
	}
	if len(registerCertificate.NodeName) == 0 || isBlankCharacter(registerCertificate.NodeName) {
		return false, errors.New("registerCertificate NodeName can't be empty or blank character")
	}
	if registerCertificate.NodeType != 0 && registerCertificate.NodeType != 1 {
		return false, errors.New("registerCertificate NodeType should be 0 or 1")
	}
	if !isLegalIP(registerCertificate.IP) {
		return false, errors.New("registerCertificate IP is illegal")
	}
	if registerCertificate.Port > 65535 {
		return false, errors.New("registerCertificate Port should be in range 0 to 65535")
	}
	if len(registerCertificate.CompanyName) == 0 || isBlankCharacter(registerCertificate.NodeName) {
		return false, errors.New("registerCertificate CompanyName can't be empty or blank character")
	}
	if len(registerCertificate.CompanyCode) == 0 || isBlankCharacter(registerCertificate.NodeName) {
		return false, errors.New("registerCertificate CompanyCode can't be empty or blank character")
	}

	return true, nil
}

// NewPeerCertificate - 初始化PeerCertificate
func (sys *System) NewPeerCertificate() *PeerCertificate {
	parseAbi, _ := abi.JSON(strings.NewReader(PeerCertificateAbiJSON))

	peerCertificate := new(PeerCertificate)
	peerCertificate.super = sys
	peerCertificate.abi = parseAbi
	return peerCertificate
}

func (peerCer *PeerCertificate) messageSignature(message, password string, keyFileData []byte, isSM2 bool) (string, string, error) {
	messageSha3 := utils.NewUtils().Sha3Raw(message)

	_, privateKey, err := peerCer.super.acc.Decrypt(keyFileData, isSM2, password)
	privKey, err := crypto.HexToECDSA(privateKey, crypto.SECP256K1)
	if err != nil {
		return "", "", err
	}

	var cryptoType crypto.CryptoType
	var t string
	if isSM2 {
		t = "00"
		cryptoType = crypto.SM2
	} else {
		cryptoType = crypto.SECP256K1
		t = "01"
	}

	messageSha3Bytes := utils.Hex2Bytes(messageSha3[2:])
	sig, err := crypto.Sign(messageSha3Bytes, privKey, cryptoType)
	if err != nil {
		return "", "", err
	}
	// r := new(big.Int).SetBytes(sig[:32])
	// s := new(big.Int).SetBytes(sig[32:64])
	// v := new(big.Int).SetBytes([]byte{sig[64] + 27})
	// fmt.Printf("r %x \n", r)
	// fmt.Printf("s %x \n", s)
	// fmt.Printf("v %x \n", v)
	// fmt.Printf("sig len is  %x \n", len(sig))
	var buf bytes.Buffer
	buf.Write([]byte{sig[64] + 27})
	buf.Write(sig[:64])
	// fmt.Printf("sig is  %s \n", t+utils.Bytes2Hex(buf.Bytes()))
	return messageSha3, t + utils.Bytes2Hex(buf.Bytes()), err
}

/*
  RegisterCertificate:
   	EN -
	CN - 为节点颁发可信证书可信
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- registerCertificate:  *dto.RegisterCertificateInfo，包含可信证书的信息
		Id          string // 节点证书的bid，必须和public_key对应，索引
		Apply       string // 申请人的bid（与Id可以相同，可以不同）
		PublicKey   string // 53个字符的公钥，也就是p2p节点id的形式
		NodeName    string // 节点名称，不含敏感词的字符串
		MessageSha3 string // 消息sha3后的16进制字符串，用于本地签名和链上验证签名，该字段不会被链保存
		Signature   string // 对上一个字段消息的签名，16进制字符串
		NodeType    uint64 // 节点类型，0企业，1个人
		Period      uint64 // 证书有效期，以年为单位的整型
		IP          string // 节点间互连的ip
		Port        uint64 // 节点间互连的端口
		CompanyName string // 公司名（如果是个人，则是个人姓名），不含敏感词的字符串
		CompanyCode string // 公司代码，不含敏感词的字符串

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以调用
*/
func (peerCer *PeerCertificate) RegisterCertificate(signTxParams *SysTxParams, registerCertificate *dto.RegisterCertificateInfo, idPassword string, idKeyFile []byte, idIsSM2 bool) (string, error) {
	ok, err := peerCer.peerRegisterPreCheck(registerCertificate)
	if !ok {
		return "", err
	}

	// encoding
	// registerCertificate is a struct we need to use the components.
	messageSha3, signature, err := peerCer.messageSignature("test", idPassword, idKeyFile, idIsSM2)
	if err != nil {
		return "", err
	}
	// fmt.Printf("messageSha3 %s, signature  %s \n", messageSha3, signature)
	// 添加messageSha3 signature
	registerCertificate.MessageSha3 = messageSha3
	registerCertificate.Signature = signature

	var values []interface{}
	values = peerCer.super.structToInterface(*registerCertificate, values)
	inputEncode, err := peerCer.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := peerCer.super.prePareSignTransaction(signTxParams, inputEncode, PeerCertificateContractAddr)
	if err != nil {
		return "", err
	}

	return peerCer.super.sendRawTransaction(signedTx)
}

/*
  RevokedCertificate:
   	EN -
	CN - 吊销节点的可信证书可信
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string，节点证书的bid

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以调用
*/
func (peerCer *PeerCertificate) RevokedCertificate(signTxParams *SysTxParams, id string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid address")
	}

	// encoding
	inputEncode, err := peerCer.abi.Pack("revokedCertificate", id)
	if err != nil {
		return "", err
	}

	signedTx, err := peerCer.super.prePareSignTransaction(signTxParams, inputEncode, PeerCertificateContractAddr)
	if err != nil {
		return "", err
	}

	return peerCer.super.sendRawTransaction(signedTx)
}

/*
  GetPeriod:
   	EN -
	CN - 查看证书有效期
  Params:
  	- id: string，节点证书的bid

  Returns:
  	- uint64，返回证书有效期，如果证书被吊销，则有效期是0
	- error

  Call permissions: Anyone
*/
func (peerCer *PeerCertificate) GetPeriod(id string) (uint64, error) {
	if !isValidHexAddress(id) {
		return 0, errors.New("id is not valid bid address")
	}

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
  GetActive:
   	EN -
	CN - 查看证书是否有效
  Params:
  	- id: string，节点证书的bid

  Returns:
  	- bool，true可用，false不可用
	- error

  Call permissions: Anyone
*/
func (peerCer *PeerCertificate) GetActive(id string) (bool, error) {
	if !isValidHexAddress(id) {
		return false, errors.New("id is not valid bid address")
	}

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
  GetPeerCertificate:
   	EN -
	CN - 查看证书的信息
  Params:
  	- id: string，节点证书的bid

  Returns:
  	- *dto.PeerCertificate
		Id          string   `json:"id"`          //唯一索引
		Issuer      string   `json:"issuer"`      //颁发者地址
		Apply       string   `json:"apply"`       //申请人bid
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
func (peerCer *PeerCertificate) GetPeerCertificate(id string) (dto.PeerCertificate, error) {
	var peerCertificate dto.PeerCertificate
	if !isValidHexAddress(id) {
		return peerCertificate, errors.New("id is not valid bid address")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_peerCertificate", params)
	if err != nil {
		return peerCertificate, err
	}

	res, err := pointer.ToPeerCertificate()
	if err != nil{
		return peerCertificate, err
	}
	return *res, nil
}

/*
  GetPeerCertificateIdList:
   	EN - Get applied certificates by bid
 	CN - 根据节点可信证书申请人的bid获取申请的证书列表
  Params:
  	- id: string，节点可信证书申请人（Apply）bid

  Returns:
  	- []string, 申请人申请的证书列表
 	- error

  Call permissions: Anyone
*/
func (peerCer *PeerCertificate) GetPeerCertificateIdList(id string) ([]string, error) {
	if !isValidHexAddress(id) {
		return nil, errors.New("id is not valid bid address")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := peerCer.super.provider.SendRequest(pointer, "peercertificate_peerCertificateIdList", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}
