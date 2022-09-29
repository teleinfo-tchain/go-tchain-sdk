package system

import (
	"errors"
	"fmt"
	"github.com/tchain/go-tchain-sdk/abi"
	"github.com/tchain/go-tchain-sdk/account"
	"github.com/tchain/go-tchain-sdk/dto"
	"math/big"
	"strings"
)

// Election - The Election Module
type Election struct {
	super *System
	abi   abi.ABI
}

// NewElection - NewElection初始化
func (sys *System) NewElection() *Election {
	parsedAbi, _ := abi.JSON(strings.NewReader(ElectionAbiJSON))

	election := new(Election)
	election.abi = parsedAbi
	election.super = sys
	return election
}

func (e *Election) judgeScheme(url string) bool {
	parts := strings.Split(url, "://")
	if len(parts) != 2 || parts[0] == "" {
		return false
	}
	return true
}

func registerTrustNodePreCheck(trustNode *dto.PeerNodeInfo) error {
	trustNode.Id = strings.TrimSpace(trustNode.Id)
	trustNode.Apply = strings.TrimSpace(trustNode.Apply)
	trustNode.PublicKey = strings.TrimSpace(trustNode.PublicKey)
	trustNode.NodeName = strings.TrimSpace(trustNode.NodeName)

	if !isValidHexAddress(trustNode.Id) {
		return errors.New("registerTrustNode id is not valid bid")
	}

	if len(trustNode.Apply) == 0 {
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is 0", "Apply")
	}

	if len(trustNode.PublicKey) != 53 {
		return fmt.Errorf("parameter is not illegal, parameter is %s, length is not 53", "")
	}

	if len(trustNode.NodeName) < 4 || len(trustNode.NodeName) > 20 {
		return errors.New("registerCertificate Website len is range in 4 to 20")
	}

	if len(trustNode.Website) < 6 || len(trustNode.Website) > 60 {
		return errors.New("registerCertificate Website len is range in 6 to 50")
	}

	if trustNode.NodeType > 1 {
		return fmt.Errorf("node type is Error, type is not 0 or 1, nodeType is %d", trustNode.NodeType)
	}

	if len(trustNode.CompanyName) == 0 || isBlankCharacter(trustNode.NodeName) {
		return errors.New("registerCertificate CompanyName can't be empty or blank character")
	}
	if len(trustNode.CompanyCode) == 0 || isBlankCharacter(trustNode.NodeName) {
		return errors.New("registerCertificate CompanyCode can't be empty or blank character")
	}

	if !isLegalIP(trustNode.Ip) {
		return errors.New("parameter is not illegal, parameter is Ip, Ip is illegal")
	}
	if trustNode.Port > 65535 {
		return errors.New("parameter is not illegal, parameter is Port, Port should be in range 0 to 65535")
	}

	return nil
}

