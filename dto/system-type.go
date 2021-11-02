package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/rlp"
	"io"
	"math/big"
)

// Message DPoS子模块所需数据结构
type Message struct {
	Code          uint64 `json:"code"`
	Msg           []byte `json:"message"`
	Address       string `json:"address"`
	Signature     []byte `json:"signature"`
	CommittedSeal []byte `json:"committedSeal"`
}

type View struct {
	Round    uint64 `json:"round"`
	Sequence uint64 `json:"sequence"`
}

type MessageSet struct {
	View     *View      `json:"view"`
	ValSet   []string   `json:"valSet"`
	Messages []*Message `json:"messages"`
}

type PrePrepare struct {
	View     *View    `json:"view"` // 序列和round
	Proposal Proposal `json:"proposal"`
}

// Proposal supports retrieving height and serialized block to be used during Istanbul consensus.
type Proposal interface {
	// Number retrieves the sequence number of this proposal.
	Number() *big.Int

	// Hash retrieves the hash of this proposal.
	Hash() utils.Hash

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error

	String() string
}

// RoundStateInfo is the information of RoundState
type RoundStateInfo struct {
	Commits     *MessageSet `json:"commits"`
	LockedHash  string      `json:"lockedHash"`
	Prepares    *MessageSet `json:"prepares"`
	Proposer    string      `json:"proposer"`
	Round       *big.Int    `json:"round"`
	Sequence    *big.Int    `json:"sequence"`
	PrePrepares *PrePrepare `json:"preprepares"`
}

func (roundStateInfo *RoundStateInfo) UnmarshalJSON(data []byte) error {
	type Alias RoundStateInfo

	temp := &struct {
		Round    string `json:"round"`
		Sequence string `json:"sequence"`
		*Alias
	}{
		Alias: (*Alias)(roundStateInfo),
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	round, success := big.NewInt(0).SetString(temp.Round[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Round))
	}

	sequence, success := big.NewInt(0).SetString(temp.Sequence[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Sequence))
	}

	roundStateInfo.Round = round
	roundStateInfo.Sequence = sequence

	return nil
}

type RoundChangeSetInfo struct {
	RoundChanges map[uint64]*MessageSet `json:"roundChanges"`
	Validates    []string               `json:"validates"`
}

// AllianceInfo 联盟合约
type AllianceInfo struct {
	Id          string `json:"id"`          // 联盟成员bid
	PublicKey   string `json:"publicKey"`   // 联盟成员公钥
	CompanyName string `json:"companyName"` // 公司名称
	CompanyCode string `json:"companyCode"` // 公司信用代码
}

type Alliance struct {
	Id           string `json:"id"`           // 联盟成员bid
	Role         uint64 `json:"role"`         // 角色 1理事，2副理事长，3理事长
	PublicKey    string `json:"publicKey"`    // 联盟成员公钥
	CompanyName  string `json:"companyName"`  // 公司名称
	CompanyCode  string `json:"companyCode"`  // 公司信用代码
	Auditor      string `json:"auditor"`      // 审核员
	AuditTime    uint64 `json:"auditTime"`    // 审核时间
	RevokeReason string `json:"revokeReason"` // 撤销理由
	Active       bool   `json:"active"`       // 是否有效，是否撤销
}

