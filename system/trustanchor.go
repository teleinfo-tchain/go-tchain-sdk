package system

import (
	"encoding/json"
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/common/hexutil"
	"github.com/bif/bif-sdk-go/complex/types"
	"regexp"
	"strings"
)

const (
	TrustAnchorContractAddr = "did:bid:00000000000000000000000c"
	AnchorStatusUnknow      = 0
	AnchorStatusOK          = 1
	AnchorStatusErr         = 2
	AnchorStatusDelete      = 3

	BaseAnchor                 = 10
	ExtendAnchor               = 11
	UnknowAnchorType           = 2
	Day                        = 86400        //24*60*60秒，每次提取积分的最小时间间隔
	Bifer                      = 100000000    //10^9
	MiniExtractAmount          = 100 * Bifer  //每次提取积分的最小额度是100
	BaseTrustAnchorPledge      = 1000 * Bifer //注册根信任锚抵押的积分
	ExtendTrustAnchorPledge    = 100 * Bifer  //注册扩展信任锚抵押的积分
	IncentivesToExtendIssueCer = 1 * Bifer    //颁发一个证书获得的奖励
	IncentivesToBaseIssueCer   = 2 * Bifer    //颁发一个证书获得的奖励
	EmptyUrl                   = ""
)

var (
	ErrIllExtracAmout    = errors.New("积分总额不足100，无法提取")
	ErrIllExtracTime     = errors.New("距离上次提取不足24小时，无法提取")
	ErrIllAnchorType     = errors.New("未知的信任锚类型，10 代表根信任锚，11代表扩展信任锚")
	ErrAnchorExist       = errors.New("该地址已存在于信任锚列表中")
	ErrAnchorNoExist     = errors.New("信任锚不存在")
	ErrIllegalAnchor     = errors.New("信任锚字段不能为空")
	ErrIllegalBalance    = errors.New("账户内积分不足，注册根信任锚需要抵押1000积分，注册扩展信任锚需要抵押100积分")
	ErrIllegalVote       = errors.New("信任锚不存在或不是基础信任锚")
	ErrIllegalRepeatVote = errors.New("同一个基础信任锚能且仅能投一票")
)

const TrustAnchorAbiJSON = `[
{"constant": false,"name":"registerTrustAnchor","inputs":[{"name":"anchor","type":"string"},{"name":"anchortype","type":"uint64"},{"name":"anchorname","type":"string"},{"name":"company","type":"string"},{"name":"companyurl","type":"string"},{"name":"website","type":"string"},{"name":"documenturl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"unRegisterTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[],"type":"function"},
{"constant": true,"name":"isTrustAnchor","inputs":[{"name":"address","type":"string"}],"outputs":[{"name":"trustanchor","type":"bool"}],"type":"function"},
{"constant": false,"name":"updateBaseAnchorInfo","inputs":[{"name":"anchor","type":"string"},{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"updateExtendAnchorInfo","inputs":[{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"extractOwnBounty","inputs":[{"name":"anchor","type":"string"}],"outputs":[],"type":"function"},
{"constant": true,"name":"queryTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[{"name":"id","type":"string"},{"name":"name","type":"string"},{"name":"company","type":"string"},{"name":"CompanyUrl","type":"string"},{"name":"website","type":"string"},{"name":"ServerUrl","type":"string"},{"name":"DocumentUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"},{"name":"TrustAnchorType","type":"uint64"},{"name":"status","type":"uint64"},{"name":"active","type":"bool"},{"name":"totalbounty","type":"uint64"},{"name":"extractedBounty","type":"uint64"},{"name":"lastExtracTime","type":"uint64"},{"name":"votecount","type":"uint64"},{"name":"stake","type":"uint64"},{"name":"createDate","type":"uint64"},{"name":"certificateAcount","type":"uint64"}],"type":"function"},
{"constant": true,"name":"queryTrustAnchorStatus","inputs":[{"name":"anchor","type":"string"}],"outputs":[{"name":"anchorstatus","type":"uint64"}],"type":"function"},
{"constant": true,"name":"queryBaseTrustAnchorList","inputs":[],"outputs":[{"name":"baseList","type":"string"}],"type":"function"},
{"constant": true,"name":"queryBaseTrustAnchorNum","inputs":[],"outputs":[{"name":"baseListNum","type":"uint64"}],"type":"function"},
{"constant": true,"name":"queryExpendTrustAnchorList","inputs":[],"outputs":[{"name":"expendList","type":"string"}],"type":"function"},
{"constant": true,"name":"queryExpendTrustAnchorNum","inputs":[],"outputs":[{"name":"expendListNum","type":"uint64"}],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"methodName","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"},{"indexed":false,"name":"time","type":"uint256"}],"name":"trustAnchorEvent","type":"event"},
{"constant": false,"name":"voteElect","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelVote","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"queryVoter","inputs":[{"name":"voterAddress","type":"string"}],"outputs":[{"name":"voterInfo","type":"string"}],"type":"function"},
{"constant": true,"name":"checkSenderAddress","inputs":[{"name":"address","type":"string"}],"outputs":[{"name":"supernode","type":"bool"}],"type":"function"}
]`

