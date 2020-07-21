package system

import (
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
	"strings"
)

const (
	CertificateContractAddr = "did:bid:00000000000000000000000b"
)

var (
	// ErrCertificateRegistered  = errors.New("certificate is already registered")
	ErrCertificateNotExist = errors.New("certificate doesn't exist")
	// ErrIllegalSender          = errors.New("illegal sender")
	// ErrIllSenderOrCerNotExist = errors.New("illegal sender Or certificate doesn't exist")
)

// 个人可信的AbiJson数据
const CertificateAbiJSON = `[
{"constant": false,"name":"registerCertificate","inputs":[{"name":"Id","type":"string"},{"name":"Context","type":"string"},{"name":"Subject","type":"string"},{"name":"Period","type":"uint64"},{"name":"IssuerAlgorithm","type":"string"},{"name":"IssuerSignature","type":"string"},{"name":"SubjectPublicKey","type":"string"},{"name":"SubjectAlgorithm","type":"string"},{"name":"SubjectSignature","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificates","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"cerdEvent","type":"event"}
]`

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

/*
RegisterCertificate: 信任锚颁发证书，如果是根信任锚颁发的证书，则证书接收者可以进行部署合约和大额转账操作

Params:
	- from: [20]byte，交易发送方地址
	- registerCertificate:  The registerCertificate object(*dto.RegisterCertificate)
		Id               string //个人可信证书bid
		Context          string //证书上下文环境，随便一个字符串，不验证
		Subject          string //证书接收者的bid，证书是颁给谁的
		Period           uint64 //证书有效期，以年为单位的整型
		IssuerAlgorithm  string // 颁发者签名算法，字符串
		IssuerSignature  string //颁发者签名值，16进制字符串
		SubjectPublicKey string // 接收者公钥，16进制字符串
		SubjectAlgorithm string //接收者签名算法，字符串
		SubjectSignature string //接收者签名值，16进制字符串

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: 只有信任锚地址可以调用

BUG(rpc):颁发证书时的Period是否必须以年为单位？？
*/
func (cer *Certificate) RegisterCertificate(from common.Address, registerCertificate *dto.RegisterCertificate) (string, error) {
	// encoding
	// registerCertificate is a struct we need to use the components.
	var values []interface{}
	values = cer.super.structToInterface(*registerCertificate, values)
	inputEncode, err := cer.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = CertificateContractAddr
	transaction.Data = types.ComplexString(hexutil.Encode(inputEncode))
	return cer.super.sendTransaction(transaction)
}

/*
RevokedCertificate: 信任锚吊销个人证书

Params:
	- from: [20]byte，交易发送方地址
	- id:  string，个人可信证书bid

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: 只有证书颁发者可以调用
*/
func (cer *Certificate) RevokedCertificate(from common.Address, id string) (string, error) {
	// encoding
	inputEncode, err := cer.abi.Pack("revokedCertificate", id)
	if err != nil {
		return "", err
	}
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = CertificateContractAddr
	transaction.Data = types.ComplexString(hexutil.Encode(inputEncode))
	return cer.super.sendTransaction(transaction)
}

/*
RevokedCertificates: 信任锚批量吊销个人证书，把自己颁发的证书全部吊销

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: 只有证书颁发者可以调用
*/
func (cer *Certificate) RevokedCertificates(from common.Address) (string, error) {
	// encoding
	inputEncode, err := cer.abi.Pack("revokedCertificates")
	if err != nil {
		return "", err
	}
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = CertificateContractAddr
	transaction.Data = types.ComplexString(hexutil.Encode(inputEncode))
	return cer.super.sendTransaction(transaction)
}

/*
GetPeriod: 查询个人证书的有效期

Params:
	- id: string,个人可信证书bid

Returns:
	- uint64，证书有效期，如果证书被吊销，则有效期是0
	- error

Call permissions: Anyone
*/
func (cer *Certificate) GetPeriod(id string) (uint64, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_period", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
GetActive: 查询证书是否可用,如果证书过期，则不可用；如果信任锚注销，则不可用

Params:
	- id: string,个人可信证书bid

Returns:
	- bool，true可用，false不可用
	- error

Call permissions: Anyone
*/
func (cer *Certificate) GetActive(id string) (bool, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_active", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
GetCertificate: 查询证书的信息，period和active没有进行逻辑判断

Params:
	- id: string,个人可信证书bid

Returns:
	- *dto.CertificateInfo
		Id             string   //凭证的hash
		Context        string   //证书所属上下文环境
		Issuer         string   //信任锚的bid
		Subject        string   //证书拥有者地址
		IssuedTime     *big.Int //颁发时间
		Period         uint64   //有效期
		IsEnable       bool     //true 凭证有效，false 凭证已撤销
		RevocationTime *big.Int //吊销时间
	- error

Call permissions: Anyone
*/
func (cer *Certificate) GetCertificate(id string) (*dto.CertificateInfo, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_certificate", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToCertificateInfo()
}

/*
GetIssuer: 查询信任锚信息，证书颁发者的信息

Params:
	- id: string,个人可信证书bid

Returns:
	- *dto.IssuerSignature
		Id        string //凭证ID
		PublicKey string // 签名公钥
		Algorithm string //签名算法
		Signature string //签名内容
	- error

Call permissions: Anyone
*/
func (cer *Certificate) GetIssuer(id string) (*dto.IssuerSignature, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_issuer", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToCertificateIssuerSignature()
}

/*
GetSubject: 查询个人(证书注册的接收者)信息，证书接收者的信息

Params:
	- id: string,个人可信证书bid

Returns:
	- *dto.SubjectSignature
		Id        string //凭证ID
		PublicKey string // 签名公钥
		Algorithm string //签名算法
		Signature string //签名内容
	- error

Call permissions: Anyone
*/
func (cer *Certificate) GetSubject(id string) (*dto.SubjectSignature, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_subject", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToCertificateSubjectSignature()
}
