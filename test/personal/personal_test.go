package personal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bif/bif-sdk-go"
	Abi "github.com/bif/bif-sdk-go/abi"
	"github.com/bif/bif-sdk-go/core/block"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"github.com/bif/bif-sdk-go/utils"
	"github.com/bif/bif-sdk-go/utils/types"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
)

type personal struct {
	provider providers.ProviderInterface
}

type wallet struct {
	Status     string    `json:"status"`
	WalletInfo string    `json:"walletInfo"`
	Accounts   []account `json:"wallet"`
}

type account struct {
	Address string `json:"address"`
}

func (acc *account) UnmarshalJSON(data []byte) error {
	type Alias account
	temp := &struct {
		*Alias
	}{
		Alias: (*Alias)(acc),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	return nil
}

func checkResponse(pointer *dto.RequestResult) error {

	if pointer.Error != nil {
		return errors.New(pointer.Error.Message)
	}

	if pointer.Result == nil {
		return dto.EMPTYRESPONSE
	}

	return nil

}

func newPersonal(provider providers.ProviderInterface) *personal {
	personal := new(personal)
	personal.provider = provider
	return personal
}

func (p *personal) personalListAccounts() ([]string, error) {
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(pointer, "personal_listAccounts", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToStringArray()
}

func (p *personal) personalListWallets() ([]*wallet, error) {
	pointer := &dto.RequestResult{}

	err := p.provider.SendRequest(pointer, "personal_listWallets", nil)

	if err != nil {
		return nil, err
	}

	if err := checkResponse(pointer); err != nil {
		return nil, err
	}

	result := (pointer).Result.([]interface{})
	wallets := make([]*wallet, len(result))
	for k, v := range result {
		result := v.(map[string]interface{})
		if len(result) == 0 {
			return nil, dto.EMPTYRESPONSE
		}
		info := &wallet{}

		marshal, err := json.Marshal(result)
		if err != nil {
			return nil, dto.EMPTYRESPONSE
		}

		err = json.Unmarshal(marshal, info)
		if err != nil {
			return nil, dto.EMPTYRESPONSE
		}

		wallets[k] = info
	}
	return wallets, err
}

func (p *personal) personalNewAccount(auth string, cryptoTypeNum int) (string, error) {
	pointer := &dto.RequestResult{}

	params := make([]interface{}, 2)
	params[0] = auth
	params[1] = cryptoTypeNum

	err := p.provider.SendRequest(pointer, "personal_newAccount", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

func (p *personal) personalImportRawKey(privateKey, password string, cryptoTypeNum uint64) (string, error) {
	pointer := &dto.RequestResult{}

	params := make([]interface{}, 3)
	params[0] = privateKey
	params[1] = password
	params[2] = cryptoTypeNum

	err := p.provider.SendRequest(pointer, "personal_importRawKey", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

func (p *personal) personalLockAccount(address string) (bool, error) {
	pointer := &dto.RequestResult{}

	params := make([]string, 1)
	params[0] = address

	err := p.provider.SendRequest(pointer, "personal_lockAccount", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

func (p *personal) personalUnLockAccount(address, password string, duration uint64) (bool, error) {
	pointer := &dto.RequestResult{}

	params := make([]interface{}, 3)
	params[0] = address
	params[1] = password
	params[2] = duration

	err := p.provider.SendRequest(pointer, "personal_unlockAccount", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

func (p *personal) personalExportToFile(address, pathDir string) (bool, error) {
	pointer := &dto.RequestResult{}

	params := make([]string, 2)
	params[0] = address
	params[1] = pathDir

	err := p.provider.SendRequest(pointer, "personal_exportToFile", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

func (p *personal) personalBatchExportToFiles(pathDir string) (bool, error) {
	pointer := &dto.RequestResult{}

	params := make([]string, 1)
	params[0] = pathDir

	err := p.provider.SendRequest(pointer, "personal_batchExportToFiles", params)
	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()
}

func (p *personal) personalImportByFile(fileName, pathDir string) (string, error) {
	pointer := &dto.RequestResult{}

	params := make([]string, 2)
	params[0] = fileName
	params[1] = pathDir

	err := p.provider.SendRequest(pointer, "personal_importByFile", params)
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

func (p *personal) personaBatchImportByFiles(pathDir string) (string, error) {
	pointer := &dto.RequestResult{}

	params := make([]string, 1)
	params[0] = pathDir

	err := p.provider.SendRequest(pointer, "personal_batchImportByFiles", params)
	if err != nil {
		return "", err
	}

	return pointer.ToString()
}

// 节点管理的账户列表
// {"jsonrpc":"2.0","method":"personal_listAccounts","params":[],"id":67}
func TestPersonalListAccounts(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	res, err := connection.personalListAccounts()
	if err != nil {
		t.Logf("errr is %s ", err)
	}
	t.Log(res)
}

// 节点管理的钱包账户
// {"jsonrpc":"2.0","method":"personal_listWallets","params":[],"id":67}
func TestPersonalListWallets(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	res, err := connection.personalListWallets()

	if err != nil {
		t.Logf("errr is %s ", err)
	}
	for _, v := range res {
		t.Logf("%#v \n ", v)
	}
}

// 生成账户地址
// {"jsonrpc":"2.0","method":"personal_newAccount","params":["teleinfo",0],"id":67}
func TestPersonaNewAccount(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	auth := "node"
	cryptoType := 0

	res, err := connection.personalNewAccount(auth, cryptoType)
	if err != nil {
		t.Logf("errr is %s ", err)
	}
	fmt.Println(res)
}

// 导入节点账户地址
// {"jsonrpc":"2.0","method":"personal_importRawKey","params":[" ","teleinfo", 0],"id":67}
func TestPersonalImportRawKey(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	for _, test := range []struct {
		privateKey       string
		password         string
		cryptoTypeNumber uint64
	}{
		{"41e46e858ea707453d8fc553805772165a4f66e6e18ca38220daa157534e0c0e", "teleinfo", 1},
	} {
		res, err := connection.personalImportRawKey(test.privateKey, test.password, test.cryptoTypeNumber)
		if err != nil {
			t.Logf("error is %s ", err)
			t.Fail()
		}
		t.Log(res)
	}

}

// lock账户
// {"jsonrpc":"2.0","method":"personal_lockAccount","params":["did:bid:c117c1794fc7a27bd301ae52"],"id":67}
func TestPersonalLockAccount(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	for _, test := range []struct {
		address string
	}{
		{"did:bid:ZFT4Y87Xdg83GEDDbiNknHLWs3Hfq58"},
	} {
		res, err := connection.personalLockAccount(test.address)
		if err != nil {
			t.Logf("error is %s ", err)
		}
		fmt.Println(res)
	}
}

// unlock账户
// {"jsonrpc":"2.0","method":"personal_unlockAccount","params":["did:bid:c117c1794fc7a27bd301ae52","teleinfo",0],"id":67}
func TestPersonalUnLockAccount(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	for _, test := range []struct {
		address  string
		password string
		duration uint64
	}{
		{"did:bid:ZFT4Y87Xdg83GEDDbiNknHLWs3Hfq58", "node", 50},
	} {
		res, err := connection.personalUnLockAccount(test.address, test.password, test.duration)
		if err != nil {
			t.Logf("error is %s ", err)
		}
		t.Log(res)
	}
}

func TestPersonalExportToFile(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	for _, test := range []struct {
		address string
	}{
		{"did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"},
	} {
		pathDir := "keyStoreTest"
		res, err := connection.personalExportToFile(test.address, pathDir)
		if err != nil {
			t.Logf("error is %s ", err)
		}
		t.Log(res)
	}
}

func TestPersonalBatchExportToFiles(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	pathDir := "keyStoreTest"
	res, err := connection.personalBatchExportToFiles(pathDir)
	if err != nil {
		t.Logf("error is %s ", err)
	}
	t.Log(res)
}

func TestPersonalImportByFile(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	for _, test := range []struct {
		fileName string
	}{
		{"did-bid-ZFT2PCJNq86pLGrsmHzSYewTxrGHPKq"},
	} {
		pathDir := "keyStoreTest"
		address, err := connection.personalImportByFile(test.fileName, pathDir)
		if err != nil {
			t.Logf("error is %s ", err)
		}
		fmt.Println(address)
	}
}

func TestPersonaBatchImportByFiles(t *testing.T) {
	var connection = newPersonal(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))

	pathDir := "keyStoreTest"
	res, err := connection.personaBatchImportByFiles(pathDir)
	if err != nil {
		t.Logf("error is %s ", err)
	}
	t.Log(res)
	// 	did:bid:ZFTExsofuzLndNueEhtiStL8QAUVcU
}

// 测试转账
func TestCoreSendTransaction(t *testing.T) {

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(generator)
	generator = "did:bid:EFTVcqqKyFR17jfPxqwEtpmRpbkvSs"
	fmt.Println(connection.Core.GetBlockNumber())

	balanceFrom, err := connection.Core.GetBalance(generator, block.LATEST)
	if err == nil {
		balBif, _ := utils.FromWei(balanceFrom)
		t.Log("fromAddress balance is ", balBif)
	}

	toAddress := "did:bid:ZFT4Y87Xdg83GEDDbiNknHLWs3Hfq58"
	balanceTo, err := connection.Core.GetBalance(toAddress, block.LATEST)
	if err == nil {
		balBif, _ := utils.FromWei(balanceTo)
		t.Log("toAddress balance is ", balBif)
	}

	transaction := new(dto.TransactionParameters)
	transaction.Sender = generator
	transaction.Recipient = toAddress
	transaction.Amount = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(1e17))
	transaction.GasLimit = 40000
	transaction.Payload = "Transfer test"

	txID, err := connection.Core.SendTransaction(transaction)

	// Wait for a block
	time.Sleep(time.Second)

	if err != nil {
		t.Errorf("Failed SendTransaction")
		t.Error(err)
		t.FailNow()
	}

	// if success, get transaction hash
	t.Log(txID)
}

// 测试部署合约(只是为了测试部署合约，实际使用contract中的Deploy)
func TestCoreSendTransactionDeployContract(t *testing.T) {
	content, err := ioutil.ReadFile("../resources/simple-contract.json")

	type Contract struct {
		Abi      string `json:"abi"`
		ByteCode string `json:"byteCode"`
	}

	var unmarshalResponse Contract

	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP00+":"+strconv.FormatUint(resources.Port, 10), 10, false))
	byteCode := unmarshalResponse.ByteCode
	parsedAbi, err := Abi.JSON(strings.NewReader(unmarshalResponse.Abi))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	inputEncode, err := parsedAbi.Pack("", big.NewInt(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction := new(dto.TransactionParameters)

	generator, err := connection.Core.GetGenerator()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	transaction.Sender = generator
	transaction.Payload = types.ComplexString(byteCode) + types.ComplexString(utils.Bytes2Hex(inputEncode))
	// estimate the gas required to deploy the contract
	gas, err := connection.Core.EstimateGas(transaction)
	if err != nil {
		t.Errorf("Failed EstimateGas")
		t.Error(err)
		t.FailNow()
	}
	t.Log("Estimate gas is ", gas)

	// transaction.GasPrice = big.NewInt(1000000)
	transaction.GasLimit = gas.Uint64()
	txHash, err := connection.Core.SendTransaction(transaction)

	if err != nil {
		t.Errorf("Failed Deploy Contract")
		t.Error(err)
		t.FailNow()
	}
	t.Log("transaction hash is ", txHash)

	var receipt *dto.TransactionReceipt

	for receipt == nil {
		time.Sleep(time.Second)
		receipt, err = connection.Core.GetTransactionReceipt(txHash)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// did:bid:ace45606ce7b19c7da1143cb
	t.Log("contract Address is ", receipt.ContractAddress)

}
