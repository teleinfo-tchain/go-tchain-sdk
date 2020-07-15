package system

import (
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"math/big"
	"strings"
)

const (
	ElectionContractAddr      = "did:bid:000000000000000000000009"
	//VoteLimit         = 64
	OneDay            = int64(24) * 3600
	VoteOrProxyOneDay = OneDay
	//VoteOrProxyOneDay = 60
	//oneWeek = OneDay * 7
	//OneYear = OneDay * 365 //代币增发周期 一年
	year    = 1559318400   // 2019-06-01 00:00:00
)

var (
	ErrCandiNameLenInvalid    = errors.New("the length of candidate's name should between [4, 20]")
	ErrCandiUrlLenInvalid     = errors.New("the length of candidate's website url should between [6, 60]")
	ErrCandiNameInvalid       = errors.New("candidate's name should consist of digits and lowercase letters")
	ErrCandiInfoDup           = errors.New("candidate's name, website url or node url is duplicated with a registered candidate")
	ErrCandiAlreadyRegistered = errors.New("candidate is already registered")
	ErrPeerNotTrust           = errors.New("peer is not apply trust")
)

var (
	nowTimeStamp = big.NewInt(year)

	// 投票周期
	unStakePeriod   = big.NewInt(VoteOrProxyOneDay)
	baseBounty      = big.NewInt(0).Mul(big.NewInt(1e+18), big.NewInt(1000))
	restTotalBounty = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1e9))

	//代币增量 初始发行量的5%
	tokenAdd = big.NewInt(0).Div(restTotalBounty, big.NewInt(20))

)

const ElectionAbiJSON = `[
{"constant": false,"name":"registerWitness","inputs":[{"name":"nodeUrl","type":"string"},{"name":"website","type":"string"},{"name":"name","type":"string"}],"outputs":[],"type":"function"}, 
{"constant": false,"name":"unregisterWitness","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"voteWitnesses","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelVote","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"startProxy","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"stopProxy","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelProxy","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"setProxy","inputs":[{"name":"proxy","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"stake","inputs":[{"name":"stakeCount","type":"uint256"}],"outputs":[],"type":"function"},
{"constant": false,"name":"unStake","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"extractOwnBounty","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"electEvent","type":"event"},
{"constant": false,"name":"issueAdditionalBounty","inputs":[],"outputs":[],"type":"function"},
{"constant": true,"name":"queryCandidates","inputs":[{"name":"candidateAddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"voteCount","type":"uint64"},{"name":"active","type":"bool"},{"name":"url","type":"string"},{"name":"totalBounty","type":"uint64"},{"name":"extractedBounty","type":"uint64"},{"name":"lastExtractTime","type":"uint64"},{"name":"website","type":"string"},{"name":"name","type":"string"}],"type":"function"},
{"constant": true,"name":"queryAllCandidates","inputs":[],"outputs":[{"name":"listNum","type":"uint64"},{"name":"candidateInfoList","type":"string"}],"type":"function"},
{"constant": true,"name":"queryVoter","inputs":[{"name":"voterAddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"isProxy","type":"bool"},{"name":"proxyVoteCount","type":"uint64"},{"name":"proxy","type":"string"},{"name":"lastVoteCount","type":"uint64"},{"name":"timestamp","type":"uint64"},{"name":"voteCandidatesList","type":"string"}],"type":"function"},
{"constant": true,"name":"queryStake","inputs":[{"name":"voterAddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"stakeCount","type":"uint64"},{"name":"timestamp","type":"uint64"}],"type":"function"},
{"constant": true,"name":"queryVoterList","inputs":[{"name":"candidateAddress","type":"string"}],"outputs":[{"name":"candidateAddress","type":"string"},{"name":"voterNum","type":"uint64"},{"name":"voterList","type":"string"}],"type":"function"}
]`


type Election struct {
	super *System
	abi   abi.ABI
}

type RegisterWitness struct {
	NodeUrl string
	Website string
	Name string
}

type candidateInfo struct {
	Owner string
	VoteCount uint64
	Active bool
	Url string
	TotalBounty uint64
	ExtractedBounty uint64
	LastExtractTime uint64
	Website string
	Name string
}

type candidatesParsing struct {
	ListNum uint64
	CandidateInfoList string
}

type candidates struct {
	Num uint64
	CandidatesList []string
}

