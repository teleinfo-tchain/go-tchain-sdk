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
	"github.com/tchain/go-tchain-sdk/account/types"
	"github.com/tchain/go-tchain-sdk/core"
	"github.com/tchain/go-tchain-sdk/debug"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/gb"
	"github.com/tchain/go-tchain-sdk/net"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/system"
	"github.com/tchain/go-tchain-sdk/test/privacycontract"
	"github.com/tchain/go-tchain-sdk/txpool"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

/*
	Bif - The Bif Module

	Bif - Bif模块
*/
type Bif struct {
	Provider        providers.ProviderInterface
	Core            *core.Core
	Gb              *gb.GB
	Account         *types.Account
	Net             *net.Net
	System          *system.System
	Debug           *debug.Debug
	TxPool          *txpool.TxPool
	PrivacyContract *privacycontract.BaseContract
}

/*
	NewBif - Web3 Module constructor to set the default provider, Core, Net and Personal

	NewBif - Web3模块构造函数，用于设置默认Provider、Core、Net和Personal
*/
func NewBif(provider providers.ProviderInterface) *Bif {
	bif := new(Bif)
	bif.Provider = provider
	bif.Core = core.NewCore(provider)
	bif.Gb = gb.NewGB(provider)
	bif.Net = net.NewNet(provider)
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

// GetCurrentAbPath 最终方案-全兼容  获取项目路径
func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()

	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
