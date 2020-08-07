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

package txpool

import (
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils/hexutil"
)

type TxPool struct {
	provider providers.ProviderInterface
}

func NewTxPool(provider providers.ProviderInterface) *TxPool {
	txPool := new(TxPool)
	txPool.provider = provider
	return txPool
}

/*
  GetStatus:
   	EN - Returns the number of pending and queued transaction in the pool
 	CN - 返回交易池中暂挂和排队的交易数
  Params:
  	- None

  Returns:
  	- map[string]hexutil.Uint
 	- error

  Call permissions: Anyone
*/
func (txPool *TxPool) GetStatus() (map[string]hexutil.Uint, error) {

	pointer := &dto.TxPoolRequestResult{}

	if err := txPool.provider.SendRequest(pointer, "txpool_status", nil); err != nil {
		return nil, err
	}

	return pointer.ToTxPoolStatus()
}

/*
  Inspect:
   	EN - Retrieves the content of the transaction pool and flattens it into an easily inspectable list
 	CN - 检索交易池的内容并将其转化为易于检查的列表
  Params:
  	- None

  Returns:
  	- map[string]map[string]map[string]string
 	- error

  Call permissions: Anyone
*/
func (txPool *TxPool) Inspect() (map[string]map[string]map[string]string, error) {

	pointer := &dto.TxPoolRequestResult{}

	if err := txPool.provider.SendRequest(pointer, "txpool_inspect", nil); err != nil {
		return nil, err
	}

	return pointer.ToTxPoolInspect()
}

/*
  Content:
   	EN - Returns the transactions contained within the transaction pool
 	CN - 返回交易池中包含的交易
  Params:
  	-

  Returns:
  	- None
 	- error

  Call permissions: Anyone
*/
func (txPool *TxPool) Content() (map[string]map[string]map[string]*dto.RPCTransaction, error) {

	pointer := &dto.TxPoolRequestResult{}

	if err := txPool.provider.SendRequest(pointer, "txpool_content", nil); err != nil {
		return nil, err
	}

	return pointer.ToTxPoolContent()

}
