package config

type (
	CryptoType     uint64
	HashType       uint64
	CodeType       uint64
	HashLengthType string
	AddressPrefix  string
)

const (
	SM2          CryptoType = 0 // 国密算法
	SECP256K1    CryptoType = 1 // 以太坊签名加密算法
	SECP256K1SM3 CryptoType = 2
	SM2SHA3      CryptoType = 3

	BASE64 CodeType = 0
	BASE58 CodeType = 1

	//BASE58String     = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	//HashLengthString = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	HashLength = 22

	HashLength20 HashLengthType = "T"
	HashLength21 HashLengthType = "U"
	HashLength22 HashLengthType = "V"
	HashLength23 HashLengthType = "W"
	HashLength24 HashLengthType = "X"
	HashLength25 HashLengthType = "Y"
	HashLength26 HashLengthType = "Z"

	ED25519_Prefix   AddressPrefix = "e"
	SM2_Prefix       AddressPrefix = "z"
	SECP256K1_Prefix AddressPrefix = "s"
	//SECP256K1SM3_Prefix AddressPrefix = "Z"
	//SM2SHA3_Prefix      AddressPrefix = "G"

	BASE32_Prefix AddressPrefix = "t"
	BASE64_Prefix AddressPrefix = "s"
	BASE58_Prefix AddressPrefix = "f"
)
