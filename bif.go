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

// bif 主文件
package bif

import (
	"github.com/bif/bif-sdk-go/core"
	"github.com/bif/bif-sdk-go/debug"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/txpool"
	"github.com/bif/bif-sdk-go/utils"
)

/*
	Bif - The Bif Module

	Bif - Bif模块
*/
type Bif struct {
	Provider providers.ProviderInterface
	Core     *core.Core
	Net      *net.Net
	Utils    *utils.Utils
	System   *system.System
	Debug    *debug.Debug
	TxPool   *txpool.TxPool
}

/*
	NewBif - Web3 Module constructor to set the default provider, Core, Net and Personal

	NewBif - Web3模块构造函数，用于设置默认Provider、Core、Net和Personal
*/
func NewBif(provider providers.ProviderInterface) *Bif {
	bif := new(Bif)
	bif.Provider = provider
	bif.Core = core.NewCore(provider)
	bif.Net = net.NewNet(provider)
	bif.Utils = utils.NewUtils()
	bif.System = system.NewSystem(provider)
	bif.TxPool = txpool.NewTxPool(provider)
	bif.Debug = debug.NewDebug(provider)
	return bif
}

/*
  ClientVersion:
	EN - Returns the current client version.
	CN - 返回当前客户端版本
  Params:
  	- None
  Returns:
  	- string - 当前客户端版本
    - error

  Call permissions: Anyone
*/
func (bif Bif) ClientVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := bif.Provider.SendRequest(pointer, "bif_clientVersion", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