type Anchor struct {
	super *System
	abi   abi.ABI
}

type TrustAnchor struct {
	Id                string
	Name              string
	Company           string
	CompanyUrl        string
	Website           string
	ServerUrl         string
	DocumentUrl       string
	Email             string
	Desc              string
	TrustAnchorType   uint64
	Status            uint64
	Active            bool
	Totalbounty       uint64
	ExtractedBounty   uint64
	LastExtracTime    uint64
	Votecount         uint64
	Stake             uint64
	CreateDate        uint64
	CertificateAcount uint64
}

type RegisterAnchor struct {
	Anchor string
	Anchortype uint64
	Anchorname string
	Company string
	ExtendAnchorInfo
}

type BaseAnchorInfo struct {
	Anchor string
	ExtendAnchorInfo
}

type ExtendAnchorInfo struct{
	Companyurl string
	Website string
	Documenturl string
	ServerUrl string
	Email string
	Desc string
}

type VoterInfo struct {
	Owner string  `json:"owner"`
	VoteCandidates []Candidates `json:"voteCandidates"`
}

type Candidates struct {
	Address string  `json:"addr"`
	Vote string  `json:"vote"`
}

func (system *System) NewTrustAnchor() (*Anchor, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(TrustAnchorAbiJSON))
	if err != nil {
		return nil, err
	}

	Anchor := new(Anchor)
	Anchor.abi = parsedAbi
	Anchor.super = system
	return Anchor, nil
}

// 是否让用户可以自定义gas和gasPrice？
//"registerTrustAnchor","inputs":[{"name":"anchor","type":"string"},{"name":"anchortype","type":"uint64"},{"name":"anchorname","type":"string"},{"name":"company","type":"string"},{"name":"companyurl","type":"string"},{"name":"website","type":"string"},{"name":"documenturl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[]
func (anchorTrust *Anchor) RegisterTrustAnchor(from common.Address, registerAnchor *RegisterAnchor) (string, error) {
	// encoding
	// RegisterAnchor is a struct we need to use the components.
	var values []interface{}
	values = anchorTrust.super.StructToInterface(*registerAnchor,values)
	inputEncode, err := anchorTrust.abi.Pack("registerTrustAnchor", values...)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	transactionHash, err := anchorTrust.super.SendTransaction(transaction)
	if err != nil {
		return "", err
	}
	return transactionHash, err
}

//"unRegisterTrustAnchor","inputs":[{"name":"anchor","type":"string"}],"outputs":[]
func (anchorTrust *Anchor) UnRegisterTrustAnchor(from common.Address, anchor string) (string, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("unRegisterTrustAnchor", anchor)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	transactionHash, err := anchorTrust.super.SendTransaction(transaction)
	if err != nil {
		return "", err
	}
	return transactionHash, err
}

//"isTrustAnchor","inputs":[{"name":"address","type":"string"}],"outputs":[{"name":"trustanchor","type":"bool"}]
// 基础信任锚和扩展信任锚的注册后，初始化时，为false;
func (anchorTrust *Anchor) IsTrustAnchor(from common.Address, address string) (bool, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("isTrustAnchor", address)
	if err != nil {
		panic(err)
	}
	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return false, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))

	var trustanchor bool
	err = anchorTrust.abi.Methods["isTrustAnchor"].Outputs.Unpack(&trustanchor, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return false, err
	}
	return trustanchor, nil
}

