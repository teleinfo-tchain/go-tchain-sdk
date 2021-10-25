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
	"github.com/bif/bif-sdk-go/account"
	block2 "github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/txpool"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
	"strconv"
	"testing"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/providers"
)

func Test_Websocket_Core_ToSyncingResponse(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

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

func Test_Websocket_Core_ToTransactionResponse(t *testing.T) {
	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))
	chainId, _ := connection.Core.GetChainId()

	var sender utils.Address
	sender = utils.StringToAddress(resources.Addr1)
	var recipient utils.Address
	recipient = utils.StringToAddress(resources.Addr2)

	tx := &account.SignTxParams{
		Sender:    &sender,
		Recipient: &recipient,
		Amount:    big.NewInt(10000),
		GasLimit:  200000,
		GasPrice:  big.NewInt(30),
		ChainId:   chainId,
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

	txFromHash, err := connection.Core.GetTransactionByHash(txHash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(txFromHash)
}

func Test_Websocket_Core_ToSignTransactionResponse(t *testing.T) {
	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	transaction := new(dto.TransactionParameters)
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

func Test_Websocket_Core_ToPendingTransactions(t *testing.T) {
	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	pendingTransactions, err := connection.Core.GetPendingTransactions()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(pendingTransactions)
}

func Test_Websocket_Core_ToBlock(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	block, err := connection.Core.GetBlockByNumber(block2.LATEST, false)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(block)

	err = connection.Provider.Close()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_Websocket_Core_ToProof(t *testing.T) {
	// todo
	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	generator, _ := connection.Core.GetGenerator()

	_, err := connection.Core.GetProof(generator, []string{"0", "1"}, block2.LATEST)

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

func Test_Websocket_Debug_ToDumpBlock(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

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

func Test_Websocket_Net_ToPeerInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
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

func Test_Websocket_Net_ToNodeInfo(t *testing.T) {
	// todo
	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
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

func Test_Websocket_Request_ToStringArray(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
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

func Test_Websocket_Request_ToUint64(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
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

func Test_Websocket_Request_ToBoolean(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
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

func Test_Websocket_Txpool_ToTxPoolStatus(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	status, err := connection.GetStatus()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(status)
}

func Test_Websocket_Txpool_ToTxPoolInspect(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	inspect, err := connection.Inspect()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(inspect)
}

func Test_Websocket_Txpool_ToTxPoolContent(t *testing.T) {

	var connection = txpool.NewTxPool(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	content, err := connection.Content()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(content)
}

func Test_Websocket_System_ToCertificateInfo(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	certificate, err := cer.GetCertificate("did:bid:sf26xSMerRVcs9E642Rkbq4TxLrhQsWzk")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(certificate)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_Websocket_System_ToCertificateIssuerSignature(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	issuer, err := cer.GetIssuer("did:bid:sf26xSMerRVcs9E642Rkbq4TxLrhQsWzk")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(issuer)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_Websocket_System_ToCertificateSubjectSignature(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cer := connection.System.NewCertificate()

	subjectSignature, err := cer.GetSubject("did:bid:sf26xSMerRVcs9E642Rkbq4TxLrhQsWzk")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(subjectSignature)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_Websocket_System_ToDocument(t *testing.T) {

	var connection = bif.NewBif(providers.NewWebSocketProvider("ws://" + resources.IP00 + ":" + strconv.FormatUint(resources.WebsocketPort, 10)))

	for index := 0; index < 2; index++ {
		var _, err = connection.ClientVersion()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	_, err := connection.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	document, err := doc.GetDocument("did:bid:sf26xSMerRVcs9E642Rkbq4TxLrhQsWzk")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(document)

	err = connection.Provider.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
