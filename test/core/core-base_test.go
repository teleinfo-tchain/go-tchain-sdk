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
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"strconv"
	"testing"
)

func TestCoreEstimateGas(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generator, err := connection.Core.GetGenerator()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.Payload = "test"
	transaction.Sender = generator
	transaction.Recipient = generator
	transaction.Amount = big.NewInt(10)
	transaction.GasLimit = uint64(40000)

	gas, err := connection.Core.EstimateGas(transaction)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(gas)

}

func TestCoreGasPrice(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	gasPrice, err := connection.Core.GetGasPrice()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if gasPrice.Int64() < 0 {
		t.Errorf("Negative gasprice")
		t.FailNow()
	}

	t.Log(gasPrice.Int64())
}

func TestCoreGenerating(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generating, err := connection.Core.Generating()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(generating)

}

func TestGetAccounts(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	accounts, err := connection.Core.GetAccounts()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(accounts)
}

func TestCoreGetBalance(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generator, _ := connection.Core.GetGenerator()

	bal, err := connection.Core.GetBalance(generator, block.PENDING)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(bal)

}

func TestGetChainId(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	chainId, err := connection.Core.GetChainId()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log(chainId)
}




func TestGetGenerator(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generator, err := connection.Core.GetGenerator()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(generator)
}

func TestCoreHashRate(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	rate, err := connection.Core.GetHashRate()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if rate.Uint64() < 0 {
		t.Errorf("Less than 0 hash rate")
		t.FailNow()
	}

	t.Log(rate)
}
func TestCoreGetTrustNumber(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	address, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	trustNumber, err := connection.Core.GetTrustNumber(address)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("trustNumber is ", trustNumber)
}

func TestCoreGetPendingTransactions(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))



	generator, err := connection.Core.GetGenerator()
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)
	transaction.Sender = generator
	transaction.Recipient = resources.Addr1
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e15))
	transaction.GasLimit = uint64(40000)
	transaction.Payload = "Transfer test"

	_, err = connection.Core.SendTransaction(transaction)
	if err != nil{
		t.FailNow()
	}
	pendingTransactions, _ := connection.Core.GetPendingTransactions()

	if len(pendingTransactions)!=0{
		t.Logf("%#v \n", pendingTransactions[0])
	}
}

func TestCoreGetStorageAt(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generator, _ := connection.Core.GetGenerator()

	value, err := connection.Core.GetStorageAt(generator, big.NewInt(0), block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(value)

}

func TestCoreGetProtocolVersion(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	version, err := connection.Core.GetProtocolVersion()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}


	t.Log(version)

}

func TestCoreSyncing(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	syncing, err := connection.Core.IsSyncing()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(syncing)

}

// todo: 测试失败
func TestCoreGetProof(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	generator, _ := connection.Core.GetGenerator()

	t.Logf("generator is %s \n", generator)

	proof, err := connection.Core.GetProof(generator, []string{"0", "1"}, block.LATEST)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", proof)

}