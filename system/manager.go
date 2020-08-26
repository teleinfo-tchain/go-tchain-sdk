package system

import (
	"errors"
	"github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/dto"
	"strings"
)

const (
	ManagerContractAddr = "did:bid:ZFT2ndSGfBuT1jKrsYBU5hLm7DmDV8u"
	ManagerAbiJSON      = `[
{"constant": false,"name":"enable","inputs":[{"name":"contract_address","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"disable","inputs":[{"name":"contract_address","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"power","inputs":[{"name":"user_address","type":"string"},{"name":"power","type":"uint64"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"superManagerEvent","type":"event"}
]`
)

// Manager - The Manager Module
type Manager struct {
	super *System
	abi   abi.ABI
}

// NewManager - 初始化Manager
func (sys *System) NewManager() *Manager {
	parseAbi, _ := abi.JSON(strings.NewReader(ManagerAbiJSON))

	manager := new(Manager)
	manager.super = sys
	manager.abi = parseAbi
	return manager
}

/*
  Enable:
   	EN -
	CN - 启用合约
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- contractAddress: string，合约地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 监管节点地址，权限包含1的地址
*/
func (manager *Manager) Enable(signTxParams *SysTxParams, contractAddress string) (string, error) {
	if !isValidHexAddress(contractAddress) {
		return "", errors.New("contractAddress is not valid address")
	}

	// encoding
	inputEncode, err := manager.abi.Pack("enable", contractAddress)
	if err != nil {
		return "", err
	}

	signedTx, err := manager.super.prePareSignTransaction(signTxParams, inputEncode, ManagerContractAddr)
	if err != nil {
		return "", err
	}

	return manager.super.sendRawTransaction(signedTx)
}

/*
  Disable:
   	EN -
	CN - 禁用合约，合约被禁用后，不能再向合约中发送send交易，但是可以发送call交易（RPC查询）
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- contractAddress: string，合约地址

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 监管节点地址，权限包含2的地址
*/
func (manager *Manager) Disable(signTxParams *SysTxParams, contractAddress string) (string, error) {
	if !isValidHexAddress(contractAddress) {
		return "", errors.New("contractAddress is not valid address")
	}

	// encoding
	inputEncode, err := manager.abi.Pack("disable", contractAddress)
	if err != nil {
		return "", err
	}

	signedTx, err := manager.super.prePareSignTransaction(signTxParams, inputEncode, ManagerContractAddr)
	if err != nil {
		return "", err
	}

	return manager.super.sendRawTransaction(signedTx)
}

/*
  SetPower:
   	EN -
	CN - 为用户授权限，使用户可以代替监管节点地址操作该合约
  Params:
  	- signTxParams *SysTxParams 系统合约构造所需参数
	- userAddress: string，用户地址
	- power: uint64，权限（1是启用，2禁用，4授权。权限和权限可以累加，类linux权限，比如3就是启用禁用权限）

  Returns:
  	- string, 交易哈希(transactionHash)，如果交易尚不可用，则为零哈希。
	- error

  Call permissions: 监管节点地址，权限包含4的地址
*/
func (manager *Manager) SetPower(signTxParams *SysTxParams, userAddress string, power uint64) (string, error) {
	if !isValidHexAddress(userAddress) {
		return "", errors.New("userAddress is not valid address")
	}
	if power != 1 && power != 2 && power != 3 && power != 4 && power != 5 && power != 6 {
		return "", errors.New("power is not illegal")
	}

	// encoding
	inputEncode, err := manager.abi.Pack("power", userAddress, power)
	if err != nil {
		return "", err
	}

	signedTx, err := manager.super.prePareSignTransaction(signTxParams, inputEncode, ManagerContractAddr)
	if err != nil {
		return "", err
	}

	return manager.super.sendRawTransaction(signedTx)
}

/*
  GetAllContracts:
   	EN -
	CN - 返回该合约管理的所有合约及合约是否启用
  Params:
  	- None

  Returns:
  	- []*dto.AllContract
	- error

  Call permissions: Anyone
*/
func (manager *Manager) GetAllContracts() ([]dto.AllContract, error) {
	pointer := &dto.SystemRequestResult{}

	err := manager.super.provider.SendRequest(pointer, "supermanager_allContracts", nil)
	if err != nil {
		return nil, err
	}

	return pointer.ToAllContract()
}

/*
  IsEnable:
   	EN -
	CN - 合约是否启用，未被该合约管理的合约，是启用状态
  Params:
  	- contractAddress string 合约地址

  Returns:
  	- bool，true启用，false禁用
	- error

  Call permissions: Anyone
*/
func (manager *Manager) IsEnable(contractAddress string) (bool, error) {
	if !isValidHexAddress(contractAddress) {
		return false, errors.New("contractAddress is not valid address")
	}

	params := make([]string, 1)
	params[0] = contractAddress

	pointer := &dto.SystemRequestResult{}

	err := manager.super.provider.SendRequest(pointer, "supermanager_isEnable", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

/*
  GetPower:
   	EN -
	CN - 查询用户的权限
  Params:
  	- userAddress string 用户地址

  Returns:
  	- uint64，1启用合约，2禁用合约，4授权  // 3=1+2启用禁用, 5=1+4启用授权, 6=2+4禁用授权, 7=1+2+4启用禁用授权 类linux权限管理
	- error

  Call permissions: Anyone
*/
func (manager *Manager) GetPower(userAddress string) (uint64, error) {
	if !isValidHexAddress(userAddress) {
		return 0, errors.New("userAddress is not valid address")
	}

	params := make([]string, 1)
	params[0] = userAddress

	pointer := &dto.SystemRequestResult{}

	err := manager.super.provider.SendRequest(pointer, "supermanager_power", params)
	if err != nil {
		return 0, err
	}

	return pointer.ToUint64()
}