type voterParsing struct {
	Owner string
	IsProxy bool
	ProxyVoteCount uint64
	Proxy string
	LastVoteCount uint64
	Timestamp uint64
	VoteCandidatesList string
}

type voter struct {
	Owner string
	IsProxy bool
	ProxyVoteCount uint64
	Proxy string
	LastVoteCount uint64
	Timestamp uint64
	VoteCandidatesList []string
}

type votersParsing struct {
	CandidateAddress string
	VoterNum uint64
	VoterList string
}

type Voters struct {
	CandidateAddress string
	VoterNum uint64
	VotersLi []string
}

type stake struct {
	Owner string
	StakeCount uint64
	Timestamp uint64
}

func (sys *System) NewElection() *Election {
	parsedAbi, _ := abi.JSON(strings.NewReader(ElectionAbiJSON))

	election := new(Election)
	election.abi = parsedAbi
	election.super = sys
	return election
}

//"registerWitness","inputs":[{"name":"nodeUrl","type":"string"},{"name":"website","type":"string"},{"name":"name","type":"string"}],"outputs":[]
func(e *Election) RegisterWitness(from common.Address, witness *RegisterWitness) (string, error){
	//encode
	// witness is a struct we need to use the components.
	var values []interface{}
	values = e.super.StructToInterface(*witness,values)
	inputEncode, err := e.abi.Pack("registerWitness", values...)
	if err != nil {
		return "", err
	}

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	return e.super.SendTransaction(transaction)
}

//"unregisterWitness","inputs":[],"outputs":[]
func(e *Election) UnRegisterWitness(from common.Address) (string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("unregisterWitness")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	return e.super.SendTransaction(transaction)
}

//"queryCandidates","inputs":[{"name":"candidateAddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"voteCount","type":"uint64"},{"name":"active","type":"bool"},{"name":"url","type":"string"},{"name":"totalBounty","type":"uint64"},{"name":"extractedBounty","type":"uint64"},{"name":"lastExtractTime","type":"uint64"},{"name":"website","type":"string"},{"name":"name","type":"string"}]
func(e *Election) GetCandidates(from common.Address, candidateAddress string)(*candidateInfo, error){
	// encoding
	inputEncode, err := e.abi.Pack("queryCandidates", candidateAddress)
	if err != nil{
		return nil, err
	}

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := e.super.Call(transaction)

	if err != nil {
		return nil, err
	}

	//fmt.Println("result is ", requestResult.Result.(string))
	var candidate candidateInfo
	err = e.abi.Methods["queryCandidates"].Outputs.Unpack(&candidate, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}

	return &candidate, err
}

//"queryAllCandidates","inputs":[],"outputs":[{"name":"listNum","type":"uint64"},{"name":"candidateInfoList","type":"string"}]
func(e *Election) GetAllCandidates(from common.Address)(*candidates, error){
	// encoding
	inputEncode, _ := e.abi.Pack("queryAllCandidates")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := e.super.Call(transaction)

	if err != nil {
		return nil, err
	}

	//fmt.Println("result is ", requestResult.Result.(string))
	var candidatesParse candidatesParsing
	err = e.abi.Methods["queryAllCandidates"].Outputs.Unpack(&candidatesParse, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}

	var candidates candidates
	for i := uint64(0);i<candidatesParse.ListNum;i++{
		candidates.CandidatesList = append(candidates.CandidatesList, candidatesParse.CandidateInfoList[i*32:(i+1)*32-1])
	}
	candidates.Num = candidatesParse.ListNum
	return &candidates, err
}

