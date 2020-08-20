package crypto

type (
	CryptoType     uint64
	HashType       uint64
	CodeType       uint64
	HashLengthType string
	AddressPrefix  string
)

const (
	SM2        CryptoType = 0 // 国密算法
	SECP256K1  CryptoType = 1 // 以太坊签名加密算法
	HashLength            = 20

	HashLength20 HashLengthType = "T"
	HashLength21 HashLengthType = "U"
	HashLength22 HashLengthType = "V"
	HashLength23 HashLengthType = "W"
	HashLength24 HashLengthType = "X"
	HashLength25 HashLengthType = "Y"
	HashLength26 HashLengthType = "Z"

	SM2_Prefix       AddressPrefix = "Z"
	SECP256K1_Prefix AddressPrefix = "E"

	BASE64_Prefix AddressPrefix = "S"
	BASE58_Prefix AddressPrefix = "F"
)

type Config struct {
	crypto_type CryptoType
}

func (c *Config) CryptoType() CryptoType {
	return c.crypto_type
}

func (c *Config) SetCryptoType(_type CryptoType) {
	c.crypto_type = _type
}
