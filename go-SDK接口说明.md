** bif-go-SDK接口说明**

**简介**：bif-go-SDK系统合约说明

**HOST**: 

**联系人**:

**Version**:1.0

**接口路径**：


# **bif-go-SDK**

## 0.System
    主要负责处理系统合约交互
```
    1 .敏感词合约交互    - 合约保存了一系列的敏感词，当其他合约有设置名字的需求时，调用该合约查询是否包含了相应敏感词
    2. 超级管理合约交互  - 管理其它的合约，可以禁用BIF链上的任何合约
    3. 节点可信合约交互  - peerCertificate.go，监管节点向节点颁发证书，节点由普通节点变成可信节点
    4. 信任锚合约交互    - trustAnchor.go，用于给个人颁发可信证书的机构被称为信任锚。信任锚又分为根信任锚和扩展信任锚，根信任锚只能通过超级节点注册，且超级节点投支持票超过2/3才可主要是由公安，学历等国家承认的机构才会注册。扩展信任锚，只要注册了就可以颁发证书，任何企业、组织都可以注册，然后为自己的员工颁发证书。
    5. 个人可信合约交互  - certificate.go，信任锚颁发给个人的证书
    6. did文档合约交互   - document.go，Did文档的操作
    7. DPoS投票合约交互  - dpos.go，选举投票人

```
### 1) 系统合约的单元测试所在路径
     ../src/bif-sdk-go/test/System

单元测试文件               | 功能说明
----------------------- | -------------
peerCertificate_test.go | 节点可信单元测试 
document_test.go        | did文档单元测试
election_test.go        | dpos投票单元测试
trustAnchor_test.go     | 信任锚单元测试
certificate_test.go     | 个人可信单元测试
dpos_test.go            | dpos属性单元测试

### 2) *SysTxParams 结构体说明
    凡是涉及到修改链的状态的交互，都会涉及到交易的构造，流程如下：
    1 (解密keyStore) -->2(本地构造交易并签署) -->3(发送链上) -->4(获取交易hash) -->5(根据hash判断执行结果)
    
    type SysTxParams struct {
    	IsSM2       bool     // 私钥生成是否使用国密，true为国密；false为非国密
    	Password    string   // 解密私钥的密码
    	KeyFileData []byte   // keystore文件内容
    	GasPrice    *big.Int // 交易的gas价格，默认是网络gas价格的平均值
    	Gas         uint64   // 交易可使用的gas，未使用的gas会退回
    	Nonce       uint64   // 从该账户发起交易的Nonce值
    	ChainId     *big.Int // 链的ChainId
    }
    
    结构体构造（传参）说明
    
    | 参数名称      | 类型      | 是否必填描述   | 描述                                               |
    | ------------| ---------| -------------| ---------------------------------------------------|
    | IsSM2       | bool     | 是           | 私钥生成是否使用国密，true为国密；false为非国密；用于解密私钥 |
    | Password    | string   | 是           | 解密私钥的密码（keyStore文件生成时的加密密码）|
    | KeyFileData | []byte   | 是           | keyStore的内容 |
    | GasPrice    | *big.Int | 否           | （选填）交易的gas价格，默认是网络gas价格的平均值   |
    | Gas         | uint64   | 是           | 交易可使用的gas，未使用的gas会退回 |
    | Nonce       | uint64   | 否           | （选填）从该账户发起交易的Nonce值   |
    | ChainId     | *big.Int | 否           | （选填）链的ChainId |
 
### 3) 请求示例
  
#### 1) 单元测试示例（send方法）   

