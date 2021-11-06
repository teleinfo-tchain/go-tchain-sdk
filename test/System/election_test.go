package System

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"testing"
	"time"
)

// 注册成为可信节点（只有理事长可以注册，即监管节点）
func TestRegisterTrustNode(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	con, sigPara, err := connectWithSig("did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	// 注册的ID（地址） 对应的keystore文件的密码
	idPassword := "teleinfo"

	for _, registerTrustNode := range []dto.PeerNodeInfo{
		{
			"did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc", "did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc",
			"16Uiu2HAmSgtVcHBHe79Ey3H3DHHxqbFBCFLL5UcEAUz8sBBxouui", "node2",
			"", "",
			"https://www.teleinfo.com", 0, "teleInfo",
			"91310000717854505W", "127.0.0.1", 44444,
			1,
		},
	} {
		// 注册的ID（地址）对应的keystore文件
		registerNodeIdFile := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
		idKeyFileData, err := ioutil.ReadFile(registerNodeIdFile)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if len(idKeyFileData) == 0 {
			t.Errorf("idKeyFileData can't be empty")
			t.FailNow()
		}

		transactionHash, err := ele.RegisterTrustNode(sigPara, &registerTrustNode, idPassword, idKeyFileData)
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

	deleteReason := "节点违规"
	transactionHash, err := ele.DeleteTrustNode(sigPara, resources.RegisterAllianceOne, deleteReason)
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
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.RegisterAllianceTwoFile
	con, sigPara, err := connectWithSig(resources.RegisterAllianceTwo, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	transactionHash, err := ele.ApplyCandidate(sigPara, resources.RegisterAllianceTwo)
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
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + resources.RegisterAllianceTwoFile
	con, sigPara, err := connectWithSig(resources.RegisterAllianceTwo, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	transactionHash, err := ele.VoteCandidate(sigPara, resources.RegisterAllianceTwo)
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
	transactionHash, err := ele.CancelConsensusNode(sigPara, resources.RegisterAllianceTwo, cancleReason)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(8 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)
	//0x57357c41bbdb60037c9b58b4995194208dcc71eb47192a76e4feaf23cb929eb6

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
	file := bif.GetCurrentAbPath()+ resources.KeyStoreFile + resources.TestAddressRegulatoryFile
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
	file := bif.GetCurrentAbPath()+ resources.KeyStoreFile + resources.TestAddressRegulatoryFile
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

	//999984574000000000000000000
	//1049984514000000000000000000
	fmt.Printf("restBounty %s \n", restBounty.String())
}

func TestGetAllTrusted(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodes, err := ele.GetAllTrusted()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("AllTrusted nodes is %+v \n", node)
	}
}

func TestGetAllCandidates(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodes, err := ele.GetAllCandidates()
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

	nodes, err := ele.GetAllConsensus()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("AllConsensus nodes is %+v \n", node)
	}
}

func TestGetAllNodes(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	nodes, err := ele.GetAllNodes()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("AllNodes node is %+v \n", node)
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

func TestGetVoteNodes(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	voter := resources.RegisterAllianceTwo
	nodes, err := ele.GetVoteNodes(voter)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, node := range nodes {
		t.Logf("VoteNodes nodes is %+v \n", node)
	}
}

func TestGetApplyNodes(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ele := con.System.NewElection()

	apply := resources.RegisterAllianceTwo
	nodes, err := ele.GetApplyNodes(apply)
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
