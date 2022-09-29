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
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/account"
	"github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/test/resources"
	"github.com/tchain/go-tchain-sdk/utils"
	"github.com/tchain/go-tchain-sdk/utils/hexutil"
	"math/big"
	"strconv"
	"testing"
	"time"
)

// 测试转账
func TestCoreSendTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	priKey := "e41219552564c956edeb0fa782c7760a6f5ade504768b3570c68dc0459a7889a"

	recipientStr := "did:bid:qwer:zftAgNtnQzLMGJHKPMdn9quPvuikNWUZ"

	balance, err := connection.Core.GetBalance(recipientStr, block.LATEST)
	if err == nil {
		balBif, _ := utils.FromWei(balance)
		t.Log("recipient balance is ", balBif)
	}

	var recipient utils.Address
	recipient = utils.StringToAddress(recipientStr)

	chainId, _ := connection.Core.GetChainId()

	nonce, _ := connection.Core.GetTransactionCount("did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", block.LATEST)

	tx := &account.SignTxParams{
		Recipient: &recipient,
		Nonce:     nonce,
		GasPrice:  big.NewInt(2000000),
		GasLimit:  uint64(41000),
		Amount:    big.NewInt(50000000000),
		Payload:   nil,
		ChainId:   chainId,
	}

	res, err := account.SignTransaction(tx, priKey, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txID, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txID)
	}

	// if success, get transaction hash
	t.Log(txID)
}
