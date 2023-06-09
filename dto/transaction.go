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
	"github.com/tchain/go-tchain-sdk/crypto"
	"github.com/tchain/go-tchain-sdk/utils/types"
	"math/big"
	"strconv"
)

// TransactionParameters GO transaction to make more easy controll the parameters
type TransactionParameters struct {
	ChainId      uint64
	Sender       string
	Recipient    string
	AccountNonce uint64
	GasPrice     *big.Int
	GasLimit     uint64
	Amount       *big.Int
	Payload      types.ComplexString
}

type TransactionCallParameters struct {
	ChainId      uint64
	Sender       string
	Recipient    string
	AccountNonce uint64
	GasPrice     *big.Int
	GasLimit     uint64

	Amount  *big.Int
	Payload types.ComplexString
}

// RequestTransactionParameters JSON
type RequestTransactionParameters struct {
	ChainId      string `json:"chainId"`
	Sender       string `json:"sender"`
	Recipient    string `json:"recipient,omitempty"`
	AccountNonce string `json:"nonce,omitempty"`
	GasPrice     string `json:"gasPrice,omitempty"`
	GasLimit     string `json:"gas,omitempty"`
	Amount       string `json:"amount,omitempty"`
	Payload      string `json:"payload,omitempty"`
}

// Transform the GO transactions parameters to json style
func (params *TransactionParameters) Transform() *RequestTransactionParameters {
	request := new(RequestTransactionParameters)
	if params.ChainId != 0 {
		request.ChainId = fmt.Sprintf("0x%x", params.ChainId)
	}
	request.Sender = params.Sender
	if params.Recipient != "" {
		request.Recipient = params.Recipient
	}
	if params.AccountNonce != 0 {
		request.AccountNonce = fmt.Sprintf("0x%x", params.AccountNonce)
	}
	if params.GasLimit != 0 {
		request.GasLimit = fmt.Sprintf("0x%x", params.GasLimit)
	}
	if params.GasPrice != nil {
		request.GasPrice = "0x" + params.GasPrice.Text(16)
	}
	if params.Amount != nil {
		request.Amount = "0x" + params.Amount.Text(16)
	}
	if params.Payload != "" {
		request.Payload = params.Payload.ToHex()
	}
	return request
}

func (params *TransactionParameters) String() string {
	return fmt.Sprintf("Sender: %s, Recipient: %s, Amount:%d, GasPrice: %d, GasLimit: %d, AccountNonce: %d, Payload: %s",
		params.Sender, params.Recipient, params.Amount, params.GasLimit, params.GasPrice, params.AccountNonce, params.Payload)
}

type SignTransactionResponse struct {
	Raw         types.ComplexString     `json:"raw"`
	Transaction SignedTransactionParams `json:"tx"`
}

type SignedTransactionParams struct {
	ChainId  uint64              `json:"chainId"`
	GasPrice *big.Int            `json:"gasPrice"`
	GasLimit uint64              `json:"gas"`
	Hash     string              `json:"hash"`
	Input    string              `json:"input"`
	Nonce    *big.Int            `json:"nonce"`
	To       string              `json:"to"`
	Amount   *big.Int            `json:"value"`
	SignNode *crypto.Signature   `json:"signNode"`
	SignUser []*crypto.Signature `json:"singUser"`
}

func (sp *SignedTransactionParams) UnmarshalJSON(data []byte) error {
	type Alias SignedTransactionParams

	temp := &struct {
		ChainId  string `json:"chainId"`
		Nonce    string `json:"nonce"`
		GasPrice string `json:"gasPrice"`
		GasLimit string `json:"gas"`
		Amount   string `json:"value"`
		*Alias
	}{
		Alias: (*Alias)(sp),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	chainId, err := strconv.ParseUint(temp.ChainId, 0, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.ChainId))
	}

	gas, err := strconv.ParseUint(temp.GasLimit, 0, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.GasLimit))
	}

	gasPrice, success := big.NewInt(0).SetString(temp.GasPrice[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.GasPrice))
	}

	nonce, success := big.NewInt(0).SetString(temp.Nonce[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Nonce))
	}

	val, success := big.NewInt(0).SetString(temp.Amount[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Amount))
	}

	sp.ChainId = chainId
	sp.GasLimit = gas
	sp.GasPrice = gasPrice
	sp.Nonce = nonce
	sp.Amount = val

	return nil
}

type TransactionResponse struct {
	ChainId          uint64              `json:"chainId"`
	BlockHash        string              `json:"blockHash"`
	BlockNumber      *big.Int            `json:"blockNumber"`
	Sender           string              `json:"sender"`
	Gas              uint64              `json:"gas"`
	GasPrice         *big.Int            `json:"gasPrice"`
	Hash             string              `json:"hash"`
	Payload          string              `json:"payload"`
	Nonce            uint64              `json:"nonce"`
	Recipient        string              `json:"recipient"`
	TransactionIndex uint                `json:"transactionIndex"`
	Amount           *big.Int            `json:"amount"`
	SignNode         *crypto.Signature   `json:"signNode"`
	SignUser         []*crypto.Signature `json:"signUser"`
}