##### A 发起交易
    
    #### 代码go（个人可信证书合约——注册）
    ```
    func TestRegisterCertificate(t *testing.T) {
        // 初始化bif
    	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
    	// 获取chainId
    	chainId, err := connection.Core.GetChainId()
    	if err != nil {
    		t.Log(err)
    		t.FailNow()
    	}
        // 获取nonce
    	nonce, err := connection.Core.GetTransactionCount(testAddress, block.LATEST)
    	if err != nil {
    		t.Log(err)
    		t.FailNow()
    	}
        
        // sysTxParams构造
    	sysTxParams := new(system.SysTxParams)
    	// 私钥生成是否使用国密，true为国密；false为非国密；用于解密私钥 
    	sysTxParams.IsSM2 = isSM2
    	// 解密私钥的密码（keyStore文件生成时的加密密码）
    	sysTxParams.Password = password
    	// keyStore的内容
    	sysTxParams.KeyFileData = keyFileData
    	// （选填）交易的gas价格，默认是网络gas价格的平均值，可以不写
    	sysTxParams.GasPrice = big.NewInt(35)
    	// 交易可使用的gas，未使用的gas会退回
    	sysTxParams.Gas = 2000000
    	// （选填）从该账户发起交易的Nonce值 ，可以不写
    	sysTxParams.Nonce = nonce.Uint64()
    	// （选填）链的ChainId，可以不写
    	sysTxParams.ChainId = big.NewInt(0).SetUint64(chainId)
    
    	cer := connection.System.NewCertificate()
        
        // 证书注册由于参数很多，构造一个结构体，参数内容不太清楚
    	registerCertificate := new(dto.RegisterCertificate)
    	registerCertificate.Id = utils.StringToAddress(testAddress).String()
    	registerCertificate.Context = "context_test"
    	registerCertificate.Subject = "did:bid:6cc796b8d6e2fbebc9b3cf9e"
    	registerCertificate.Period = 3
    	registerCertificate.IssuerAlgorithm = ""
    	registerCertificate.IssuerSignature = ""
    	registerCertificate.SubjectPublicKey = ""
    	registerCertificate.SubjectAlgorithm = ""
    	registerCertificate.SubjectSignature = ""
    	// registerCertificate
    	registerCertificateHash, err := cer.RegisterCertificate(sysTxParams, registerCertificate)
    	if err != nil {
    		t.Error(err)
    		t.FailNow()
    	}
    	// 获取交易的hash
    	t.Log(registerCertificateHash, err)
    }
    ```
    
##### B 交易结果查验
    #### 代码go（个人可信证书合约——注册）
    ```
    // 测试系统合约的执行结果
    func TestSystemLogDecode(t *testing.T) {
    	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
    	_, err := connection.Core.GetCoinBase()
    	if err != nil {
    		t.Error(err)
    		t.FailNow()
    	}
        
        txHash := "0x2105effdc5d303391a3f194f0b29d1d1d2755b0cc239c03e3dd88581830116d0"
    	log, err := connection.System.SystemLogDecode()
    
    	if err != nil {
    		t.Errorf("err log : %v ", err)
    		t.FailNow()
    	}
    
    	t.Log("method is ", log.Method)
    	t.Log("status is ", log.Status)
    	t.Log("result is ", log.Result)
    }
    ```
   **LogData结构体**
   ```
   type LogData struct {
       Method string  // 交易调用的方法
       Status bool    // 交易执行成功还是失败
       Result string  // 交易执行的结果（如果失败返回错误原因；）
    }

   connection.System.SystemLogDecode() 返回LogData结构体指针
   通过其内部的三个参数获取交易的执行情况
   ```
   
#### 2) 单元测试示例（rpc方法） 
    
     #### 代码go（个人可信证书合约——查看有效期）
        ```
        func TestGetPeriod(t *testing.T) {
            // 初始化bif
            var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
            // 获取coinBase
            coinBase, err := connection.Core.GetCoinBase()
            // 错误判断
            if err != nil {
                t.Error(err)
                t.FailNow()
            }
            
            // 初始化 个人可信证书
            cer := connection.System.NewCertificate()
            
            // 查询个人证书的有效期
            period, err := cer.GetPeriod(coinBase)
            if err != nil {
                t.Error(err)
                t.FailNow()
            }
            // 打印结果
            t.Log(period)
        }
        ```

#### 3) 单元测试组测试示例（用于相同类型的多个输入测试）

     #### 代码go（个人可信证书合约——查看有效期）
        ```
        func TestHashMessage(t *testing.T) {
            // 通过结构体构建多个相同用例
        	for _, test := range []struct {
        		message     string
        		isSm2       bool
        		hashMessage string
        	}{
        		{"Hello World", false, "0xa1de988600a42c4b4ab089b619297c17d53cffae5d5120d82d8a92d0bb3b78f2"},
        		{utils.NewUtils().Utf8ToHex("Hello World"), false, "0xa1de988600a42c4b4ab089b619297c17d53cffae5d5120d82d8a92d0bb3b78f2"},
        	} {
        		var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
        		res := connection.Account.HashMessage(test.message, test.isSm2)
        		if res != test.hashMessage {
        			t.Errorf("hash error, input is %s, result is %s, expect is %s \n", test.message, res, test.hashMessage)
        			t.FailNow()
        		}
        
        	}
        }
        ```
        
