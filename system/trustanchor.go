package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
	"strings"
)

const (
	TrustAnchorContractAddr = "did:bid:00000000000000000000000c"
	//AnchorStatusUnknow      = 0
	//AnchorStatusOK          = 1
	//AnchorStatusErr         = 2
	//AnchorStatusDelete      = 3
	//
	//BaseAnchor                 = 10
	//ExtendAnchor               = 11
	//UnknowAnchorType           = 2
	//Day                        = 86400        //24*60*60秒，每次提取积分的最小时间间隔
	//Bifer                      = 100000000    //10^9
	//MiniExtractAmount          = 100 * Bifer  //每次提取积分的最小额度是100
	//BaseTrustAnchorPledge      = 1000 * Bifer //注册根信任锚抵押的积分
	//ExtendTrustAnchorPledge    = 100 * Bifer  //注册扩展信任锚抵押的积分
	//IncentivesToExtendIssueCer = 1 * Bifer    //颁发一个证书获得的奖励
	//IncentivesToBaseIssueCer   = 2 * Bifer    //颁发一个证书获得的奖励
	//EmptyUrl                   = ""
)

//var (
//	ErrIllExtracAmout    = errors.New("积分总额不足100，无法提取")
//	ErrIllExtracTime     = errors.New("距离上次提取不足24小时，无法提取")
//	ErrIllAnchorType     = errors.New("未知的信任锚类型，10 代表根信任锚，11代表扩展信任锚")
//	ErrAnchorExist       = errors.New("该地址已存在于信任锚列表中")
//	ErrAnchorNoExist     = errors.New("信任锚不存在")
//	ErrIllegalAnchor     = errors.New("信任锚字段不能为空")
//	ErrIllegalBalance    = errors.New("账户内积分不足，注册根信任锚需要抵押1000积分，注册扩展信任锚需要抵押100积分")
//	ErrIllegalVote       = errors.New("信任锚不存在或不是基础信任锚")
//	ErrIllegalRepeatVote = errors.New("同一个基础信任锚能且仅能投一票")
//)

const TrustAnchorAbiJSON = `[
{"constant": false,"name":"registerTrustAnchor","inputs":[{"name":"anchor","type":"string"},{"name":"anchorType","type":"uint64"},{"name":"anchorName","type":"string"},{"name":"company","type":"string"},{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"unRegisterTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"updateAnchorInfo","inputs":[{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"extractOwnBounty","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"trustAnchorEvent","type":"event"},
{"constant": false,"name":"voteElect","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelVote","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
]`

type Anchor struct {
	super *System
	abi   abi.ABI
}

func (sys *System) NewTrustAnchor() *Anchor {
	parsedAbi, _ := abi.JSON(strings.NewReader(TrustAnchorAbiJSON))

	anchor := new(Anchor)
	anchor.abi = parsedAbi
	anchor.super = sys
	return anchor
}

// 是否让用户可以自定义gas和gasPrice？
//"registerTrustAnchor","inputs":[{"name":"anchor","type":"string"},{"name":"anchorType","type":"uint64"},{"name":"anchorName","type":"string"},{"name":"company","type":"string"},{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[]
func (anc *Anchor) RegisterTrustAnchor(from common.Address, registerAnchor *dto.RegisterAnchor) (string, error) {
	// encoding
	// RegisterAnchor is a struct we need to use the components.
	var values []interface{}
	values = anc.super.StructToInterface(*registerAnchor,values)
	inputEncode, err := anc.abi.Pack("registerTrustAnchor", values...)
	if err != nil {
		return "", err
	}

	transaction := anc.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.SendTransaction(transaction)
}

//"unRegisterTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[]
func (anc *Anchor) UnRegisterTrustAnchor(from common.Address, anchor string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("unRegisterTrustAnchor", anchor)
	if err != nil {
		return "", err
	}

	transaction := anc.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.SendTransaction(transaction)
}

// 基础信任锚，初始化时，为false;扩展信任锚的注册后,即生效，为true
func (anc *Anchor) IsTrustAnchor(address string) (bool, error) {
	params := make([]string, 1)
	params[0] = address

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_isTrustAnchor", params)
	if err != nil {
		return false, err
	}

	return pointer.ToIsTrustAnchor()
}

func (anc *Anchor) UpdateAnchorInfo(from common.Address, extendAnchorInfo *dto.UpdateAnchorInfo) (string, error) {
	// encoding
	// extendAnchorInfo is a struct we need to use the components.
	var values []interface{}
	values = anc.super.StructToInterface(*extendAnchorInfo, values)
	inputEncode, err := anc.abi.Pack("updateAnchorInfo", values...)
	if err != nil {
		return "", err
	}

	transaction := anc.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.SendTransaction(transaction)
}

func (anc *Anchor) ExtractOwnBounty(from common.Address) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("extractOwnBounty")
	if err != nil {
		return "", err
	}

	transaction := anc.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return  anc.super.SendTransaction(transaction)
}

func (anc *Anchor) GetTrustAnchor(anchor string) (*dto.TrustAnchor, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_trustAnchor", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToTrustAnchor()
}

func (anc *Anchor) GetTrustAnchorStatus(anchor string) (uint64, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_trustAnchorStatus", params)
	if err != nil{
		return 0, err
	}

	return pointer.ToTrustAnchorStatus()

}

func (anc *Anchor) GetCertificateList(anchor string)([]string, error){
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_certificateList", params)
	if err != nil{
		return nil, err
	}

	return pointer.ToTrustAnchorCertificateList()

}

func (anc *Anchor) GetBaseList() ([]string, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_baseList", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBaseTrustAnchor()
}

func (anc *Anchor) GetBaseNum() (uint64, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_baseNumber", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToBaseTrustAnchorNumber()
}

func (anc *Anchor) GetExpendList() ([]string, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_expendList", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToExpendTrustAnchor()
}

func (anc *Anchor) GetExpendNum() (uint64, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_enpendNumber", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToExpendTrustAnchorNumber()
}

//向信任锚投支持票
func (anc *Anchor) VoteElect(from common.Address, candidate string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("voteElect", candidate)
	if err != nil {
		return "", err
	}

	transaction := anc.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.SendTransaction(transaction)
}

//向信任锚投反对票
func (anc *Anchor) CancelVote(from common.Address, candidate string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("cancelVote", candidate)
	if err != nil {
		return "", err
	}

	transaction := anc.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.SendTransaction(transaction)
}

//查询投票人信息  !!!!!!!(存在投的人的列表，需要解析下)
func (anc *Anchor) GetVoter(voterAddress string) (*dto.TrustAnchorVoter, error) {
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_voter", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToTrustAnchorVoter()
}
