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

// bif ���ļ�
package bif

import (
	"github.com/bif/bif-sdk-go/core"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/personal"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/utils"
)

// Web3 - The Web3 Module
// Web3 - Web3ģ��
type Bif struct {
	Provider providers.ProviderInterface
	Core     *core.Core
	Net      *net.Net
	Personal *personal.Personal
	Utils    *utils.Utils
	System   *system.System
}

/*
	NewBif - Web3 Module constructor to set the default provider, Core, Net and Personal

	NewBif - Web3ģ�鹹�캯������������Ĭ��Provider��Core��Net��Personal
*/
func NewBif(provider providers.ProviderInterface) *Bif {
	bif := new(Bif)
	bif.Provider = provider
	bif.Core = core.NewEth(provider)
	bif.Net = net.NewNet(provider)
	bif.Personal = personal.NewPersonal(provider)
	bif.Utils = utils.NewUtils()
	bif.System = system.NewSystem(provider)
	return bif
}

// ClientVersion - Returns the current client version.
//
// ClientVersion - ���ص�ǰ�ͻ��˰汾.
//
// Returns
//	- String - The current client version
func (bif Bif) ClientVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := bif.Provider.SendRequest(pointer, "bif_clientVersion", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
