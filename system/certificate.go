package system

import (
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils"
	"strings"
)

// Certificate - The Certificate Module
type Certificate struct {
	super *System
	abi   abi.ABI
}

// NewCertificate - NewCertificate初始化
func (sys *System) NewCertificate() *Certificate {
	parsedAbi, _ := abi.JSON(strings.NewReader(CertificateAbiJSON))

	cer := new(Certificate)
	cer.abi = parsedAbi
	cer.super = sys
	return cer
}

func (cer *Certificate) registerCertificatePreCheck(registerCertificate dto.RegisterCertificate) (bool, error) {
	if !isValidHexAddress(registerCertificate.Id) {
		return false, errors.New("registerCertificate Id is not valid bid")
	}

	if len(registerCertificate.Context) == 0 || isBlankCharacter(registerCertificate.Context) {
		return false, errors.New("registerCertificate Context can't be empty")
	}

	if !isValidHexAddress(registerCertificate.Subject) {
		return false, errors.New("registerCertificate Subject is not valid bid")
	}

	// 检测公钥，且判断是否与subject的地址相匹配
	ok, err := account.CheckPublicKeyToAccount(registerCertificate.Subject, registerCertificate.SubjectPublicKey)
	if !ok {
		return false, err
	}


	if len(registerCertificate.IssuerAlgorithm) == 0 || isBlankCharacter(registerCertificate.IssuerAlgorithm) {
		return false, errors.New("registerCertificate IssuerAlgorithm can't be empty")
	}

	if len(registerCertificate.IssuerSignature) == 0 || isBlankCharacter(registerCertificate.IssuerSignature) {
		return false, errors.New("registerCertificate IssuerSignature can't be empty")
	}

	if len(registerCertificate.SubjectAlgorithm) == 0 || isBlankCharacter(registerCertificate.SubjectAlgorithm) {
		return false, errors.New("registerCertificate SubjectAlgorithm can't be empty")
	}

	if len(registerCertificate.SubjectSignature) == 0 || isBlankCharacter(registerCertificate.SubjectSignature) {
		return false, errors.New("registerCertificate SubjectSignature can't be empty")
	}
	return true, nil
}

/*
  RegisterCertificate:
   	EN -
	CN - 信任锚颁发证书，如果是根信任锚颁发的证书，则证书接收者可以进行部署合约和大额转账操作
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- registerCertificate:  The registerCertificate object(*dto.RegisterCertificate)
		Id               string // 个人可信证书bid
		Context          string // 证书上下文环境，随便一个字符串，不验证
		Subject          string // 证书接收者的bid，证书是颁给谁的
		Period           uint64 // 证书有效期，以年为单位的整型
		IssuerAlgorithm  string // 颁发者签名算法，字符串
		IssuerSignature  string // 颁发者签名值，16进制字符串
		SubjectPublicKey string // 接收者公钥，16进制字符串
		SubjectAlgorithm string // 接收者签名算法，字符串
		SubjectSignature string // 接收者签名值，16进制字符串

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有信任锚地址可以调用
*/
func (cer *Certificate) RegisterCertificate(signTxParams *SysTxParams, registerCertificate *dto.RegisterCertificate) (string, error) {
	ok, err := cer.registerCertificatePreCheck(*registerCertificate)
	if !ok {
		return "", err
	}

	if utils.Has0xPrefix(registerCertificate.SubjectPublicKey){
		registerCertificate.SubjectPublicKey = registerCertificate.SubjectPublicKey[2:]
	}

	// registerCertificate is a struct we need to use the components.
	var values []interface{}
	values = cer.super.structToInterface(*registerCertificate, values)
	inputEncode, err := cer.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := cer.super.prePareSignTransaction(signTxParams, inputEncode, CertificateContract)
	if err != nil {
		return "", err
	}

	return cer.super.sendRawTransaction(signedTx)
}

