package dto

import "math/big"

// RoundStateInfoResponse is the information of RoundState
// 数据类型定义是否合适？？
type RoundStateInfo struct {
	Commits          *stageInfo     `json:"commits"`
	LockedHash       string         `json:"lockedHash"`
	Prepares         *stageInfo     `json:"prepares"`
	Proposer         string         `json:"proposer"`
	Round            *big.Int       `json:"round"`
	Sequence         *big.Int       `json:"sequence"`
	View             string         `json:"view"`      // string的数据类型待确认!!!!!!
}

type stageInfo struct {
	Messages        []string        `json:"messages"`  // message的数据类型待确认!!!!!!
	ValSet          []string        `json:"valSet"`
	View            *view           `json:"view"`
}

type view struct {
	Round        	*big.Int        `json:"round"`
	Sequence     	*big.Int        `json:"sequence"`
}


type RoundChangeSetInfo struct {
	RoundChanges     interface{}     `json:"roundChanges"` // message的数据结构待确认!!!!!!
	Validates        []string         `json:"validates"`
}

// 数据结构待确认
type Message struct {
}