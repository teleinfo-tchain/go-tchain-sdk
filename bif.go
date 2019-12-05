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
 * @file bif.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package bif

import (
	"github.com/bif/bifGo/core"
	"github.com/bif/bifGo/dto"
	"github.com/bif/bifGo/net"
	"github.com/bif/bifGo/personal"
	"github.com/bif/bifGo/providers"
	"github.com/bif/bifGo/utils"
)

// Coin - Ethereum value unity value
const (
	Coin float64 = 1000000000000000000
)

// Web3 - The Web3 Module
type Bif struct {
	Provider providers.ProviderInterface
	Core     *core.Core
	Net      *net.Net
	Personal *personal.Personal
	Utils    *utils.Utils
}

// NewBif - Web3 Module constructor to set the default provider, Core, Net and Personal
func NewBif(provider providers.ProviderInterface) *Bif {
	bif := new(Bif)
	bif.Provider = provider
	bif.Core = core.NewEth(provider)
	bif.Net = net.NewNet(provider)
	bif.Personal = personal.NewPersonal(provider)
	bif.Utils = utils.NewUtils(provider)
	return bif
}

// ClientVersion - Returns the current client version.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion
// Parameters:
//    - none
// Returns:
// 	  - String - The current client version
func (bif Bif) ClientVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := bif.Provider.SendRequest(pointer, "bif_clientVersion", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
