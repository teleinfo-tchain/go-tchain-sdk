package dto

import (
	"encoding/json"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	customerror "github.com/bif/bif-sdk-go/constants"
	"math/big"
)

type Voter struct {
	Owner             common.Address   `json:"owner"`             // 投票人的地址
	IsProxy           bool             `json:"isProxy"`           // 是否是代理人
	ProxyVoteCount    *big.Int         `json:"proxyVoteCount"`    // 收到的代理的票数
	Proxy             common.Address   `json:"proxy"`             // 该节点设置的代理人
	LastVoteCount     *big.Int         `json:"lastVoteCount"`     // 上次投的票数
	LastVoteTimeStamp *big.Int         `json:"lastVoteTimeStamp"` // 上次投票时间戳
	VoteCandidates    []common.Address `json:"voteCandidates"`    // 投了哪些人
}

type Stake struct {
	Owner              common.Address `json:"owner"`              // 抵押代币的所有人
	StakeCount         *big.Int       `json:"stakeCount"`         // 抵押的代币数量
	LastStakeTimeStamp *big.Int       `json:"lastStakeTimeStamp"` // 上次抵押时间戳
}

func (pointer *RequestResult) ToPeerPeriod() (uint64, error) {

	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	return uint64(result.(hexutil.Uint64)), nil
}

func (pointer *RequestResult) ToPeerActive() (bool, error) {

	if err := pointer.checkResponse(); err != nil {
		return false, err
	}

	result := (pointer).Result.(interface{})

	return result.(bool), nil
}

func (pointer *RequestResult) ToPeerCertificate() (*PeerCertificate, error) {

	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	peerCertificate := &PeerCertificate{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, peerCertificate)

	return peerCertificate, err
}

func (pointer *RequestResult) ToIsBaseTrustAnchor() (bool, error) {
	if err := pointer.checkResponse(); err != nil {
		return false, err
	}

	result := (pointer).Result.(interface{})

	return result.(bool), nil
}

func (pointer *RequestResult) ToIsTrustAnchor() (bool, error) {
	if err := pointer.checkResponse(); err != nil {
		return false, err
	}

	result := (pointer).Result.(interface{})

	return result.(bool), nil
}

func (pointer *RequestResult) ToTrustAnchor() (*TrustAnchor, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	trustAnchor := &TrustAnchor{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, trustAnchor)

	return trustAnchor, err
}

func (pointer *RequestResult) ToTrustAnchorStatus() (uint64, error) {
	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	return uint64(result.(hexutil.Uint64)), nil
}

func (pointer *RequestResult) ToTrustAnchorCertificateList() ([]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	stringLi := make([]string, len(result))
	for i, v := range result {
		stringLi[i] = v.(string)
	}

	return stringLi, nil
}

func (pointer *RequestResult) ToBaseTrustAnchor() ([]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	stringLi := make([]string, len(result))
	for i, v := range result {
		stringLi[i] = v.(string)
	}

	return stringLi, nil
}

func (pointer *RequestResult) ToBaseTrustAnchorNumber() (uint64, error) {
	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	return uint64(result.(hexutil.Uint64)), nil
}

func (pointer *RequestResult) ToExpendTrustAnchor() ([]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})

	stringLi := make([]string, len(result))
	for i, v := range result {
		stringLi[i] = v.(string)
	}

	return stringLi, nil
}

func (pointer *RequestResult) ToExpendTrustAnchorNumber() (uint64, error) {
	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	return uint64(result.(hexutil.Uint64)), nil
}

// 解析测试注意！！！
func (pointer *RequestResult) ToTrustAnchorVoter() ([]TrustAnchorVoter, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]map[string]interface{})

	trustAnchorVoterLi := make([]TrustAnchorVoter, len(result))

	for i, v := range result {
		marshal, err := json.Marshal(v)

		if err != nil {
			return nil, customerror.UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, &trustAnchorVoterLi[i])
		return nil, err
	}

	return trustAnchorVoterLi, nil
}

func (pointer *RequestResult) ToCertificatePeriod() (uint64, error) {
	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}

	result := (pointer).Result.(interface{})

	return uint64(result.(hexutil.Uint64)), nil
}

func (pointer *RequestResult) ToCertificateActive() (bool, error) {
	if err := pointer.checkResponse(); err != nil {
		return false, err
	}

	result := (pointer).Result.(interface{})

	return result.(bool), nil
}

func (pointer *RequestResult) ToCertificateInfo() (*CertificateInfo, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	certificateInfo := &CertificateInfo{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, certificateInfo)
	return certificateInfo, err
}

func (pointer *RequestResult) ToCertificateIssuerSignature() (*IssuerSignature, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	issuerSignature := &IssuerSignature{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, issuerSignature)
	return issuerSignature, err
}

func (pointer *RequestResult) ToCertificateSubjectSignature() (*SubjectSignature, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	subjectSignature := &SubjectSignature{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, subjectSignature)
	return subjectSignature, err
}

func (pointer *RequestResult) ToDocument() (*Document, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	document := &Document{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, document)
	return document, err
}

func (pointer *RequestResult) ToDocIsEnable() (bool, error) {
	if err := pointer.checkResponse(); err != nil {
		return false, err
	}

	result := (pointer).Result.(interface{})

	return result.(bool), nil
}

func (pointer *RequestResult) ToElectionRestBIFBounty() (*big.Int, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(interface{})

	return result.(*big.Int), nil
}

func (pointer *RequestResult) ToElectionCandidate() (*Candidate, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	candidate := &Candidate{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, candidate)

	return candidate, err
}

func (pointer *RequestResult) ToElectionCandidates() ([]Candidate, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]map[string]interface{})

	candidateLi := make([]Candidate, len(result))

	for i, v := range result {
		marshal, err := json.Marshal(v)

		if err != nil {
			return nil, customerror.UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, &candidateLi[i])
		return nil, err
	}
	return candidateLi, nil
}

func (pointer *RequestResult) ToElectionVoter() (*Voter, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	voter := &Voter{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, voter)
	return voter, err
}

func (pointer *RequestResult) ToElectionVoterList() ([]Voter, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]map[string]interface{})

	voterLi := make([]Voter, len(result))

	for i, v := range result {
		marshal, err := json.Marshal(v)

		if err != nil {
			return nil, customerror.UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, &voterLi[i])
		return nil, err
	}
	return voterLi, nil
}

func (pointer *RequestResult) ToElectionStake() (*Stake, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}

	stake := &Stake{}

	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, stake)

	return stake, err
}
