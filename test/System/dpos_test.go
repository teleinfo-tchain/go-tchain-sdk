package System

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils/hexutil"
	"strconv"
	"testing"
)

func TestGetValidators(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	// 距离最新区块间隔数超过126会报错，因此测试最新区块
	blockNumber, err := connection.Core.GetBlockNumber()
	// fmt.Println("blockNumber is ", blockNumber)
	// blockNumberOld := new(big.Int)
	// fmt.Println("blockNumberOld is ", blockNumberOld.Sub(blockNumber, big.NewInt(127)))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	DPoS := connection.System.NewDPoS()
	validators, err := DPoS.GetValidators(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(validators)
}

func TestGetValidatorsAtHash(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	// 距离最新区块间隔数超过126会报错，因此测试最新区块
	blockNumber, err := connection.Core.GetBlockNumber()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	block, err := connection.Core.GetBlockByNumber(hexutil.EncodeBig(blockNumber), false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	DPoS := connection.System.NewDPoS()
	// input block hash
	validators, err := DPoS.GetValidatorsAtHash(block.(*dto.BlockNoDetails).Hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(validators)
}

func TestRoundStateInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	DPoS := connection.System.NewDPoS()
	roundStateInfo, err := DPoS.RoundStateInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("roundStateInfo is %#v \n", roundStateInfo)
}

func TestRoundChangeSetInfo(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	DPoS := connection.System.NewDPoS()
	roundChangeSetInfo, err := DPoS.RoundChangeSetInfo()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("roundStateInfo is %#v \n", roundChangeSetInfo)
}

func TestBacklogs(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	DPoS := connection.System.NewDPoS()
	backlogs, err := DPoS.Backlogs()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("roundStateInfo is %#v \n", backlogs)
}
