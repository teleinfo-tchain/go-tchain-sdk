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
	ElectionContractAddr = "did:bid:000000000000000000000009"
	// VoteLimit         = 64
	// OneDay            = int64(24) * 3600
	// VoteOrProxyOneDay = OneDay
	// VoteOrProxyOneDay = 60
	// oneWeek = OneDay * 7
	// OneYear = OneDay * 365 //代币增发周期 一年
	// year    = 1559318400   // 2019-06-01 00:00:00
)

// var (
// 	ErrCandiNameLenInvalid    = errors.New("the length of candidate's name should between [4, 20]")
// 	ErrCandiUrlLenInvalid     = errors.New("the length of candidate's website url should between [6, 60]")
// 	ErrCandiNameInvalid       = errors.New("candidate's name should consist of digits and lowercase letters")
// 	ErrCandiInfoDup           = errors.New("candidate's name, website url or node url is duplicated with a registered candidate")
// 	ErrCandiAlreadyRegistered = errors.New("candidate is already registered")
// 	ErrPeerNotTrust           = errors.New("peer is not apply trust")
// )

// var (
// 	nowTimeStamp = big.NewInt(year)
//
// 	// 投票周期
// 	unStakePeriod   = big.NewInt(VoteOrProxyOneDay)
// 	baseBounty      = big.NewInt(0).Mul(big.NewInt(1e+18), big.NewInt(1000))
// 	restTotalBounty = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1e9))
//
// 	//代币增量 初始发行量的5%
// 	tokenAdd = big.NewInt(0).Div(restTotalBounty, big.NewInt(20))
//
// )

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

Params:
	- from: [20]byte，交易发送方地址
	- witness: *dto.RegisterWitness，注册的见证人信息
		NodeUrl string
		Website string
		Name    string

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) RegisterWitness(from common.Address, witness *dto.RegisterWitness) (string, error) {
	// encode
	// witness is a struct we need to use the components.
	var values []interface{}
	values = e.super.structToInterface(*witness, values)
	inputEncode, err := e.abi.Pack("registerWitness", values...)
	if err != nil {
		return "", err
	}

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	return e.super.sendTransaction(transaction)
}

/*
UnRegisterWitness: 取消成为见证人

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: 只能自己取消自己
*/
func (e *Election) UnRegisterWitness(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("unregisterWitness")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	return e.super.sendTransaction(transaction)
}

/*
GetCandidate: 查询候选人

Params:
	- candidateAddress: string，候选人的地址

Returns:
	- *dto.Candidate
		Owner           string       `json:"owner"`           // 候选人地址
		Name            string       `json:"name"`            // 候选人名称
		Active          bool         `json:"active"`          // 当前是否是候选人
		Url             string       `json:"url"`             // 节点的URL
		VoteCount       *hexutil.Big `json:"voteCount"`       // 收到的票数
		TotalBounty     *hexutil.Big `json:"totalBounty"`     // 总奖励金额
		ExtractedBounty *hexutil.Big `json:"extractedBounty"` // 已提取奖励金额
		LastExtractTime *hexutil.Big `json:"lastExtractTime"` // 上次提权时间
		Website         string       `json:"website"`         // 见证人网站
	- error

Call permissions: Anyone
*/
func (e *Election) GetCandidate(candidateAddress string) (*dto.Candidate, error) {
	params := make([]string, 1)
	params[0] = candidateAddress

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_candidate", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionCandidate()
}

/*
GetAllCandidates: 查询所有候选人

Params:
	- None

Returns:
	- []dto.Candidate，列表内为候选人信息，参考GetCandidate的候选人信息
	- error

Call permissions: Anyone
*/
func (e *Election) GetAllCandidates() ([]*dto.Candidate, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allCandidates", nil)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionCandidates()
}

