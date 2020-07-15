package system

import (
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
)

type DPoS struct {
	super *System
}

func (sys *System) NewDPoS() *DPoS {
	dPoS := new(DPoS)
	dPoS.super = sys
	return dPoS
}

func (dp *DPoS) GetValidators(blockNumber *big.Int) ([]string, error) {
	params := make([]interface{}, 1)
	params[0] = utils.IntToHex(blockNumber)

	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_getValidators", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToValidators()
}

func (dp *DPoS) GetValidatorsAtHash(hash string) ([]string, error) {
	params := make([]interface{}, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_getValidatorsAtHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToValidatorsAtHash()
}

func (dp *DPoS) RoundStateInfo() (*dto.RoundStateInfo, error) {
	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_roundStateInfo", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToRoundStateInfo()
}

func (dp *DPoS) RoundChangeSetInfo() (*dto.RoundChangeSetInfo, error) {
	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_roundChangeSetInfo", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToRoundChangeSetInfo()
}

func (dp *DPoS) Backlogs() (map[string][]*dto.Message, error) {
	pointer := &dto.RequestResult{}

	err := dp.super.provider.SendRequest(pointer, "dpos_backlogs", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBacklogs()
}
