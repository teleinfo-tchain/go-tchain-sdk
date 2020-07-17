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

// 信任锚的AbiJson数据
const TrustAnchorAbiJSON = `[
{"constant": false,"name":"registerTrustAnchor","inputs":[{"name":"anchor","type":"string"},{"name":"anchorType","type":"uint64"},{"name":"anchorName","type":"string"},{"name":"company","type":"string"},{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"unRegisterTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"updateAnchorInfo","inputs":[{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"extractOwnBounty","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"trustAnchorEvent","type":"event"},
{"constant": false,"name":"voteElect","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelVote","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
]`

// Anchor - The Anchor Module
type Anchor struct {
	super *System
	abi   abi.ABI
}

// NewTrustAnchor - 初始化Anchor
func (sys *System) NewTrustAnchor() *Anchor {
	parsedAbi, _ := abi.JSON(strings.NewReader(TrustAnchorAbiJSON))

	anchor := new(Anchor)
	anchor.abi = parsedAbi
	anchor.super = sys
	return anchor
}

/*
 RegisterTrustAnchor: 注册信任锚，刚刚注册的信任锚都是扩展信任锚，但是如果10类型的信任锚，经过超级节点投票，大于2/3的超级节点同意，可以变成根信任锚。根信任锚需要抵押1000积分，扩展信任锚需要抵押100积分。

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 如果是注册根信任锚，必须是超级节点才可以注册
*/
func (anc *Anchor) RegisterTrustAnchor(from common.Address, registerAnchor *dto.RegisterAnchor) (string, error) {
	// encoding
	// RegisterAnchor is a struct we need to use the components.
	var values []interface{}
	values = anc.super.structToInterface(*registerAnchor, values)
	inputEncode, err := anc.abi.Pack("registerTrustAnchor", values...)
	if err != nil {
		return "", err
	}

	transaction := anc.super.prePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.sendTransaction(transaction)
}

/*
 UnRegisterTrustAnchor: 注销自己的信任锚，自动退回抵押。但是，需要手动批量吊销自己颁发的证书，如果存在未吊销的证书，则抵押不退回。

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 只能注销自己的信任锚
*/
func (anc *Anchor) UnRegisterTrustAnchor(from common.Address, anchor string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("unRegisterTrustAnchor", anchor)
	if err != nil {
		return "", err
	}

	transaction := anc.super.prePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.sendTransaction(transaction)
}

/*
 IsBaseTrustAnchor: 查询bid地址是否为根信任锚

 Returns： bool，true为是根信任锚，false为不是根信任锚
*/
func (anc *Anchor) IsBaseTrustAnchor(address string) (bool, error) {
	params := make([]string, 1)
	params[0] = address

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_isBaseTrustAnchor", params)
	if err != nil {
		return false, err
	}

	return pointer.ToIsBaseTrustAnchor()
}

/*
 IsTrustAnchor: 查询bid地址是否为信任锚

 Returns：  bool，true为是信任锚，false为不是信任锚
*/
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

/*
 UpdateAnchorInfo: 更新信任锚数据

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 只能修改自己的
*/
func (anc *Anchor) UpdateAnchorInfo(from common.Address, extendAnchorInfo *dto.UpdateAnchorInfo) (string, error) {
	// encoding
	// extendAnchorInfo is a struct we need to use the components.
	var values []interface{}
	values = anc.super.structToInterface(*extendAnchorInfo, values)
	inputEncode, err := anc.abi.Pack("updateAnchorInfo", values...)
	if err != nil {
		return "", err
	}

	transaction := anc.super.prePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.sendTransaction(transaction)
}

/*
 ExtractOwnBounty: 提取信任锚激励，只有超过100积分，且24小时内只能提取一次

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 只能提取自己的
*/
func (anc *Anchor) ExtractOwnBounty(from common.Address) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("extractOwnBounty")
	if err != nil {
		return "", err
	}

	transaction := anc.super.prePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.sendTransaction(transaction)
}

/*
 GetTrustAnchor: 查询信任锚信息

 Returns： *dto.TrustAnchor
*/
func (anc *Anchor) GetTrustAnchor(anchor string) (*dto.TrustAnchor, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_trustAnchor", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToTrustAnchor()
}

/*
 GetTrustAnchorStatus: 查询信任锚状态

 Returns： uint64，0未知，1可用，2错误，3删除
*/
func (anc *Anchor) GetTrustAnchorStatus(anchor string) (uint64, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_trustAnchorStatus", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToTrustAnchorStatus()

}

/*
 GetCertificateList: 查询信任锚颁发的证书列表

 Returns： []string， 证书列表
*/
func (anc *Anchor) GetCertificateList(anchor string) ([]string, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_certificateList", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToTrustAnchorCertificateList()

}

/*
 GetBaseList: 查询根信任锚列表

 Returns： []string，根信任锚列表
*/
func (anc *Anchor) GetBaseList() ([]string, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_baseList", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBaseTrustAnchor()
}

/*
 GetBaseNum: 查询根信任锚个数

 Returns： uint64， 根信任锚个数
*/
func (anc *Anchor) GetBaseNum() (uint64, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_baseNumber", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToBaseTrustAnchorNumber()
}

/*
 GetExpendList: 查询扩展信任锚列表

 Returns： []string， 扩展信任锚列表
*/
func (anc *Anchor) GetExpendList() ([]string, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_expendList", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToExpendTrustAnchor()
}

/*
 GetExpendList: 查询扩展信任锚个数

 Returns： uint64， 扩展信任锚个数
*/
func (anc *Anchor) GetExpendNum() (uint64, error) {
	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_enpendNumber", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToExpendTrustAnchorNumber()
}

/*
 VoteElect: 向信任锚投支持票

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 只有超级节点才可投票
*/
func (anc *Anchor) VoteElect(from common.Address, candidate string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("voteElect", candidate)
	if err != nil {
		return "", err
	}

	transaction := anc.super.prePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.sendTransaction(transaction)
}

/*
 CancelVote: 向信任锚投反对票

 Returns： transactionHash，32 Bytes - 交易哈希，如果交易尚不可用，则为零哈希。

 Call permissions: 只有超级节点才可投票
*/
func (anc *Anchor) CancelVote(from common.Address, candidate string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("cancelVote", candidate)
	if err != nil {
		return "", err
	}

	transaction := anc.super.prePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))

	return anc.super.sendTransaction(transaction)
}

/*
 GetVoter: 查询投票人信息

 Returns： []dto.TrustAnchorVoter， 投票人信息
*/
func (anc *Anchor) GetVoter(voterAddress string) ([]dto.TrustAnchorVoter, error) {
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.RequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_voter", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToTrustAnchorVoter()
}