//"updateBaseAnchorInfo","inputs":[{"name":"anchor","type":"string"},{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[]
func (anchorTrust *Anchor) UpdateBaseAnchorInfo(from common.Address, updateAnchor *BaseAnchorInfo) (string, error) {
	// encoding
	// updateAnchor is a struct we need to use the components.
	var values []interface{}
	values = anchorTrust.super.StructToInterface(*updateAnchor, values)
	inputEncode, err := anchorTrust.abi.Pack("updateBaseAnchorInfo", values...)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	transactionHash, err := anchorTrust.super.SendTransaction(transaction)
	if err != nil {
		return "", err
	}
	return transactionHash, err
}

//"inputs":[{"name":"companyUrl","type":"string"},{"name":"website","type":"string"},{"name":"documentUrl","type":"string"},{"name":"serverUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"}],"outputs":[]
func (anchorTrust *Anchor) UpdateExtendAnchorInfo(from common.Address, extendAnchorInfo *ExtendAnchorInfo) (string, error) {
	// encoding
	// extendAnchorInfo is a struct we need to use the components.
	var values []interface{}
	values = anchorTrust.super.StructToInterface(*extendAnchorInfo, values)
	inputEncode, err := anchorTrust.abi.Pack("updateExtendAnchorInfo", values...)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	transactionHash, err := anchorTrust.super.SendTransaction(transaction)
	if err != nil {
		return "", err
	}
	return transactionHash, err
}

//"extractOwnBounty","inputs":[{"name":"anchor","type":"string"}],"outputs":[]
func (anchorTrust *Anchor) ExtractOwnBounty(from common.Address, anchor string) (string, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("extractOwnBounty", anchor)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	transactionHash, err := anchorTrust.super.SendTransaction(transaction)
	if err != nil {
		return "", err
	}
	return transactionHash, err
}


//"inputs":[{"name":"anchor","type":"string"}],"outputs":[{"name":"id","type":"string"},{"name":"name","type":"string"},{"name":"company","type":"string"},{"name":"CompanyUrl","type":"string"},{"name":"website","type":"string"},{"name":"ServerUrl","type":"string"},{"name":"DocumentUrl","type":"string"},{"name":"email","type":"string"},{"name":"desc","type":"string"},{"name":"TrustAnchorType","type":"uint64"},{"name":"status","type":"uint64"},{"name":"active","type":"bool"},{"name":"totalbounty","type":"uint64"},{"name":"extractedBounty","type":"uint64"},{"name":"lastExtracTime","type":"uint64"},{"name":"votecount","type":"uint64"},{"name":"stake","type":"uint64"},{"name":"createDate","type":"uint64"},{"name":"certificateAcount","type":"uint64"}]
func (anchorTrust *Anchor) GetTrustAnchor(from common.Address, anchor string) (*TrustAnchor, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("queryTrustAnchor", anchor)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return nil, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var trustAnchor TrustAnchor
	err = anchorTrust.abi.Methods["queryTrustAnchor"].Outputs.Unpack(&trustAnchor, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}
	if trustAnchor.Id == "0000000000000000000000000000000000000000"{
		return nil, ErrCertificateNotExist
	}
	return &trustAnchor, err
}

//"name":"queryTrustAnchorStatus","inputs":[{"name":"anchor","type":"string"}],"outputs":[{"name":"anchorstatus","type":"uint64"}]
func (anchorTrust *Anchor) GetTrustAnchorStatus(from common.Address, anchor string) (uint64, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("queryTrustAnchorStatus", anchor)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return 0, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var trustAnchorStatus uint64
	err = anchorTrust.abi.Methods["queryTrustAnchorStatus"].Outputs.Unpack(&trustAnchorStatus, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return 0, err
	}

	return trustAnchorStatus, err
}

//"name":"queryBaseTrustAnchorList","inputs":[],"outputs":[{"name":"baseList","type":"string"}]
func (anchorTrust *Anchor) GetBaseTrustAnchorList(from common.Address) ([]string, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("queryBaseTrustAnchorList")
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return nil, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var baseList string
	err = anchorTrust.abi.Methods["queryBaseTrustAnchorList"].Outputs.Unpack(&baseList, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}
	anchorCount := len(baseList)/32
	if anchorCount == 0{
		return nil, nil
	}
	anchorList := make([]string,anchorCount)
	for i:=0;i<anchorCount;i++{
		anchorList = append(anchorList, baseList[i*32:(i+1)*32-1])
	}
	return anchorList, err
}

