package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	registerBaseAnchor := new(system.RegisterAnchor)
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

	registerExtendAnchor := new(system.RegisterAnchor)
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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transactionHash, err := anchor.UnRegisterTrustAnchor(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestIsTrustAnchor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	trustAnchor, err := anchor.IsTrustAnchor(common.StringToAddress(coinBase), coinBase)
	if err != nil && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}
	if err == system.ErrCertificateNotExist {
		t.Log("not ", err)
	}
	t.Log(trustAnchor)
}

func TestUpdateBaseAnchorInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	baseAnchorInfo := new(system.BaseAnchorInfo)
	baseAnchorInfo.Anchor = BaseAnchorAddr
	baseAnchorInfo.CompanyUrl = "www.teleinfoTest.cn"
	baseAnchorInfo.Website = "www.server.teleinfoTest.cn"
	baseAnchorInfo.DocumentUrl = "www.doc.teleinfoTest.cn"
	baseAnchorInfo.ServerUrl = "2.2.2.2"
	baseAnchorInfo.Email = "www.email.teleinfoTest.cn"
	baseAnchorInfo.Desc = "info test Test"

	transactionHash, err := anchor.UpdateBaseAnchorInfo(common.StringToAddress(coinBase), baseAnchorInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash)
}

func TestUpdateExtendAnchorInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	extendAnchorInfo := new(system.ExtendAnchorInfo)
	extendAnchorInfo.CompanyUrl = "0www.teleinfo.cn"
	extendAnchorInfo.Website = "0www.server.teleinfo.cn"
	extendAnchorInfo.DocumentUrl = "0www.doc.teleinfo.cn"
	extendAnchorInfo.ServerUrl = "01.1.1.1"
	extendAnchorInfo.Email = "0www.email.teleinfo.cn"
	extendAnchorInfo.Desc = "0info test"

	transactionHash, err := anchor.UpdateExtendAnchorInfo(common.StringToAddress(coinBase), extendAnchorInfo)
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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transactionHash, err := anchor.ExtractOwnBounty(common.StringToAddress(coinBase), coinBase)
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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	trustAnchor, err := anchor.GetTrustAnchor(common.StringToAddress(coinBase), coinBase)
	if err != nil && err != system.ErrCertificateNotExist {
		t.Error(err)
		t.FailNow()
	}
	if err == system.ErrCertificateNotExist {
		t.Log(err)
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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	trustAnchorStatus, err := anchor.GetTrustAnchorStatus(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(trustAnchorStatus)
}

func TestGetBaseTrustAnchorList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	baseList, err := anchor.GetBaseTrustAnchorList(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(baseList)
}

func TestGetBaseTrustAnchorNum(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	baseListNum, err := anchor.GetBaseTrustAnchorNum(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(baseListNum)
}

func TestGetExpendTrustAnchorList(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expendList, err := anchor.GetExpendTrustAnchorList(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(expendList)
}

func TestGetExpendTrustAnchorNum(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expendListNum, err := anchor.GetExpendTrustAnchorNum(common.StringToAddress(coinBase))
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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

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
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	voterInfo, err := anchor.GetVoter(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(voterInfo)
}

func TestCheckSenderAddress(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	anchor, err := connection.System.NewTrustAnchor()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	superNode, err := anchor.CheckSenderAddress(common.StringToAddress(coinBase), coinBase)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(superNode)
}

