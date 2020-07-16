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
	//ErrCertificateRegistered  = errors.New("certificate is already registered")
	ErrCertificateNotExist    = errors.New("certificate doesn't exist")
	//ErrIllegalSender          = errors.New("illegal sender")
	//ErrIllSenderOrCerNotExist = errors.New("illegal sender Or certificate doesn't exist")
)

const CertificateAbiJSON = `[
{"constant": false,"name":"registerCertificate","inputs":[{"name":"Id","type":"string"},{"name":"Context","type":"string"},{"name":"Subject","type":"string"},{"name":"Period","type":"uint64"},{"name":"IssuerAlgorithm","type":"string"},{"name":"IssuerSignature","type":"string"},{"name":"SubjectPublicKey","type":"string"},{"name":"SubjectAlgorithm","type":"string"},{"name":"SubjectSignature","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificates","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"cerdEvent","type":"event"}
]`

type Certificate struct {
	super *System
	abi   abi.ABI
}

func (sys *System) NewCertificate() *Certificate {
	parsedAbi, _ := abi.JSON(strings.NewReader(CertificateAbiJSON))

	cer := new(Certificate)
	cer.abi = parsedAbi
	cer.super = sys
	return cer
}

//信任锚颁发证书，如果是根信任锚颁发的证书，则证书接收者可以进行部署合约和大额转账操作
func (cer *Certificate) RegisterCertificate(from common.Address, registerCertificate *dto.RegisterCertificate) (string, error) {
	// encoding
	// registerCertificate is a struct we need to use the components.
	var values []interface{}
	values = cer.super.StructToInterface(*registerCertificate,values)
	inputEncode, err := cer.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = CertificateContractAddr
	transaction.Data = types.ComplexString(hexutil.Encode(inputEncode))
	return cer.super.SendTransaction(transaction)
}

//信任锚吊销个人证书
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
	return cer.super.SendTransaction(transaction)
}

//批量吊销个人证书，把自己颁发的证书全部吊销
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
	return cer.super.SendTransaction(transaction)
}

//证书有效期，如果证书被吊销，则有效期是0
func (cer *Certificate) GetPeriod(id string) (uint64, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_period", params)
	if err != nil{
		return 0, err
	}

	return pointer.ToCertificatePeriod()
}

func (cer *Certificate) GetActive(id string) (bool, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_active", params)
	if err != nil{
		return false, err
	}

	return pointer.ToCertificateActive()
}

func (cer *Certificate) GetCertificate(id string) (*dto.CertificateInfo, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_certificate", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToCertificateInfo()
}

//信任锚信息，证书颁发者的信息
func (cer *Certificate) GetIssuer(id string) (*dto.IssuerSignature, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_issuer", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToCertificateIssuerSignature()
}


func (cer *Certificate) GetSubject(id string) (*dto.SubjectSignature, error) {
	params := make([]string, 1)
	params[0] = id

	pointer := &dto.RequestResult{}

	err := cer.super.provider.SendRequest(pointer, "certificate_subject", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToCertificateSubjectSignature()
}

