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

package test

import (
	"fmt"
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/account"
	block2 "github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/test/resources"
	"github.com/tchain/go-tchain-sdk/txpool"
	"github.com/tchain/go-tchain-sdk/utils"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func Test_HTTP_Core_ToSyncingResponse(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	_, err = connection.Core.IsSyncing()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = connection.Provider.Close()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Core_ToTransactionResponse(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount(resources.Addr1, block2.LATEST)

	var sender utils.Address
	sender = utils.StringToAddress(resources.Addr1)
	var recipient utils.Address
	recipient = utils.StringToAddress(resources.Addr2)

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(200),
		GasLimit:  200000,
		Sender:    &sender,
		Recipient: &recipient,
		Amount:    big.NewInt(10000),
	}

	res, err := account.SignTransaction(tx, resources.Addr1Pri, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txHash, err := connection.Core.SendRawTransaction(res.Raw.String())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(time.Second * 8)

	txResponse, err := connection.Core.GetTransactionByHash(txHash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Printf("txRes is %+v \n", txResponse)
}

func Test_HTTP_Core_ToSignTransactionResponse(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	chainId, _ := connection.Core.GetChainId()

	transaction := new(dto.TransactionParameters)
	transaction.ChainId = chainId
	transaction.AccountNonce = uint64(2)
	transaction.Sender = resources.Addr1
	transaction.Recipient = resources.Addr2
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(5), big.NewInt(1e17))
	transaction.GasLimit = uint64(50000)
	transaction.GasPrice = big.NewInt(1)
	transaction.Payload = "Sign Transfer bif test"

	txID, err := connection.Core.SignTransaction(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(txID)
}

func Test_HTTP_Core_ToPendingTransactions(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	pendingTransactions, err := connection.Core.GetPendingTransactions()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(pendingTransactions)
}

func Test_HTTP_Core_ToBlock(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	_, err = connection.Core.GetBlockByNumber(block2.LATEST, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = connection.Provider.Close()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Core_ToProof(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	generator, _ := connection.Core.GetGenerator()

	proof, err := connection.Core.GetProof(generator, []string{"0", "1"}, block2.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(proof)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Debug_ToDumpBlock(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	dumpBlock, err := connection.Debug.DumpBlock("latest")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(dumpBlock)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Net_ToPeerInfo(t *testing.T) {
	// todo
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	peers, err := connection.Net.GetPeers()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(peers)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Net_ToNodeInfo(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	nodeInfo, err := connection.Net.GetNodeInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(nodeInfo)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Request_ToStringArray(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	accounts, err := connection.Core.GetAccounts()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(accounts)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Request_ToUint64(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	chainId, err := connection.Core.GetChainId()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(chainId)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Request_ToBigInt(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	generator, _ := connection.Core.GetGenerator()
	nonce, err := connection.Core.GetTransactionCount(generator, block2.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(nonce)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Request_ToBoolean(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	var _, err = connection.ClientVersion()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	generator, err := connection.Core.GetGenerator()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(generator)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_HTTP_Request_ToTransactionReceipt(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	from := "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	fromPriKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	to := "did:bid:qwer:zftAgNtnQzLMGJHKPMdn9quPvuikNWUZ"

	chainId, _ := connection.Core.GetChainId()

	// sender := utils.StringToAddress(from)
	var recipient utils.Address
	recipient = utils.StringToAddress(to)

	nonce, err := connection.Core.GetTransactionCount(from, block2.LATEST)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tx := &account.SignTxParams{
		ChainId:   chainId,
		Nonce:     nonce,
		GasPrice:  big.NewInt(0),
		GasLimit:  21000,
		Recipient: &recipient,
		Amount:    big.NewInt(500000000),
		Payload:   nil,
	}

	res, err := account.SignTransaction(tx, fromPriKey, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txHash, err := connection.Core.SendRawTransaction(res.Raw.String())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(time.Second * 8)

	txReceipt, err := connection.Core.GetTransactionReceipt(txHash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(txReceipt)
}

func Test_HTTP_Txpool_ToTxPoolStatus(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	status, err := connection.GetStatus()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(status)
}

func Test_HTTP_Txpool_ToTxPoolInspect(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	inspect, err := connection.Inspect()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(inspect)
}

func Test_HTTP_Txpool_ToTxPoolContent(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	content, err := connection.Content()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(content)
}

