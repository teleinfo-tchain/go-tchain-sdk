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

/**
 * @file utils.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package utils

import (
	"github.com/bif/bif-sdk-go/complex/types"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
)

// Utils - The Utils Module
type Utils struct {
	provider providers.ProviderInterface
}

// NewUtils - Utils Module constructor to set the default provider
func NewUtils(provider providers.ProviderInterface) *Utils {
	utils := new(Utils)
	utils.provider = provider
	return utils
}

// Sha3 - Returns Keccak-256 (not the standardized SHA3-256) of the given data.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_sha3
//    - DATA - the data to convert into a SHA3 hash
// Returns:
// 	  - DATA - The SHA3 result of the given string.
func (utils *Utils) Sha3(data types.ComplexString) (string, error) {

	params := make([]string, 1)
	params[0] = data.ToHex()

	pointer := &dto.RequestResult{}

	err := utils.provider.SendRequest(pointer, "bif_sha3", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
