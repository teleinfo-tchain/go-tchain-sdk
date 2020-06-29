package dpos

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"math/big"
	"testing"
)


func TestGetValidatorsAtHash(t *testing.T){
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	blockNumber := big.NewInt(157532)

	block, err := connection.Core.GetBlockByNumber(blockNumber, false)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	//input block hash
	validators, err := connection.System.Dpos.GetValidatorsAtHash(block.(*dto.BlockNoDetails).Hash)
	//missing trie node ???? 最初测试可行，现在为什么报错？？？

	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(validators)
}

