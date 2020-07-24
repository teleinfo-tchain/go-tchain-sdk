package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"testing"
)

// 注册成为候选者节点
func TestRegisterWitness(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	keyJson, err := ioutil.ReadFile("../resources/superNodeKeyStore/UTC--172.17.6.53--did-bid-c935bd29a90fbeea87badf3e")
	if err != nil{
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = ""
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyJson
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	registerWitness := new(dto.RegisterWitness)
	registerWitness.NodeUrl = "127.0.0.1/test"
	registerWitness.Website = "www.tele.info.com"
	registerWitness.Name = "BeiJing"

	// registerWitness
	registerWitnessHash, err := elect.RegisterWitness(sysTxParams, registerWitness)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(registerWitnessHash, err)
}

// 取消成为候选者节点
func TestUnRegisterWitness(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	unRegisterWitnessHash, err := elect.UnRegisterWitness(sysTxParams)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(unRegisterWitnessHash, err)
}

// 获取候选节点基本信息
func TestGetCandidate(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP52+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	candidate, err := elect.GetCandidate(coinBase)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%#v \n", candidate)
}

// 获取所有的见证候选节点
func TestGetAllCandidates(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	allCandidates, err := elect.GetAllCandidates()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(allCandidates)
	// "did:bid:6cc796b8d6e2fbebc9b3cf9",
	// "did:bid:13803fb30b7e95d57103c2d",
	// "did:bid:c117c1794fc7a27bd301ae5",
	// "did:bid:590ed37615bdfefa496224c",
	// "did:bid:2e4f6a140ed099d15177d32",
	// "did:bid:136993a53a0f3b4e5b5dc0d",
	// "did:bid:55460d49c4364da777f6170",
	// "did:bid:d8e55affc8ac5016185ddc6",
	// "did:bid:0fcb3f265375cea7a57cceb",
	// "did:bid:ade588198aeb5991541c850",
	// "did:bid:c935bd29a90fbeea87badf3",
	// "did:bid:cd68ab895f5a238ce60b98f",
	// "did:bid:ee4f5aa444dc5ec815a25c7",
	// "did:bid:d423816e6984eaa9bafb4bb",
	// "did:bid:ae05edab896300f583747fd"
}

func TestVoteWitnesses(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	candidate := "did:bid:ee4f5aa444dc5ec815a25c7"
	voteWitnessHash, err := elect.VoteWitnesses(sysTxParams, candidate)
	if err != nil && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}
	if err == system.ErrCertificateNotExist {
		t.Log(err)
	}
	t.Log(voteWitnessHash)
}

func TestElectCancelVote(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	cancelVoteHash, err := elect.CancelVote(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(cancelVoteHash)
}

func TestStartProxy(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	startProxyHash, err := elect.StartProxy(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(startProxyHash)
}

func TestStopProxy(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	stopProxyHash, err := elect.StopProxy(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(stopProxyHash)
}

// not set proxy
func TestCancelProxy(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	cancelProxyHash, err := elect.CancelProxy(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(cancelProxyHash)
}

// account registered as a proxy is not allowed to use a proxy??
func TestSetProxy(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	proxy := resources.Address51
	setProxyHash, err := elect.SetProxy(sysTxParams, proxy)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(setProxyHash)
}

// 权益抵押
func TestStake(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	stakeHash, err := elect.Stake(sysTxParams, big.NewInt(20))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(stakeHash)
}

// 权益抵押撤销
func TestUnStake(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	unStakeHash, err := elect.UnStake(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(unStakeHash)
}

// 获取某个地址的权益抵押信息
func TestGetStake(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP51+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	stake, err := elect.GetStake(coinBase)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", stake)
	// {Owner:"did:bid:c935bd29a90fbeea87badf3e", StakeCount:0xea60, Timestamp:0x5f0d2529}
	// {Owner:"did:bid:590ed37615bdfefa496224c7", StakeCount:0xea60, Timestamp:0x5ef5b639}
	// {Owner:"did:bid:6cc796b8d6e2fbebc9b3cf9e", StakeCount:0xea60, Timestamp:0x5ef5b55d}
	// {Owner:"did:bid:13803fb30b7e95d57103c2dc", StakeCount:0xea60, Timestamp:0x5ef5b5f5}
	// {Owner:"did:bid:c117c1794fc7a27bd301ae52", StakeCount:0xea60, Timestamp:0x5ef5b619}
}

func TestGetRestBIFBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP52+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	restBounty, err := elect.GetRestBIFBounty()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(restBounty)
}

func TestElectionExtractOwnBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	extractHash, err := elect.ExtractOwnBounty(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(extractHash)
}

// ？？？？？？？？？？？这个没有输出日志吗？？？？？
func TestIssueAdditionalBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(coinBase, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = common.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	elect := connection.System.NewElection()

	issueAdditionalBountyHash, err := elect.IssueAdditionalBounty(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(issueAdditionalBountyHash)
}

// 获取指定候选者的基本信息
func TestElectGetVoter(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	voter, err := elect.GetVoter(coinBase)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", voter)
	// {Owner:"did:bid:c935bd29a90fbeea87badf3e", IsProxy:true, ProxyVoteCount:0x0, Proxy:"0000000000000000000000000000000000000000", LastVoteCount:0x0, Timestamp:0x5f0d20f9, VoteCandidatesList:""}
	// {Owner:"did:bid:590ed37615bdfefa496224c7", IsProxy:false, ProxyVoteCount:0x0, Proxy:"0000000000000000000000000000000000000000", LastVoteCount:0x1e7df, Timestamp:0x5ef5b6a9, VoteCandidatesList:"did:bid:590ed37615bdfefa496224c7"}
	// {Owner:"did:bid:6cc796b8d6e2fbebc9b3cf9e", IsProxy:false, ProxyVoteCount:0x0, Proxy:"0000000000000000000000000000000000000000", LastVoteCount:0x1e7df, Timestamp:0x5ef5b735, VoteCandidatesList:"did:bid:6cc796b8d6e2fbebc9b3cf9e"}
	// {Owner:"did:bid:13803fb30b7e95d57103c2dc", IsProxy:false, ProxyVoteCount:0x0, Proxy:"0000000000000000000000000000000000000000", LastVoteCount:0x1e7df, Timestamp:0x5ef5b6f1, VoteCandidatesList:"did:bid:13803fb30b7e95d57103c2dc"}
	// {Owner:"did:bid:c117c1794fc7a27bd301ae52", IsProxy:false, ProxyVoteCount:0x0, Proxy:"0000000000000000000000000000000000000000", LastVoteCount:0x1e7df, Timestamp:0x5ef5b6d1, VoteCandidatesList:"did:bid:c117c1794fc7a27bd301ae52"}

}

// 获取指定的候选者的投票列表
func TestGetVoterList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	voterList, err := elect.GetVoterList(resources.Address52)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", voterList)
}
