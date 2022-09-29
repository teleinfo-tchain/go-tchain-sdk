package System

import (
	"errors"
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/system"
	"github.com/tchain/go-tchain-sdk/test/resources"
	"io/ioutil"
	"math/big"
	"strconv"
)

func connectWithSig(sigAddr string, singAddrFile string) (*bif.Bif, *system.SysTxParams, error) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
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
	sysTxParams.Password = resources.SystemPassword
	sysTxParams.KeyFileData = keyFileData
	sysTxParams.GasPrice = big.NewInt(45)
	sysTxParams.Gas = 200000
	sysTxParams.Nonce = nonce
	sysTxParams.ChainId = chainId

	return connection, sysTxParams, nil
}

func connectBif() (*bif.Bif, error) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	_, err := connection.Core.GetChainId()
	if err != nil {
		return nil, err
	}
	return connection, nil
}