### 4) 其它可能用到的接口

##### 1） 如果是core的接口，参考core模块，其单元测试在 ../src/bif-sdk-go/test/core

##### 2） 涉及到账户生成和解密的参照account模块，其单元测试在 ../src/bif-sdk-go/test/account

## 1. 节点可信合约
### 1)  GetActive(id string) (bool, error)
```
  GetActive:
   	EN -
	CN - 查看证书是否有效
  Params:
  	- id: string，节点证书的bid

  Returns:
  	- bool，true可用，false不可用
	- error

  Call permissions: Anyone
```

### 2)  GetPeerCertificate(id string) (*dto.PeerCertificate, error)
```
  GetPeerCertificate:
   	EN -
	CN - 查看证书的信息
  Params:
  	- id: string，节点证书的bid

  Returns:
  	- *dto.PeerCertificate
		Id          string   `json:"id"`          //唯一索引
		Issuer      string   `json:"issuer"`      //颁发者地址
		Apply       string   `json:"apply"`       // 申请人bid
		PublicKey   string   `json:"publicKey"`   //节点公钥
		NodeName    string   `json:"nodeName"`    //节点名称
		Signature   string   `json:"signature"`   //节点签名内容
		NodeType    uint64   `json:"nodeType"`    //节点类型0企业，1个人
		CompanyName string   `json:"companyName"` //公司名称
		CompanyCode string   `json:"companyCode"` //公司信用代码
		IssuedTime  *big.Int `json:"issuedTime"`  //颁发时间
		Period      uint64   `json:"period"`      //有效期
		IsEnable    bool     `json:"isEnable"`    //true 凭证有效，false 凭证已撤销
	- error

  Call permissions: Anyone
```

### 3)  GetPeerCertificateIdList(id string) ([]string, error)
```
 GetPeerCertificateIdList:
  	EN - Get applied certificates by bid
	CN - 根据节点可信证书申请人的bid获取申请的证书列表
 Params:
 	- id: string，节点证书的bid

 Returns:
 	- []string, 申请人申请的证书列表
	- error

 Call permissions: Anyone
```

### 4)  GetPeriod(id string) (uint64, error)
```
  GetPeriod:
   	EN -
	CN - 查看证书有效期
  Params:
  	- id: string，节点证书的bid

  Returns:
  	- uint64，返回证书有效期，如果证书被吊销，则有效期是0
	- error

  Call permissions: Anyone
```

### 5)  RegisterCertificate(signTxParams *SysTxParams, registerCertificate *dto.RegisterCertificateInfo) (string, error)
```
  RegisterCertificate:
   	EN -
	CN - 为节点颁发可信证书可信
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- registerCertificate:  *dto.RegisterCertificateInfo，包含可信证书的信息
		Id          string //节点证书的bid,，必须和public_key相同
		Apply       string
		PublicKey   string // 53个字符的公钥
		NodeName    string // 节点名称，不含敏感词的字符串
		MessageSha3 string // 消息sha3后的16进制字符串
		Signature   string // 对上一个字段消息的签名，16进制字符串
		NodeType    uint64 // 节点类型，0企业，1个人
		Period      uint64 // 证书有效期，以年为单位的整型
		IP          string // ip
		Port        uint64 // port
		CompanyName string // 公司名（如果是个人，则是个人姓名）
		CompanyCode string // 公司代码

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以调用
```

### 6)  RevokedCertificate(signTxParams *SysTxParams, id string) (string, error)
```
  RevokedCertificate:
   	EN -
	CN - 吊销节点的可信证书可信
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string，节点证书的bid

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以调用
```

## 2. did文档合约
### 1)   AddAuth(signTxParams *SysTxParams, id string, auth string) (string, error)
```
AddAuth:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: ??
```

### 2)   AddExtra(signTxParams *SysTxParams, id string, extra string) (string, error)
```
AddExtra:
   	EN -
	CN - 添加用户的基本信息
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string, bid
	- extra: string, ??

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
```

### 3)   AddProof(signTxParams *SysTxParams, id string, proofType string, proofCreator string, proofSign string) (string, error)
```
 AddProof:
   	EN -
	CN - 增加证明
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- proofType: string, ??
	- proofCreator: string, ??
	- proofSign: string, ??

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
```

