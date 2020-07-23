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

// 可让您与bif节点的帐户进行交互。
package personal

import (
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
)

// Personal - The Personal Module
// Personal - Personal 模块
type Personal struct {
	provider providers.ProviderInterface
}

// NewPersonal - Personal Module constructor to set the default provider
// NewPersonal - 初始化Personal
func NewPersonal(provider providers.ProviderInterface) *Personal {
	personal := new(Personal)
	personal.provider = provider
	return personal
}

 /*
  ListAccounts:
   	EN - Return a list of addresses for accounts this node manages
 	CN -  返回此节点管理的帐户的地址列表

   Params:
  	- None

  Returns:
  	- []string, Array,A list of 20 byte account identifiers.
 	- error
 
  Call permissions: Anyone
  */
func (personal *Personal) ListAccounts() ([]string, error) {

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_listAccounts", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()

}

 /*
  NewAccount:
   	EN - Creates new account and returns the address for the new account
 	CN - 创建一个新账户,并返回新账户地址
  Params:
  	- password， string, 新账户的密码
 	- isSm2, true 使用国密的加密方式生成账户；false使用非国密的方式生成账户

  Returns:
  	- string, address , 20 Bytes, 新账户
 	- error

  Call permissions: Anyone， 切勿通过不安全的Websocket或HTTP提供程序调用此函数，因为您的密码将以纯文本形式发送
  */
func (personal *Personal) NewAccount(password string, isSm2 bool) (string, error) {
	params := make([]interface{}, 2)
	params[0] = password
	if isSm2{
		params[1] = uint64(0)
	}else {
		params[1] = uint64(1)
	}

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(&pointer, "personal_newAccount", params)

	if err != nil {
		return "", err
	}

	response, err := pointer.ToString()

	return response, err

}

 /*
   SendTransaction:
	EN - Create a transaction from the given arguments and tries to sign it with the key associated with args.From
  	CN - 对给定参数创建交易，并使用与交易发起方（From）的相关密钥进行签名
   Params:
   	- transaction: 要发送的交易对象(*dto.TransactionParameters)
 		from: string，20 Bytes - 指定的发送者的地址。
 		to: string，20 Bytes - （可选）交易消息的目标地址，如果是合约创建，则不填.
 		gas: *big.Int - （可选）默认是自动，交易可使用的gas，未使用的gas会退回。
 		gasPrice: *big.Int - （可选）默认是自动确定，交易的gas价格，默认是网络gas价格的平均值 。
 		data: string - （可选）或者包含相关数据的字节字符串，如果是合约创建，则是初始化要用到的代码。
 		value: *big.Int - （可选）交易携带的货币量，以bifer为单位。如果合约创建交易，则为初始的基金
 		nonce: *big.Int - （可选）整数，使用此值，可以允许你覆盖你自己的相同nonce的，待pending中的交易
 	- password, string, 当前帐户的密码来解锁账户

   Returns:
   	- string, transactionHash，32 Bytes，交易哈希，如果交易尚不可用，则为零哈希
  	- error

   Call permissions: Anyone， 注意，通过不安全的HTTP RPC连接发送帐户密码非常不安全。
  */
func (personal *Personal) SendTransaction(transaction *dto.TransactionParameters, password string) (string, error) {

	params := make([]interface{}, 2)

	transactionParameters := transaction.Transform()

	params[0] = transactionParameters
	params[1] = password

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_sendTransaction", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  UnlockAccount:
   	EN - Unlock the account associated with the given address and the given password
 	CN - 解锁与给定地址和给定密码关联的帐户
  Params:
  	- address,string 20 Bytes. 要解锁的账户地址
   	- password,string ,解锁该账户的密码
   	- duration,uint64 ,帐户保持解锁状态的持续时间（？？？？如果是0，就是永久解锁？？）

  Returns:
  	- bool, true是解锁成功， false解锁失败
 	- error

  Call permissions: Anyone
*/
func (personal *Personal) UnlockAccount(address string, password string, duration uint64) (bool, error) {

	params := make([]interface{}, 3)
	params[0] = address
	params[1] = password
	params[2] = duration

	pointer := &dto.RequestResult{}

	err := personal.provider.SendRequest(pointer, "personal_unlockAccount", params)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}
