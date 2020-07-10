package system

import (
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
)

type Dpos struct {
	super *System
}

func (system *System) NewDpos() *Dpos {
	Dpos := new(Dpos)
	Dpos.super = system
	return Dpos
}

func (dpos *Dpos) GetValidators(blockNumber *big.Int) ([]string, error) {
	params := make([]interface{}, 1)
	params[0] = utils.IntToHex(blockNumber)

	pointer := &dto.RequestResult{}

	err := dpos.super.provider.SendRequest(pointer, "dpos_getValidators", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToValidators()
}

func (dpos *Dpos) GetValidatorsAtHash(hash string) ([]string, error) {
	params := make([]interface{}, 1)
	params[0] = hash

	pointer := &dto.RequestResult{}

	err := dpos.super.provider.SendRequest(pointer, "dpos_getValidatorsAtHash", params)

	if err != nil {
		return nil, err
	}

	return pointer.ToValidatorsAtHash()
}

func (dpos *Dpos) RoundStateInfo() (*dto.RoundStateInfo, error) {
	pointer := &dto.RequestResult{}

	err := dpos.super.provider.SendRequest(pointer, "dpos_roundStateInfo", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToRoundStateInfo()
}

func (dpos *Dpos) RoundChangeSetInfo() (*dto.RoundChangeSetInfo, error) {
	pointer := &dto.RequestResult{}

	err := dpos.super.provider.SendRequest(pointer, "dpos_roundChangeSetInfo", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToRoundChangeSetInfo()
}

func (dpos *Dpos) Backlogs() (map[string][]*dto.Message, error) {
	pointer := &dto.RequestResult{}

	err := dpos.super.provider.SendRequest(pointer, "dpos_backlogs", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBacklogs()
}
