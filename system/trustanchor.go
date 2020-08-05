package system

import (
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils"
	"strings"
)

const (
	TrustAnchorContractAddr = "did:bid:00000000000000000000000c"
	// AnchorStatusUnknow      = 0
	// AnchorStatusOK          = 1
	// AnchorStatusErr         = 2
	// AnchorStatusDelete      = 3
	//
	// BaseAnchor                 = 10
	// ExtendAnchor               = 11
	// UnknowAnchorType           = 2
	// Day                        = 86400        //24*60*60秒，每次提取积分的最小时间间隔
	// Bifer                      = 100000000    //10^9
	// MiniExtractAmount          = 100 * Bifer  //每次提取积分的最小额度是100
	// BaseTrustAnchorPledge      = 1000 * Bifer //注册根信任锚抵押的积分
	// ExtendTrustAnchorPledge    = 100 * Bifer  //注册扩展信任锚抵押的积分
	// IncentivesToExtendIssueCer = 1 * Bifer    //颁发一个证书获得的奖励
	// IncentivesToBaseIssueCer   = 2 * Bifer    //颁发一个证书获得的奖励
	// EmptyUrl                   = ""
)

// var (
//	ErrIllExtracAmout    = errors.New("积分总额不足100，无法提取")
//	ErrIllExtracTime     = errors.New("距离上次提取不足24小时，无法提取")
//	ErrIllAnchorType     = errors.New("未知的信任锚类型，10 代表根信任锚，11代表扩展信任锚")
//	ErrAnchorExist       = errors.New("该地址已存在于信任锚列表中")
//	ErrAnchorNoExist     = errors.New("信任锚不存在")
//	ErrIllegalAnchor     = errors.New("信任锚字段不能为空")
//	ErrIllegalBalance    = errors.New("账户内积分不足，注册根信任锚需要抵押1000积分，注册扩展信任锚需要抵押100积分")
//	ErrIllegalVote       = errors.New("信任锚不存在或不是基础信任锚")
//	ErrIllegalRepeatVote = errors.New("同一个基础信任锚能且仅能投一票")
// )