func (alliance *Alliance) UnmarshalJSON(data []byte) error {
	type Alli Alliance

	temp := &struct {
		*Alli
	}{
		Alli: (*Alli)(alliance),
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	return nil
}

type Weights struct {
	DirectorWeights        uint64 `json:"directorWeights"`        // 理事权重
	ViceWeights            uint64 `json:"viceWeights"`            // 副理事长权重
	DirectorGeneralWeights uint64 `json:"directorGeneralWeights"` // 理事长权重
	Score                  uint64 `json:"score"`                  // 总票数
}

// RegisterCertificate 可信认证
type RegisterCertificate struct {
	Id               string // 个人可信证书bid
	Context          string // 证书上下文环境，随便一个字符串，不验证
	Subject          string // 证书接收者的bid，证书是颁给谁的
	Period           uint64 // 证书有效期，以年为单位的整型
	IssuerAlgorithm  string // 颁发者签名算法，字符串
	IssuerSignature  string // 颁发者签名值，16进制字符串
	SubjectPublicKey string // 接收者公钥，16进制字符串
	SubjectAlgorithm string // 接收者签名算法，字符串
	SubjectSignature string // 接收者签名值，16进制字符串
}

type CertificateInfo struct {
	Id             string // 凭证的hash
	Context        string // 证书所属上下文环境
	Issuer         string // 信任锚的bid
	Subject        string // 证书拥有者地址
	IssuedTime     uint64 // 颁发时间
	Period         uint64 // 有效期
	IsEnable       bool   // true 凭证有效，false 凭证已撤销
	RevocationTime uint64 // 吊销时间
}

func (certificateInfo *CertificateInfo) UnmarshalJSON(data []byte) error {
	type Alias CertificateInfo

	temp := &struct {
		// TransactionIndex string `json:"transactionIndex"`
		*Alias
	}{
		Alias: (*Alias)(certificateInfo),
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	return nil
}

type IssuerSignature struct {
	Id        string // 凭证ID
	PublicKey string // 签名公钥
	Algorithm string // 签名算法
	Signature string // 签名内容
}

type SubjectSignature struct {
	Id        string // 凭证ID
	PublicKey string // 签名公钥
	Algorithm string // 签名算法
	Signature string // 签名内容
}

// PublicKey did文档合约
type PublicKey struct {
	Id string `json:"id"`
	// KeyId      string `json:"key_id"`
	Type       string `json:"type"`
	Controller string `json:"controller"`
	Authority  string `json:"authority"` // 公钥权限
	PublicKey  string `json:"publicKey"`
}

type Authentication struct {
	Id        utils.Address `json:"id"`
	ProofId   []byte        `json:"proofId"`
	Issuer    utils.Address `json:"type"`
	PublicKey []byte        `json:"public_key"`
}

type Attribute struct {
	Id       utils.Address `json:"id"`
	AttrType []byte        `json:"attr_type"`
	Value    []byte        `json:"value"`
}

type Document struct {
	Id              string       `json:"id"` // bid
	Contexts        string       `json:"context"`
	Name            string       `json:"name"`           // bid标识符昵称
	Type            uint64       `json:"type"`           // bid的类型，包括0：普通用户,1:智能合约以及设备，2：企业或者组织，BID类型一经设置，永不能变
	PublicKeys      []*PublicKey `json:"publicKey"`      // 用户用于身份认证的公钥信息
	Authentications []string     `json:"authentication"` // 用户身份认证列表信息
	Services        []*Service   `json:"service"`        // 用户填写的服务端点信息
	Proof           *Proof       `json:"proof"`          // 用户填写的证明信息值
	Extra           string       `json:"extra"`          // 用户填写的备注
	IsEnable        bool         `json:"isEnable"`       // 该BID是否启用
	CreateTime      string       `json:"created"`
	UpdateTime      string       `json:"updated"`
}

type Service struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Endpoint string `json:"serviceEndpoint"`
}

type Proof struct {
	Type       string `json:"type"`
	CreateTime string `json:"created"`
	Creator    string `json:"creator"`
	Signature  string `json:"signatureValue"`
}

type PeerNodeInfo struct {
	Id          string `json:"id"`          // 唯一索引，节点bid
	Apply       string `json:"apply"`       // m申请人地址（联盟成员的地址）
	PublicKey   string `json:"publicKey"`   // 节点公钥
	NodeName    string `json:"nodeName"`    // 节点名称
	MessageSha3 string `json:"messageSha3"` // 要签名的内容sha3的hash
	Signature   string `json:"signature"`   // 节点签名内容
	Url         string `json:"url"`         // 节点URL
	Website     string `json:"website"`     // 节点网站地址
	NodeType    uint64 `json:"nodeType"`    // 节点类型0企业，1个人
	CompanyName string `json:"companyName"` // 公司名称
	CompanyCode string `json:"companyCode"` // 公司信用代码
	Ip          string `json:"ip"`          // ip
	Port        uint64 `json:"port"`        // 端口
}

type PeerNodeDetail struct {
	Id              string   `json:"id"`              // 唯一索引，节点bid
	Issuer          string   `json:"issuer"`          // 颁发者地址(理事长地址)
	Apply           string   `json:"apply"`           // 申请人地址（联盟成员的地址）
	PublicKey       string   `json:"publicKey"`       // 节点公钥
	NodeName        string   `json:"nodeName"`        // 节点名称
	Signature       string   `json:"signature"`       // 节点签名内容
	Url             string   `json:"url"`             // 节点URL
	Website         string   `json:"website"`         // 节点网站地址
	NodeType        uint64   `json:"nodeType"`        // 节点类型0企业，1个人
	CompanyName     string   `json:"companyName"`     // 公司名称
	CompanyCode     string   `json:"companyCode"`     // 公司信用代码
	Role            uint64   `json:"role"`            // 可信节点、候选节点、共识节点
	Active          bool     `json:"active"`          // 是否有效
	CreateTime      uint64   `json:"create"`          // 创建时间
	RevokeReason    string   `json:"revokeReason"`    // 节点删除理由
	ConsensusRevoke string   `json:"consensusRevoke"` // 共识撤销理由
	StartTime       uint64   `json:"startTime"`       // 投票开始时间
	Score           uint64   `json:"score"`           // 投票得分
	VoterList       []string `json:"voterList"`       // 投票给这个节点的投票人列表
}

type AllContract struct {
	ContractName string `json:"contractName"` // 用户地址
	IsEnable     bool   `json:"isEnable"`     // 用户权限,1启用合约，2禁用合约，4授权  // 3=1+2, 5=1+4, 6=2+4, 7=1+2+4 类linux权限管理
}

type SubChainInfo struct {
	Id             string `json:"id"`             // 子链的bid
	Apply          string `json:"apply"`          // 子链申请者的bid
	SubChainName   string `json:"subChainName"`   // 子链名
	ChainCode      string `json:"chainCode"`      // 子链AC码
	ChainIndustry  string `json:"chainIndustry"`  // 子链所属行业
	ChainFramework string `json:"chainFramework"` // 子链架构
	Consensus      string `json:"consensus"`      // 子链共识
	ChainMsgHash   string `json:"chainMsgHash"`   // 子链其他信息存储哈希
}

type SubChainDetail struct {
	Id             string   `json:"id"`             // 子链的bid
	Apply          string   `json:"apply"`          // 子链申请者的bid
	Auditor        string   `json:"auditor"`        // 子链审核的理事长bid
	SubChainName   string   `json:"subChainName"`   // 子链名
	ChainCode      string   `json:"chainCode"`      // 子链AC码
	ChainIndustry  string   `json:"chainIndustry"`  // 子链所属行业
	ChainFramework string   `json:"chainFramework"` // 子链架构
	Consensus      string   `json:"consensus"`      // 子链共识
	ChainMsgHash   string   `json:"chainMsgHash"`   // 子链其他信息存储哈希
	Status         uint64   `json:"status"`         // 该子链的状态, 1为申请投票中, 2为成功 3 失效（投票时间过期/撤销）
	RevokeReason   string   `json:"revokeReason"`   // 撤销原由
	StartTime      uint64   `json:"startTime"`      // 投票开始时间
	Score          uint64   `json:"score"`          // 投票得分
	VoterList      []string `json:"voterList"`      // 投票给这个节点的投票人列表
}