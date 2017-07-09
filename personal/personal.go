/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

/**
 * @file personal.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package personal

import (
	"github.com/regcostajr/go-web3/complex/types"
	"github.com/regcostajr/go-web3/dto"
	"github.com/regcostajr/go-web3/providers"
)

// Personal - The Personal Module
type Personal struct {
	provider providers.ProviderInterface
}

// NewPersonal - Personal Module constructor to set the default provider
func NewPersonal(provider providers.ProviderInterface) *Personal {
	personal := new(Personal)
	personal.provider = provider
	return personal
}

// ListAccounts - Lists all stored accounts.
// Reference: https://github.com/paritytech/parity/wiki/JSONRPC-personal-module#personal_listaccounts
// Parameters:
//    - none
// Returns:
//    - Array - A list of 20 byte account identifiers.
func (personal *Personal) ListAccounts() ([]types.Address, error) {

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_listAccounts", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToAddressArray()

}

// NewAccount - Creates new account.
// Note: it becomes the new current unlocked account. There can only be one unlocked account at a time.
// Reference: https://github.com/paritytech/parity/wiki/JSONRPC-personal-module#personal_newaccount
// Parameters:
//    - String - Password for the new account.
// Returns:
//	  - Address - 20 Bytes - The identifier of the new account.
func (personal *Personal) NewAccount(password string) (types.Address, error) {

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(&pointer, "personal_newAccount", nil)

	if err != nil {
		return "", err
	}

	response, err := pointer.ToString()

	return types.Address(response), err

}

// SendTransaction - Sends transaction and signs it in a single call. The account does not need to be unlocked to make this call, and will not be left unlocked after.
// Reference: https://github.com/paritytech/parity/wiki/JSONRPC-personal-module#personal_sendtransaction
// Parameters:
//    1. Object - The transaction object
//        - from: Address - 20 Bytes - The address the transaction is send from.
//        - to: Address - (optional) 20 Bytes - The address the transaction is directed to.
//        - gas: Quantity - (optional) Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.
//        - gasPrice: Quantity - (optional) Integer of the gas price used for each paid gas.
//        - value: Quantity - (optional) Integer of the value sent with this transaction.
//        - data: Data - (optional) 4 byte hash of the method signature followed by encoded parameters. For details see Ethereum Contract ABI.
//        - nonce: Quantity - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
//        - condition: Object - (optional) Conditional submission of the transaction. Can be either an integer block number { block: 1 } or UTC timestamp (in seconds) { time: 1491290692 } or null.
//    2. String - Passphrase to unlock the from account.
// Returns:
//    - Data - 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available
func (personal *Personal) SendTransaction(from types.Address, to types.Address, value types.ComplexIntParameter, hexData types.ComplexString, password string) (string, error) {

	params := make([]interface{}, 1)

	transactionParameters := &dto.TransactionParameters{}
	transactionParameters.From = from
	transactionParameters.To = to
	transactionParameters.Value = value.ToHex()
	transactionParameters.Data = hexData.ToHex()

	params[0] = transactionParameters
	params[1] = password

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

// UnlockAccount - Unlocks specified account for use.
// Reference: https://github.com/paritytech/parity/wiki/JSONRPC-personal-module#personal_unlockaccount
// If permanent unlocking is disabled (the default) then the duration argument will be ignored,
// and the account will be unlocked for a single signing. With permanent locking enabled, the
// duration sets the number of seconds to hold the account open for. It will default to 300 seconds.
// Passing 0 unlocks the account indefinitely.
// There can only be one unlocked account at a time.
// Parameters:
//    - Address - 20 Bytes - The address of the account to unlock.
//    - String - Passphrase to unlock the account.
//    - Quantity - (default: 300) Integer or null - Duration in seconds how long the account should remain unlocked for.
// Returns:
// 	   - Boolean - whether the call was successful
func (personal *Personal) UnlockAccount(address types.Address, password string, duration types.ComplexIntParameter) (bool, error) {

	params := make([]string, 3)
	params[0] = string(address)
	params[1] = password
	params[2] = duration.ToHex()

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_unlockAccount", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}
