package privacycontract

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tchain/go-tchain-sdk/abi"
	"github.com/tchain/go-tchain-sdk/crypto/privacyutils"
	"github.com/tchain/go-tchain-sdk/utils"
	"math/big"
	"strings"
)

type Root struct {
	super *BaseContract
	abi   abi.ABI
	bin   []byte
}

func (base *BaseContract) NewRoot() *Root {
	parseAbi, _ := abi.JSON(strings.NewReader(RootAbiJSON))

	Root := new(Root)
	Root.super = base
	Root.abi = parseAbi
	Root.bin, _ = hex.DecodeString(RootBin)
	return Root
}

// 部署合約
func (base *Root) RootDeploy(signTxParams *BaseTxParams) (string, error) {
	signedTx, err := base.super.prePareSignDeployTransaction(signTxParams, base.bin)
	if err != nil {
		return "", err
	}
	return base.super.sendRawTransaction(signedTx)
}

//發行接口
// 日志：Issue
func (base *Root) RootIssue(signTxParams *BaseTxParams, contractAddress string) (string, error) {
	if !isValidHexAddress(contractAddress) {
		return "", errors.New("contractAddress is not valid address")
	}
	var base_blind string = "1"
	var blind string = strings.Repeat(base_blind, 64)
	fmt.Println(blind)
	// 1.币的名称
	name := "spf"
	// 2.币的单位
	symbol := "s"
	// 3.構建ptr參數
	var ptr uint64 = privacyutils.InitPrivacy()
	// 4.構造發行commit參數
	commit, err := privacyutils.CreatePedersenCommit(ptr, 500, blind)
	var commitA []byte = []byte(commit)
	if err != nil {
		return "", err
	}
	// 5.構造發行rangeproof參數
	rangeproof, err := privacyutils.BpRangeproofProve(ptr, blind, 500)
	var rangeproofA []byte = []byte(rangeproof)
	if err != nil {
		return "", err
	}

	var from_pubkey string = "cabf47bb65af26e719678d4f16606dcd22c368"

	// encoding
	inputEncode, err := base.abi.Pack("issue", name, symbol, from_pubkey, commitA, rangeproofA)
	if err != nil {
		return "", err
	}

	signedTx, err := base.super.prePareSignTransaction(signTxParams, inputEncode, FinacialContract)
	if err != nil {
		return "", err
	}

	return base.super.sendRawTransaction(signedTx)
}

//轉移接口
func (base *Root) RootTransfer(signTxParams *BaseTxParams, contractAddress string) (string, error) {
	if !isValidHexAddress(contractAddress) {
		return "", errors.New("contractAddress is not valid address")
	}
	var base_blind string = "1"
	var blind string = strings.Repeat(base_blind, 64)
	// 1.構建ptr參數
	var ptr uint64 = privacyutils.InitPrivacy()
	// 2.構造轉移commit參數
	commit, err := privacyutils.CreatePedersenCommit(ptr, 500, blind)
	if err != nil {
		return "", err
	}
	var commits [1]string
	commits[0] = commit
	// 3.構造轉移rangeproof參數
	rangeproof, err := privacyutils.BpRangeproofProve(ptr, blind, 500)
	if err != nil {
		return "", err
	}
	var rangeproofs [1]string
	rangeproofs[0] = rangeproof
	// 4.構造轉移encryptvalue
	var encryptvalue = [1]string{"63542b9f52c15f79192d88e833cabf47bb65af26e719678d4f16606dcd22c368"}
	// 5.構造轉移from_pubkey
	var fromPubkeys = [1]string{"cabf47bb65af26e719678d4f16606dcd22c368"}
	// 6.構造inputId
	var inputId [1]*big.Int
	inputId[0] = big.NewInt(1)
	// 7.代筆類別
	var tokenId *big.Int = big.NewInt(1)
	// 8.接收人地址
	b := utils.StringToAddress("did:bid:qwer:zf21JocUsuvd4gMNFuK3qkCRS1AigY46Y").Bytes()
	fmt.Println("地址為：", b)
	//[122 102 214 212 58 80 181 176 117 79 97 66 78 232 69 163 152 254 150 2 250 95 108 213]
	toadds := [1]utils.Address{{122, 102, 214, 212, 58, 80, 181, 176, 117, 79, 97, 66, 78, 232, 69, 163, 152, 254, 150, 2, 250, 95, 108, 213}}

	inputEncode, err := base.abi.Pack("transfer", tokenId, inputId, commits, fromPubkeys, rangeproofs, toadds, encryptvalue)
	if err != nil {
		return "", err
	}

	signedTx, err := base.super.prePareSignTransaction(signTxParams, inputEncode, FinacialContract)
	if err != nil {
		return "", err
	}

	return base.super.sendRawTransaction(signedTx)
}

//查詢代筆名接口
func (base *Root) RootQueryName(signTxParams *BaseTxParams, contractAddress string) (string, error) {
	if !isValidHexAddress(contractAddress) {
		return "", errors.New("contractAddress is not valid address")
	}

	// encoding
	inputEncode, err := base.abi.Pack("queryName")
	if err != nil {
		return "", err
	}

	signedTx, err := base.super.prePareSignTransaction(signTxParams, inputEncode, FinacialContract)
	if err != nil {
		return "", err
	}

	return base.super.sendRawTransaction(signedTx)
}

//查詢代筆標志接口
func (base *Root) RootQuerySymbol(signTxParams *BaseTxParams, contractAddress string) (string, error) {
	if !isValidHexAddress(contractAddress) {
		return "", errors.New("contractAddress is not valid address")
	}

	// encoding
	inputEncode, err := base.abi.Pack("querySymbol")
	if err != nil {
		return "", err
	}

	signedTx, err := base.super.prePareSignTransaction(signTxParams, inputEncode, FinacialContract)
	if err != nil {
		return "", err
	}

	return base.super.sendRawTransaction(signedTx)
}
