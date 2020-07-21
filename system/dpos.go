package system

import (
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
)

// DPoS - The DPoS Module
type DPoS struct {
	super *System
}

// NewDPoS - NewDPoS
func (sys *System) NewDPoS() *DPoS {
	dPoS := new(DPoS)
	dPoS.super = sys
	return dPoS
}

/*
GetValidators: 根据区块数查询验证人

Params:
	- blockNumber: *big.Int，区块数

Returns:
	- []string
	- error

Call permissions: Anyone

BUG(agl): rpc接收的类型为*rpc.BlockNumber？？？
*/
func (dp *DPoS) GetValidators(blockNumber *big.Int) ([]string, error) {
	params := make([]interface{}, 1)
	params[0] = utils.IntToHex(blockNumber)

	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_getValidators", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}

/*
GetValidatorsAtHash: 根据区块hash查询验证人

Params:
	- hash: string，区块hash

Returns:
	- []string
	- error

Call permissions: Anyone
*/
func (dp *DPoS) GetValidatorsAtHash(hash string) ([]string, error) {
	params := make([]interface{}, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_getValidatorsAtHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}

/*
RoundStateInfo: 获取bft周期状态信息

Params:
	- None

Returns:
	- *dto.RoundStateInfo
		Commits    *MessageSet `json:"commits"`
		LockedHash common.Hash `json:"lockedHash"`
		Prepares   *MessageSet `json:"prepares"`
		Proposer   string      `json:"proposer"`
		Round      *big.Int    `json:"round"`
		Sequence   *big.Int    `json:"sequence"`
		View       *View       `json:"view"`
	- error

Call permissions: Anyone
*/
func (dp *DPoS) RoundStateInfo() (*dto.RoundStateInfo, error) {
	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_roundStateInfo", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToRoundStateInfo()
}

/*
RoundChangeSetInfo: 获取bft周期切换信息

Params:
	- None

Returns:
	- *dto.RoundChangeSetInfo
		RoundChanges map[uint64]*MessageSet `json:"roundChanges"`
		Validates    []string               `json:"validates"`
	- error

Call permissions: Anyone
*/
func (dp *DPoS) RoundChangeSetInfo() (*dto.RoundChangeSetInfo, error) {
	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_roundChangeSetInfo", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToRoundChangeSetInfo()
}

/*
Backlogs: 获取bft未来数据

Params:
	- None

Returns:
	- map[string][]*dto.Message
	- error

Call permissions: Anyone
*/
func (dp *DPoS) Backlogs() (map[string][]*dto.Message, error) {
	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_backlogs", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBacklogs()
}