### 4)   AddPublic(signTxParams *SysTxParams, id string, publicType string, publicAuth string, publicKey string) (string, error)
```
 AddPublic:
   	EN -
	CN - 增加用户did身份认证
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- publicType: string, 公钥类型（secp256k1、SM2）
	- publicAuth: string, 公钥权限（all、update、ban）
	- publicKey: string,  公钥（十六进制的字符串）

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ？？
```

### 5)   AddService(signTxParams *SysTxParams, id string, serviceId string, serviceType string, serviceEndpoint string) (string, error)
```
 AddService:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: Anyone
```

### 6)   DelAuth(signTxParams *SysTxParams, id string, auth string) (string, error)
```
  DelAuth:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: ??
```

### 7)   DelExtra(signTxParams *SysTxParams, id string) (string, error)
```
  DelExtra:
   	EN -
	CN - 删除用户的基本信息
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string, bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: Anyone
```

### 8)   DelProof(signTxParams *SysTxParams, id string) (string, error)
```
  DelProof:
   	EN -
	CN - 删除证明
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ？？
```

### 9)   DelPublic(signTxParams *SysTxParams, id string, publicKey string) (string, error)
```
  DelPublic:
   	EN -
	CN - 删除用户did身份认证
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- publicKey: string,   公钥

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
```

### 10)  DelService(signTxParams *SysTxParams, id string, serviceId string) (string, error)
```
  DelService:
   	EN -
	CN -
  Params:
  	-
  	-

  Returns:
  	-
	- error

  Call permissions: Anyone
```

### 11)  Disable(signTxParams *SysTxParams, id string) (string, error)
```
  Disable:
   	EN -
	CN - 使用户的Did身份不可用
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
  	- id string bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
```

### 12)  Enable(signTxParams *SysTxParams, id string) (string, error)
```
  Enable:
   	EN -
	CN - 使用户的Did身份可用
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
  	- id string bid

  Returns:
  	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
	- error

  Call permissions: ??
```

### 13)  GetDocument(isDidAddress bool, did string) (*dto.Document, error)
```
  GetDocument:
   	EN -
	CN - 查询文档的信息
  Params:
  	- isDidAddress: bool，如果是true则第二个参数传did的地址；否则传入bidName
	- did: string，did文档的地址或bidName

  Returns:
  	- *dto.Document
		Id              utils.Address `json:"id"` // bid
		Contexts        []byte        `json:"context"`
		Name            []byte        `json:"name"`            // bid标识符昵称
		Type            []byte        `json:"type"`            // bid的类型，包括0: 普通用户,1:智能合约以及设备，2: 企业或者组织，BID类型一经设置，永不能变
		PublicKeys      []byte        `json:"publicKeys"`      // 用户用于身份认证的公钥信息
		Authentications []byte        `json:"authentications"` // 用户身份认证列表信息
		Attributes      []byte        `json:"attributes"`      // 用户填写的个人信息值
		IsEnable        []byte        `json:"is_enable"`       // 该BID是否启用
		CreateTime      time.Time     `json:"createTime"`
		UpdateTime      time.Time     `json:"updateTime"`
	- error

  Call permissions: ？？
```

### 14)  Init(signTxParams *SysTxParams, bidType uint64) (string, error)
```
  init:
   	EN -
	CN - did文档初始化
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- bidType: uint64，0: 普通用户,1:智能合约以及设备，2: 企业或者组织，BID类型一经设置，永不能变

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
```

### 15)  IsEnable(isDidAddress bool, did string) (bool, error)
```
  IsEnable:
   	EN -
	CN - 查询文档是否可用
  Params:
  	- isDidAddress: bool，如果是true则第二个参数传did的地址；否则传入bidName
	- did: string，did文档的地址或bidName

  Returns:
  	- bool, true可用，false不可用
	- error

  Call permissions: Anyone
```

### 16)  SetBidName(signTxParams *SysTxParams, id string, bidName string) (string, error)
```
  SetBidName:
   	EN -
	CN - 设置昵称
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id: string      bid
	- bidName: string 昵称(字符串长度6~20)

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
```

## 3. DPoS投票合约
### 1)   CancelProxy(signTxParams *SysTxParams) (string, error)
```
  CancelProxy:
   	EN -
	CN - 取消代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 2)   CancelVote(signTxParams *SysTxParams) (string, error)
```
 CancelVote:
   	EN -
	CN - 撤销投票
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 3)   ExtractOwnBounty(signTxParams *SysTxParams) (string, error)
```
  ExtractOwnBounty:
   	EN -
	CN - 取出自身的赏金
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
```