//"voteWitnesses","inputs":[{"name":"candidate","type":"string"}],"outputs":[]
func(e *Election) VoteWitnesses(from common.Address, candidate string)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("voteWitnesses", candidate)

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"cancelVote","inputs":[],"outputs":[]
func(e *Election) CancelVote(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("cancelVote")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"startProxy","inputs":[],"outputs":[]
func(e *Election) StartProxy(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("startProxy")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"stopProxy","inputs":[],"outputs":[]
func(e *Election) StopProxy(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("stopProxy")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"cancelProxy","inputs":[],"outputs":[]
func(e *Election) CancelProxy(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("cancelProxy")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"setProxy","inputs":[{"name":"proxy","type":"string"}],"outputs":[]
func(e *Election) SetProxy(from common.Address, proxy string)(string, error){
	// encoding
	inputEncode, err := e.abi.Pack("setProxy", proxy)
	if err!= nil{
		return "", err
	}

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"stake","inputs":[{"name":"stakeCount","type":"uint256"}],"outputs":[]
func(e *Election) Stake(from common.Address, stakeCount *big.Int)(string, error){
	// encoding
	inputEncode, err := e.abi.Pack("stake", stakeCount)
	if err!= nil{
		return "", err
	}

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"unStake","inputs":[],"outputs":[]
func(e *Election) UnStake(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("unStake")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"queryStake","inputs":[{"name":"voterAddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"stakeCount","type":"uint64"},{"name":"timestamp","type":"uint64"}]
func(e *Election) GetStake(from common.Address, voterAddress string)(*stake, error){
	// encoding
	inputEncode, err := e.abi.Pack("queryStake", voterAddress)
	if err != nil{
		return nil, err
	}

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := e.super.Call(transaction)

	if err != nil {
		return nil, err
	}

	//fmt.Println("result is ", requestResult.Result.(string))
	var stakeInfo stake
	err = e.abi.Methods["queryStake"].Outputs.Unpack(&stakeInfo, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}

	return &stakeInfo, err
}

//"extractOwnBounty","inputs":[],"outputs":[]
func(e *Election) ExtractOwnBounty(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("extractOwnBounty")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//"issueAdditionalBounty","inputs":[],"outputs":[]
func(e *Election) IssueAdditionalBounty(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("issueAdditionalBounty")

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.SendTransaction(transaction)
}

//queryVoter","inputs":[{"name":"voterAddress","type":"string"}],"outputs":[{"name":"owner","type":"string"},{"name":"isProxy","type":"bool"},{"name":"proxyVoteCount","type":"uint64"},{"name":"proxy","type":"string"},{"name":"lastVoteCount","type":"uint64"},{"name":"timestamp","type":"uint64"},{"name":"voteCandidatesList","type":"string"}]
func(e *Election) GetVoter(from common.Address, voterAddress string)(*voter, error){
	// encoding
	inputEncode, err := e.abi.Pack("queryVoter", voterAddress)
	if err != nil{
		return nil, err
	}

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := e.super.Call(transaction)

	if err != nil {
		return nil, err
	}

	//fmt.Println("result is ", requestResult.Result.(string))
	var voterInfo voterParsing
	fmt.Println("voter ", voterInfo)
	err = e.abi.Methods["queryVoter"].Outputs.Unpack(&voterInfo, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}

	var voter voter
	voter.Owner = voterInfo.Owner
	voter.IsProxy = voterInfo.IsProxy
	voter.ProxyVoteCount = voterInfo.ProxyVoteCount
	voter.Proxy = voterInfo.Proxy
	voter.LastVoteCount = voterInfo.LastVoteCount
	voter.Timestamp = voterInfo.Timestamp
	for i := 0;i<len(voterInfo.VoteCandidatesList);i++{
		voter.VoteCandidatesList = append(voter.VoteCandidatesList, voterInfo.VoteCandidatesList[i*32:(i+1)*32-1])
	}

	return &voter, err
}

//"queryVoterList","inputs":[{"name":"candidateAddress","type":"string"}],"outputs":[{"name":"candidateAddress","type":"string"},{"name":"voterNum","type":"uint64"},{"name":"voterList","type":"string"}]
func(e *Election) GetVoterList(from common.Address, candidateAddress string)(*Voters, error){
	// encoding
	inputEncode, err := e.abi.Pack("queryVoterList", candidateAddress)
	if err != nil{
		return nil, err
	}

	transaction := e.super.PrePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := e.super.Call(transaction)

	if err != nil {
		return nil, err
	}

	//fmt.Println("result is ", requestResult.Result.(string))
	var voterList votersParsing
	err = e.abi.Methods["queryVoterList"].Outputs.Unpack(&voterList, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}

	var voters Voters
	for i := uint64(0);i<voterList.VoterNum;i++{
		voters.VotersLi = append(voters.VotersLi, voterList.VoterList[i*32:(i+1)*32-1])
	}

	voters.CandidateAddress = voterList.CandidateAddress
	voters.VoterNum = voterList.VoterNum
	return &voters, err
}
