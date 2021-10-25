package resources

// base config
const (
	IP00 string = "39.99.132.122"
	IP51 string = "172.17.6.51"
	IP52 string = "172.17.6.52"
	IP54 string = "172.17.6.54"
	IP55 string = "172.17.6.55"

	Port          uint64 = 44002
	WebsocketPort uint64 = 5003

	ChainCode string = "tele"
	PassWord  string = "tele"
)

// accounts
const (
	Addr1    string = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	Addr1Hex string = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	Addr1Pri string = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	Addr2    string = "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"
	Addr2Hex string = "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"
	Addr2Pri string = "78a0fc8f2e8440e1cc13eb12e5eb0a76c70e4cb0b864dfcc4d9530832f259363"
)

const KeyStoreFile string = "/test/resources/keystore/"

// system contract
const (
	IsSM2          = true
	NotSm2         = false
	SystemPassword = "tele"

	TestAddressRegulatory     = "did:bid:tele:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	TestAddressRegulatoryFile = "UTC--2021-10-18T03-00-25.788254467Z--did:bid:tele:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"

	TestAddressAlliance     = "did:bid:tele:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	TestAddressAllianceFile = "UTC--2021-10-18T03-00-25.788254467Z--did:bid:tele:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"

	// RegisterAllianceOnePubKey 联盟合约配置
	RegisterAllianceOnePubKey = "16Uiu2HAm3z3rBzpH5tpFkdTxf7CU2JSdEDT4A6JH78ieKc69Aotp"
	RegisterAllianceOnePriKey = "78a0fc8f2e8440e1cc13eb12e5eb0a76c70e4cb0b864dfcc4d9530832f259363"
	RegisterAllianceOne       = "did:bid:tele:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"
	RegisterAllianceOneFile   = "UTC--2021-10-18T03-00-25.866164055Z--did:bid:tele:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY"

	RegisterAllianceTwoPubKey = "16Uiu2HAmPR6ruDPZnoAJEo8PzJXPBfaRT8ureQtkhabXgdMyuttc"
	RegisterAllianceTwoPriKey = "344d73bf98e1b9afdd6451ce7bb531d54b3c5cdbe6611d9252f3826b458ef5e9"
	RegisterAllianceTwo       = "did:bid:tele:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n"
	RegisterAllianceTwoFile   = "UTC--2021-10-18T03-00-25.911112270Z--did:bid:tele:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n"

	RegisterAllianceThreePubKey = "16Uiu2HAmAomJHsiKfnYsBAdKSbzqFdoZetoyKPEdG5g8vdDKPHny"
	RegisterAllianceThreePriKey = "7499f87d504b672f0671eaa7dcde51581bd3b1e9f5c7fcde1bade03766dbcdfa"
	RegisterAllianceThree       = "did:bid:tele:sfCXQusR8SEWgp8fQ9BQu61riWdDLCMN"
	RegisterAllianceThreeFile   = "UTC--2021-10-18T03-00-25.953464329Z--did:bid:tele:sfCXQusR8SEWgp8fQ9BQu61riWdDLCMN"

	// TestAddressElect 选举合约配置

	// PersonCertificate 证书合约配置
	PersonCertificate = ""

	// PersonCertificatePublicKey 与personCertificate对应的公钥
	PersonCertificatePublicKey = ""
	TestAddressCertificate     = ""
	TestAddressCertificateFile = ""

	// TestAddressDoc did文档合约配置
	TestAddressDoc = ""

	TestAddressDocPublicKey = ""
	TestAddressDocFile      = ""

	// TestAddressManager 管理合约配置
	TestAddressManager = ""

	TestAddressManagerFile = ""
	TestContractAddress    = ""

	// TestAddressSen 敏感词合约配置
	TestAddressSen     = ""
	TestAddressSenFile = ""
)

//node0
//rm -rf node1/gbif/


//node1
//rm -rf node1/gbif/
//./gbifBone init --datadir node1 --chain.code tele genesis.json
//./gbifBone --syncmode full --datadir node1 --verbosity 3 --port 5101 --rpc --rpcaddr 0.0.0.0 --rpcport 5102 --rpcapi core,gb,admin,personal,db,net,txpool,bif,dpos,alliance,election,certificate,document,sensitive,supermanager --unlock did:bid:tele:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY --password ./password.txt --chain.code tele --generate console

//node2
//rm -rf ./node2/gbif/
//./gbifBone init --datadir ./node2/ --chain.code tele genesis.json
//./gbifBone --syncmode full --datadir node2 --verbosity 3 --port 5201 --rpc --rpcaddr 0.0.0.0 --rpcport 5202 --rpcapi core,gb,admin,personal,db,net,txpool,bif,dpos,alliance,election,certificate,document,sensitive,supermanager --unlock did:bid:tele:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n --password ./password.txt --chain.code tele --generate console

//node3
//rm -rf node3/gbif/
//./gbifBone init --datadir ./node3/ --chain.code tele genesis.json
//./gbifBone --syncmode full --datadir node3 --verbosity 3 --port 5301 --rpc --rpcaddr 0.0.0.0 --rpcport 5302 --rpcapi core,gb,admin,personal,db,net,txpool,bif,dpos,alliance,election,certificate,document,sensitive,supermanager --unlock did:bid:tele:sfCXQusR8SEWgp8fQ9BQu61riWdDLCMN --password ./password.txt --chain.code tele --generate console
