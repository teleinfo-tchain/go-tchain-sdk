package system

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/tchain/go-tchain-sdk/abi"
	"github.com/tchain/go-tchain-sdk/crypto"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/utils"
	libp2pcorecrypto "github.com/libp2p/go-libp2p-core/crypto"
	libp2pcorepeer "github.com/libp2p/go-libp2p-core/peer"
	"strings"
)

// Alliance - The Alliance Module
type Alliance struct {
	super *System
	abi   abi.ABI
}

// NewAlliance - NewAlliance初始化
func (sys *System) NewAlliance() *Alliance {
	parsedAbi, _ := abi.JSON(strings.NewReader(AllianceAbiJSON))

	ali := new(Alliance)
	ali.abi = parsedAbi
	ali.super = sys
	return ali
}

func registerDirectorPreCheck(directorInfo dto.AllianceInfo) error {
	// PublicKey length must be 53
	if len(directorInfo.PublicKey) != 53 {
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is not 53", "PublicKey")
	}

	// CompanyName can't be null
	if len(directorInfo.CompanyName) == 0 {
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "CompanyName")
	}

	// CompanyCode can't be null
	if len(directorInfo.CompanyCode) == 0 {
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "CompanyCode")
	}

	// Id can't be null
	address := utils.StringToAddress(directorInfo.Id)
	if address == utils.EmptyAddress {
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "Id")
	}

	_, Id, err := pub2addr(directorInfo.PublicKey)
	if err != nil {
		return err
	}

	//  Verify the address by PublicKey
	if address != Id {
		return fmt.Errorf("alliance | registerDirector id is not match the address by PublicKey")
	}
	return nil
}

func (ali *Alliance) RegisterDirector(signTxParams *SysTxParams, directorInfo *dto.AllianceInfo) (string, error) {
	err := registerDirectorPreCheck(*directorInfo)
	if err != nil {
		return "", err
	}

	// RegisterDirector is a struct we need to use the components.
	var values []interface{}
	values = ali.super.structToInterface(*directorInfo, values)
	inputEncode, err := ali.abi.Pack("registerDirector", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := ali.super.prePareSignTransaction(signTxParams, inputEncode, AllianceContract)
	if err != nil {
		return "", err
	}

	return ali.super.sendRawTransaction(signedTx)
}

func (ali *Alliance) UpgradeDirector(signTxParams *SysTxParams, director string) (string, error) {
	if !isValidHexAddress(director) {
		return "", errors.New("id is not valid bid")
	}

	// encoding
	inputEncode, err := ali.abi.Pack("upgradeDirector", director)
	if err != nil {
		return "", err
	}

	signedTx, err := ali.super.prePareSignTransaction(signTxParams, inputEncode, AllianceContract)
	if err != nil {
		return "", err
	}

	return ali.super.sendRawTransaction(signedTx)
}

func (ali *Alliance) Revoke(signTxParams *SysTxParams, member string, revokeReason string) (string, error) {
	if !isValidHexAddress(member) {
		return "", errors.New("member is not valid bid")
	}

	// Revoke is a struct we need to use the components.
	var values []interface{}
	type RevokeInfo struct {
		Member       string
		RevokeReason string
	}
	var revokeInfo = RevokeInfo{member, revokeReason}

	values = ali.super.structToInterface(revokeInfo, values)
	inputEncode, err := ali.abi.Pack("revoke", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := ali.super.prePareSignTransaction(signTxParams, inputEncode, AllianceContract)
	if err != nil {
		return "", err
	}

	return ali.super.sendRawTransaction(signedTx)
}

func (ali *Alliance) SetWeights(signTxParams *SysTxParams, directorWeights, viceWeights, directorGeneralWeights uint64) (string, error) {
	// Revoke is a struct we need to use the components.
	var values []interface{}
	type Weights struct {
		DirectorWeights        uint64
		ViceWeights            uint64
		DirectorGeneralWeights uint64
	}
	var weights = Weights{directorWeights, viceWeights, directorGeneralWeights}

	values = ali.super.structToInterface(weights, values)
	inputEncode, err := ali.abi.Pack("setWeights", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := ali.super.prePareSignTransaction(signTxParams, inputEncode, AllianceContract)
	if err != nil {
		return "", err
	}

	return ali.super.sendRawTransaction(signedTx)
}

func (ali *Alliance) AllDirectors() ([]*dto.Alliance, error) {
	pointer := &dto.SystemRequestResult{}

	err := ali.super.provider.SendRequest(pointer, "alliance_directors", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToAllianceDirectors()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ali *Alliance) AllVices() ([]*dto.Alliance, error) {
	pointer := &dto.SystemRequestResult{}

	err := ali.super.provider.SendRequest(pointer, "alliance_vices", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToAllianceVices()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ali *Alliance) AllDirectorGenerals() ([]*dto.Alliance, error) {
	pointer := &dto.SystemRequestResult{}

	err := ali.super.provider.SendRequest(pointer, "alliance_directorGenerals", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToAllianceDirectors()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ali *Alliance) AllAlliances() ([]*dto.Alliance, error) {
	pointer := &dto.SystemRequestResult{}

	err := ali.super.provider.SendRequest(pointer, "alliance_alliances", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToAllianceDirectors()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ali *Alliance) GetAlliance(id string) (*dto.Alliance, error) {
	if !isValidHexAddress(id) {
		return nil, errors.New("id is not valid bid")
	}

	params := make([]string, 1)
	params[0] = id

	pointer := &dto.SystemRequestResult{}

	err := ali.super.provider.SendRequest(pointer, "alliance_alliance", params)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToAlliance()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ali *Alliance) GetWeights() (*dto.Weights, error) {
	pointer := &dto.SystemRequestResult{}

	err := ali.super.provider.SendRequest(pointer, "alliance_weights", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToWeights()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func pub2addr(publicKey string) (*ecdsa.PublicKey, utils.Address, error) {
	peerID, err := libp2pcorepeer.Decode(publicKey)
	if err != nil {
		return nil, utils.Address{}, err
	}
	pk, err := peerID.ExtractPublicKey()
	if a, ok := pk.(*libp2pcorecrypto.Secp256k1PublicKey); ok {
		key := (*ecdsa.PublicKey)(a)
		id := crypto.PubkeyToAddress(*key)
		return key, id, nil
	}
	return nil, utils.Address{}, err
}
