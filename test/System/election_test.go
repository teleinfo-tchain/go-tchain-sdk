package System

import (
	"fmt"
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/test/resources"
	"io/ioutil"
	"testing"
	"time"
)

// 注册成为可信节点（只有理事长可以注册，即监管节点）
func TestRegisterTrustNode(t *testing.T) {
	// 签名的节点是联盟成员
	// file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	// con, sigPara, err := connectWithSig("did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc", file)
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-10-19T05-33-49.419105162Z--did_bid_qwer_sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	con, sigPara, err := connectWithSig("did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	// 注册的ID（地址） 对应的keystore文件的密码
	idPassword := "tele"

	for _, registerTrustNode := range []struct {
		*dto.PeerNodeInfo
		KeyStoreFile string
	}{
		{
			PeerNodeInfo: &dto.PeerNodeInfo{
				Id:        "did:bid:qwer:sfbg3EBvpfCkqzh2dweGH217Gm758w9w",
				Apply:     "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y",
				PublicKey: "16Uiu2HAmVDhHxPFetMHXZxHeZSxmPX95AKbD7Y1PFzFZbeDJi3tx", NodeName: "node4",
				Website: "https://www.teleinfo4.com", CompanyName: "teleInfo",
				CompanyCode: "91310000717854505W", Ip: "192.168.1.30", Port: 44402,
			},
			KeyStoreFile: "UTC--2021-11-15T02-13-59.857639300Z--did-bid-qwer-sfbg3EBvpfCkqzh2dweGH217Gm758w9w",
		},
		// {
		// 	PeerNodeInfo: &dto.PeerNodeInfo{
		// 		Id:        "did:bid:qwer:sf8MD76bSGkxCzFXnyDNvtrGgCdodhBF",
		// 		Apply:     "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y",
		// 		PublicKey: "16Uiu2HAmDbHWA7MdEh9EQvbAPKb4AcrFMp9RwTevuNqzxLU57goD", NodeName: "node5",
		// 		Website: "https://www.teleinfo5.com", CompanyName: "teleInfo",
		// 		CompanyCode: "91310000717854505W", Ip: "139.198.15.29", Port: 44502,
		// 	},
		// 	KeyStoreFile: "UTC--2021-11-15T02-14-01.409690300Z--did-bid-qwer-sf8MD76bSGkxCzFXnyDNvtrGgCdodhBF",
		// },
		// {
		// 	PeerNodeInfo: &dto.PeerNodeInfo{Id: "did:bid:qwer:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY",
		// 		Apply:     "did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y",
		// 		PublicKey: "16Uiu2HAm3z3rBzpH5tpFkdTxf7CU2JSdEDT4A6JH78ieKc69Aotp",
		// 		NodeName:  "node2",
		// 		Website:   "https://www.teleinfo4.com", CompanyName: "teleInfo",
		// 		CompanyCode: "91310000717854505W", Ip: "139.198.15.29", Port: 44101},
		// 	KeyStoreFile: "UTC--2021-10-19T05-33-49.486361230Z--did_bid_qwer_sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY",
		// },
	} {

		fmt.Println(con.Core.GetBalance(registerTrustNode.Id, block.LATEST))
		// 注册的ID（地址）对应的keystore文件
		registerNodeIdFile := bif.GetCurrentAbPath() + resources.KeyStoreFile + registerTrustNode.KeyStoreFile
		idKeyFileData, err := ioutil.ReadFile(registerNodeIdFile)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if len(idKeyFileData) == 0 {
			t.Errorf("idKeyFileData can't be empty")
			t.FailNow()
		}

		transactionHash, err := ele.RegisterTrustNode(sigPara, registerTrustNode.PeerNodeInfo, idPassword, idKeyFileData)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Log(transactionHash, err)

		time.Sleep(8 * time.Second)

		log, err := con.System.SystemLogDecode(transactionHash)

		if err != nil {
			t.Errorf("err log : %v ", err)
			t.FailNow()
		}

		if !log.Status {
			t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
		}
	}
}

// 删除可信节点（只有理事长可以注册，即监管节点）
func TestDeleteTrustNode(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.TestAddressRegulatoryFile
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	deleteReason := "节点信息注册错误"
	transactionHash, err := ele.DeleteTrustNode(sigPara, "did:bid:qwer:sf8MD76bSGkxCzFXnyDNvtrGgCdodhBF", deleteReason)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

// 申请成为候选节点
func TestApplyCandidate(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-10-19T05-33-49.419105162Z--did_bid_qwer_sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	con, sigPara, err := connectWithSig("did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	transactionHash, err := ele.ApplyCandidate(sigPara, "did:bid:qwer:sfbg3EBvpfCkqzh2dweGH217Gm758w9w")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

// 取消成为候选节点
func TestCancelCandidate(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.TestAddressRegulatoryFile
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	transactionHash, err := ele.CancelCandidate(sigPara, resources.RegisterAllianceTwo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestVoteCandidate(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-10-19T05-33-49.419105162Z--did_bid_qwer_sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y"
	con, sigPara, err := connectWithSig("did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	transactionHash, err := ele.VoteCandidate(sigPara, "did:bid:qwer:sfbg3EBvpfCkqzh2dweGH217Gm758w9w")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestCancelConsensusNode(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.TestAddressRegulatoryFile
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	cancleReason := "违背规则"
	transactionHash, err := ele.CancelConsensusNode(sigPara, "did:bid:qwer:sfH7tKf9ohrsWGacPsr8TNvFdLEbvcej", cancleReason)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)
	// 0x57357c41bbdb60037c9b58b4995194208dcc71eb47192a76e4feaf23cb929eb6

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

// 设置投票期限  只能监管节点设置
func TestSetDeadline(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.TestAddressRegulatoryFile
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	// 以妙为单位
	var deadline uint64 = 1 * 24 * 3600
	transactionHash, err := ele.SetDeadline(sigPara, deadline)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestExtractOwnBounty(t *testing.T) {
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.RegisterAllianceOneFile
	con, sigPara, err := connectWithSig(resources.RegisterAllianceOne, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	transactionHash, err := ele.ExtractOwnBounty(sigPara)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestIssueAdditionalBounty(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.TestAddressRegulatoryFile
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	transactionHash, err := ele.IssueAdditionalBounty(sigPara)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestGetRestBIFBounty(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	restBounty, err := ele.GetRestBIFBounty()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// 999984574000000000000000000
	// 1049984514000000000000000000
	fmt.Printf("restBounty %s \n", restBounty.String())
}

func TestAllTrusted(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodes, err := ele.AllTrusted()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("AllTrusted nodes is %+v \n", node)
	}
}

func TestAllCandidates(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodes, err := ele.AllCandidates()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("AllCandidates nodes is %+v \n", node)
	}
}

func TestGetAllConsensus(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodes, err := ele.AllConsensus()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("AllConsensus nodes is %+v \n", node)
	}

}

func TestAllNodes(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodes, err := ele.AllNodes()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("AllNodes node is %s , %d\n", node.Id, node.Role)
	}
}

func TestGetPeerNode(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	peerNodeId := resources.RegisterAllianceTwo

	node, err := ele.GetPeerNode(peerNodeId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("PeerNode is %+v \n", node)
}

func TestVoteNodes(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	voter := resources.RegisterAllianceTwo
	nodes, err := ele.VoteNodes(voter)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("VoteNodes nodes is %+v \n", node)
	}
}

func TestApplyNodes(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	apply := resources.RegisterAllianceTwo
	nodes, err := ele.ApplyNodes(apply)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("ApplyNodes nodes is %+v \n", node)
	}
}

func TestGetDeadline(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	deadline, err := ele.GetDeadline()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("deadline is %v \n", deadline)
}

func TestNodeBounty(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodeBounty, err := ele.NodeBounty("did:bid:qwer:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("nodeBounty is %v \n", nodeBounty)
}
