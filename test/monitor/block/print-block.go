package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/hexutil"
)

func main() {
	var connection = bif.NewBif(providers.NewHTTPProvider(os.Args[1], 10, false))

	for {
		block, err := connection.Core.GetBlockByNumber("latest", false)

		if err != nil {
			fmt.Println(err)
		}
		b := block.(*dto.BlockNoDetails)
		printInfo(b)
		period, err := strconv.ParseInt(os.Args[2], 10, 64)
		time.Sleep(time.Duration(period) * time.Second)
	}
}

func printInfo(b *dto.BlockNoDetails) {
	extra, e := hexutil.Decode(b.ExtraData)
	if e != nil {
		fmt.Println(e)
	}
	istanbulExtra, e := utils.ExtractIstanbulExtra(extra)
	if e != nil {
		fmt.Println(e)
	}

	fmt.Printf("%-10d %-20s %-50s || %2d\t",
		b.Number,
		time.Unix(int64(b.Timestamp), 0).Format("2006-01-02 15:04:05"),
		b.Generator,
		len(istanbulExtra.Validators),
	)
	for _, validator := range istanbulExtra.Validators {
		fmt.Printf("%s,", validator.String(""))
	}
	fmt.Println()
}
