package System

import (
	"errors"
	"fmt"
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/system"
	"github.com/tchain/go-tchain-sdk/test/resources"
	"io/ioutil"
	"math/big"
	"testing"
	"time"
)

func connectWithSigNow(sigAddr string, singAddrFile string) (*bif.Bif, *system.SysTxParams, error) {
	var connection = bif.NewBif(providers.NewHTTPProvider("127.0.0.1:55550", 10, false))
	chainId, err := connection.Core.GetChainId()
	if err != nil {
		return nil, nil, err
	}

	nonce, err := connection.Core.GetTransactionCount(sigAddr, block.LATEST)
	if err != nil {
		return nil, nil, err
	}

	// keyFileData 还可以进一步校验
	keyFileData, err := ioutil.ReadFile(singAddrFile)
	if err != nil {
		return nil, nil, err
	}
	if len(keyFileData) == 0 {
		return nil, nil, errors.New("keyFileData can't be empty")
	}

	sysTxParams := new(system.SysTxParams)
	sysTxParams.From = sigAddr
	sysTxParams.IsSM2 = resources.NotSm2
	sysTxParams.Password = "teleinfo"
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(45)
	sysTxParams.Gas = 200000
	sysTxParams.Nonce = nonce
	sysTxParams.ChainId = chainId

	return connection, sysTxParams, nil
}

func connectBifNow() (*bif.Bif, error) {
	var connection = bif.NewBif(providers.NewHTTPProvider("127.0.0.1:55550", 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func TestApplySubChain(t *testing.T) {
	// 签名的节点是联盟成员
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	con, sigPara, err := connectWithSigNow("did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	sub := con.System.NewSubChain()

	for _, subChain := range []dto.SubChainInfo{
		{
			"did:bid:llj1:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", "did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc",
			"node1",
			"llj1", "煤炭12",
			"公链",
			"hotStuff", "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79",
		},
		// {
		// 	"did:bid:llj1:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY", "did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc",
		// 	"node2",
		// 	"qwer", "煤炭12",
		// 	"公链",
		// 	"hotStuff", "0xc1912fee45d61c87cc5ea59dae311904cd86b84fee17cc96966216f811ce6a79",
		// },
	} {
		fmt.Println("now apply is ", subChain.Id)
		// 注册的ID（地址）对应的keystore文件
		applySubChainIdFile := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
		idKeyFileData, err := ioutil.ReadFile(applySubChainIdFile)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if len(idKeyFileData) == 0 {
			t.Errorf("idKeyFileData can't be empty")
			t.FailNow()
		}
		transactionHash, err := sub.ApplySubChain(sigPara, &subChain)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Log(transactionHash, err)

		time.Sleep(10 * time.Second)

		// t.Log(sigPara, subChain, sub)
		// log, err := con.System.SystemLogDecode("0x2e36de42eb2956e8dc0a9173e4b17c80590959cc13fb69788a2770fdc18987bd")
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

func TestRevokeSubChain(t *testing.T) {
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	con, sigPara, err := connectWithSigNow("did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	subchain := con.System.NewSubChain()

	deleteReason := "违规"
	transactionHash, err := subchain.Revoke(sigPara, "did:bid:llj1:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y", deleteReason)
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

func TestVoteSubChain(t *testing.T) {
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	con, sigPara, err := connectWithSigNow("did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	subChain := con.System.NewSubChain()

	transactionHash, err := subChain.VoteSubChain(sigPara, "did:bid:llj1:sf2BX7RNbmdtGgyYuD3HL7H7w1XmGSTFY")
	// transactionHash, err := subChain.VoteSubChain(sigPara, "did:bid:llj1:sf25XGBQU8E8wGFo9wGKo95jUgtYPM24Y")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(transactionHash, err)

	time.Sleep(10 * time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestSetSubDeadline(t *testing.T) {
	file := bif.GetCurrentAbPath() + resources.KeyStoreFile + "UTC--2021-08-26T09-21-33.005300071Z--did-bid-llj1-sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"
	con, sigPara, err := connectWithSigNow("did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc", file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	subChain := con.System.NewSubChain()

	// 以秒为单位
	var deadline uint64 = 1 * 24 * 3600
	transactionHash, err := subChain.SetDeadline(sigPara, deadline)
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

func TestAllSubChains(t *testing.T) {
	con, err := connectBifNow()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	subChain := con.System.NewSubChain()

	subChains, err := subChain.AllSubChains()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, sub := range subChains {
		t.Logf("subChain is %+v \n", sub)
	}
}

func TestGetSubChain(t *testing.T) {
	var con = bif.NewBif(providers.NewHTTPProvider("127.0.0.1:55550", 10, false))
	_, err := con.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	subChain := con.System.NewSubChain()

	subChainId := "did:bid:llj1:sfrVXK5LxB6ZYrqXsaqp6g3izMkm2r8n"

	sub, err := subChain.GetSubChain(subChainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("SubChain is %+v \n", sub)

}

func TestGetSubDeadline(t *testing.T) {
	var con = bif.NewBif(providers.NewHTTPProvider("127.0.0.1:55550", 10, false))
	_, err := con.Core.GetChainId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	subChain := con.System.NewSubChain()

	deadline, err := subChain.GetDeadline()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("deadline is %v \n", deadline)
}
