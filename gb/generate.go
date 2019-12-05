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
 * @file net.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package gb

import (
	"github.com/bif/bifGo/dto"
	"github.com/bif/bifGo/providers"
)

// Net - The Net Module
type GB struct {
	provider providers.ProviderInterface
}

// NewNet - Net Module constructor to set the default provider
func NewGB(provider providers.ProviderInterface) *GB {
	gb := new(GB)
	gb.provider = provider
	return gb
}

func (gb *GB) Start() error {

	pointer := &dto.RequestResult{}

	err := gb.provider.SendRequest(pointer, "gb_start", nil)

	return err
}

func (gb *GB) Stop() error {

	pointer := &dto.RequestResult{}

	err := gb.provider.SendRequest(pointer, "gb_stop", nil)

	return err
}
