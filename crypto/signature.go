package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/tchain/go-tchain-sdk/utils/hexutil"

	"github.com/tchain/go-tchain-sdk/crypto/config"
	"github.com/tchain/go-tchain-sdk/crypto/secp"
	"github.com/tchain/go-tchain-sdk/crypto/sm2"
	"github.com/tchain/go-tchain-sdk/utils"

	"github.com/btcsuite/btcd/btcec"

	gmsm2 "github.com/tchain/go-tgmsm/sm2"
)

type Signature struct {
	PublicKey  []byte `json:"publicKey"    gencodec:"required"`  // 公钥，33字节，第1个字节是类型0, 1, 2，3，后32字节是公钥的x
	CryptoType []byte `json:"cryptoType"    gencodec:"required"` // 签名类型，0是sm2，1是secp
	Signature  []byte `json:"signature"    gencodec:"required"`  // 签名，64字字，前32字节是r，后32字节是s
}

func (s *Signature) String() string {
	return fmt.Sprintf("{cryptoType: %s, publicKey: %s, signature}", hexutil.Encode(s.CryptoType), hexutil.Encode(s.PublicKey), hexutil.Encode(s.Signature))
}

func NewSignature(hash []byte, prv *ecdsa.PrivateKey, cryptoType config.CryptoType) (*Signature, error) {
	signature := &Signature{
		PublicKey:  make([]byte, 0, 33),
		CryptoType: make([]byte, 0, 1),
		Signature:  make([]byte, 0, 64),
	}

	sig, err := sign(hash, prv, cryptoType)
	if err != nil {
		return signature, err
	}

	var pubkey []byte
	var ct []byte

	switch cryptoType {
	case config.SM2:
		pubkey = sm2.CompressPubkeySm2(&prv.PublicKey)
		ct = []byte{0}
	case config.SECP256K1:
		pubkey = secp.CompressPubkeyBtc(&prv.PublicKey)
		ct = []byte{1}
	default:
		pubkey = sm2.CompressPubkeySm2(&prv.PublicKey)
		ct = []byte{0}
	}

	signature.PublicKey = pubkey
	signature.CryptoType = ct
	signature.Signature = sig[:64]

	return signature, nil
}

func sign(hash []byte, prv *ecdsa.PrivateKey, cryptoType config.CryptoType) (sig []byte, err error) {
	switch cryptoType {
	case config.SM2:
		//start := time.Now()
		sig, err = sm2.SignSm2(hash, prv)
		//log.Debug("Sign SM2", "nanoseconds", utils.PrettyDuration(time.Since(start)))
	case config.SECP256K1:
		//start := time.Now()
		sig, err = secp.SignBtc(hash, prv)
		//log.Debug("Sign SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
	default:
		//start := time.Now()
		sig, err = secp.SignBtc(hash, prv)
		//log.Debug("Sign SECP256K1", "nanoseconds", utils.PrettyDuration(time.Since(start)))
	}

	if bytes.Count(sig[:32], []byte{0}) == 32 {
		log.Error("Sign, r is nil", "sig", sig, "err", err)
		return nil, errors.New("r of sig of Sign is nil")
	}

	if bytes.Count(sig[32:64], []byte{0}) == 32 {
		log.Error("Sign, s is nil", "sig", sig, "err", err)
		return nil, errors.New("s of sig of Sign is nil")
	}

	return sig, err
}

func Sign2Address(hash []byte, signature *Signature) (addr utils.Address, pubKey []byte, err error) {

	var pub *ecdsa.PublicKey
	switch signature.CryptoType[0] {
	case 0:
		pub, err = sm2.DecompressPubkeySm2(signature.PublicKey)
		if err != nil {
			return utils.Address{}, nil, fmt.Errorf("decompress public is error, err is %s", err.Error())
		}
		if !sm2.VerifySignatureSm2(pub, hash, signature.Signature) {
			return utils.Address{}, nil, fmt.Errorf("signature verify is error")
		}
		pubKey = (*gmsm2.PublicKey)(pub).SerializeUncompressed()
	case 1:
		pub, err = secp.DecompressPubkeyBtc(signature.PublicKey)
		if err != nil {
			return utils.Address{}, nil, fmt.Errorf("decompress public is error, err is %s", err.Error())
		}
		if !secp.VerifySignatureBtc(signature.PublicKey, hash, signature.Signature) {
			return utils.Address{}, nil, fmt.Errorf("signature verify is error")
		}
		pubKey = (*btcec.PublicKey)(pub).SerializeUncompressed()
	default:
		pub, err = sm2.DecompressPubkeySm2(signature.PublicKey)
		if err != nil {
			return utils.Address{}, nil, fmt.Errorf("decompress public is error, err is %s", err.Error())
		}
		if !sm2.VerifySignatureSm2(pub, hash, signature.Signature) {
			return utils.Address{}, nil, fmt.Errorf("signature verify is error")
		}
		pubKey = (*gmsm2.PublicKey)(pub).SerializeUncompressed()
	}
	if err != nil {
		return utils.Address{}, nil, err
	}

	return PubkeyToAddress(*pub), pubKey, nil
}

// S256 returns an instance of the secp256k1 curve.
func S256(cryptoType config.CryptoType) elliptic.Curve {
	switch cryptoType {
	case config.SM2:
		return sm2.S256Sm2()
	case config.SECP256K1:
		return secp.S256Btc()
	default:
		return sm2.S256Sm2()
	}
}