//"queryBaseTrustAnchorNum","inputs":[],"outputs":[{"name":"baseListNum","type":"uint64"}]
func (anchorTrust *Anchor) GetBaseTrustAnchorNum(from common.Address) (uint64, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("queryBaseTrustAnchorNum")
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return 0, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var baseListNum uint64
	err = anchorTrust.abi.Methods["queryBaseTrustAnchorNum"].Outputs.Unpack(&baseListNum, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return 0, err
	}

	return baseListNum, err
}

//"name":"queryExpendTrustAnchorList","inputs":[],"outputs":[{"name":"expendList","type":"string"}]
func (anchorTrust *Anchor) GetExpendTrustAnchorList(from common.Address) ([]string, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("queryExpendTrustAnchorList")
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return nil, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var expendList string
	err = anchorTrust.abi.Methods["queryExpendTrustAnchorList"].Outputs.Unpack(&expendList, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}
	anchorCount := len(expendList)/32
	if anchorCount == 0{
		return nil, nil
	}
	expendAnchorList := make([]string,anchorCount)
	for i:=0;i<anchorCount;i++{
		expendAnchorList = append(expendAnchorList, expendList[i*32:(i+1)*32-1])
	}
	return expendAnchorList, err
}

//"queryExpendTrustAnchorNum","inputs":[],"outputs":[{"name":"expendListNum","type":"uint64"}]
func (anchorTrust *Anchor) GetExpendTrustAnchorNum(from common.Address) (uint64, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("queryExpendTrustAnchorNum")
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return 0, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var expendListNum uint64
	err = anchorTrust.abi.Methods["queryExpendTrustAnchorNum"].Outputs.Unpack(&expendListNum, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return 0, err
	}

	return expendListNum, err
}

//"voteElect","inputs":[{"name":"candidate","type":"string"}],"outputs":[]
func (anchorTrust *Anchor) VoteElect(from common.Address, candidate string) (string, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("voteElect", candidate)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	transactionHash, err := anchorTrust.super.SendTransaction(transaction)
	if err != nil {
		return "", err
	}
	return transactionHash, err
}


//"cancelVote","inputs":[{"name":"candidate","type":"string"}],"outputs":[]
func (anchorTrust *Anchor) CancelVote(from common.Address, candidate string) (string, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("cancelVote", candidate)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	transactionHash, err := anchorTrust.super.SendTransaction(transaction)
	if err != nil {
		return "", err
	}
	return transactionHash, err
}

//"queryVoter","inputs":[{"name":"voterAddress","type":"string"}],"outputs":[{"name":"voterInfo","type":"string"}]
func (anchorTrust *Anchor) GetVoter(from common.Address, voterAddress string) (*VoterInfo, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("queryVoter", voterAddress)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return nil, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var voterInfo string
	err = anchorTrust.abi.Methods["queryVoter"].Outputs.Unpack(&voterInfo, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return nil, err
	}
	re1, _ := regexp.Compile(`"`)
	re2, _ := regexp.Compile(`'`)
	rep := re2.ReplaceAllString(re1.ReplaceAllString(voterInfo, ``), `"`)
	voteInfo := &VoterInfo{}
	err = json.Unmarshal([]byte(rep), voteInfo)
	if err != nil {
		return nil, err
	}
	return voteInfo, err
}


//"checkSenderAddress","inputs":[{"name":"address","type":"string"}],"outputs":[{"name":"supernode","type":"bool"}]
func (anchorTrust *Anchor) CheckSenderAddress(from common.Address, address string) (bool, error) {
	// encoding
	inputEncode, err := anchorTrust.abi.Pack("checkSenderAddress", address)
	if err != nil {
		panic(err)
	}

	transaction := anchorTrust.super.PrePareTransaction(from, TrustAnchorContractAddr, types.ComplexString(hexutil.Encode(inputEncode)))
	requestResult, err := anchorTrust.super.Call(transaction)
	if err != nil {
		return false, err
	}
	//fmt.Println("result is ", requestResult.Result.(string))
	var superNode bool
	err = anchorTrust.abi.Methods["checkSenderAddress"].Outputs.Unpack(&superNode, common.FromHex(requestResult.Result.(string)))
	// 解码不应该出错，除非底层逻辑变更
	if err != nil {
		return false, err
	}

	return superNode, err
}
