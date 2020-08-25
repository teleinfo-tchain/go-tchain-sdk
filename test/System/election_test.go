package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"testing"
)

const (
	isSM2Elect       = false
	passwordElect    = "teleinfo"
	testAddressElect = "did:bid:EFTTQWPMdtghuZByPsfQAUuPkWkWYb"
	testAddressFile  = "../resources/superNodeKeyStore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb"
)

// 注册成为候选者节点
func TestRegisterWitness(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	elect := connection.System.NewElection()

	registerWitness := new(dto.RegisterWitness)
	// 但是需要
	registerWitness.NodeUrl = "/ip4/169.254.187.66/tcp/30303/p2p/16Uiu2HAkwviNXPoPHBkZxpg8nURQPiNVeCB9HrocfhXTRCs8j34z"
	registerWitness.Website = "https://www.tele.com"
	registerWitness.Name = "test"

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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	candidate, err := elect.GetCandidate(testAddressElect)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%#v \n", candidate)
}

// 获取所有的见证候选节点
func TestGetAllCandidates(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
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

	t.Logf("%#v \n", allCandidates)
}

func TestVoteWitnesses(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	elect := connection.System.NewElection()

	candidate := "did:bid:6cc796b8d6e2fbebc9b3cf9e"
	voteWitnessHash, err := elect.VoteWitnesses(sysTxParams, candidate)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(voteWitnessHash)
}

func TestElectCancelVote(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	elect := connection.System.NewElection()

	proxy := resources.NewAddrZ
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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	stake, err := elect.GetStake(testAddressElect)

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

//  这个接口暂时不用测试，链的后台需要修改
func TestGetRestBIFBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	elect := connection.System.NewElection()

	extractHash, err := elect.ExtractOwnBounty(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(extractHash)
}

func TestIssueAdditionalBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Elect
	sysTxParams.Password = passwordElect
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()
	voter, err := elect.GetVoter(testAddressElect)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", voter)
}

// 获取指定的候选者的投票列表
func TestGetVoterList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	elect := connection.System.NewElection()

	voterList, err := elect.GetVoterList(testAddressElect)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", voterList)
}
