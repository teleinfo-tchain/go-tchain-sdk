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
	BaseAnchorAddr         = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	ExtendAnchorAddr       = "did:bid:EFTTQWPMdtghuZByPsfQAUuPkWkWYb"
	BaseAnchorType         = 10
	ExtendAnchorType       = 11
	TrustAnchorName        = "trustAnchor"
	TrustAnchorCompany     = "teleinfo"
	TrustAnchorCompanyUrl  = "www.teleinfo.cn"
	TrustAnchorWebsite     = "www.server.teleinfo.cn"
	TrustAnchorDocumentUrl = "www.doc.teleinfo.cn"
	TrustAnchorServerUrl   = "1.1.1.1"
	TrustAnchorEmail       = "j203@163.com"
	TrustAnchorDesc        = "info test"
)

const (
	isSM2Trust = false
	passwordTrust = "teleinfo"
	testAddressTrust = resources.TestAddr
	testAddressTrustFile = "../resources/superNodeKeyStore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
)

func TestRegisterBaseTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressTrust, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressTrustFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Trust
	sysTxParams.Password = passwordTrust
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	anchor := connection.System.NewTrustAnchor()

	registerBaseAnchor := new(dto.RegisterAnchor)
	registerBaseAnchor.Anchor = BaseAnchorAddr
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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressTrust, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressTrustFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Trust
	sysTxParams.Password = passwordTrust
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(55)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	anchor := connection.System.NewTrustAnchor()

	registerExtendAnchor := new(dto.RegisterAnchor)
	registerExtendAnchor.Anchor = ExtendAnchorAddr
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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressTrust, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressTrustFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Trust
	sysTxParams.Password = passwordTrust
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	baseAnchor, err := anchor.IsBaseTrustAnchor(ExtendAnchorAddr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(baseAnchor)
}

func TestIsTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	trustAnchor, err := anchor.IsTrustAnchor(ExtendAnchorAddr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchor)
}

func TestUpdateAnchorInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressTrust, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressTrustFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Trust
	sysTxParams.Password = passwordTrust
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	anchor := connection.System.NewTrustAnchor()

	extendAnchorInfo := new(dto.UpdateAnchorInfo)
	extendAnchorInfo.CompanyUrl = "0www.teleinfo.cn"
	extendAnchorInfo.Website = "0www.server.teleinfo.cn"
	extendAnchorInfo.DocumentUrl = "0www.doc.teleinfo.cn"
	extendAnchorInfo.ServerUrl = "01.1.1.1"
	extendAnchorInfo.Email = "j23@163.com"
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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressTrust, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressTrustFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Trust
	sysTxParams.Password = passwordTrust
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

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
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	trustAnchor, err := anchor.GetTrustAnchor(testAddressTrust)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchor)
}

func TestGetTrustAnchorStatus(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	trustAnchorStatus, err := anchor.GetTrustAnchorStatus(testAddressTrust)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchorStatus)
}

func TestGetCertificateList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	certificateLi, err := anchor.GetCertificateList(testAddressTrust)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(certificateLi)
}

func TestGetBaseList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
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
	_, err := connection.Core.GetChainId()
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
	_, err := connection.Core.GetChainId()
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

	// [did:bid:6cc796b8d6e2fbebc9b3cf9e did:bid:c935bd29a90fbeea87badf3e]
	t.Log(expendList)
}

func TestGetExpendNum(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
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
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressTrust, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressTrustFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Trust
	sysTxParams.Password = passwordTrust
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	anchor := connection.System.NewTrustAnchor()

	transactionHash, err := anchor.VoteElect(sysTxParams, testAddressTrust)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestCancelVote(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	nonce, err := connection.Core.GetTransactionCount(testAddressTrust, block.LATEST)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(testAddressTrustFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(keyFileData) == 0 {
		t.Errorf("keyFileData can't be empty")
		t.FailNow()
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.IsSM2 = isSM2Trust
	sysTxParams.Password = passwordTrust
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(35)
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
	sysTxParams.ChainId = chainId
	sysTxParams.Version = 1

	anchor := connection.System.NewTrustAnchor()
	transactionHash, err := anchor.CancelVote(sysTxParams, testAddressTrust)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestGetVoter(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()
	trustAnchorVoterLi, err := anchor.GetVoter(testAddressTrust)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchorVoterLi)
}