func (e *Election) RegisterTrustNode(signTxParams *SysTxParams, trustNode *dto.PeerNodeInfo, idPassword string, idKeyFile []byte) (string, error) {
	err := registerTrustNodePreCheck(trustNode)
	if err != nil {
		return "", err
	}

	messageSha3, signature, err := account.MessageSignatureBtc("registerTrustNode", idPassword, idKeyFile)
	if err != nil {
		return "", err
	}
	// 添加messageSha3 signature
	trustNode.MessageSha3 = messageSha3
	trustNode.Signature = signature

	// encode
	// trustNode is a struct we need to use the components.
	var values []interface{}
	values = e.super.structToInterface(*trustNode, values)
	inputEncode, err := e.abi.Pack("registerTrustNode", values...)

	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) DeleteTrustNode(signTxParams *SysTxParams, trustNodeId string, revokeReason string) (string, error) {
	if !isValidHexAddress(trustNodeId) {
		return "", errors.New("trustNodeId is not valid bid")
	}

	// Revoke is a struct we need to use the components.
	var values []interface{}
	type nodeInfo struct {
		TrustNode    string
		RevokeReason string
	}
	var node = nodeInfo{trustNodeId, revokeReason}

	values = e.super.structToInterface(node, values)
	inputEncode, err := e.abi.Pack("deleteTrustNode", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) ApplyCandidate(signTxParams *SysTxParams, candidateAddress string) (string, error) {
	if !isValidHexAddress(candidateAddress) {
		return "", errors.New("candidateAddress is not valid bid")
	}

	// encoding
	inputEncode, err := e.abi.Pack("applyCandidate", candidateAddress)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) CancelCandidate(signTxParams *SysTxParams, candidateAddress string) (string, error) {
	if !isValidHexAddress(candidateAddress) {
		return "", errors.New("candidateAddress is not valid bid")
	}

	// encoding
	inputEncode, err := e.abi.Pack("cancelCandidate", candidateAddress)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) VoteCandidate(signTxParams *SysTxParams, candidates string) (string, error) {
	candis := strings.TrimSpace(candidates)
	if len(candis) == 0 {
		return "", errors.New("candidateAddress is not valid bid")
	}

	// encoding
	inputEncode, err := e.abi.Pack("voteCandidate", candis)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) CancelConsensusNode(signTxParams *SysTxParams, consensusNode string, cancelConsensusReason string) (string, error) {
	if !isValidHexAddress(consensusNode) {
		return "", errors.New("consensusNode is not valid bid")
	}

	var values []interface{}
	type nodeInfo struct {
		Consensus             string
		CancelConsensusReason string
	}
	var consensusNodeInfo = nodeInfo{consensusNode, cancelConsensusReason}

	values = e.super.structToInterface(consensusNodeInfo, values)
	inputEncode, err := e.abi.Pack("cancelConsensusNode", values...)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) SetDeadline(signTxParams *SysTxParams, deadline uint64) (string, error) {
	// encoding
	inputEncode, err := e.abi.Pack("setDeadline", deadline)
	if err != nil {
		return "", err
	}

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) ExtractOwnBounty(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("extractOwnBounty")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) IssueAdditionalBounty(signTxParams *SysTxParams) (string, error) {
	// encoding
	inputEncode, _ := e.abi.Pack("issueAdditionalBounty")

	signedTx, err := e.super.prePareSignTransaction(signTxParams, inputEncode, ElectionContract)
	if err != nil {
		return "", err
	}

	return e.super.sendRawTransaction(signedTx)
}

func (e *Election) GetRestBIFBounty() (*big.Int, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_restBIFBounty", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionRestBIFBounty()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) AllTrusted() ([]*dto.PeerNodeDetail, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allTrusted", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionPeerNodes()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) AllCandidates() ([]*dto.PeerNodeDetail, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allCandidates", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionPeerNodes()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) AllConsensus() ([]*dto.PeerNodeDetail, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allConsensus", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionPeerNodes()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) AllNodes() ([]*dto.PeerNodeDetail, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_allNodes", nil)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionPeerNodes()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) GetPeerNode(peerNodeId string) (*dto.PeerNodeDetail, error) {
	if !isValidHexAddress(peerNodeId) {
		return nil, errors.New("peerNodeId is not valid bid")
	}

	params := make([]string, 1)
	params[0] = peerNodeId

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_peerNode", params)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionPeerNode()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) VoteNodes(voter string) ([]*dto.PeerNodeDetail, error) {
	if !isValidHexAddress(voter) {
		return nil, errors.New("voter is not valid bid")
	}

	params := make([]string, 1)
	params[0] = voter

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_voteNodes", params)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionPeerNodes()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) ApplyNodes(apply string) ([]*dto.PeerNodeDetail, error) {
	if !isValidHexAddress(apply) {
		return nil, errors.New("apply is not valid bid")
	}

	params := make([]string, 1)
	params[0] = apply

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_applyNodes", params)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionPeerNodes()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Election) GetDeadline() (uint64, error) {
	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_deadline", nil)
	if err != nil {
		return 0, err
	}

	res, err := pointer.ToUint64()
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (e *Election) NodeBounty(peerNodeId string) (*dto.PeerNodeBounty, error) {
	if !isValidHexAddress(peerNodeId) {
		return nil, errors.New("peerNodeId is not valid bid")
	}

	params := make([]string, 1)
	params[0] = peerNodeId

	pointer := &dto.SystemRequestResult{}

	err := e.super.provider.SendRequest(pointer, "election_nodeBounty", params)
	if err != nil {
		return nil, err
	}

	res, err := pointer.ToElectionNodeBounty()
	if err != nil {
		return nil, err
	}

	return res, nil
}