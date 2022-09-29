package main

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils/math"
	"os"
	"strconv"
	"time"
)

func main() {
	address := os.Args[1] + ":" + os.Args[2]

	fmt.Println("http://" + address)

	var connection = bif.NewBif(providers.NewHTTPProvider(address, 10, false))

	var lastBlock *dto.BlockNoDetails
	var lastLastBlock *dto.BlockNoDetails
	for true {
		number, _ := connection.Core.GetBlockNumber()
		block, _ := connection.Core.GetBlockByNumber(number.String(), false)
		b := block.(*dto.BlockNoDetails)
		txLen := len(b.Transactions)
		if lastBlock == nil {
			lastBlock = b
			lastLastBlock = b
		}
		if b.Number.Uint64() == lastBlock.Number.Uint64() {
			lastBlock = lastLastBlock
		}
		sub := b.Timestamp - lastBlock.Timestamp
		if sub == 0 {
			sub = math.MaxUint64
		}
		fmt.Printf("number:%d, timestamp:%d, txCount:%d, tps:%f\n", number.Uint64(), b.Timestamp, txLen, float64(txLen)/float64(sub))

		if b.Number.Uint64() != lastBlock.Number.Uint64() {
			lastLastBlock = lastBlock
			lastBlock = b
		}
		second, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(time.Duration(second) * time.Second)
	}
}