### 4)   GetAllCandidates() ([]*dto.Candidate, error)
```
  GetAllCandidates:
   	EN -
	CN - 查询所有候选人
  Params:
  	- None

  Returns:
  	- []dto.Candidate，列表内为候选人信息，参考GetCandidate的候选人信息
	- error

  Call permissions: Anyone
```

### 5)   GetCandidate(candidateAddress string) (*dto.Candidate, error)
```
  GetCandidate:
   	EN -
	CN -  查询候选人
  Params:
  	- candidateAddress: string，候选人的地址

  Returns:
  	- *dto.Candidate
		Owner           string       `json:"owner"`           // 候选人地址
		Name            string       `json:"name"`            // 候选人名称
		Active          bool         `json:"active"`          // 当前是否是候选人
		Url             string       `json:"url"`             // 节点的URL
		VoteCount       *hexutil.Big `json:"voteCount"`       // 收到的票数
		TotalBounty     *hexutil.Big `json:"totalBounty"`     // 总奖励金额
		ExtractedBounty *hexutil.Big `json:"extractedBounty"` // 已提取奖励金额
		LastExtractTime *hexutil.Big `json:"lastExtractTime"` // 上次提权时间
		Website         string       `json:"website"`         // 见证人网站
	- error

  Call permissions: Anyone
```

### 6)   GetRestBIFBounty() (*big.Int, error)
```
  GetRestBIFBounty:
   	EN -
	CN - 查询剩余的Bif总激励
  Params:
  	- None

  Returns:
  	- *big.Int
	- error

  Call permissions: Anyone
```

### 7)   GetStake(voterAddress string) (*dto.Stake, error)
```
  GetStake:
   	EN -
	CN - 查询抵押权益
  Params:
  	- voterAddress: string，投票者的地址

  Returns:
  	- *dto.Stake
		Owner              common.Address `json:"owner"`              // 抵押代币的所有人
		StakeCount         *big.Int       `json:"stakeCount"`         // 抵押的代币数量
		LastStakeTimeStamp *big.Int       `json:"lastStakeTimeStamp"` // 上次抵押时间戳
	- error

  Call permissions: Anyone
```

### 8)   GetVoter(voterAddress string) (*dto.Voter, error)
```
  GetVoter:
   	EN -
	CN - 查询投票人信息
  Params:
  	- voterAddress: string，投票者的地址

  Returns:
  	- *dto.Voter
		Owner             common.Address   `json:"owner"`             // 投票人的地址
		IsProxy           bool             `json:"isProxy"`           // 是否是代理人
		ProxyVoteCount    *big.Int         `json:"proxyVoteCount"`    // 收到的代理的票数
		Proxy             common.Address   `json:"proxy"`             // 该节点设置的代理人
		LastVoteCount     *big.Int         `json:"lastVoteCount"`     // 上次投的票数
		LastVoteTimeStamp *big.Int         `json:"lastVoteTimeStamp"` // 上次投票时间戳
		VoteCandidates    []common.Address `json:"voteCandidates"`    // 投了哪些人
	- error

  Call permissions: Anyone
```

### 9)   GetVoterList(voterAddress string) ([]*dto.Voter, error)
```
  GetVoterList:
   	EN -
	CN - 查询所有投票人信息
  Params:
  	- voterAddress: string，投票者的地址

  Returns:
  	- []dto.Voter，投票人的详细信息，参考GetVoter
	- error

  Call permissions: Anyone
```

### 10)  IssueAdditionalBounty(signTxParams *SysTxParams) (string, error)
```
  IssueAdditionalBounty:
   	EN -
	CN - ??
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？
```

### 11)  RegisterWitness(signTxParams *SysTxParams, witness *dto.RegisterWitness) (string, error)
```
  RegisterWitness:
   	EN -
	CN - 注册成为候选
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- witness: *dto.RegisterWitness，注册的见证人信息
		NodeUrl string
		Website string
		Name    string

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 12)  SetProxy(signTxParams *SysTxParams, proxy string) (string, error)
```
  SetProxy:
   	EN -
	CN - 设置代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- proxy: string，???

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 13)  Stake(signTxParams *SysTxParams, stakeCount *big.Int) (string, error)
```
  Stake:
   	EN -
	CN - 权益抵押
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- stakeCount: *big.Int，抵押的权益数量

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 14)  StartProxy(signTxParams *SysTxParams) (string, error)
```
  StartProxy:
   	EN -
	CN - 开启代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 15)  StopProxy(signTxParams *SysTxParams) (string, error)
