package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
	"testing"
)

const (
	BaseAnchorAddr         = "did:bid:a3fa9bb1b84e722f30dbda8c"
	ExtendAnchorAddr       = "did:bid:c935bd29a90fbeea87badf3e"
	BaseAnchorType         = 10
	ExtendAnchorType       = 11
	TrustAnchorName        = "trustAnchor"
	TrustAnchorCompany     = "teleinfo"
	TrustAnchorCompanyUrl  = "www.teleinfo.cn"
	TrustAnchorWebsite     = "www.server.teleinfo.cn"
	TrustAnchorDocumentUrl = "www.doc.teleinfo.cn"
	TrustAnchorServerUrl   = "1.1.1.1"
	TrustAnchorEmail       = "www.email.teleinfo.cn"
	TrustAnchorDesc        = "info test"
)

func TestRegisterBaseTrustAnchor(t *testing.T) {
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
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	anchor := connection.System.NewTrustAnchor()

	registerBaseAnchor := new(dto.RegisterAnchor)
	registerBaseAnchor.Anchor = utils.StringToAddress(BaseAnchorAddr).String()
	registerBaseAnchor.AnchorType = BaseAnchorType
	registerBaseAnchor.AnchorName = TrustAnchorName
	registerBaseAnchor.Company = TrustAnchorCompany
	registerBaseAnchor.CompanyUrl = TrustAnchorCompanyUrl
	registerBaseAnchor.Website = TrustAnchorWebsite
	registerBaseAnchor.DocumentUrl = TrustAnchorDocumentUrl
	registerBaseAnchor.ServerUrl = TrustAnchorServerUrl
	registerBaseAnchor.Email = TrustAnchorEmail
	registerBaseAnchor.Desc = TrustAnchorDesc
	// registerBaseAnchor
	registerBaseAnchorHash, err := anchor.RegisterTrustAnchor(sysTxParams, registerBaseAnchor)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(registerBaseAnchorHash)
}

func TestRegisterExtendTrustAnchor(t *testing.T) {
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
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	anchor := connection.System.NewTrustAnchor()

	registerExtendAnchor := new(dto.RegisterAnchor)
	registerExtendAnchor.Anchor = utils.StringToAddress(ExtendAnchorAddr).String()
	registerExtendAnchor.AnchorType = ExtendAnchorType
	registerExtendAnchor.AnchorName = TrustAnchorName
	registerExtendAnchor.Company = TrustAnchorCompany
	registerExtendAnchor.CompanyUrl = TrustAnchorCompanyUrl
	registerExtendAnchor.Website = TrustAnchorWebsite
	registerExtendAnchor.DocumentUrl = TrustAnchorDocumentUrl
	registerExtendAnchor.ServerUrl = TrustAnchorServerUrl
	registerExtendAnchor.Email = TrustAnchorEmail
	registerExtendAnchor.Desc = TrustAnchorDesc
	// registerExtendAnchor
	registerExtendAnchorHash, err := anchor.RegisterTrustAnchor(sysTxParams, registerExtendAnchor)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(registerExtendAnchorHash)
}

func TestUnRegisterTrustAnchor(t *testing.T) {
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
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	anchor := connection.System.NewTrustAnchor()

	transactionHash, err := anchor.UnRegisterTrustAnchor(sysTxParams)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestIsBaseTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	baseAnchor, err := anchor.IsBaseTrustAnchor(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(baseAnchor)
}

func TestIsTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	trustAnchor, err := anchor.IsTrustAnchor(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchor)
}

func TestUpdateAnchorInfo(t *testing.T) {
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
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	anchor := connection.System.NewTrustAnchor()

	extendAnchorInfo := new(dto.UpdateAnchorInfo)
	extendAnchorInfo.CompanyUrl = "0www.teleinfo.cn"
	extendAnchorInfo.Website = "0www.server.teleinfo.cn"
	extendAnchorInfo.DocumentUrl = "0www.doc.teleinfo.cn"
	extendAnchorInfo.ServerUrl = "01.1.1.1"
	extendAnchorInfo.Email = "0www.email.teleinfo.cn"
	extendAnchorInfo.Desc = "0info test"

	transactionHash, err := anchor.UpdateAnchorInfo(sysTxParams, extendAnchorInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestExtractOwnBounty(t *testing.T) {
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
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	anchor := connection.System.NewTrustAnchor()

	transactionHash, err := anchor.ExtractOwnBounty(sysTxParams)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestGetTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	trustAnchor, err := anchor.GetTrustAnchor(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchor)
}

func TestGetTrustAnchorStatus(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	trustAnchorStatus, err := anchor.GetTrustAnchorStatus(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchorStatus)
}

func TestGetCertificateList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	certificateLi, err := anchor.GetCertificateList(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(certificateLi)
}

func TestGetBaseList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	baseList, err := anchor.GetBaseList()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(baseList)
}

func TestGetBaseTrustAnchorNum(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	baseListNum, err := anchor.GetBaseNum()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(baseListNum)
}

func TestGetExpendList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	expendList, err := anchor.GetExpendList()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(expendList)
}

func TestGetExpendNum(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	expendListNum, err := anchor.GetExpendNum()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(expendListNum)
}

// 投票超过2/3才能激活信任锚(现在有5个超级节点，超过2/3就是需要有至少4个投票)
func TestVoteElect(t *testing.T) {
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
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	anchor := connection.System.NewTrustAnchor()

	transactionHash, err := anchor.VoteElect(sysTxParams, coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestCancelVote(t *testing.T) {
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
	sysTxParams.From = utils.StringToAddress(coinBase)
	sysTxParams.PrivateKey = resources.CoinBasePriKey
	sysTxParams.Gas = 2000000
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)

	anchor := connection.System.NewTrustAnchor()
	transactionHash, err := anchor.CancelVote(sysTxParams, coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestGetVoter(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()
	trustAnchorVoterLi, err := anchor.GetVoter(coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchorVoterLi)
}
