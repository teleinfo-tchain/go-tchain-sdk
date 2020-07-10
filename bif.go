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
	"fmt"
	"github.com/bif/bif-sdk-go/core"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/net"
	"github.com/bif/bif-sdk-go/personal"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Coin - Ethereum value unity value
const (
	Coin float64 = 1000000000000000000
)

//profile variables
type conf struct {
	Version  string            `yaml:"version"`
	Core     map[string]string `yaml:"core"`
	Gb       map[string]string `yaml:"gb"`
	Net      map[string]string `yaml:"net"`
	Admin    map[string]string `yaml:"admin"`
	Personal map[string]string `yaml:"personal"`
	Txpool   map[string]string `yaml:"txpool"`
	Bif      map[string]string `yaml:"bif"`
}

var ConfigPathName = "rpc-method-config.yaml"

//var yamlPath string
//flag.StringVar(&yamlPath, "yaml-path", "yaml-path", "config rpc yaml path")
//flag.Parse()
//fmt.Printf("yaml path: %s", yamlPath)

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile(ConfigPathName)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

//var RpcCof conf
//func init() {
//	RpcCof.getConf()
//}

// Web3 - The Web3 Module
type Bif struct {
	Provider providers.ProviderInterface
	Core     *core.Core
	Net      *net.Net
	Personal *personal.Personal
	Utils    *utils.Utils
	System   *system.System
}

// NewBif - Web3 Module constructor to set the default provider, Core, Net and Personal
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

//var instance *bif
//var once sync.Once
//// NewBif - Bif Module constructor to set the default provider, Core, Net and Personal
//// Singleton pattern
//func NewBif(provider providers.ProviderInterface) *bif {
//	once.Do(func() {
//		instance = new(bif)
//		instance.Provider = provider
//		instance.Core = core.NewEth(provider)
//		instance.Net = net.NewNet(provider)
//		instance.Personal = personal.NewPersonal(provider)
//		instance.Utils = utils.NewUtils()
//	})
//	return instance
//}

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
