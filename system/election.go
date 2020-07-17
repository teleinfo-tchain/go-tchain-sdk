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
	//OneDay            = int64(24) * 3600
	//VoteOrProxyOneDay = OneDay
	//VoteOrProxyOneDay = 60
	//oneWeek = OneDay * 7
	//OneYear = OneDay * 365 //代币增发周期 一年
	//year    = 1559318400   // 2019-06-01 00:00:00
)

//var (
//	ErrCandiNameLenInvalid    = errors.New("the length of candidate's name should between [4, 20]")
//	ErrCandiUrlLenInvalid     = errors.New("the length of candidate's website url should between [6, 60]")
//	ErrCandiNameInvalid       = errors.New("candidate's name should consist of digits and lowercase letters")
//	ErrCandiInfoDup           = errors.New("candidate's name, website url or node url is duplicated with a registered candidate")
//	ErrCandiAlreadyRegistered = errors.New("candidate is already registered")
//	ErrPeerNotTrust           = errors.New("peer is not apply trust")
//)

//var (
//	nowTimeStamp = big.NewInt(year)
//
//	// 投票周期
//	unStakePeriod   = big.NewInt(VoteOrProxyOneDay)
//	baseBounty      = big.NewInt(0).Mul(big.NewInt(1e+18), big.NewInt(1000))
//	restTotalBounty = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1e9))
//
//	//代币增量 初始发行量的5%
//	tokenAdd = big.NewInt(0).Div(restTotalBounty, big.NewInt(20))
//
//)

// 见证人选举的AbiJson数据
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

// Election - The Election Module
type Election struct {
	super *System
	abi   abi.ABI
}

// NewElection - NewElection初始化
func (sys *System) NewElection() *Election {
	parsedAbi, _ := abi.JSON(strings.NewReader(ElectionAbiJSON))

	election := new(Election)
	election.abi = parsedAbi
	election.super = sys
	return election
}

/*
 RegisterWitness: 注册成为见证人

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) RegisterWitness(from common.Address, witness *dto.RegisterWitness) (string, error){
	//encode
	// witness is a struct we need to use the components.
	var values []interface{}
	values = e.super.structToInterface(*witness,values)
	inputEncode, err := e.abi.Pack("registerWitness", values...)
	if err != nil {
		return "", err
	}

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	return e.super.sendTransaction(transaction)
}

/*
 UnRegisterWitness: 取消成为见证人

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) UnRegisterWitness(from common.Address) (string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("unregisterWitness")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	return e.super.sendTransaction(transaction)
}

/*
 GetCandidate: 查询候选人

 Returns：*dto.Candidate
*/
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

/*
 GetAllCandidates: 查询所有候选人

 Returns：[]dto.Candidate，列表内为候选人
*/
func(e *Election) GetAllCandidates()([]dto.Candidate, error){
	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allCandidates", nil)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionCandidates()
}

/*
 VoteWitnesses: 给见证人投票

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) VoteWitnesses(from common.Address, candidate string)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("voteWitnesses", candidate)

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 CancelVote: 撤销投票？？还是说是投反对票？？

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) CancelVote(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("cancelVote")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 StartProxy: 开启代理

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) StartProxy(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("startProxy")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 StopProxy: 关闭代理

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) StopProxy(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("stopProxy")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 CancelProxy: 取消代理

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) CancelProxy(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("cancelProxy")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 SetProxy: 设置代理

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) SetProxy(from common.Address, proxy string)(string, error){
	// encoding
	inputEncode, err := e.abi.Pack("setProxy", proxy)
	if err!= nil{
		return "", err
	}

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 Stake: 权益抵押

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) Stake(from common.Address, stakeCount *big.Int)(string, error){
	// encoding
	inputEncode, err := e.abi.Pack("stake", stakeCount)
	if err!= nil{
		return "", err
	}

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 UnStake: 撤销权益抵押

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？？
*/
func(e *Election) UnStake(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("unStake")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 GetStake: 查询抵押权益

 Returns：*dto.Stake
*/
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

/*
 GetRestBIFBounty: 获取剩余的Bif总激励

 Returns：*big.Int
*/
func (e *Election) GetRestBIFBounty()(*big.Int, error){
	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_restBIFBounty", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToElectionRestBIFBounty()
}

/*
 ExtractOwnBounty: 取出自身的赏金

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？
*/
func(e *Election) ExtractOwnBounty(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("extractOwnBounty")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 IssueAdditionalBounty: ？？？？？？？？

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: ？？
*/
func(e *Election) IssueAdditionalBounty(from common.Address)(string, error){
	// encoding
	inputEncode, _ := e.abi.Pack("issueAdditionalBounty")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
 GetVoter: 查询投票人信息

 Returns：*dto.Voter
*/
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

/*
 GetVoterList: 查询所有投票人信息

 Returns：[]string
*/
func(e *Election) GetVoterList(voterAddress string)([]dto.Voter, error){
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.RequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_voterList", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToElectionVoterList()
}