// 信任锚的AbiJson数据
const TrustAnchorAbiJSON = `[
{"constant": false,"name":"registerTrustAnchor","inputs":[{"name":"anchor","type":"string"},{"name":"anchorType","type":"uint64"},{"name":"anchorName","type":"string"},{"name":"company","type":"string"},{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"unRegisterTrustAnchor","inputs":[],"outputs":[],"type":"function"},
{"constant": false,"name":"updateAnchorInfo","inputs":[{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"extractOwnBounty","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"trustAnchorEvent","type":"event"}]
{"constant": false,"name":"voteElect","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelVote","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"}
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
  RegisterTrustAnchor:
   	EN -
	CN - 注册信任锚，刚刚注册的信任锚都是扩展信任锚，但是如果10类型的信任锚，经过超级节点投票，大于2/3的超级节点同意，可以变成根信任锚。根信任锚需要抵押1000积分，扩展信任锚需要抵押100积分。
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- registerAnchor: *dto.RegisterAnchor，包含注册信任锚的信息
		Anchor      string // 信任锚bid
		AnchorType  uint64 // 信任锚的类型，10为根信任锚，11为扩展信任锚
		AnchorName  string // 信任锚名称，不含敏感词的字符串
		Company     string // 公司名
		CompanyUrl  string // 公司网址
		Website     string // 信任锚网址
		DocumentUrl string // 信任锚接口字段文档
		ServerUrl   string // 服务链接
		Email       string // 邮箱地址 email没有做格式校验，在sdk中做？？
		Desc        string // 描述

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 如果是注册根信任锚，必须是超级节点才可以注册
*/
func (anc *Anchor) RegisterTrustAnchor(signTxParams *SysTxParams, registerAnchor *dto.RegisterAnchor) (string, error) {
	// encoding
	// RegisterAnchor is a struct we need to use the components.
	var values []interface{}
	values = anc.super.structToInterface(*registerAnchor, values)
	inputEncode, err := anc.abi.Pack("registerTrustAnchor", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := anc.super.prePareSignTransaction(signTxParams, inputEncode, utils.StringToAddress(TrustAnchorContractAddr))
	if err != nil {
		return "", err
	}

	return anc.super.sendRawTransaction(signedTx)
}

/*
  UnRegisterTrustAnchor:
   	EN -
	CN - 注销自己的信任锚，自动退回抵押。但是，需要手动批量吊销自己颁发的证书，如果存在未吊销的证书，则抵押不退回。
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只能注销自己的信任锚
*/
func (anc *Anchor) UnRegisterTrustAnchor(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := anc.abi.Pack("unRegisterTrustAnchor")

	signedTx, err := anc.super.prePareSignTransaction(signTxParams, inputEncode, utils.StringToAddress(TrustAnchorContractAddr))
	if err != nil {
		return "", err
	}

	return anc.super.sendRawTransaction(signedTx)
}

/*
  IsBaseTrustAnchor:
   	EN -
	CN - 查询bid地址是否为根信任锚
  Params:
  	- anchor: string，信任锚bid

  Returns:
  	- bool，true为是根信任锚，false为不是根信任锚
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) IsBaseTrustAnchor(anchor string) (bool, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_isBaseTrustAnchor", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
  IsTrustAnchor:
   	EN -
	CN - 查询bid地址是否为信任锚
  Params:
  	- anchor: string，信任锚bid

  Returns:
  	- bool，true为是信任锚，false为不是信任锚
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) IsTrustAnchor(anchor string) (bool, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_isTrustAnchor", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
  UpdateAnchorInfo:
   	EN -
	CN - 更新信任锚数据
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- extendAnchorInfo: *dto.UpdateAnchorInfo，更新的信任锚数据信息

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。。
	- error

  Call permissions: 只能修改自己的
*/
func (anc *Anchor) UpdateAnchorInfo(signTxParams *SysTxParams, extendAnchorInfo *dto.UpdateAnchorInfo) (string, error) {
	// encoding
	// extendAnchorInfo is a struct we need to use the components.
	var values []interface{}
	values = anc.super.structToInterface(*extendAnchorInfo, values)
	inputEncode, err := anc.abi.Pack("updateAnchorInfo", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := anc.super.prePareSignTransaction(signTxParams, inputEncode, utils.StringToAddress(TrustAnchorContractAddr))
	if err != nil {
		return "", err
	}

	return anc.super.sendRawTransaction(signedTx)
}

/*
  ExtractOwnBounty:
   	EN -
	CN - 提取信任锚激励，只有超过100积分，且24小时内只能提取一次
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只能提取自己的
*/
func (anc *Anchor) ExtractOwnBounty(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := anc.abi.Pack("extractOwnBounty")

	signedTx, err := anc.super.prePareSignTransaction(signTxParams, inputEncode, utils.StringToAddress(TrustAnchorContractAddr))
	if err != nil {
		return "", err
	}

	return anc.super.sendRawTransaction(signedTx)
}

/*
  GetTrustAnchor:
   	EN -
	CN - 查询信任锚信息
  Params:
  	- anchor: string，信任锚bid

  Returns:
  	- *dto.TrustAnchor
		Id               string   `json:"id"              gencodec:"required"`   //信任锚BID地址
		Name             string   `json:"name"            gencodec:"required"`   //信任锚名称
		Company          string   `json:"company"         gencodec:"required"`   //信任锚所属公司
		CompanyUrl       string   `json:"company_url"     gencodec:"required"`   //公司网址
		Website          string   `json:"website"         gencodec:"required"`   //信任锚网址
		ServerUrl        string   `json:"server_url"      gencodec:"required"`   //服务链接
		DocumentUrl      string   `json:"document_url"    gencodec:"required"`   //信任锚接口字段文档
		Email            string   `json:"email"           gencodec:"required"`   //信任锚客服邮箱
		Desc             string   `json:"desc" gencodec:"required"`              //描述
		TrustAnchorType  uint64   `json:"type"            gencodec:"required"`   //信任锚类型
		Status           uint64   `json:"status"          gencodec:"required"`   //服务状态
		Active           bool     `json:"active"          gencodec:"required"`   //是否是根信任锚
		TotalBounty      *big.Int `json:"totalBounty"     gencodec:"required"`   //总激励
		ExtractedBounty  *big.Int `json:"extractedBounty" gencodec:"required"`   //已提取激励
		LastExtractTime  *big.Int `json:"lastExtractTime" gencodec:"required"`   //上次提取时间
		VoteCount        *big.Int `json:"vote_count" gencodec:"required"`        //得票数
		Stake            *big.Int `json:"stake" gencodec:"required"`             //抵押
		CreateDate       *big.Int `json:"create_date" gencodec:"required"`       //创建时间
		CertificateCount *big.Int `json:"certificate_count" gencodec:"required"` //证书总数
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetTrustAnchor(anchor string) (*dto.TrustAnchor, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_trustAnchor", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToTrustAnchor()
}

/*
  GetTrustAnchorStatus:
   	EN -
	CN - 查询信任锚状态
  Params:
  	- anchor: string，信任锚bid

  Returns:
  	- uint64，0未知，1可用，2错误，3删除
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetTrustAnchorStatus(anchor string) (uint64, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_trustAnchorStatus", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()

}

/*
  GetCertificateList:
   	EN -
	CN - 查询信任锚颁发的证书列表
  Params:
  	- anchor: string，信任锚bid
  	-

  Returns:
  	- []string， 证书列表
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetCertificateList(anchor string) ([]string, error) {
	params := make([]string, 1)
	params[0] = anchor

	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_certificateList", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()

}

/*
  GetBaseList:
   	EN -
	CN - 查询根信任锚列表
  Params:
  	- None

  Returns:
  	- []string，根信任锚列表
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetBaseList() ([]string, error) {
	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_baseList", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}

/*
  GetBaseNum:
   	EN -
	CN - 查询根信任锚个数
  Params:
  	- None

  Returns:
  	- uint64， 根信任锚个数
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetBaseNum() (uint64, error) {
	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_baseNumber", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
  GetExpendList:
   	EN -
	CN - 查询扩展信任锚列表
  Params:
  	- None
  	-

  Returns:
  	- []string， 扩展信任锚列表
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetExpendList() ([]string, error) {
	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_expendList", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}

/*
  GetExpendNum:
   	EN -
	CN - 查询扩展信任锚个数
  Params:
  	- None

  Returns:
  	- uint64， 扩展信任锚个数
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetExpendNum() (uint64, error) {
	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_enpendNumber", nil)

	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}

/*
  VoteElect:
   	EN -
	CN - 向信任锚投支持票
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- candidate: string，信任锚地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有超级节点才可投票
*/
func (anc *Anchor) VoteElect(signTxParams *SysTxParams, candidate string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("voteElect", candidate)
	if err != nil {
		return "", err
	}

	signedTx, err := anc.super.prePareSignTransaction(signTxParams, inputEncode, utils.StringToAddress(TrustAnchorContractAddr))
	if err != nil {
		return "", err
	}

	return anc.super.sendRawTransaction(signedTx)
}

/*
  CancelVote:
   	EN -
	CN - 向信任锚投反对票
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- candidate: string，信任锚地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有超级节点才可投票
*/
func (anc *Anchor) CancelVote(signTxParams *SysTxParams, candidate string) (string, error) {
	// encoding
	inputEncode, err := anc.abi.Pack("cancelVote", candidate)
	if err != nil {
		return "", err
	}

	signedTx, err := anc.super.prePareSignTransaction(signTxParams, inputEncode, utils.StringToAddress(TrustAnchorContractAddr))
	if err != nil {
		return "", err
	}

	return anc.super.sendRawTransaction(signedTx)
}

/*
  GetVoter:
   	EN -
	CN - 查询投票人信息
  Params:
  	- voterAddress: 投票人地址（也就是超级节点地址，因为只有超级节点才可以投票）

  Returns:
  	- []dto.TrustAnchorVoter， 投票人信息
	- error

  Call permissions: Anyone
*/
func (anc *Anchor) GetVoter(voterAddress string) ([]*dto.TrustAnchorVoter, error) {
	params := make([]string, 1)
	params[0] = voterAddress

	pointer := &dto.SystemRequestResult{}

	err := anc.super.provider.SendRequest(pointer, "trustanchor_voter", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToTrustAnchorVoter()
}
