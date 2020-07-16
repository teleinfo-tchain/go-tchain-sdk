package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
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

//var (
//	ErrCandiNameLenInvalid    = errors.New("the length of candidate's name should between [4, 20]")
//	ErrCandiUrlLenInvalid     = errors.New("the length of candidate's website url should between [6, 60]")
//	ErrCandiNameInvalid       = errors.New("candidate's name should consist of digits and lowercase letters")
//	ErrCandiInfoDup           = errors.New("candidate's name, website url or node url is duplicated with a registered candidate")
//	ErrCandiAlreadyRegistered = errors.New("candidate is already registered")
//	ErrPeerNotTrust           = errors.New("peer is not apply trust")
//)

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
]`


type Election struct {
	super *System
	abi   abi.ABI
}

func (sys *System) NewElection() *Election {
	parsedAbi, _ := abi.JSON(strings.NewReader(ElectionAbiJSON))

	election := new(Election)
	election.abi = parsedAbi
	election.super = sys
	return election
}

//"registerWitness","inputs":[{"name":"nodeUrl","type":"string"},{"name":"website","type":"string"},{"name":"name","type":"string"}],"outputs":[]
func(e *Election) RegisterWitness(from common.Address, witness *dto.RegisterWitness) (string, error){
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

//查询候选人
func(e *Election) GetCandidate(candidateAddress string)(*dto.Candidate, error){
	params := make([]string, 1)
	params[0] = candidateAddress

	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_candidate", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToElectionCandidate()
}

func(e *Election) GetAllCandidates()([]string, error){
	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allCandidates", nil)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionCandidates()
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

func(e *Election) GetStake(voterAddress string)(*dto.Stake, error){
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_stake", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionStake()
}

func (e *Election) GetRestBIFBounty()(*big.Int, error){
	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_restBIFBounty", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToElectionRestBIFBounty()
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

func(e *Election) GetVoter(voterAddress string)(*dto.Voter, error){
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_voter", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToElectionVoter()
}

//投票人列表
func(e *Election) GetVoterList(voterAddress string)([]string, error){
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_voterList", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToElectionVoterList()
}