/*
VoteWitnesses: 给见证人投票

Params:
	- from: [20]byte，交易发送方地址
	- candidate: string，候选人的地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) VoteWitnesses(from common.Address, candidate string) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("voteWitnesses", candidate)

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
CancelVote: 撤销投票？？还是说是投反对票？？

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) CancelVote(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("cancelVote")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
StartProxy: 开启代理

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) StartProxy(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("startProxy")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
StopProxy: 关闭代理

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) StopProxy(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("stopProxy")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
CancelProxy: 取消代理

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) CancelProxy(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("cancelProxy")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
SetProxy: 设置代理

Params:
	- from: [20]byte，交易发送方地址
	- proxy: string，???

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) SetProxy(from common.Address, proxy string) (string, error) {
	// encoding
	inputEncode, err := e.abi.Pack("setProxy", proxy)
	if err != nil {
		return "", err
	}

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
Stake: 权益抵押

Params:
	- from: [20]byte，交易发送方地址
	- stakeCount: *big.Int，抵押的权益数量

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) Stake(from common.Address, stakeCount *big.Int) (string, error) {
	// encoding
	inputEncode, err := e.abi.Pack("stake", stakeCount)
	if err != nil {
		return "", err
	}

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
UnStake: 撤销权益抵押

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？？
*/
func (e *Election) UnStake(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("unStake")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
GetStake: 查询抵押权益

Params:
	- voterAddress: string，投票者的地址

Returns:
	- *dto.Stake
		Owner              common.Address `json:"owner"`              // 抵押代币的所有人
		StakeCount         *big.Int       `json:"stakeCount"`         // 抵押的代币数量
		LastStakeTimeStamp *big.Int       `json:"lastStakeTimeStamp"` // 上次抵押时间戳
	- error

Call permissions: Anyone
*/
func (e *Election) GetStake(voterAddress string) (*dto.Stake, error) {
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_stake", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionStake()
}

/*
GetRestBIFBounty: 查询剩余的Bif总激励

Params:
	- None

Returns:
	- *big.Int
	- error

Call permissions: Anyone
*/
func (e *Election) GetRestBIFBounty() (*big.Int, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_restBIFBounty", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()
}

/*
ExtractOwnBounty: 取出自身的赏金

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？
*/
func (e *Election) ExtractOwnBounty(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("extractOwnBounty")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
IssueAdditionalBounty: ？？？？？？？？

Params:
	- from: [20]byte，交易发送方地址

Returns:
	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

Call permissions: ？？
*/
func (e *Election) IssueAdditionalBounty(from common.Address) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("issueAdditionalBounty")

	transaction := e.super.prePareTransaction(from, ElectionContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return e.super.sendTransaction(transaction)
}

/*
GetVoter: 查询投票人信息

Params:
	- voterAddress: string，投票者的地址

Returns:
	- *dto.Voter
		Owner             common.Address   `json:"owner"`             // 投票人的地址
		IsProxy           bool             `json:"isProxy"`           // 是否是代理人
		ProxyVoteCount    *big.Int         `json:"proxyVoteCount"`    // 收到的代理的票数
		Proxy             common.Address   `json:"proxy"`             // 该节点设置的代理人
		LastVoteCount     *big.Int         `json:"lastVoteCount"`     // 上次投的票数
		LastVoteTimeStamp *big.Int         `json:"lastVoteTimeStamp"` // 上次投票时间戳
		VoteCandidates    []common.Address `json:"voteCandidates"`    // 投了哪些人
	- error

Call permissions: Anyone
*/
func (e *Election) GetVoter(voterAddress string) (*dto.Voter, error) {
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_voter", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionVoter()
}

/*
GetVoterList: 查询所有投票人信息

Params:
	- voterAddress: string，投票者的地址

Returns:
	- []dto.Voter，投票人的详细信息，参考GetVoter
	- error

Call permissions: Anyone
*/
func (e *Election) GetVoterList(voterAddress string) ([]*dto.Voter, error) {
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_voterList", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionVoterList()
}
