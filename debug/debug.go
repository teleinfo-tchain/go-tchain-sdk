package debug

import (
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/providers"
)

type Debug struct {
	provider providers.ProviderInterface
}

func NewDebug(provider providers.ProviderInterface) *Debug {
	debug := new(Debug)
	debug.provider = provider
	return debug
}

/*
  DumpBlock:
   	EN - retrieves the entire state of the database at a given block
 	CN - 根据指定的区块检索数据库的整个状态
  Params:
  	- blockNumber, string, options are:
	 (1) HEX String - an integer block number
	 (2) String "latest" - for the latest mined block
	 (3) String "pending" - for the pending state/transactions

  Returns:
  	- *dto.Dump
 	- error

  Call permissions: Anyone
*/
func (debug *Debug) DumpBlock(blockNumber string) (*dto.Dump, error) {
	params := make([]string, 1)
	params[0] = blockNumber

	pointer := &dto.DebugRequestResult{}
	err := debug.provider.SendRequest(pointer, "debug_dumpBlock", params)
	if err != nil {
		return nil, err
	}

	return pointer.ToDumpBlock()

}

/*
  GetBlockRlp:
   	EN - retrieves the RLP encoded for of a single block
 	CN - 获取指定区块的RLP编码
  Params:
  	- number,uint64

  Returns:
  	- string
 	- error

  Call permissions: Anyone
*/
func (debug *Debug) GetBlockRlp(number uint64) (string, error) {
	params := make([]uint64, 1)
	params[0] = number

	pointer := &dto.DebugRequestResult{}
	err := debug.provider.SendRequest(pointer, "debug_getBlockRlp", params)
	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

/*
  PrintBlock:
   	EN - retrieves a block and returns its pretty printed form
 	CN - 根据指定的区块号返回区块的表示格式
  Params:
   	- number,uint64

  Returns:
   	- string
 	- error

  Call permissions: Anyone
*/
func (debug *Debug) PrintBlock(number uint64) (string, error) {
	params := make([]uint64, 1)
	params[0] = number

	pointer := &dto.DebugRequestResult{}
	err := debug.provider.SendRequest(pointer, "debug_printBlock", params)
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}