```
  StopProxy:
   	EN -
	CN - 关闭代理
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 16)  UnRegisterWitness(signTxParams *SysTxParams) (string, error)
```
  UnRegisterWitness:
   	EN -
	CN -  取消成为候选
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只能自己取消自己
```

### 17)  UnStake(signTxParams *SysTxParams) (string, error)
```
  UnStake:
   	EN -
	CN - 撤销权益抵押
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```

### 18)  VoteWitnesses(signTxParams *SysTxParams, candidate string) (string, error)
```
  VoteWitnesses:
   	EN -
	CN - 给见证人投票
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- candidate: string，候选人的地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: ？？？
```
## 4. DPoS属性
### 1) Backlogs() (map[string][]*dto.Message, error)
```
Backlogs: 获取bft未来数据

Params:

- None
Returns:

- map[string][]*dto.Message
- error
```

### 2) GetValidators(blockNumber *big.Int) ([]string, error)
```
GetValidators: 根据区块数查询验证人

Params:
    - blockNumber: *big.Int，区块数

Returns:
    - []string
    - error
```


### 3) GetValidatorsAtHash(hash string) ([]string, error)
```
GetValidatorsAtHash: 根据区块hash查询验证人

Params:
    - hash: string，区块hash

Returns:
    - []string
    - error
```

### 4) RoundChangeSetInfo() (*dto.RoundChangeSetInfo, error)
```
RoundChangeSetInfo: 获取bft周期切换信息

Params:
    - None  

Returns:
    - *dto.RoundChangeSetInfo
        RoundChanges map[uint64]*MessageSet `json:"roundChanges"`
        Validates    []string               `json:"validates"`
    - error
Call permissions: Anyone
```

### 5) RoundStateInfo() (*dto.RoundStateInfo, error)
```
RoundStateInfo: 获取bft周期状态信息

Params:
    - None

Returns:
    - *dto.RoundStateInfo
        Commits    *MessageSet `json:"commits"`
        LockedHash common.Hash `json:"lockedHash"`
        Prepares   *MessageSet `json:"prepares"`
        Proposer   string      `json:"proposer"`
        Round      *big.Int    `json:"round"`
        Sequence   *big.Int    `json:"sequence"`
        View       *View       `json:"view"`
    - error

Call permissions: Anyone
```

## 5. 信任锚合约
### 1)   CancelVote(signTxParams *SysTxParams, candidate string) (string, error)
```
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
```

### 2)   ExtractOwnBounty(signTxParams *SysTxParams) (string, error)
```
  ExtractOwnBounty:
   	EN -
	CN - 提取信任锚激励，只有超过100积分，且24小时内只能提取一次
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只能提取自己的
```

### 3)   GetBaseList() ([]string, error)
```
  GetBaseList:
   	EN -
	CN - 查询根信任锚列表
  Params:
  	- None

  Returns:
  	- []string，根信任锚列表
	- error

  Call permissions: Anyone
```

### 4)   GetBaseNum() (uint64, error)
```
  GetBaseNum:
   	EN -
	CN - 查询根信任锚个数
  Params:
  	- None

  Returns:
  	- uint64， 根信任锚个数
	- error

  Call permissions: Anyone
```

### 5)   GetCertificateList(anchor string) ([]string, error)
```
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
```

### 6)   GetExpendList() ([]string, error)
```
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
```

### 7)   GetExpendNum() (uint64, error)
```
  GetExpendNum:
   	EN -
	CN - 查询扩展信任锚个数
  Params:
  	- None

  Returns:
  	- uint64， 扩展信任锚个数
	- error

  Call permissions: Anyone
```

### 8)   GetTrustAnchor(anchor string) (*dto.TrustAnchor, error)
```
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
```

### 9)   GetTrustAnchorStatus(anchor string) (uint64, error)
```
  GetTrustAnchorStatus:
   	EN -
	CN - 查询信任锚状态
  Params:
  	- anchor: string，信任锚bid

  Returns:
  	- uint64，0未知，1可用，2错误，3删除
	- error

  Call permissions: Anyone
