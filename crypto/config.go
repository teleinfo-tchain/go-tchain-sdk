package crypto

type CryptoType uint64

const (
	SM2       CryptoType = 0 // 国密算法
	SECP256K1 CryptoType = 1 // 以太坊签名加密算法
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
