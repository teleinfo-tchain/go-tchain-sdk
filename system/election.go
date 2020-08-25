package system

import (
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/dto"
	"math/big"
	"strings"
)

const (
	ElectionContractAddr = "did:bid:ZFT6dswwQWN63BF8pup729HuxLTRhy"
)

// todo:这几种错误对应的判断，需要填入前置判断
// var (
// 	ErrCandiNameLenInvalid    = errors.New("the length of candidate's name should between [4, 20]")
// 	ErrCandiUrlLenInvalid     = errors.New("the length of candidate's website url should between [6, 60]")
// 	ErrCandiNameInvalid       = errors.New("candidate's name should consist of digits and lowercase letters")
// 	ErrCandiInfoDup           = errors.New("candidate's name, website url or node url is duplicated with a registered candidate")
// 	ErrCandiAlreadyRegistered = errors.New("candidate is already registered")
// 	ErrPeerNotTrust           = errors.New("peer is not apply trust")
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
{"constant": false,"name":"issueAddtitionalBounty","inputs":[],"outputs":[],"type":"function"}
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

func (e *Election) judgeScheme(url string) bool {
	parts := strings.Split(url, "://")
	if len(parts) != 2 || parts[0] == "" {
		return false
	}
	return true
}

func (e *Election) registerWitnessPreCheck(witness *dto.RegisterWitness) (bool, error) {
	if len(witness.NodeUrl) == 0 || isBlankCharacter(witness.NodeUrl) {
		return false, errors.New("witness NodeUrl can't be empty or blank character")
	}

	if len(witness.Website) == 0 || isBlankCharacter(witness.Website) {
		return false, errors.New("witness Website can't be empty or blank character")
	}

	if !e.judgeScheme(witness.Website) {
		return false, errors.New("witness Website protocol scheme missing")
	}

	if len(witness.Name) == 0 || isBlankCharacter(witness.Name) {
		return false, errors.New("witness Name can't be empty or blank character")
	}

	return true, nil
}

/*
  RegisterWitness:
   	EN -
	CN - 注册成为候选
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- witness: *dto.RegisterWitness，注册的见证人信息
		NodeUrl string
		Website string
		Name    string

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) RegisterWitness(signTxParams *SysTxParams, witness *dto.RegisterWitness) (string, error) {
	ok, err := e.registerWitnessPreCheck(witness)
	if !ok {
		return "", err
	}

	// encode
	// witness is a struct we need to use the components.
	var values []interface{}
	values = e.super.structToInterface(*witness, values)
	inputEncode, err := e.abi.Pack("registerWitness", values...)

	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  UnRegisterWitness:
   	EN -
	CN -  取消成为候选
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只能自己取消自己
*/
func (e *Election) UnRegisterWitness(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("unregisterWitness")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  GetCandidate:
   	EN -
	CN -  查询候选人
  Params:
  	- candidateAddress: string，候选人的地址

  Returns:
  	- dto.Candidate
		Owner           string       `json:"owner"`           // 候选人地址
		Name            string       `json:"name"`            // 候选人名称
		Active          bool         `json:"active"`          // 当前是否是候选人
		Url             string       `json:"url"`             // 节点的URL
		VoteCount       *hexutil.Big `json:"voteCount"`       // 收到的票数
		TotalBounty     *hexutil.Big `json:"totalBounty"`     // 总奖励金额
		ExtractedBounty *hexutil.Big `json:"extractedBounty"` // 已提取奖励金额
		LastExtractTime uint64       `json:"lastExtractTime"` // 上次提权时间
		Website         string       `json:"website"`         // 见证人网站
	- error

  Call permissions: Anyone
*/
func (e *Election) GetCandidate(candidateAddress string) (dto.Candidate, error) {
	var candidate dto.Candidate
	if !isValidHexAddress(candidateAddress) {
		return candidate, errors.New("candidateAddress is not valid bid")
	}

	params := make([]string, 1)
	params[0] = candidateAddress

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_candidate", params)
	if err != nil {
		return candidate, err
	}

	res, err := pointer.ToElectionCandidate()
	if err != nil{
		return candidate, err
	}
	if res.Owner == ""{
		return candidate, errors.New("无候选者")
	}
	return *res,nil
}

/*
  GetAllCandidates:
   	EN -
	CN - 查询所有候选人
  Params:
  	- None

  Returns:
  	- []dto.Candidate，列表内为候选人信息，参考GetCandidate的候选人信息
	- error

  Call permissions: Anyone
*/
func (e *Election) GetAllCandidates() ([]dto.Candidate, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allCandidates", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionCandidates()

	return res,nil
}

/*
  VoteWitnesses:
   	EN -
	CN - 给见证人投票
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- candidate: string，候选人的地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) VoteWitnesses(signTxParams *SysTxParams, candidate string) (string, error) {
	if !isValidHexAddress(candidate) {
		return "", errors.New("candidate is not valid bid")
	}

	// encoding
	inputEncode, err := e.abi.Pack("voteWitnesses", candidate)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  CancelVote:
   	EN -
	CN - 撤销投票
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) CancelVote(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("cancelVote")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  StartProxy:
   	EN -
	CN - 开启代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) StartProxy(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("startProxy")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  StopProxy:
   	EN -
	CN - 关闭代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) StopProxy(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("stopProxy")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  CancelProxy:
   	EN -
	CN - 取消代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) CancelProxy(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("cancelProxy")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  SetProxy:
   	EN -
	CN - 设置代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- proxy: string，是个did的地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) SetProxy(signTxParams *SysTxParams, proxy string) (string, error) {
	// todo:加上proxy是否与签名方或者发送方的地址是否相同的判断
	if !isValidHexAddress(proxy) {
		return "", errors.New("proxy is not valid bid")
	}

	// encoding
	inputEncode, err := e.abi.Pack("setProxy", proxy)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  Stake:
   	EN -
	CN - 权益抵押
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- stakeCount: *big.Int，抵押的权益数量

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) Stake(signTxParams *SysTxParams, stakeCount *big.Int) (string, error) {
	// encoding
	inputEncode, err := e.abi.Pack("stake", stakeCount)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  UnStake:
   	EN -
	CN - 撤销权益抵押
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
*/
func (e *Election) UnStake(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("unStake")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  GetStake:
   	EN -
	CN - 查询抵押权益
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
	if !isValidHexAddress(voterAddress) {
		return nil, errors.New("voterAddress is not valid bid")
	}

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
  GetRestBIFBounty:
   	EN -
	CN - 查询剩余的Bif总激励
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

	return pointer.ToRestBIFBounty()
}

/*
  ExtractOwnBounty:
   	EN -
	CN - 取出自身的赏金
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
*/
func (e *Election) ExtractOwnBounty(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("extractOwnBounty")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  IssueAdditionalBounty:
   	EN -
	CN - ??
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
*/
func (e *Election) IssueAdditionalBounty(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("issueAddtitionalBounty")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContractAddr)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

/*
  GetVoter:
   	EN -
	CN - 查询投票人信息
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
	if !isValidHexAddress(voterAddress) {
		return nil, errors.New("voterAddress is not valid bid")
	}

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
  GetVoterList:
   	EN -
	CN - 查询所有投票人信息
  Params:
  	- voterAddress: string，投票者的地址

  Returns:
  	- []dto.Voter，投票人的详细信息，参考GetVoter
	- error

  Call permissions: Anyone
*/
func (e *Election) GetVoterList(voterAddress string) ([]dto.Voter, error) {
	if !isValidHexAddress(voterAddress) {
		return nil, errors.New("voterAddress is not valid bid")
	}

	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_voterList", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToElectionVoterList()
}