```

### 10)  GetVoter(voterAddress string) ([]*dto.TrustAnchorVoter, error)
```
  GetVoter:
   	EN -
	CN - 查询投票人信息
  Params:
  	- voterAddress: 投票人地址（也就是超级节点地址，因为只有超级节点才可以投票）

  Returns:
  	- []dto.TrustAnchorVoter， 投票人信息
	- error

  Call permissions: Anyone
```

### 11)  IsBaseTrustAnchor(anchor string) (bool, error)
```
  IsBaseTrustAnchor:
   	EN -
	CN - 查询bid地址是否为根信任锚
  Params:
  	- anchor: string，信任锚bid

  Returns:
  	- bool，true为是根信任锚，false为不是根信任锚
	- error

  Call permissions: Anyone
```

### 12)  IsTrustAnchor(anchor string) (bool, error)
```
  IsTrustAnchor:
   	EN -
	CN - 查询bid地址是否为信任锚
  Params:
  	- anchor: string，信任锚bid

  Returns:
  	- bool，true为是信任锚，false为不是信任锚
	- error

  Call permissions: Anyone
```

### 13)  RegisterTrustAnchor(signTxParams *SysTxParams, registerAnchor *dto.RegisterAnchor) (string, error)
```
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
```

### 14)  UnRegisterTrustAnchor(signTxParams *SysTxParams) (string, error)
```
  UnRegisterTrustAnchor:
   	EN -
	CN - 注销自己的信任锚，自动退回抵押。但是，需要手动批量吊销自己颁发的证书，如果存在未吊销的证书，则抵押不退回。
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只能注销自己的信任锚
```

### 15)  UpdateAnchorInfo(signTxParams *SysTxParams, extendAnchorInfo *dto.UpdateAnchorInfo) (string, error)
```
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
```

### 16)  VoteElect(signTxParams *SysTxParams, candidate string) (string, error)
```
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
```

## 6. 个人可信合约
### 1) GetActive(id string) (bool, error)
```
GetActive:
   	EN -
	CN -  查询证书是否可用,如果证书过期，则不可用；如果信任锚注销，则不可用
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- bool，true可用，false不可用
	- error

  Call permissions: Anyone
```

### 2) GetCertificate(id string) (*dto.CertificateInfo, error)
```
GetCertificate:
   	EN -
	CN -  查询证书的信息，period和active没有进行逻辑判断
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- *dto.CertificateInfo
		Id             string   //凭证的hash
		Context        string   //证书所属上下文环境
		Issuer         string   //信任锚的bid
		Subject        string   //证书拥有者地址
		IssuedTime     *big.Int //颁发时间
		Period         uint64   //有效期
		IsEnable       bool     //true 凭证有效，false 凭证已撤销
		RevocationTime *big.Int //吊销时间
	- error

  Call permissions:
```

### 3) GetIssuer(id string) (*dto.IssuerSignature, error)
```
GetIssuer:
   	EN -
	CN - 查询信任锚信息，证书颁发者的信息
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- *dto.IssuerSignature
		Id        string //凭证ID
		PublicKey string // 签名公钥
		Algorithm string //签名算法
		Signature string //签名内容
	- error

  Call permissions: Anyone
```

### 4) GetPeriod(id string) (uint64, error)
```
 GetPeriod:
   	EN -
	CN -  查询个人证书的有效期
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- uint64，证书有效期，如果证书被吊销，则有效期是0
	- error

  Call permissions: Anyone
```

### 5) GetSubject(id string) (*dto.SubjectSignature, error)
```
GetSubject:
   	EN -
	CN - 查询个人(证书注册的接收者)信息，证书接收者的信息
  Params:
  	- id: string,个人可信证书bid

  Returns:
  	- *dto.SubjectSignature
		Id        string //凭证ID
		PublicKey string // 签名公钥
		Algorithm string //签名算法
		Signature string //签名内容
	- error

  Call permissions: Anyone
