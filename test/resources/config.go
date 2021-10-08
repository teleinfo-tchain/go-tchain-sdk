package resources

// base config
const (
	IP00 string = "127.0.0.1"
	IP51 string = "172.17.6.51"
	IP52 string = "172.17.6.52"
	IP54 string = "172.17.6.54"
	IP55 string = "172.17.6.55"

	Port          uint64 = 44002
	WebsocketPort uint64 = 44003

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
	SystemPassword = "teleInfo" 

	// PersonCertificate 注册证书的bid
	PersonCertificate = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	// PersonCertificatePublicKey 与personCertificate对应的公钥
	PersonCertificatePublicKey = "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"
	TestAddressCertificate     = "did:bid:EFTTQWPMdtghuZByPsfQAUuPkWkWYb"
	TestAddressCertificateFile = "../resources/keystore/UTC--2020-08-19T05-48-46.004537900Z--did-bid-EFTTQWPMdtghuZByPsfQAUuPkWkWYb"

	TestAddressDoc          = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressDocPublicKey = "0x04647f729afb309e4cd20f4b186a7883e1cd23b245e9fb6eb939ad74e47cc16c55e60aa12f20ed21bee8d23291aae377ad319b166604dec1a81dfb2b008bdc3c68"
	TestAddressDocFile      = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"

	TestAddressManager = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressManagerFile    = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestContractAddress = ""

	TestAddressSen = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressSenFile    = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	
	TestAddressElect = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	TestAddressFile  = "../resources/keystore/UTC--2020-08-20T05-28-39.403642600Z--did-bid-EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
)
