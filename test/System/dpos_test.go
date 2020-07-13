package System

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestGetValidators(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	// 距离最新区块间隔数超过126会报错，因此测试最新区块
	blockNumber, err := connection.Core.GetBlockNumber()
	//fmt.Println("blockNumber is ", blockNumber)
	//blockNumberOld := new(big.Int)
	//fmt.Println("blockNumberOld is ", blockNumberOld.Sub(blockNumber, big.NewInt(127)))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	Dpos := connection.System.NewDpos()
	validators, err := Dpos.GetValidators(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(validators)
}

func TestGetValidatorsAtHash(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	// 距离最新区块间隔数超过126会报错，因此测试最新区块
	blockNumber, err := connection.Core.GetBlockNumber()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	block, err := connection.Core.GetBlockByNumber(blockNumber, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	Dpos := connection.System.NewDpos()
	//input block hash
	validators, err := Dpos.GetValidatorsAtHash(block.(*dto.BlockNoDetails).Hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(validators)
}

func TestRoundStateInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	Dpos := connection.System.NewDpos()
	roundStateInfo, err := Dpos.RoundStateInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("roundStateInfo is %#v \n", roundStateInfo)
}

func TestRoundChangeSetInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))

	Dpos := connection.System.NewDpos()
	roundChangeSetInfo, err := Dpos.RoundChangeSetInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("roundStateInfo is %#v \n", roundChangeSetInfo)
}

func TestBacklogs(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	Dpos := connection.System.NewDpos()
	backlogs, err := Dpos.Backlogs()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("roundStateInfo is %#v \n", backlogs)
}
