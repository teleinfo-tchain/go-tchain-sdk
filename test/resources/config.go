package resources

// base config
const (
	IP00 string = "127.0.0.1"
	IP51 string = "172.17.6.51"
	IP52 string = "172.17.6.52"
	IP54 string = "172.17.6.54"
	IP55 string = "172.17.6.55"

	Port          uint64 = 5002
	WebsocketPort uint64 = 5003

	ChainCode string = "tele"
	PassWord  string = "teleInfo"

)

// accounts
const (
	Addr1    string = "did:bid:sf26xSMerRVcs9E642Rkbq4TxLrhQsWzk"
	Addr1Hex string = ""
	Addr1Pri string = "d37df84af4156fe9ab65a5642418cd7bd9e9371acb5ae1bd282d1d473bcb1f13"

	Addr2    string = "did:bid:sf26xSMerRVcs9E642Rkbq4TxLrhQsWzk"
	Addr2Hex string = ""
	Addr2Pri string = "d37df84af4156fe9ab65a5642418cd7bd9e9371acb5ae1bd282d1d473bcb1f13"

	Addr3    string = "did:bid:sf26xSMerRVcs9E642Rkbq4TxLrhQsWzk"
	Addr3Hex string = ""
	Addr3Pri string = "d37df84af4156fe9ab65a5642418cd7bd9e9371acb5ae1bd282d1d473bcb1f13"

)

// system contract
const (
	IsSM2 = true
	NotSm2 = false
	SystemPassword = "teleinfo"

	FilePath = "../resources/keystore/"
	TestAddressAlliance     = "did:bid:tele:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	TestAddressAllianceFile = FilePath+"UTC--2021-10-11T02-48-43.619738164Z--did:bid:tele:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"

	// RegisterAllianceOnePubKey 联盟合约配置
	RegisterAllianceOnePubKey = "16Uiu2HAm5wRVEuPzTjjq18tXFNxUepnMAxPFA1UQyMxCingPHEAP"
	RegisterAllianceOnePriKey = "f50eca2610b6ec2a2cc7a781174f9051a5313813e66168970fa986613df796f0"
	RegisterAllianceOne = "did:bid:tele:sfx2ee4QMyUsJEUJZgC9GuoGLwGYMtAk"
	RegisterAllianceOneFile = FilePath+"UTC--2021-10-11T02-48-43.660543959Z--did:bid:tele:sfx2ee4QMyUsJEUJZgC9GuoGLwGYMtAk"

	RegisterAllianceTwoPubKey = "16Uiu2HAmGxSMc6uyMD44ofU6M2r4C4eL5Qs3LpRDWdtg6i7SG3E7"
	RegisterAllianceTwoPriKey = "b08e9b9c87ddd936213e55e9db005e2218fea7ee3cf752f3b0cc480871ead389"
	RegisterAllianceTwo = "did:bid:tele:sfyqv1R7kMYKxLcZMeS5ZQeez91ae1kt"
	RegisterAllianceTwoFile = FilePath+"UTC--2021-10-11T02-48-43.701398271Z--did:bid:tele:sfyqv1R7kMYKxLcZMeS5ZQeez91ae1kt"

	RegisterAllianceThreePubKey = "16Uiu2HAmBYGHwpUuBJEHfwisKqdzdbHwKADVtuT3Db8AiSsXNiJH"
	RegisterAllianceThreePriKey = "3c4a99e1600447ecd937574679a1a8fada4026bdc7aaa3e6f7909c99712b6b2a"
	RegisterAllianceThree = "did:bid:tele:sftkV9fC1fdcDaW1972fPs3CatUmc5SF"
	RegisterAllianceThreeFile = FilePath+"UTC--2021-10-11T02-48-43.740088399Z--did:bid:tele:sftkV9fC1fdcDaW1972fPs3CatUmc5SF"

	// TestAddressElect 选举合约配置
	TestAddressElect = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressFile  = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"

	RegisterTrustNodePubKey = "16Uiu2HAmBYGHwpUuBJEHfwisKqdzdbHwKADVtuT3Db8AiSsXNiJH"
	RegisterTrustNodePriKey = "3c4a99e1600447ecd937574679a1a8fada4026bdc7aaa3e6f7909c99712b6b2a"
	RegisterTrustNodeOne = "did:bid:tele:sftkV9fC1fdcDaW1972fPs3CatUmc5SF"
	RegisterTrustNodeOneFile = FilePath+"UTC--2021-10-11T02-48-43.740088399Z--did:bid:tele:sftkV9fC1fdcDaW1972fPs3CatUmc5SF"

	// PersonCertificate 证书合约配置
	PersonCertificate = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	// PersonCertificatePublicKey 与personCertificate对应的公钥
	PersonCertificatePublicKey = "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"
	TestAddressCertificate     = "did:bid:tele:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	TestAddressCertificateFile = FilePath+"UTC--2021-10-11T02-48-43.619738164Z--did:bid:tele:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"

	// TestAddressDoc did文档合约配置
	TestAddressDoc          = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressDocPublicKey = "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"
	TestAddressDocFile      = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"

	// TestAddressManager 管理合约配置
	TestAddressManager = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressManagerFile    = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestContractAddress = ""

	// TestAddressSen 敏感词合约配置
	TestAddressSen = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressSenFile    = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
)
