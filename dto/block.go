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
	"strconv"
)

type BlockDetails struct {
	Number           *big.Int              `json:"number"`
	Hash             string                `json:"hash"`
	ParentHash       string                `json:"parentHash"`
	LogsBloom        string                `json:"logsBloom"`
	StateRoot        string                `json:"stateRoot"`
	Generator        string                `json:"generator"`
	Regulatory       string                `json:"regulatory"`
	ExtraData        string                `json:"extraData"`
	Size             uint64                `json:"size"`
	Timestamp        uint64                `json:"timestamp"`
	TransactionsRoot string                `json:"transactionsRoot"`
	ReceiptsRoot     string                `json:"receiptsRoot"`
	Transactions     []TransactionResponse `json:"transactions"`
}

type BlockNoDetails struct {
	Number           *big.Int `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	LogsBloom        string   `json:"logsBloom"`
	StateRoot        string   `json:"stateRoot"`
	Generator        string   `json:"generator"`
	Regulatory       string   `json:"regulatory"`
	ExtraData        string   `json:"extraData"`
	Size             uint64   `json:"size"`
	Timestamp        uint64   `json:"timestamp"`
	TransactionsRoot string   `json:"transactionsRoot"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
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

	size, err := strconv.ParseUint(temp.Size, 0, 64)

	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.Size))
	}

	timestamp, err := strconv.ParseUint(temp.Timestamp, 0, 64)

	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.Timestamp))
	}

	b.Number = num
	b.Size = size
	b.Timestamp = timestamp

	return nil
}

func (b *BlockNoDetails) UnmarshalJSON(data []byte) error {
	type Alias BlockNoDetails
	temp := &struct {
		Number    string `json:"number"`
		Size      string `json:"size"`
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

	size, err := strconv.ParseUint(temp.Size, 0, 64)

	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.Size))
	}

	timestamp, err := strconv.ParseUint(temp.Timestamp, 0, 64)

	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.Timestamp))
	}

	b.Number = num
	b.Size = size
	b.Timestamp = timestamp

	return nil
}
