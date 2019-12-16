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

package txpool

import (
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/go-bif/common/hexutil"
)

type Txpool struct {
	provider providers.ProviderInterface
}

func NewTxpool(provider providers.ProviderInterface) *Txpool {
	txpool := new(Txpool)
	txpool.provider = provider
	return txpool
}

func (txpool *Txpool) Status() (map[string]hexutil.Uint, error) {

	pointer := &dto.RequestResult{}

	if err := txpool.provider.SendRequest(pointer, "txpool_status", nil); err != nil {
		return nil , err
	}

	result, err := pointer.ToTxpoolStatus()
	if err != nil {
		return nil , err
	}

	return result, nil
}

func (txpool *Txpool) Inspect() (map[string]map[string]map[string]string ,error) {

	pointer := &dto.RequestResult{}

	if err := txpool.provider.SendRequest(pointer, "txpool_inspect", nil); err != nil {
		return nil, err
	}

	result, err := pointer.ToTxpoolInspect()
	if err != nil {
		return nil , err
	}

	return result, nil
}
