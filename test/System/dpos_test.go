package System

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)

//missing trie node????
func TestGetValidators(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	blockNumber := big.NewInt(157532)
	Dpos := connection.System.NewDpos()
	validators, err := Dpos.GetValidators(blockNumber)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(validators)
}

//missing trie node ???? 最初测试可行，现在为什么报错？？？
func TestGetValidatorsAtHash(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	blockNumber := big.NewInt(157532)

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
	fmt.Println(validators)
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
