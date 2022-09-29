package system

import (
	"errors"
	"fmt"
	"github.com/tchain/go-tchain-sdk/abi"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/prometheus/common/log"
	"regexp"
	"strings"
)

type SubChain struct {
	super *System
	abi   abi.ABI
}

// NewElection - NewElection初始化
func (sys *System) NewSubChain() *SubChain {
	parsedAbi, _ := abi.JSON(strings.NewReader(SubChainAbiJSON))

	subChain := new(SubChain)
	subChain.abi = parsedAbi
	subChain.super = sys
	return subChain
}

var (
	AddressRegexp   = regexp.MustCompile(`^(did:bid:(([a-z0-9]{4}):)?)?([zes][sft])([1-9a-km-zA-HJ-NP-Z]+)$`)
	ChainCodeRegexp = regexp.MustCompile(`^[a-z0-9]{4}$`)
)

func applySubChainPreCheck(subChainInfo *dto.SubChainInfo) error {
	subChainInfo.Id = strings.TrimSpace(subChainInfo.Id)
	subChainInfo.SubChainName = strings.TrimSpace(subChainInfo.SubChainName)
	subChainInfo.ChainCode = strings.TrimSpace(subChainInfo.ChainCode)
	subChainInfo.ChainIndustry = strings.TrimSpace(subChainInfo.ChainIndustry)
	subChainInfo.ChainFramework = strings.TrimSpace(subChainInfo.ChainFramework)
	subChainInfo.Consensus = strings.TrimSpace(subChainInfo.Consensus)
	subChainInfo.ChainMsgHash = strings.TrimSpace(subChainInfo.ChainMsgHash)
	subChainInfo.Apply = strings.TrimSpace(subChainInfo.Apply)

	switch {
	case len(subChainInfo.SubChainName) == 0:
		log.Error("subChain | parameter is not illegal", "parameter", "SubChainName")
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "SubChainName")
	case len(subChainInfo.ChainIndustry) == 0:
		log.Error("subChain | parameter is not illegal", "parameter", "ChainIndustry")
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "ChainIndustry")
	case len(subChainInfo.ChainFramework) == 0:
		log.Error("subChain | parameter is not illegal", "parameter", "ChainFramework")
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "ChainFramework")
	case len(subChainInfo.Consensus) == 0:
		log.Error("subChain | parameter is not illegal", "parameter", "Consensus")
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "Consensus")
	case len(subChainInfo.ChainMsgHash) == 0: // TODO: 对哈希增加正则校验
		log.Error("subChain | parameter is not illegal", "parameter", "ChainMsgHash")
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "ChainMsgHash")
	case !ChainCodeRegexp.MatchString(subChainInfo.ChainCode):
		log.Error("subChain | parameter is not illegal", "parameter", "ChainCode")
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is not 4 or invalid", "ChainCode")
	case !AddressRegexp.MatchString(subChainInfo.Id):
		log.Error("subChain | parameter is not illegal", "parameter", "Id")
		return fmt.Errorf("parameter is not illegal, parameter is %s", "Id")
	case !AddressRegexp.MatchString(subChainInfo.Apply):
		log.Error("subChain | parameter is not illegal", "parameter", "Apply")
		return fmt.Errorf("parameter is not illegal, parameter is %s", "Apply")
	}

	return nil
}

func (sc *SubChain) ApplySubChain(signTxParams *SysTxParams, subChain *dto.SubChainInfo) (string, error) {
	err := applySubChainPreCheck(subChain)
	if err != nil {
		return "", err
	}

	// encode
	var values []interface{}
	values = sc.super.structToInterface(*subChain, values)
	inputEncode, err := sc.abi.Pack("applySubChain", values...)

	if err != nil {
		return "", err
	}

	signedTx, err := sc.super.prePareSignTransaction(signTxParams, inputEncode, SubChainContract)
	if err != nil {
		return "", err
	}

	return sc.super.sendRawTransaction(signedTx)
}

func (sc *SubChain) VoteSubChain(signTxParams *SysTxParams, candidates string) (string, error) {
	candis := strings.TrimSpace(candidates)
	if len(candis) == 0 {
		return "", errors.New("candidateAddress is not valid bid")
	}

	// encoding
	inputEncode, err := sc.abi.Pack("voteSubChain", candis)
	if err != nil {
		return "", err
	}

	signedTx, err := sc.super.prePareSignTransaction(signTxParams, inputEncode, SubChainContract)
	if err != nil {
		return "", err
	}

	return sc.super.sendRawTransaction(signedTx)
}

func (sc *SubChain) SetDeadline(signTxParams *SysTxParams, deadline uint64) (string, error) {
	// encoding
	inputEncode, err := sc.abi.Pack("setDeadline", deadline)
	if err != nil {
		return "", err
	}

	signedTx, err := sc.super.prePareSignTransaction(signTxParams, inputEncode, SubChainContract)
	if err != nil {
		return "", err
	}

	return sc.super.sendRawTransaction(signedTx)
}

func (sc *SubChain) Revoke(signTxParams *SysTxParams, subChainId string, revokeReason string) (string, error) {
	if !isValidHexAddress(subChainId) {
		return "", errors.New("subChainId is not valid bid")
	}

	// Revoke is a struct we need to use the components.
	var values []interface{}
	type subChainInfo struct {
		SubChainId   string
		RevokeReason string
	}
	var node = subChainInfo{subChainId, revokeReason}

	values = sc.super.structToInterface(node, values)
	inputEncode, err := sc.abi.Pack("revoke", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := sc.super.prePareSignTransaction(signTxParams, inputEncode, SubChainContract)
	if err != nil {
		return "", err
	}

	return sc.super.sendRawTransaction(signedTx)
}

func (sc *SubChain) AllSubChains() ([]*dto.SubChainDetail, error) {
	pointer := &dto.SystemRequestResult{}

	err := sc.super.provider.SendRequest(pointer, "subchain_allSubChains", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToSubChains()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sc *SubChain) GetSubChain(subChainId string) (*dto.SubChainDetail, error) {
	if !isValidHexAddress(subChainId) {
		return nil, errors.New("subChainId is not valid bid")
	}

	params := make([]string, 1)
	params[0] = subChainId

	pointer := &dto.SystemRequestResult{}

	err := sc.super.provider.SendRequest(pointer, "subchain_subChain", params)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToSubChain()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sc *SubChain) GetDeadline() (uint64, error) {
	pointer := &dto.SystemRequestResult{}

	err := sc.super.provider.SendRequest(pointer, "subchain_deadline", nil)
	if err != nil {
		return 0, err
	}

	res, err := pointer.ToUint64()
	if err != nil {
		return 0, err
	}

	return res, nil
}
