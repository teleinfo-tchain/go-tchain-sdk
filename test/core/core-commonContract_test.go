/********************************************************************************
   This file is part of go-bif.
   go-bif is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-bif is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-bif.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/account"
	"github.com/bif/bif-sdk-go/compiler"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	bif "go-sdk"
	"math/big"
	"path"
	"runtime"
	"strconv"
	"testing"
	"time"
)

//  部署合约及合约的call方法调用
func TestCoreContract(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	solcDir := path.Join(bif.GetCurrentAbPath(), "compiler", "tmp")
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows" {
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	} else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solFile := path.Join(bif.GetCurrentAbPath(), "compiler", "contract", "simple-token.sol")

	solc := compiler.NewSolidityCompiler(solcDir)

	output, err := solc.Compile(solFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	type Contract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}
	var contract Contract

	for _, v := range output.Contracts {
		contract.Abi = v.Abi
		contract.ByteCode = v.Bin
	}

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(singAddr)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  200000,
		Sender:    &sender,
		Recipient: nil,
		Amount:    nil,
	}

	hash, err := contractObj.Deploy(tx, false, singAddrPriKey, contract.ByteCode)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txHash is ", hash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Contract Address: ", receipt.ContractAddress)

	transaction := new(dto.TransactionParameters)
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// txHash is  0x0e57cc103a20ebdbd3ac36530e440e95af68b4e6b3ce2f55d845218c7a7b1add
	// Contract Address:  did:bid:qwer:sfzxtY15dmcYSKN9r5sWs3GDBFvCU2Sq

	// chainId, _ := connection.Core.GetChainId()
	// nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)
	transaction.ChainId = chainId
	transaction.Sender = generator
	transaction.Recipient = receipt.ContractAddress
	transaction.AccountNonce = nonce.Uint64()

	result, err := contractObj.Call(transaction, "getName")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		val, _ := result.ToComplexString()
		fmt.Println(val)
	}
}

type Contract struct {
	Abi      string `json:"abi"`
	ByteCode string `json:"byteCode"`
}

func ballotContract() Contract {
	solcDir := path.Join(bif.GetCurrentAbPath(), "compiler", "tmp")
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows" {
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	} else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solFile := path.Join(bif.GetCurrentAbPath(), "compiler", "contract", "ballot.sol")

	solc := compiler.NewSolidityCompiler(solcDir)

	output, err := solc.Compile(solFile)
	if err != nil {
		panic(err)
	}

	var contract Contract

	for _, v := range output.Contracts {
		contract.Abi = v.Abi
		contract.ByteCode = v.Bin
	}

	return contract
}

func ballotDeploy(t *testing.T) string {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	contract := ballotContract()

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(singAddr)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  2000000,
		Sender:    &sender,
		Recipient: nil,
		Amount:    nil,
	}

	var key1 [32]byte
	var key2 [32]byte
	copy(key1[:], []byte("test1"))
	copy(key1[:], []byte("test2"))
	proposalNames := []utils.Hash{key1, key2}

	hash, err := contractObj.Deploy(tx, false, singAddrPriKey, contract.ByteCode, proposalNames)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("txHash is ", hash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Contract Address: ", receipt.ContractAddress)

	return hash
}

// 合约交互 查询 call 方法
func ballotCallWinnerName(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	transaction := new(dto.TransactionParameters)
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contractAddress := "did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)
	transaction.ChainId = chainId
	transaction.Sender = generator
	transaction.Recipient = contractAddress
	transaction.AccountNonce = nonce.Uint64()

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	result, err := contractObj.Call(transaction, "winnerName")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		val, _ := result.ToComplexString()
		fmt.Println(val)
	}
}

func ballotCallVoters(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	transaction := new(dto.TransactionParameters)
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contractAddress := "did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)
	transaction.ChainId = chainId
	transaction.Sender = generator
	transaction.Recipient = contractAddress
	transaction.AccountNonce = nonce.Uint64()

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	voter := utils.StringToAddress("did:bid:qwer:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n")

	result, err := contractObj.Call(transaction, "voters", voter)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		val, _ := result.ToComplexString()
		fmt.Println(val)
	}
}

// 合约交互 交易 send 方法
func ballotSend(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractAddress := "did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)

	sender := utils.StringToAddress(singAddr)
	recipient := utils.StringToAddress(contractAddress)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  2000000,
		Sender:    &sender,
		Recipient: &recipient,
		Amount:    nil,
	}

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// voter := utils.StringToAddress("did:bid:qwer:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n")
	// voter := utils.StringToAddress("did:bid:qwer:sfCXQusR8SEWgp8fQ9BQu61riWdDLCMN")
	voter := utils.StringToAddress("did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY")

	txHash, err := contractObj.Send(tx, false, singAddrPriKey, "giveRightToVote", voter)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("txHash is ", txHash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txHash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试合约的交互(只是为了测试合约交互，实际使用contract中的Call或者Send)
func TestCoreSendTransactionInteractContract(t *testing.T) {
	// 合约部署
	// txHash := ballotDeploy(t)
	// txHash is  0x8e17880962519fa1421f1eea19a3503290758db10ce37524239030b9b7aa17c3
	// Contract Address:  did:bid:qwer:sfjGbVtUc3RNBNhMdPBJXrmRN2tzTCH8

	// // 合约交互 查询 call 方法
	// ballotCallWinnerName(t)

	// // 合约交互 交易 send 方法
	// ballotSend(t)
	//
	// call voters
	ballotCallVoters(t)
}

func TestEVMQue(t *testing.T){
	// 编译器的问题，telChain的链已经做了其它的修改
	// var singAddr = "did:bid:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	// // 0x000000000000000000000000b39229c015c595850c16948d72f940303c8f3e11
	// // 0x66323558474251553845387747466f3977474b6f39356a55677459504d323459
	// // 0x66323558474251553845387747466f3977474b6f39356a55677459504d323459
	// 0xadFfd5A56Fc679aa9E8BfD89a7438b614f317ad1
	// fmt.Println(utils.StringToHash(singAddr).String())

	// // res, err := connection.Core.GetTransactionReceipt("0x86401d8ffac5ce804cef47ff89f5297a33ad33f3e9dad5bfb270352b330e4169")
	// res, err := connection.Core.GetTransactionReceipt("0x8451fa14eb1e301a40e5eb01ae075b7186cc6e0747a13600c44ae550394c72b0")
	// fmt.Println(res)
	// fmt.Println(err)
	// res, _ := hexutil.Decode("0xb39229c015c595850c16948d72f940303c8f3e11")
	// fmt.Println(res)
	// var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	// //
	// sender := utils.StringToAddress(singAddr)
	// fmt.Println(sender)
}

func TestTelChainEVM(t *testing.T){
	// // 基于新的合约编译器进行部署
	// ballotDeployTelChainEVM(t)
	// txHash 0x0d122db927efe4cf8247390486141a0b2e44d0b4fea2c70d61711589ebc47330
	// contractAddr did:bid:qwer:sfWkwvBdBaeAmApBJ6hn82DpAb9gc3M4

	// var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	// recp, err := connection.Core.GetTransactionReceipt("0x0d122db927efe4cf8247390486141a0b2e44d0b4fea2c70d61711589ebc47330")
	// fmt.Println(recp, err)

	// // 合约交互 查询 call 方法
	ballotCallVotersTelChainEVM(t)

	// // 0x00000000000000007366e65bb39229c015c595850c16948d72f940303c8f3e11
	// res, _ := hexutil.Decode("0x7366e65bb39229c015c595850c16948d72f940303c8f3e11")
	// fmt.Println(res)

	// 合约交互 交易 send 方法
	// ballotSendTelChainEVM(t)

}

func ballotDeployTelChainEVM(t *testing.T) string {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	var contract Contract

	contract.Abi = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"proposalNames\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"chairperson\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"delegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"giveRightToVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposal\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"voters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"voted\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"vote\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"winnerName\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"winnerName_\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"winningProposal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"winningProposal_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	contract.ByteCode = "608060405234801561001057600080fd5b5060405161088a38038061088a8339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825186602082028301116401000000008211171561008557600080fd5b82525081516020918201928201910280838360005b838110156100b257818101518382015260200161009a565b50505050919091016040908152600080546001600160c01b03191633178082556001600160c01b03168152600160208190529181209190915593505050505b8151811015610153576002604051806040016040528084848151811061011357fe5b60209081029190910181015182526000918101829052835460018181018655948352918190208351600290930201918255919091015190820155016100f1565b5050610726806101646000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063609ff1bd1161005b578063609ff1bd1461012c5780639e7b8d6114610146578063a3ec138d1461016c578063e2ba53f0146101c057610088565b80630121b93f1461008d578063013cf08b146100ac5780632e4176cf146100e25780635c19a95c14610106575b600080fd5b6100aa600480360360208110156100a357600080fd5b50356101c8565b005b6100c9600480360360208110156100c257600080fd5b50356102b3565b6040805192835260208301919091528051918290030190f35b6100ea6102e1565b604080516001600160c01b039092168252519081900360200190f35b6100aa6004803603602081101561011c57600080fd5b50356001600160c01b03166102f0565b6101346104df565b60408051918252519081900360200190f35b6100aa6004803603602081101561015c57600080fd5b50356001600160c01b0316610546565b6101926004803603602081101561018257600080fd5b50356001600160c01b0316610641565b6040805194855292151560208501526001600160c01b03909116838301526060830152519081900360800190f35b610134610675565b3360009081526001602052604090208054610221576040805162461bcd60e51b8152602060048201526014602482015273486173206e6f20726967687420746f20766f746560601b604482015290519081900360640190fd5b600181015460ff161561026c576040805162461bcd60e51b815260206004820152600e60248201526d20b63932b0b23c903b37ba32b21760911b604482015290519081900360640190fd5b6001818101805460ff191690911790556002808201839055815481549091908490811061029557fe5b60009182526020909120600160029092020101805490910190555050565b600281815481106102c357600080fd5b60009182526020909120600290910201805460019091015490915082565b6000546001600160c01b031681565b3360009081526001602081905260409091209081015460ff1615610350576040805162461bcd60e51b81526020600482015260126024820152712cb7ba9030b63932b0b23c903b37ba32b21760711b604482015290519081900360640190fd5b6001600160c01b0382163314156103ae576040805162461bcd60e51b815260206004820152601e60248201527f53656c662d64656c65676174696f6e20697320646973616c6c6f7765642e0000604482015290519081900360640190fd5b6001600160c01b038281166000908152600160208190526040909120015461010090041615610458576001600160c01b039182166000908152600160208190526040909120015461010090049091169033821415610453576040805162461bcd60e51b815260206004820152601960248201527f466f756e64206c6f6f7020696e2064656c65676174696f6e2e00000000000000604482015290519081900360640190fd5b6103ae565b6001818101805460ff19168217610100600160c81b0319166101006001600160c01b0386169081029190911790915560009081526020829052604090209081015460ff16156104d2578154600282810154815481106104b357fe5b60009182526020909120600160029092020101805490910190556104da565b815481540181555b505050565b600080805b6002548110156105415781600282815481106104fc57fe5b9060005260206000209060020201600101541115610539576002818154811061052157fe5b90600052602060002090600202016001015491508092505b6001016104e4565b505090565b6000546001600160c01b0316331461058f5760405162461bcd60e51b81526004018080602001828103825260288152602001806106a36028913960400191505060405180910390fd5b6001600160c01b0381166000908152600160208190526040909120015460ff1615610601576040805162461bcd60e51b815260206004820152601860248201527f54686520766f74657220616c726561647920766f7465642e0000000000000000604482015290519081900360640190fd5b6001600160c01b0381166000908152600160205260409020541561062457600080fd5b6001600160c01b0316600090815260016020819052604090912055565b600160208190526000918252604090912080549181015460029091015460ff82169161010090046001600160c01b03169084565b600060026106816104df565b8154811061068b57fe5b90600052602060002090600202016000015490509056fe4f6e6c79206368616972706572736f6e2063616e206769766520726967687420746f20766f74652ea2646970667358221220359927aa0db07eccbb352dca8aa7dc06bce5775ed527cc251219bfaacfb1b62d64736f6c637828302e372e362d646576656c6f702e323032302e31322e31382b636f6d6d69742e37333166666535340059"

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(singAddr)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  2000000,
		Sender:    &sender,
		Recipient: nil,
		Amount:    nil,
	}

	var key1 [32]byte
	var key2 [32]byte
	copy(key1[:], []byte("test1"))
	copy(key1[:], []byte("test2"))
	proposalNames := []utils.Hash{key1, key2}

	hash, err := contractObj.Deploy(tx, false, singAddrPriKey, contract.ByteCode, proposalNames)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("txHash is ", hash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(hash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Contract Address: ", receipt.ContractAddress)

	return hash
}

func ballotCallVotersTelChainEVM(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	transaction := new(dto.TransactionParameters)
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contractAddress := "did:bid:qwer:sfWkwvBdBaeAmApBJ6hn82DpAb9gc3M4"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(singAddr, block.LATEST)
	transaction.ChainId = chainId
	transaction.Sender = generator
	transaction.Recipient = contractAddress
	transaction.AccountNonce = nonce.Uint64()

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	voter := utils.StringToAddress("did:bid:qwer:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n")

	result, err := contractObj.Call(transaction, "voters", voter)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result != nil {
		val, _ := result.ToComplexString()
		fmt.Println(val)
	}
}

// 合约交互 交易 send 方法
func ballotSendTelChainEVM(t *testing.T) {
	var singAddr = "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	var singAddrPriKey = "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	contractAddress := "did:bid:qwer:sfWkwvBdBaeAmApBJ6hn82DpAb9gc3M4"
	chainId, _ := connection.Core.GetChainId()
	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block.LATEST)

	sender := utils.StringToAddress(singAddr)
	recipient := utils.StringToAddress(contractAddress)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  2000000,
		Sender:    &sender,
		Recipient: &recipient,
		Amount:    nil,
	}

	// 创建合约实例
	contract := ballotContract()
	contractObj, err := connection.Core.NewContract(contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	voter := utils.StringToAddress("did:bid:qwer:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n")
	// voter := utils.StringToAddress("did:bid:qwer:sfCXQusR8SEWgp8fQ9BQu61riWdDLCMN")
	// voter := utils.StringToAddress("did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY")

	txHash, err := contractObj.Send(tx, false, singAddrPriKey, "giveRightToVote", voter)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("txHash is ", txHash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txHash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}