```

### 6) RegisterCertificate(signTxParams *SysTxParams, registerCertificate *dto.RegisterCertificate) (string, error)
```
RegisterCertificate:
   	EN -
	CN - 信任锚颁发证书，如果是根信任锚颁发的证书，则证书接收者可以进行部署合约和大额转账操作
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- registerCertificate:  The registerCertificate object(*dto.RegisterCertificate)
		Id               string //个人可信证书bid
		Context          string //证书上下文环境，随便一个字符串，不验证
		Subject          string //证书接收者的bid，证书是颁给谁的
		Period           uint64 //证书有效期，以年为单位的整型
		IssuerAlgorithm  string // 颁发者签名算法，字符串
		IssuerSignature  string //颁发者签名值，16进制字符串
		SubjectPublicKey string // 接收者公钥，16进制字符串
		SubjectAlgorithm string //接收者签名算法，字符串
		SubjectSignature string //接收者签名值，16进制字符串

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有信任锚地址可以调用
```

### 7) RevokedCertificate(signTxParams *SysTxParams, id string) (string, error)
```
 RevokedCertificate:
   	EN -
	CN - 信任锚吊销个人证书
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- id:  string，个人可信证书bid

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有证书颁发者可以调用
```

### 8) RevokedCertificates(signTxParams *SysTxParams) (string, error)
```
 RevokedCertificates:
   	EN -
	CN - 信任锚批量吊销个人证书，把自己颁发的证书全部吊销
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error
```

## 7. 敏感词合约
### 1) AddWord(signTxParams *SysTxParams, word string) (string, error)
```
  AddWord:
   	EN -
	CN - 向合约中添加敏感词
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- word: string，敏感词

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以操作
```


### 2) AddWords(signTxParams *SysTxParams, wordsLi []string) (string, error)
```
  AddWords:
   	EN -
	CN - 批量向合约中添加敏感词
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- wordsLi: []string，敏感词列表

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以操作
```

### 3) DelWord(signTxParams *SysTxParams, word string) (string, error)
```
  DelWord:
   	EN -
	CN - 删除合约中的字符串，只能一个一个地删除
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- word: string，敏感词

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 只有监管节点地址可以操作
```

### 4) GetAllWords() ([]string, error)
```
  GetAllWords:
   	EN -
	CN - 返回合约中保存的所有敏感词
  Params:
  	- None

  Returns:
  	- []string，返回敏感词列表
	- error

  Call permissions: Anyone
```

### 5) IsContainWord(word string) (bool, error)
```
  IsContainWord:
   	EN -
	CN - 查询词语是否包含敏感词
  Params:
  	- word: string，敏感词

  Returns:
  	- bool，true包含，false不包含
	- error

  Call permissions: Anyone
```

## 8. 超级管理合约
### 1) Disable(signTxParams *SysTxParams, contractAddress string) (string, error)
```
  Disable:
   	EN -
	CN - 禁用合约，合约被禁用后，不能再向合约中发送send交易，但是可以发送call交易（RPC查询）
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- contractAddress: string，合约地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 监管节点地址，权限包含2的地址
```

### 2) Enable(signTxParams *SysTxParams, contractAddress string) (string, error)
```
  Enable:
   	EN -
	CN - 启用合约
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- contractAddress: string，合约地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 监管节点地址，权限包含1的地址
```

### 3) GetAllContracts() ([]*dto.AllContract, error)
```
  GetAllContracts:
   	EN -
	CN - 返回该合约管理的所有合约及合约是否启用
  Params:
  	- None

  Returns:
  	- []*dto.AllContract
	- error

  Call permissions: Anyone
```

### 4) GetPower(userAddress string) (uint64, error)
```
  GetPower:
   	EN -
	CN - 查询用户的权限
  Params:
  	- userAddress string 用户地址

  Returns:
  	- uint64，1启用合约，2禁用合约，4授权  // 3=1+2启用禁用, 5=1+4启用授权, 6=2+4禁用授权, 7=1+2+4启用禁用授权 类linux权限管理
	- error

  Call permissions: Anyone
```

### 5) IsEnable(contractAddress string) (bool, error)
```
  IsEnable:
   	EN -
	CN - 合约是否启用，未被该合约管理的合约，是启用状态
  Params:
  	- contractAddress string 合约地址

  Returns:
  	- bool，true启用，false禁用
	- error

  Call permissions: Anyone
```

### 6) SetPower(signTxParams *SysTxParams, userAddress string, power uint64) (string, error)
```
  SetPower:
   	EN -
	CN - 为用户授权限，使用户可以代替监管节点地址操作该合约
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- userAddress: string，用户地址
	- power: uint64，权限（1是启用，2禁用，4授权。权限和权限可以累加，类linux权限，比如3就是启用禁用权限）

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 监管节点地址，权限包含4的地址
```
