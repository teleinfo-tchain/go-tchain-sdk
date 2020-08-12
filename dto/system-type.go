package dto

import (
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"math/big"
)

/*
	DPoS子模块所需数据结构
*/
type Message struct {
	Code          uint64 `json:"code"`
	Msg           []byte `json:"message"`
	Address       string `json:"address"`
	Signature     []byte `json:"signature"`
	CommittedSeal []byte `json:"committedSeal"`
}

type View struct {
	Round    *big.Int `json:"round"`
	Sequence *big.Int `json:"sequence"`
}

type MessageSet struct {
	View     *View      `json:"view"`
	ValSet   []string   `json:"valSet"`
	Messages []*Message `json:"messages"`
}

// RoundStateInfoResponse is the information of RoundState
type RoundStateInfo struct {
	Commits    *MessageSet `json:"commits"`
	LockedHash utils.Hash  `json:"lockedHash"`
	Prepares   *MessageSet `json:"prepares"`
	Proposer   string      `json:"proposer"`
	Round      *big.Int    `json:"round"`
	Sequence   *big.Int    `json:"sequence"`
	View       *View       `json:"view"`
}

type RoundChangeSetInfo struct {
	RoundChanges map[uint64]*MessageSet `json:"roundChanges"`
	Validates    []string               `json:"validates"`
}

/*
	节点可信合约
*/

type RegisterCertificateInfo struct {
	Id          string // 节点证书的bid，必须和public_key对应，索引
	Apply       string // 申请人的bid（与Id可以相同，可以不同）
	PublicKey   string // 53个字符的公钥，也就是p2p节点id的形式
	NodeName    string // 节点名称，不含敏感词的字符串
	MessageSha3 string // 消息sha3后的16进制字符串，用于本地签名和链上验证签名，该字段不会被链保存
	Signature   string // 对上一个字段消息的签名，16进制字符串
	NodeType    uint64 // 节点类型，0企业，1个人
	Period      uint64 // 证书有效期，以年为单位的整型
	IP          string // 节点间互连的ip
	Port        uint64 // 节点间互连的端口
	CompanyName string // 公司名（如果是个人，则是个人姓名），不含敏感词的字符串
	CompanyCode string // 公司代码，不含敏感词的字符串
}

type PeerCertificate struct {
	Id          string   `json:"id"`          // 唯一索引
	Issuer      string   `json:"issuer"`      // 颁发者地址
	Apply       string   `json:"apply"`       // 申请人bid
	PublicKey   string   `json:"publicKey"`   // 节点公钥
	NodeName    string   `json:"nodeName"`    // 节点名称
	Signature   string   `json:"signature"`   // 节点签名内容
	NodeType    uint64   `json:"nodeType"`    // 节点类型0企业，1个人
	CompanyName string   `json:"companyName"` // 公司名称
	CompanyCode string   `json:"companyCode"` // 公司信用代码
	IssuedTime  *big.Int `json:"issuedTime"`  // 颁发时间
	Period      uint64   `json:"period"`      // 有效期
	IsEnable    bool     `json:"isEnable"`    // true 凭证有效，false 凭证已撤销
}

/*
	信任锚合约
*/

type RegisterAnchor struct {
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
}

type UpdateAnchorInfo struct {
	CompanyUrl  string
	Website     string
	DocumentUrl string
	ServerUrl   string
	Email       string
	Desc        string
}

type TrustAnchor struct {
	Id               string   `json:"id"              gencodec:"required"`   // 信任锚BID地址
	Name             string   `json:"name"            gencodec:"required"`   // 信任锚名称
	Company          string   `json:"company"         gencodec:"required"`   // 信任锚所属公司
	CompanyUrl       string   `json:"company_url"     gencodec:"required"`   // 公司网址
	Website          string   `json:"website"         gencodec:"required"`   // 信任锚网址
	ServerUrl        string   `json:"server_url"      gencodec:"required"`   // 服务链接
	DocumentUrl      string   `json:"document_url"    gencodec:"required"`   // 信任锚接口字段文档
	Email            string   `json:"email"           gencodec:"required"`   // 信任锚客服邮箱
	Desc             string   `json:"desc" gencodec:"required"`              // 描述
	TrustAnchorType  uint64   `json:"type"            gencodec:"required"`   // 信任锚类型
	Status           uint64   `json:"status"          gencodec:"required"`   // 服务状态
	Active           bool     `json:"active"          gencodec:"required"`   // 是否是根信任锚
	TotalBounty      *big.Int `json:"totalBounty"     gencodec:"required"`   // 总激励
	ExtractedBounty  *big.Int `json:"extractedBounty" gencodec:"required"`   // 已提取激励
	LastExtractTime  *big.Int `json:"lastExtractTime" gencodec:"required"`   // 上次提取时间
	VoteCount        *big.Int `json:"vote_count" gencodec:"required"`        // 得票数
	Stake            *big.Int `json:"stake" gencodec:"required"`             // 抵押
	CreateDate       *big.Int `json:"create_date" gencodec:"required"`       // 创建时间
	CertificateCount *big.Int `json:"certificate_count" gencodec:"required"` // 证书总数
}

type TrustAnchorVoter struct {
	Id             string `json:"owner"`          // 投票人地址
	VoteCandidates string `json:"voteCandidates"` // 投的人
	Votes          uint64 `json:"Votes"`          // 得票人的票数
}

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
	Id             string   // 凭证的hash
	Context        string   // 证书所属上下文环境
	Issuer         string   // 信任锚的bid
	Subject        string   // 证书拥有者地址
	IssuedTime     *big.Int // 颁发时间
	Period         uint64   // 有效期
	IsEnable       bool     // true 凭证有效，false 凭证已撤销
	RevocationTime *big.Int // 吊销时间
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

/*
	did文档合约
*/
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
	ServiceId string `json:"id"`
	Type      string `json:"type"`
	Endpoint  string `json:"serviceEndpoint"`
}

type Proof struct {
	// Id         string   `json:"id"`
	Type       string `json:"type"`
	CreateTime string `json:"created"`
	Creator    string `json:"creator"`
	Signature  string `json:"signatureValue"`
}

/*
	dpos投票合约
*/

type RegisterWitness struct {
	NodeUrl string
	Website string
	Name    string
}

type Candidate struct {
	Owner           string       `json:"owner"`           // 候选人地址
	Name            string       `json:"name"`            // 候选人名称
	Active          bool         `json:"active"`          // 当前是否是候选人
	Url             string       `json:"url"`             // 节点的URL
	VoteCount       *hexutil.Big `json:"voteCount"`       // 收到的票数
	TotalBounty     *hexutil.Big `json:"totalBounty"`     // 总奖励金额
	ExtractedBounty *hexutil.Big `json:"extractedBounty"` // 已提取奖励金额
	LastExtractTime *hexutil.Big `json:"lastExtractTime"` // 上次提权时间
	Website         string       `json:"website"`         // 见证人网站
}
