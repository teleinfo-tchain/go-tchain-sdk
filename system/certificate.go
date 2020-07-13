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
	ErrCertificateRegistered  = errors.New("Certificate is already registered")
	ErrCertificateNotExist    = errors.New("Certificate doesn't exist")
	ErrIllegalSender          = errors.New("Illegal Sender")
	ErrIllSenderOrCerNotExist = errors.New("Illegal Sender Or Certificate doesn't exist")
)

const CertificateAbiJSON = `[
{"constant": false,"name":"registerCertificate","inputs":[{"name":"Id","type":"string"},{"name":"Context","type":"string"},{"name":"Subject","type":"string"},{"name":"Period","type":"uint64"},{"name":"IssuerAlgorithm","type":"string"},{"name":"IssuerSignature","type":"string"},{"name":"SubjectPublicKey","type":"string"},{"name":"SubjectAlgorithm","type":"string"},{"name":"SubjectSignature","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": true,"name":"queryPeriod","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"period","type":"uint64"}],"type":"function"},
{"constant": true,"name":"queryActive","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"isEnable","type":"bool"}],"type":"function"},
{"constant": true,"name":"queryIssuer","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Issuer","type":"string"}],"type":"function"},
{"constant": true,"name":"queryIssuerSignature","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Id","type":"string"},{"name":"PublicKey","type":"string"},{"name":"Algorithm","type":"string"},{"name":"Signature","type":"string"}],"type":"function"},
{"constant": true,"name":"querySubjectSignature","inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Id","type":"string"},{"name":"PublicKey","type":"string"},{"name":"Algorithm","type":"string"},{"name":"Signature","type":"string"}],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"cerdEvent","type":"event"}
]`

type Certificate struct {
	super *System
	abi   abi.ABI
}

type Signature struct {
	Id        string
	PublicKey string
	Algorithm string
	Signature string
}

type RegisterCertificate struct {
	Id string
	Context string
	Subject string
	Period uint64
	IssuerAlgorithm string
	IssuerSignature string
	SubjectPublicKey string
	SubjectAlgorithm string
	SubjectSignature string
}
func (system *System) NewCertificate() (*Certificate, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(CertificateAbiJSON))
	if err != nil {
		return nil, err
	}

	certificate := new(Certificate)
	certificate.abi = parsedAbi
	certificate.super = system
	return certificate, nil
}

//"inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"period","type":"uint64"}]
func (certificate *Certificate) GetPeriod(from common.Address, id string) (uint64, error) {
	// encoding
	inputEncode, err := certificate.abi.Pack("queryPeriod", id)
	if err != nil {
		return 0, err
	}

	transaction := certificate.super.PrePareTransaction(from, CertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := certificate.super.Call(transaction)
	if err != nil {
		return 0, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var period uint64
	err = certificate.abi.Methods["queryPeriod"].Outputs.Unpack(&period, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return 0, err
	}

	if period == 0 {
		return 0, ErrCertificateNotExist
	}
	return period, err
}

//"inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"isEnable","type":"bool"}]
func (certificate *Certificate) GetActive(from common.Address, id string) (bool, error) {
	// encoding
	inputEncode, err := certificate.abi.Pack("queryActive", id)
	if err != nil {
		return false, err
	}

	transaction := certificate.super.PrePareTransaction(from, CertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := certificate.super.Call(transaction)
	if err != nil {
		return false, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var isEnable bool
	err = certificate.abi.Methods["queryActive"].Outputs.Unpack(&isEnable, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return false, err
	}

	return isEnable, err
}

//"inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Issuer","type":"string"}]
func (certificate *Certificate) GetIssuer(from common.Address, id string) (string, error) {
	// encoding
	inputEncode, err := certificate.abi.Pack("queryIssuer", id)
	if err != nil {
		return "", err
	}

	transaction := certificate.super.PrePareTransaction(from, CertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := certificate.super.Call(transaction)
	if err != nil {
		return "", err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var issuer string
	err = certificate.abi.Methods["queryIssuer"].Outputs.Unpack(&issuer, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return "", err
	}
	return issuer, err
}

//"inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Id","type":"string"},{"name":"PublicKey","type":"string"},{"name":"Algorithm","type":"string"},{"name":"Signature","type":"string"}]
func (certificate *Certificate) GetIssuerSignature(from common.Address, id string) (*Signature, error) {
	// encoding
	inputEncode, err := certificate.abi.Pack("queryIssuerSignature", id)
	if err != nil {
		return nil, err
	}

	transaction := certificate.super.PrePareTransaction(from, CertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := certificate.super.Call(transaction)
	if err != nil {
		return nil, err
	}
	//fmt.Println("result string is ", requestResult.Result.(string))
	var issuerSignature Signature
	err = certificate.abi.Methods["queryIssuerSignature"].Outputs.Unpack(&issuerSignature, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil{
		return nil, err
	}
	if issuerSignature.Id == "0000000000000000000000000000000000000000"{
		return nil, ErrCertificateNotExist
	}
	return &issuerSignature, err
}

//"inputs":[{"name":"id", "type":"string"}],"outputs":[{"name":"Id","type":"string"},{"name":"PublicKey","type":"string"},{"name":"Algorithm","type":"string"},{"name":"Signature","type":"string"}]
func (certificate *Certificate) GetSubjectSignature(from common.Address, id string) (*Signature, error) {
	// encoding
	inputEncode, err := certificate.abi.Pack("querySubjectSignature", id)
	if err != nil {
		return nil, err
	}

	transaction := certificate.super.PrePareTransaction(from, CertificateContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := certificate.super.Call(transaction)
	if err != nil {
		return nil, err
	}

	//fmt.Println("result string is ", requestResult.Result.(string))
	var subjectSignature Signature
	err = certificate.abi.Methods["querySubjectSignature"].Outputs.Unpack(&subjectSignature, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}
	if subjectSignature.Id == "0000000000000000000000000000000000000000"{
		return nil, ErrCertificateNotExist
	}
	return &subjectSignature, err
}

//"inputs":[{"name":"Id","type":"string"},{"name":"Context","type":"string"},{"name":"Subject","type":"string"},{"name":"Period","type":"uint64"},{"name":"IssuerAlgorithm","type":"string"},{"name":"IssuerSignature","type":"string"},{"name":"SubjectPublicKey","type":"string"},{"name":"SubjectAlgorithm","type":"string"},{"name":"SubjectSignature","type":"string"}],"outputs":[]
func (certificate *Certificate) RegisterCertificate(from common.Address, registerCertificate *RegisterCertificate) (string, error) {
	// encoding
	// registerCertificate is a struct we need to use the components.
	var values []interface{}
	values = certificate.super.StructToInterface(*registerCertificate,values)
	inputEncode, err := certificate.abi.Pack("registerCertificate", values...)
	if err != nil {
		return "", err
	}
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = CertificateContractAddr
	transaction.Data = types.ComplexString(hexutil.Encode(inputEncode))
	return certificate.super.SendTransaction(transaction)
}

//inputs":[{"name":"id","type":"string"}],"outputs":[]
func (certificate *Certificate) RevokedCertificate(from common.Address, id string) (string, error) {
	// encoding
	inputEncode, err := certificate.abi.Pack("revokedCertificate", id)
	if err != nil {
		return "", err
	}
	transaction := new(dto.TransactionParameters)
	transaction.From = from.String()
	transaction.To = CertificateContractAddr
	transaction.Data = types.ComplexString(hexutil.Encode(inputEncode))
	return certificate.super.SendTransaction(transaction)
}
