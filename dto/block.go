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

package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

type BlockDetails struct {
	Number           *big.Int              `json:"number"`
	Hash             string                `json:"hash"`
	ParentHash       string                `json:"parentHash"`
	Author           string                `json:"author,omitempty"`
	Miner            string                `json:"miner,omitempty"`
	Size             *big.Int              `json:"size"`
	GasLimit         *big.Int              `json:"gasLimit"`
	GasUsed          *big.Int              `json:"gasUsed"`
	Nonce            *big.Int              `json:"nonce"`
	Timestamp        *big.Int              `json:"timestamp"`
	ExtraData        string                `json:extraData`
	LogsBloom        string                `json:logsBloom`
	MixHash          string                `json:mixHash`
	ReceiptsRoot     string                `json:receiptsRoot`
	StateRoot        string                `json:stateRoot`
	TransactionsRoot string                `json:transactionsRoot`
	Transactions     []TransactionResponse `json:"transactions"`
}

type BlockNoDetails struct {
	Number           *big.Int `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Author           string   `json:"author,omitempty"`
	Miner            string   `json:"miner,omitempty"`
	Size             *big.Int `json:"size"`
	GasLimit         *big.Int `json:"gasLimit"`
	GasUsed          *big.Int `json:"gasUsed"`
	Nonce            *big.Int `json:"nonce"`
	Timestamp        *big.Int `json:"timestamp"`
	ExtraData        string   `json:extraData`
	LogsBloom        string   `json:logsBloom`
	MixHash          string   `json:mixHash`
	ReceiptsRoot     string   `json:receiptsRoot`
	StateRoot        string   `json:stateRoot`
	TransactionsRoot string   `json:transactionsRoot`
	Transactions     []string `json:"transactions"`
}

/**
 * How to un-marshal the block struct using the Big.Int rather than the
 * `complexReturn` type.
 */
func (b *BlockDetails) UnmarshalJSON(data []byte) error {
	type Alias BlockDetails
	temp := &struct {
		Number    string `json:"number"`
		Size      string `json:"size"`
		GasLimit  string `json:"gasLimit"`
		GasUsed   string `json:"gasUsed"`
		Nonce     string `json:"nonce"`
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	num, success := big.NewInt(0).SetString(temp.Number[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Number))
	}

	size, success := big.NewInt(0).SetString(temp.Size[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Size))
	}

	gasLimit, success := big.NewInt(0).SetString(temp.GasLimit[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.GasLimit))
	}

	gas, success := big.NewInt(0).SetString(temp.GasUsed[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.GasUsed))
	}

	nonce, success := big.NewInt(0).SetString(temp.Nonce[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Nonce))
	}

	timestamp, success := big.NewInt(0).SetString(temp.Timestamp[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Timestamp))
	}

	b.Number = num
	b.Size = size
	b.GasLimit = gasLimit
	b.GasUsed = gas
	b.Nonce = nonce
	b.Timestamp = timestamp

	return nil
}

func (b *BlockNoDetails) UnmarshalJSON(data []byte) error {
	type Alias BlockNoDetails
	temp := &struct {
		Number    string `json:"number"`
		Size      string `json:"size"`
		GasLimit  string `json:"gasLimit"`
		GasUsed   string `json:"gasUsed"`
		Nonce     string `json:"nonce"`
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	num, success := big.NewInt(0).SetString(temp.Number[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Number))
	}

	size, success := big.NewInt(0).SetString(temp.Size[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Size))
	}

	gasLimit, success := big.NewInt(0).SetString(temp.GasLimit[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.GasLimit))
	}

	gas, success := big.NewInt(0).SetString(temp.GasUsed[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.GasUsed))
	}

	nonce, success := big.NewInt(0).SetString(temp.Nonce[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Nonce))
	}

	timestamp, success := big.NewInt(0).SetString(temp.Timestamp[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Timestamp))
	}

	b.Number = num
	b.Size = size
	b.GasLimit = gasLimit
	b.GasUsed = gas
	b.Nonce = nonce
	b.Timestamp = timestamp

	return nil
}
