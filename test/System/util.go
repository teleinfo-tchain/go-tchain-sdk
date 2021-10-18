package System

import (
	"errors"
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/system"
	"github.com/bif/bif-sdk-go/test/resources"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
	sysTxParams.Gas = 2000000
	sysTxParams.Nonce = nonce.Uint64()
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

// GetCurrentAbPath 最终方案-全兼容  获取项目路径
func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()

	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
