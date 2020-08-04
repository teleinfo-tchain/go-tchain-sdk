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
	Abi "github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/utils/types"
	"strings"
)

// Contract ...
type Contract struct {
	super *Core
	abi   Abi.ABI
}

// NewContract - Contract abstraction
func (core *Core) NewContract(abi string) (*Contract, error) {
	parsedAbi, _ := Abi.JSON(strings.NewReader(abi))

	contract := new(Contract)
	contract.abi = parsedAbi
	contract.super = core
	return contract, nil
}

// prepareTransaction ...
func (contract *Contract) prepareTransaction(transaction *dto.TransactionParameters, functionName string, args []interface{}) (*dto.TransactionParameters, error) {
	inputEncode, err := contract.abi.Pack(functionName, args...)
	if err != nil {
		return nil, err
	}
	transaction.Data = types.ComplexString(inputEncode)
	return transaction, nil

}

func (contract *Contract) Call(transaction *dto.TransactionParameters, functionName string, args ...interface{}) (*dto.RequestResult, error) {
	transaction, err := contract.prepareTransaction(transaction, functionName, args)

	if err != nil {
		return nil, err
	}

	return contract.super.Call(transaction)

}

func (contract *Contract) Send(transaction *dto.TransactionParameters, functionName string, args ...interface{}) (string, error) {

	transaction, err := contract.prepareTransaction(transaction, functionName, args)

	if err != nil {
		return "", err
	}

	return contract.super.SendTransaction(transaction)

}

func (contract *Contract) Deploy(transaction *dto.TransactionParameters, byteCode string, args ...interface{}) (string, error) {
	inputEncode, err := contract.abi.Pack("", args...)
	if err != nil {
		return "", err
	}
	transaction.Data = types.ComplexString(byteCode) + types.ComplexString(inputEncode)
	return contract.super.SendTransaction(transaction)

}

