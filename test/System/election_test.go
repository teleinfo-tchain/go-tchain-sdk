package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"strconv"
	"testing"
)

//// 注册成为候选者节点
//func TestRegisterWitness(t *testing.T) {
//	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
//	chainId, err := connection.Core.GetChainId()
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressElect, block.LATEST)
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	// keyFileData 还可以进一步校验
//	keyFileData, err := ioutil.ReadFile(resources.TestAddressFile)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	if len(keyFileData) == 0 {
//		t.Errorf("keyFileData can't be empty")
//		t.FailNow()
//	}
//
//	sysTxParams := new(system.SysTxParams)
//	sysTxParams.IsSM2 = resources.NotSm2
//	sysTxParams.Password = resources.SystemPassword
//	sysTxParams.KeyFileData = keyFileData
//	sysTxParams.GasPrice = big.NewInt(35)
//	sysTxParams.Gas = 2000000
//	sysTxParams.Nonce = nonce.Uint64()
//	sysTxParams.ChainId = chainId
//
//	elect := connection.System.NewElection()
//
//	registerWitness := new(dto.RegisterWitness)
//	// 但是需要
//	registerWitness.NodeUrl = "/ip4/169.254.187.66/tcp/30303/p2p/16Uiu2HAkwviNXPoPHBkZxpg8nURQPiNVeCB9HrocfhXTRCs8j34z"
//	registerWitness.Website = "https://www.tele.com"
//	registerWitness.Name = "test"
//
//	// registerWitness
//	registerWitnessHash, err := elect.RegisterWitness(sysTxParams, registerWitness)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	t.Log(registerWitnessHash, err)
//}
//
//// 取消成为候选者节点
//func TestUnRegisterWitness(t *testing.T) {
//	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
//	_, err := connection.Core.GetChainId()
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressElect, block.LATEST)
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	// keyFileData 还可以进一步校验
//	keyFileData, err := ioutil.ReadFile(resources.TestAddressFile)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	if len(keyFileData) == 0 {
//		t.Errorf("keyFileData can't be empty")
//		t.FailNow()
//	}
//
//	sysTxParams := new(system.SysTxParams)
//	sysTxParams.IsSM2 = resources.NotSm2
//	sysTxParams.Password = resources.SystemPassword
//	sysTxParams.KeyFileData = keyFileData
//	sysTxParams.GasPrice = big.NewInt(35)
//	sysTxParams.Gas = 2000000
//	sysTxParams.Nonce = nonce.Uint64()
//
//	elect := connection.System.NewElection()
//
//	unRegisterWitnessHash, err := elect.UnRegisterWitness(sysTxParams)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	t.Log(unRegisterWitnessHash, err)
//}
//
//// 获取候选节点基本信息
//func TestGetCandidate(t *testing.T) {
//	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
//	_, err := connection.Core.GetChainId()
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//
//	elect := connection.System.NewElection()
//
//	candidate, err := elect.GetCandidate(resources.TestAddressElect)
//
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	t.Logf("%#v \n", candidate)
//}
//
//// 获取所有的见证候选节点
//func TestGetAllCandidates(t *testing.T) {
//	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
//	_, err := connection.Core.GetChainId()
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//
//	elect := connection.System.NewElection()
//
//	allCandidates, err := elect.GetAllCandidates()
//
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//
//	t.Logf("%#v \n", allCandidates)
//}
//
//func TestVoteWitnesses(t *testing.T) {
//	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
//	chainId, err := connection.Core.GetChainId()
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressElect, block.LATEST)
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	// keyFileData 还可以进一步校验
//	keyFileData, err := ioutil.ReadFile(resources.TestAddressFile)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	if len(keyFileData) == 0 {
//		t.Errorf("keyFileData can't be empty")
//		t.FailNow()
//	}
//
//	sysTxParams := new(system.SysTxParams)
//	sysTxParams.IsSM2 = resources.NotSm2
//	sysTxParams.Password = resources.SystemPassword
//	sysTxParams.KeyFileData = keyFileData
//	sysTxParams.Gas = 2000000
//	sysTxParams.Nonce = nonce.Uint64()
//	sysTxParams.ChainId = chainId
//
//	elect := connection.System.NewElection()
//
//	candidate := "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
//	voteWitnessHash, err := elect.VoteWitnesses(sysTxParams, candidate)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	t.Log(voteWitnessHash)
//}
//
//func TestElectCancelVote(t *testing.T) {
//	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
//	chainId, err := connection.Core.GetChainId()
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressElect, block.LATEST)
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//
//	// keyFileData 还可以进一步校验
//	keyFileData, err := ioutil.ReadFile(resources.TestAddressFile)
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	if len(keyFileData) == 0 {
//		t.Errorf("keyFileData can't be empty")
//		t.FailNow()
//	}
//
//	sysTxParams := new(system.SysTxParams)
//	sysTxParams.IsSM2 = resources.NotSm2
//	sysTxParams.Password = resources.SystemPassword
//	sysTxParams.KeyFileData = keyFileData
//	sysTxParams.GasPrice = big.NewInt(35)
//	sysTxParams.Gas = 2000000
//	sysTxParams.Nonce = nonce.Uint64()
//	sysTxParams.ChainId = chainId
//
//	elect := connection.System.NewElection()
//
//	cancelVoteHash, err := elect.CancelVote(sysTxParams)
//
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//
//	t.Log(cancelVoteHash)
//}

func TestGetRestBIFBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
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
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = resources.NotSm2
	sysTxParams.Password = resources.SystemPassword
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	elect := connection.System.NewElection()

	extractHash, err := elect.ExtractOwnBounty(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(extractHash)
}

func TestIssueAdditionalBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(resources.TestAddressElect, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(resources.TestAddressFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = resources.NotSm2
	sysTxParams.Password = resources.SystemPassword
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId

	elect := connection.System.NewElection()

	issueAdditionalBountyHash, err := elect.IssueAdditionalBounty(sysTxParams)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(issueAdditionalBountyHash)
}

