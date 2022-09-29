package debug

import (
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/test/resources"
	"strconv"
	"testing"
)

func TestDumpBlock(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	_, err := connection.Core.GetBlockNumber()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	dumpBlock, err := connection.Debug.DumpBlock("latest")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v \n", dumpBlock)
	// t.Logf("%#v \n", dumpBlock.Accounts)
}

func TestGetBlockRlp(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	_, err := connection.Core.GetBlockNumber()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	blockRlp, err := connection.Debug.GetBlockRlp(12)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(blockRlp)
}

func TestPrintBlock(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	_, err := connection.Core.GetBlockNumber()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	printBlock, err := connection.Debug.PrintBlock(12)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(printBlock)
}