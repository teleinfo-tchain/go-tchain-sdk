package account

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/bif/bif-sdk-go/crypto"
	"github.com/bif/bif-sdk-go/crypto/config"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"github.com/bif/bif-sdk-go/utils/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type SignTxParams struct {
	ChainId      uint64
	AccountNonce uint64
	GasPrice     *big.Int
	GasLimit     uint64
	Sender       *utils.Address
	Recipient    *utils.Address
	Amount       *big.Int
	Payload      []byte
}

type txData struct {
	ChainId      uint64         `json:"chainId"    gencodec:"required"`
	AccountNonce uint64         `json:"nonce"    gencodec:"required"`
	GasPrice     *big.Int       `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64         `json:"gas"      gencodec:"required"`
	Sender       *utils.Address `json:"sender"     rlp:"nil"` // nil means contract creation
	Recipient    *utils.Address `json:"recipient"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int       `json:"amount"    gencodec:"required"`
	Payload      []byte         `json:"payload"    gencodec:"required"`

	// account Signature values
	SignUser []byte `json:"signUser"    gencodec:"required"` // 33 + 1 + 32 + 32  即  公钥 + 类型 + r + s

	// This is only used when marshaling to JSON.
	Hash *utils.Hash `json:"hash" rlp:"-"`
}

type SignTransactionResult struct {
	Raw hexutil.Bytes `json:"raw"`
	Tx  *txData       `json:"tx"`
}

func rlpHash(x interface{}) (h utils.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type BIFSigner struct{}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (bfs BIFSigner) Hash(tx *txData) utils.Hash {
	return rlpHash([]interface{}{
		tx.ChainId,
		tx.AccountNonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.Sender,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
	})
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *txData, s BIFSigner, prv *ecdsa.PrivateKey, cryptoType config.CryptoType) (*txData, error) {
	h := s.Hash(tx)
	var err error
	Signature, err := crypto.GenSignature(h[:], prv, cryptoType)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	b.Write(Signature.PublicKey)
	b.Write(Signature.CryptoType)
	b.Write(Signature.Signature)
	tx.SignUser = b.Bytes()
	return tx, nil
}
