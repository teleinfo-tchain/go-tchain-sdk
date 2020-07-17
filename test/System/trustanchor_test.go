package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)
const (
	BaseAnchorAddr = "did:bid:a3fa9bb1b84e722f30dbda8c"
	ExtendAnchorAddr = "did:bid:c935bd29a90fbeea87badf3e"
	BaseAnchorType = 10
	ExtendAnchorType = 11
	TrustAnchorName = "trustAnchor"
	TrustAnchorCompany = "teleinfo"
	TrustAnchorCompanyUrl = "www.teleinfo.cn"
	TrustAnchorWebsite =  "www.server.teleinfo.cn"
	TrustAnchorDocumentUrl = "www.doc.teleinfo.cn"
	TrustAnchorServerUrl = "1.1.1.1"
	TrustAnchorEmail = "www.email.teleinfo.cn"
	TrustAnchorDesc = "info test"
)

func TestRegisterTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	registerBaseAnchor := new(dto.RegisterAnchor)
	registerBaseAnchor.Anchor = common.StringToAddress(BaseAnchorAddr).String()
	registerBaseAnchor.AnchorType = BaseAnchorType
	registerBaseAnchor.AnchorName = TrustAnchorName
	registerBaseAnchor.Company = TrustAnchorCompany
	registerBaseAnchor.CompanyUrl = TrustAnchorCompanyUrl
	registerBaseAnchor.Website = TrustAnchorWebsite
	registerBaseAnchor.DocumentUrl = TrustAnchorDocumentUrl
	registerBaseAnchor.ServerUrl = TrustAnchorServerUrl
	registerBaseAnchor.Email = TrustAnchorEmail
	registerBaseAnchor.Desc = TrustAnchorDesc
	//registerBaseAnchor
	registerBaseAnchorHash, err := anchor.RegisterTrustAnchor(common.StringToAddress(coinBase), registerBaseAnchor)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(registerBaseAnchorHash)

	registerExtendAnchor := new(dto.RegisterAnchor)
	registerExtendAnchor.Anchor = common.StringToAddress(ExtendAnchorAddr).String()
	registerExtendAnchor.AnchorType = ExtendAnchorType
	registerExtendAnchor.AnchorName = TrustAnchorName
	registerExtendAnchor.Company = TrustAnchorCompany
	registerExtendAnchor.CompanyUrl = TrustAnchorCompanyUrl
	registerExtendAnchor.Website = TrustAnchorWebsite
	registerExtendAnchor.DocumentUrl = TrustAnchorDocumentUrl
	registerExtendAnchor.ServerUrl = TrustAnchorServerUrl
	registerExtendAnchor.Email = TrustAnchorEmail
	registerExtendAnchor.Desc = TrustAnchorDesc
	//registerExtendAnchor
	registerExtendAnchorHash, err := anchor.RegisterTrustAnchor(common.StringToAddress(coinBase), registerExtendAnchor)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(registerExtendAnchorHash)
}

func TestUnRegisterTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	transactionHash, err := anchor.UnRegisterTrustAnchor(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestIsBaseTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	baseAnchor, err := anchor.IsBaseTrustAnchor(coinBase)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Log(baseAnchor)
}

func TestIsTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	trustAnchor, err := anchor.IsTrustAnchor(coinBase)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchor)
}

func TestUpdateAnchorInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	extendAnchorInfo := new(dto.UpdateAnchorInfo)
	extendAnchorInfo.CompanyUrl = "0www.teleinfo.cn"
	extendAnchorInfo.Website = "0www.server.teleinfo.cn"
	extendAnchorInfo.DocumentUrl = "0www.doc.teleinfo.cn"
	extendAnchorInfo.ServerUrl = "01.1.1.1"
	extendAnchorInfo.Email = "0www.email.teleinfo.cn"
	extendAnchorInfo.Desc = "0info test"

	transactionHash, err := anchor.UpdateAnchorInfo(common.StringToAddress(coinBase), extendAnchorInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestExtractOwnBounty(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	transactionHash, err := anchor.ExtractOwnBounty(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestGetTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
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
	coinBase, err := connection.Core.GetCoinbase()
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
	coinBase, err := connection.Core.GetCoinbase()
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
	_, err := connection.Core.GetCoinbase()
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
	_, err := connection.Core.GetCoinbase()
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
	_, err := connection.Core.GetCoinbase()
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
	_, err := connection.Core.GetCoinbase()
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
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()

	transactionHash, err := anchor.VoteElect(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestCancelVote(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor := connection.System.NewTrustAnchor()
	transactionHash, err := anchor.CancelVote(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestGetVoter(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
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

