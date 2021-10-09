package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

type SystemRequestResult struct {
	RequestResult
}

func (pointer *SystemRequestResult) ToAllianceDirectors() ([]*Alliance, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("联盟成员为空")
		}
		return nil, err
	}

	directorLi := (pointer).Result.([]interface{})
	directors := make([]*Alliance, len(directorLi))

	for i, v := range directorLi {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}

		info := &Alliance{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, info)
		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		directors[i] = info

	}
	return directors, nil
}

func (pointer *SystemRequestResult) ToAllianceVices() ([]*Alliance, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("联盟成员为空")
		}
		return nil, err
	}

	viceLi := (pointer).Result.([]interface{})
	vices := make([]*Alliance, len(viceLi))

	for i, v := range viceLi {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}

		info := &Alliance{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, info)
		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		vices[i] = info

	}
	return vices, nil
}

func (pointer *SystemRequestResult) ToAlliance() (*Alliance, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("联盟成员为空")
		}
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	alliance := &Alliance{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, alliance)
	return alliance, err
}

func (pointer *SystemRequestResult) ToWeights() (*Weights, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("联盟成员权重为空")
		}
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	weights := &Weights{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, weights)
	return weights, err
}

func (pointer *SystemRequestResult) ToCertificateInfo() (*CertificateInfo, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("可信证书不存在")
		}
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	certificateInfo := &CertificateInfo{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, certificateInfo)
	return certificateInfo, err
}

func (pointer *SystemRequestResult) ToCertificateIssuerSignature() (*IssuerSignature, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("证书颁发者为空")
		}
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	issuerSignature := &IssuerSignature{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, issuerSignature)
	return issuerSignature, err
}

func (pointer *SystemRequestResult) ToCertificateSubjectSignature() (*SubjectSignature, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("证书接收者为空")
		}
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	subjectSignature := &SubjectSignature{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, subjectSignature)
	return subjectSignature, err
}

func (pointer *SystemRequestResult) ToDocument() (*Document, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("did 文档未初始化")
		}
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	document := &Document{}
	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, document)
	return document, err
}

func (pointer *SystemRequestResult) ToRestBIFBounty() (*big.Int, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	res := (pointer).Result.(interface{})
	ret, success := big.NewInt(0).SetString(res.(string)[2:], 16)

	if !success {
		return nil, errors.New(fmt.Sprintf("Failed to convert %s to BigInt", res.(string)))
	}

	return ret, nil
}

func (pointer *SystemRequestResult) ToElectionCandidate() (*Candidate, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if result["owner"] == "" {
		return nil, errors.New("无候选者")
	}

	candidate := &Candidate{}

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, candidate)

	return candidate, err
}

func (pointer *SystemRequestResult) ToElectionCandidates() ([]Candidate, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	candidateLi := (pointer).Result.([]interface{})
	candidates := make([]Candidate, len(candidateLi))

	for i, v := range candidateLi {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}

		info := &Candidate{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, info)
		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		candidates[i] = *info

	}

	return candidates, nil
}

func (pointer *SystemRequestResult) ToElectionVoter() (*Voter, error) {
	if err := pointer.checkResponse(); err != nil {
		if err == EMPTYRESPONSE {
			return nil, errors.New("投票人信息为空")
		}
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	voter := &Voter{}
	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, voter)
	return voter, err
}

func (pointer *SystemRequestResult) ToElectionVoterList() ([]Voter, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	resultVoters := (pointer).Result.([]interface{})

	voters := make([]Voter, len(resultVoters))

	for i, v := range resultVoters {

		result := v.(map[string]interface{})
		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}

		info := &Voter{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, info)
		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		voters[i] = *info

	}

	return voters, nil
}

func (pointer *SystemRequestResult) ToRoundStateInfo() (*RoundStateInfo, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	roundStateInfo := &RoundStateInfo{}

	marshal, err := json.Marshal(result)

	err = json.Unmarshal(marshal, roundStateInfo)

	return roundStateInfo, err
}

func (pointer *SystemRequestResult) ToRoundChangeSetInfo() (*RoundChangeSetInfo, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	roundChangeSetInfo := &RoundChangeSetInfo{}

	marshal, err := json.Marshal(result)

	err = json.Unmarshal(marshal, roundChangeSetInfo)

	return roundChangeSetInfo, err
}

func (pointer *SystemRequestResult) ToBacklogs() (map[string][]*Message, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	result := (pointer).Result.(map[string]interface{})

	if len(result) == 0 {
		return nil, EMPTYRESPONSE
	}

	backlogs := make(map[string][]*Message, len(result))

	marshal, err := json.Marshal(result)

	if err != nil {
		return nil, UNPARSEABLEINTERFACE
	}

	err = json.Unmarshal(marshal, &backlogs)

	return backlogs, err
}

func (pointer *SystemRequestResult) ToAllContract() ([]AllContract, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}

	resultContracts := (pointer).Result.([]interface{})

	contracts := make([]AllContract, len(resultContracts))

	for i, v := range resultContracts {

		result := v.(map[string]interface{})

		if len(result) == 0 {
			return nil, EMPTYRESPONSE
		}

		info := &AllContract{}

		marshal, err := json.Marshal(result)

		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		err = json.Unmarshal(marshal, info)
		if err != nil {
			return nil, UNPARSEABLEINTERFACE
		}

		contracts[i] = *info

	}

	return contracts, nil
}
