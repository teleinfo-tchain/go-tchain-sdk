package dto

import (
	"github.com/bif/bif-sdk-go/common"
	"math/big"
)

type Message struct {
	Code          uint64 `json:"code"`
	Msg           []byte `json:"message"`
	Address       string `json:"address"`
	Signature     []byte `json:"signature"`
	CommittedSeal []byte `json:"committedSeal"`
}

type View struct {
	Round    *big.Int `json:"round"`
	Sequence *big.Int `json:"sequence"`
}

type MessageSet struct {
	View     *View      `json:"view"`
	ValSet   []string   `json:"valSet"`
	Messages []*Message `json:"messages"`
}

// RoundStateInfoResponse is the information of RoundState
type RoundStateInfo struct {
	Commits    *MessageSet `json:"commits"`
	LockedHash common.Hash `json:"lockedHash"`
	Prepares   *MessageSet `json:"prepares"`
	Proposer   string      `json:"proposer"`
	Round      *big.Int    `json:"round"`
	Sequence   *big.Int    `json:"sequence"`
	View       *View       `json:"view"`
}

type RoundChangeSetInfo struct {
	RoundChanges map[uint64]*MessageSet `json:"roundChanges"`
	Validates    []string               `json:"validates"`
}