/*
  RevokedCertificate:
   	EN -
	CN - 信任锚吊销个人证书
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id:  string，个人可信证书bid

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有证书颁发者可以调用
*/
func (cer *Certificate) RevokedCertificate(signTxParams *SysTxParams, id string) (string, error) {
	if !isValidHexAddress(id) {
		return "", errors.New("id is not valid bid")
	}

	// encoding
	inputEncode, err := cer.abi.Pack("revokedCertificate", id)
	if err != nil {
		return "", err
	}

	signedTx, err := cer.super.prePareSignTransaction(signTxParams, inputEncode, CertificateContract)
	if err != nil {
		return "", err
	}

	return cer.super.sendRawTransaction(signedTx)
}

/*
  RevokedCertificates:
   	EN -
	CN - 信任锚批量吊销个人证书，把自己颁发的证书全部吊销
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有证书颁发者可以调用
*/
func (cer *Certificate) RevokedCertificates(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := cer.abi.Pack("revokedCertificates")

	signedTx, err := cer.super.prePareSignTransaction(signTxParams, inputEncode, CertificateContract)
	if err != nil {
		return "", err
	}

	return cer.super.sendRawTransaction(signedTx)
}

/*
  GetPeriod:
   	EN -
	CN -  查询个人证书的有效期
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- uint64，证书有效期，如果证书被吊销，则有效期是0
	- error

  Call permissions: Anyone
*/
func (cer *Certificate) GetPeriod(id string) (uint64, error) {
	if !isValidHexAddress(id) {
		return 0, errors.New("id is not valid bid")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_period", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
  GetActive:
   	EN -
	CN -  查询证书是否可用,如果证书过期，则不可用；如果信任锚注销，则不可用
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- bool，true可用，false不可用
	- error

  Call permissions: Anyone
*/
func (cer *Certificate) GetActive(id string) (bool, error) {
	if !isValidHexAddress(id) {
		return false, errors.New("id is not valid bid")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_active", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
  GetCertificate:
   	EN -
	CN -  查询证书的信息，period和active没有进行逻辑判断
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- dto.CertificateInfo
		Id             string // 凭证的hash
		Context        string // 证书所属上下文环境
		Issuer         string // 信任锚的bid
		Subject        string // 证书拥有者地址
		IssuedTime     uint64 // 颁发时间
		Period         uint64 // 有效期
		IsEnable       bool   // true 凭证有效，false 凭证已撤销
		RevocationTime uint64 // 吊销时间
	- error

  Call permissions:
*/
func (cer *Certificate) GetCertificate(id string) (dto.CertificateInfo, error) {
	var certificate dto.CertificateInfo
	if !isValidHexAddress(id) {
		return certificate, errors.New("id is not valid bid")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_certificate", params)
	if err != nil {
		return certificate, err
	}

	res, err := pointer.ToCertificateInfo()
	if err != nil{
		return certificate, err
	}

	return *res, nil
}

/*
  GetIssuer:
   	EN -
	CN - 查询信任锚信息，证书颁发者的信息
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- dto.IssuerSignature
		Id        string // 凭证ID
		PublicKey string // 签名公钥
		Algorithm string // 签名算法
		Signature string // 签名内容
	- error

  Call permissions: Anyone
*/
func (cer *Certificate) GetIssuer(id string) (dto.IssuerSignature, error) {
	var issuerSignature dto.IssuerSignature
	if !isValidHexAddress(id) {
		return issuerSignature, errors.New("id is not valid bid")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_issuer", params)
	if err != nil {
		return issuerSignature, err
	}

	res, err := pointer.ToCertificateIssuerSignature()
	if err != nil{
		return issuerSignature, err
	}

	return *res, nil
}

/*
  GetSubject:
   	EN -
	CN - 查询个人(证书注册的接收者)信息，证书接收者的信息
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- dto.SubjectSignature
		Id        string //凭证ID
		PublicKey string // 签名公钥
		Algorithm string //签名算法
		Signature string //签名内容
	- error

  Call permissions: Anyone
*/
func (cer *Certificate) GetSubject(id string) (dto.SubjectSignature, error) {
	var subjectSignature dto.SubjectSignature
	if !isValidHexAddress(id) {
		return subjectSignature, errors.New("id is not valid bid")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_subject", params)
	if err != nil {
		return subjectSignature, err
	}

	res, err := pointer.ToCertificateSubjectSignature()
	if err != nil{
		return subjectSignature, err
	}

	return *res, nil
}
