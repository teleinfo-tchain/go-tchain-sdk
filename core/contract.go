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

package core

import (
	Abi "github.com/tchain/go-tchain-sdk/abi"
	"github.com/tchain/go-tchain-sdk/account"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/utils"
	"github.com/tchain/go-tchain-sdk/utils/types"
	"strings"
)

// Contract ...
type Contract struct {
	super *Core
	abi   Abi.ABI
}

// NewContract - Contract abstraction
func (core *Core) NewContract(abi string) (*Contract, error) {
	parsedAbi, err := Abi.JSON(strings.NewReader(abi))
	if err != nil {
		return nil, err
	}

	contract := new(Contract)
	contract.abi = parsedAbi
	contract.super = core
	return contract, nil
}

func (contract *Contract) Call(transaction *dto.TransactionParameters, functionName string, args ...interface{}) (*dto.RequestResult, error) {
	inputEncode, err := contract.abi.Pack(functionName, args...)
	if err != nil {
		return nil, err
	}
	transaction.Payload = types.ComplexString("0x" + utils.Bytes2Hex(inputEncode))

	return contract.super.Call(transaction)

}

func (contract *Contract) Send(tx *account.SignTxParams, isSM2 bool, signPriKey, functionName string, args ...interface{}) (string, error) {
	inputEncode, err := contract.abi.Pack(functionName, args...)
	if err != nil {
		return "", err
	}

	tx.Payload = inputEncode

	signTx, err := account.SignTransaction(tx, signPriKey, isSM2)
	if err != nil {
		return "", err
	}

	return contract.super.SendRawTransaction(signTx.Raw.String())

}

func (contract *Contract) Deploy(tx *account.SignTxParams, isSM2 bool, signPriKey, byteCode string, args ...interface{}) (string, error) {
	inputEncode, err := contract.abi.Pack("", args...)
	if err != nil {
		return "", err
	}

	tx.Payload = append(utils.Hex2Bytes(byteCode), inputEncode...)

	signTx, err := account.SignTransaction(tx, signPriKey, isSM2)
	if err != nil {
		return "", err
	}

	return contract.super.SendRawTransaction(signTx.Raw.String())
}