func (t *TransactionResponse) UnmarshalJSON(data []byte) error {
	type Alias TransactionResponse

	temp := &struct {
		ChainId          string `json:"chainId"`
		BlockNumber      string `json:"blockNumber"`
		Gas              string `json:"gas"`
		GasPrice         string `json:"gasPrice"`
		Nonce            string `json:"nonce"`
		TransactionIndex string `json:"transactionIndex"`
		Amount           string `json:"amount"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	chainId, err := strconv.ParseUint(temp.ChainId, 0, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.ChainId))
	}

	if len(temp.BlockNumber) == 0 {
		temp.BlockNumber = "0x0"
	}
	blockNum, success := big.NewInt(0).SetString(temp.BlockNumber[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.BlockNumber))
	}

	gas, err := strconv.ParseUint(temp.Gas, 0, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.Gas))
	}

	gasPrice, success := big.NewInt(0).SetString(temp.GasPrice[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.GasPrice))
	}

	nonce, err := strconv.ParseUint(temp.Nonce, 0, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.Nonce))
	}

	if len(temp.TransactionIndex) == 0 {
		temp.TransactionIndex = "0x0"
	}
	txIndex := big.NewInt(0).Uint64()
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.TransactionIndex))
	}

	amount, success := big.NewInt(0).SetString(temp.Amount[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Amount))
	}

	t.ChainId = chainId
	t.BlockNumber = blockNum
	t.Gas = gas
	t.GasPrice = gasPrice
	t.Nonce = nonce
	t.TransactionIndex = uint(txIndex)
	t.Amount = amount

	return nil
}

type TransactionReceipt struct {
	TransactionHash   string            `json:"transactionHash"`
	TransactionIndex  *big.Int          `json:"transactionIndex"`
	BlockHash         string            `json:"blockHash"`
	BlockNumber       *big.Int          `json:"blockNumber"`
	From              string            `json:"sender"`
	To                string            `json:"recipient"`
	CumulativeGasUsed uint64            `json:"cumulativeGasUsed"`
	GasUsed           uint64            `json:"gasUsed"`
	ContractAddress   string            `json:"contractAddress"`
	Logs              []TransactionLogs `json:"logs"`
	LogsBloom         string            `json:"logsBloom"`
	Status            bool              `json:"status"`
}

func (r *TransactionReceipt) UnmarshalJSON(data []byte) error {
	type Alias TransactionReceipt

	temp := &struct {
		TransactionIndex  string `json:"transactionIndex"`
		BlockNumber       string `json:"blockNumber"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		GasUsed           string `json:"gasUsed"`
		Status            string `json:"status"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	blockNum, success := big.NewInt(0).SetString(temp.BlockNumber[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.BlockNumber))
	}

	if len(temp.TransactionIndex) == 0 {
		temp.TransactionIndex = "0x0"
	}
	txIndex, success := big.NewInt(0).SetString(temp.TransactionIndex[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.TransactionIndex))
	}

	gasUsed, err := strconv.ParseUint(temp.GasUsed, 0, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.GasUsed))
	}

	cumulativeGas, err := strconv.ParseUint(temp.CumulativeGasUsed, 0, 64)

	if err != nil {
		return errors.New(fmt.Sprintf("Error converting %s to uint64", temp.CumulativeGasUsed))
	}

	status, success := big.NewInt(0).SetString(temp.Status[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Status))
	}

	r.TransactionIndex = txIndex
	r.BlockNumber = blockNum
	r.CumulativeGasUsed = cumulativeGas
	r.GasUsed = gasUsed
	r.Status = false
	if status.Cmp(big.NewInt(1)) == 0 {
		r.Status = true
	}

	return nil
}

type TransactionLogs struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      *big.Int `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         *big.Int `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

func (r *TransactionLogs) UnmarshalJSON(data []byte) error {
	type Alias TransactionLogs

	log := &struct {
		TransactionIndex string `json:"transactionIndex"`
		BlockNumber      string `json:"blockNumber"`
		LogIndex         string `json:"logIndex"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &log); err != nil {
		return err
	}

	blockNumLog, success := big.NewInt(0).SetString(log.BlockNumber[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", log.BlockNumber))
	}

	if len(log.TransactionIndex) == 0 {
		log.TransactionIndex = "0x0"
	}
	txIndexLogs, success := big.NewInt(0).SetString(log.TransactionIndex[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", log.TransactionIndex))
	}

	logIndex, success := big.NewInt(0).SetString(log.LogIndex[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", log.LogIndex))
	}

	r.BlockNumber = blockNumLog
	r.TransactionIndex = txIndexLogs
	r.LogIndex = logIndex
	return nil

}

type CandidateResponse struct {
	Owner           string `json:"owner"`           // 候选人地址
	Name            string `json:"name"`            // 候选人名称
	Active          bool   `json:"active"`          // 当前是否是候选人
	Url             string `json:"url"`             // 节点的URL
	VoteCount       int64  `json:"voteCount"`       // 收到的票数
	TotalBounty     int64  `json:"totalBounty"`     // 总奖励金额
	ExtractedBounty int64  `json:"extractedBounty"` // 已提取奖励金额
	LastExtractTime int64  `json:"lastExtractTime"` // 上次提权时间
	Website         string `json:"website"`         // 见证人网站
}

// Voter is the information of who has vote witness candidate
type VoterResponse struct {
	Owner             string   `json:"owner"`             // 投票人的地址
	IsProxy           bool     `json:"isProxy"`           // 是否是代理人
	ProxyVoteCount    int64    `json:"proxyVoteCount"`    // 收到的代理的票数
	Proxy             string   `json:"proxy"`             // 该节点设置的代理人
	LastVoteCount     int64    `json:"lastVoteCount"`     // 上次投的票数
	LastVoteTimeStamp int64    `json:"lastVoteTimeStamp"` // 上次投票时间戳
	VoteCandidates    []string `json:"voteCandidates"`    // 投了哪些人
}

// Stake is the information of a user
type StakeResponse struct {
	Owner              string `json:"owner"`              // 抵押代币的所有人
	StakeCount         int64  `json:"stakeCount"`         // 抵押的代币数量
	LastStakeTimeStamp int64  `json:"lastStakeTimeStamp"` // 上次抵押时间戳
}
