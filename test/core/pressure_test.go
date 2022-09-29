package test

import (
	"fmt"
	"github.com/tchain/go-tchain-sdk"
	"github.com/tchain/go-tchain-sdk/account"
	"github.com/tchain/go-tchain-sdk/core/block"
	"github.com/tchain/go-tchain-sdk/crypto"
	"github.com/tchain/go-tchain-sdk/crypto/config"
	"github.com/tchain/go-tchain-sdk/dto"
	"github.com/tchain/go-tchain-sdk/providers"
	"github.com/tchain/go-tchain-sdk/utils"
	"github.com/tchain/go-tchain-sdk/utils/hexutil"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	ip               = "172.17.6.51"
	port      uint64 = 44002
	chainId   uint64 = 222
	chainCode        = "gjdw"

	privateKey, cryptoType = "37d4e59ab4b95e8e4c81bef0ea83b5deced9a484555ea6fa1c8c8b7c25e7a83b", config.SECP256K1
)

var (
	tos = []string{
		"269d17aa16a3a61024d104bd1740284e5fe6f4106c8c6c9045d1fefb0e2ad168",
		"628c70ae613a0444e3f0e1a187778749a281dd6dad0d23aedbe8d20fae5e3ec3",
		"8c02e470b4d821b159e80a54d4341e12b0ebc9a7a5973fcfad5caacd000a2d8e",
		"a991b28a1b52f0ec8171c112011de3a61493a4efb53f1558318f0b4b19955907",
		"f7fc175d860da83b2879a5b272d9a8f6b1fc943b93b3222a9f2addad5e6049ad",
		"52318d09f5d9e804c66f4674d253da62127add5dc9a1a3b7823358f4b2a524fd",
		"e684393be5302dba261e0cbccb5fdc1f0bef02d5080bf05c4b022acc20dc4836",
		"187ea33825a0e84df18b18980f6b59c742a38df30d03286f89e273a1516deff8",
		"fb63a20ab3998892e3a83e601a125fb562df40df41276d91b39cf1d489862cd5",
		"e8003f36d40ab1961894a9d4552b14a0c05c02c1064b8402c1e45adfc2826ec1",
		"c2a09101aa95fb4203456f2035288882b338525bb742c9e209d30688eb4fa698",
		"fc7c37be767b779ce5ef088ed5a1fa74114a59e751444bd08de79eef7a7471fc",
		"e776ad0cc4c0352e9e7d9b33f1b771189e1cc2c73532ca05b3fc551d4f2a0b3b",
		"664d58ace1c9e6cb0a49595f63b6e4984979f062cec106b5393008cd6f57ab09",
		"6fce4b120d476038878fd3799b711cecc2fcd8491054455ca987345cfc3c0d31",
		"fec38a08ebe7488a813f96224200abb0a0c119aad8adb5d929339745e0f60ec6",
		"fc88253f89f55991ce32be18fceb12b8a85083c183d202356adc9f3d65155efa",
		"07193496b311e29cdf5a2082757091aecfd51ffcbea5641ac112e691c88290b8",
		"0c12065962e875bea2f8450e38f6c193346c4e480c67498596eff372c24e186f",
		"8652e7c58cef21618b19ba562f21e0dc573f8afa420b5d1024c58ef5fefb601d",
		"b6ec8fe59cbba93ba10ac2faeeceded0159bb956b8405f669bec41695a34a0fd",
		"9554e6d5541309f090a59e8bde56eb0d4f3b11bb897e77a0bd153fede0590557",
		"4d9d505c653a27f3f680706bd96906c949412b70c7cbb336246c49c85b36e11e",
		"5f4c0239f8e5cbca5265f861c9bb7d3ad24ffa5427313e47086f9ef14f86d9e1",
		"d38e0acf0c6b74d87c904b6ea4e0846d8c0e6855cce2cb8c8bc5eecff9ab7df1",
		"d97884a40ba5e2ceee9c071e9275c93ae14c00fce18e3f4aaa8aa6c974def57d",
		"dcbe638a2ed0444a6e6dc71eba7439e64847dc005d9f3a6e8e9e30a1d1061f4c",
		"59503e38c7ac3b515fcd9a809afe2fb35240dae070e7bf8bfa648fcf017662e2",
		"1d784f57ff3348d760c14663c05e8be118e2c70c490fea2509fffeb88cb106bd",
		"c9910cac7219264a04d4f363768127afc25061289609194e8e0e87c081a1b7e6",
		"ecdd9b3167aa59de19ef0839f7b571320667ef9416e9c60a59602a50fb480164",
		"2af650a93ca60316bebc189475c074085d26c56a1fdd8785ca985e46c33af7c0",
		"1190d1106d6d22a84fd48bb4a36a61230e1f90f8ed93d9bf9ce33ee04ea77dc4",
		"369f832984bc2ae7c9b9965042ff4ebd3a26e32e20734304b967324ca6791c39",
		"1ba4af72694d4717eac4e98988d2d6a8148c889f0b767d8c67575653e7c6f079",
		"b803995a1ff56fb208f7bf3f3bffb7d138ac8ac75e22e15dcd84b971a10a4668",
		"4d8933da974f0694f91b5a3195e21f61900fe5544c281c73a730dd23cd7cb66e",
		"fddabb157bce9e3dd62a534d6fe97848c9caa0744dc74547126ee1f184461ed4",
		"3593ebfcb6e20029f3083f5dccef8e37218bab0bcfee8797451d9e114966e3b2",
		"ae16df6a7b4d79f35b6834733812734c72c2c14170eb19d9c36fdc8f9fdd4b75",
		"90b98cd3cc6a838640cadd657af0fc494f511cfacda2785d83d6c057248e452b",
		"24c326c87264753ac6b5b669bf8b2b2c9ce01be423c35b3ae6e481396ccb464a",
		"8fd2da73dbaae9b9500c2e4cae316e38e1faf28276bda500324f9b65da3e82de",
		"ef94ea75403435ae186a059864385e9cc702135840edb8c2910bd82a776d7798",
		"427e7a4844c65872fb667e6be71d414796df9563352bd540943e9fae94c3b76d",
		"aa39d7c11958ca8ed274c1fede2dc2937c875c08e77e8d7080231a8dfee293d2",
		"419e311a30951f1ccc797b402aae6eeabed0349bffee44718060a6c2a63c536d",
		"e76abcac190bf198cdc0fdc1a8615eed5c739997a135a2710adc0ed06c72ed79",
		"def627a0c4c13b43f4af5b05d5839ca1457c07e494473d89da730592850c5036",
		"0cd01e50a96af52ca3b131d5dc6afd291e8425b450ee256b7e949bb77138c9b1",
	}
)

func TestGenerate(t *testing.T) {
	for i := 0; i < 50; i++ {
		generateKey, _ := crypto.GenerateKey(config.SECP256K1)
		fmt.Println(hexutil.Encode(generateKey.D.Bytes()))
	}
}

func TestMonitor(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(ip+":"+strconv.FormatUint(port, 10), 10, false))

	for true {
		number, _ := connection.Core.GetBlockNumber()
		block, _ := connection.Core.GetBlockByNumber(number.String(), false)
		b := block.(*dto.BlockNoDetails)
		fmt.Printf("number:%d, timestamp:%d, txLength:%d\n", number.Uint64(), b.Timestamp, len(b.Transactions))
		time.Sleep(4 * time.Second)
	}
}

func TestPressure(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(ip+":"+strconv.FormatUint(port, 10), 10, false))

	generateKey, _ := crypto.GenerateKey(config.SECP256K1)

	recipient := crypto.PubkeyToAddress(generateKey.PublicKey)

	txObjs := make([]*account.SignTxParams, len(tos))
	for i, to := range tos {
		txObj := tx(connection, to, cryptoType, recipient, big.NewInt(0))

		txObjs[i] = txObj

		res, err := account.SignTransaction(txObj, to, cryptoType == config.SM2)
		if err != nil {
			fmt.Println(err)
		}

		txID, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(txID)
	}

	for j := 0; j < 100; j++ {
		for i, to := range tos {
			txObjs[i].Nonce.Add(txObjs[i].Nonce, utils.Big1)

			res, err := account.SignTransaction(txObjs[i], to, cryptoType == config.SM2)
			if err != nil {
				fmt.Println(err)
			}

			txID, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(txID)
		}
	}
}

func TestInit(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(ip+":"+strconv.FormatUint(port, 10), 10, false))

	for _, to := range tos {
		txObj := tx(connection, privateKey, cryptoType, utils.StringToAddress(to), big.NewInt(1000000000000000000))

		res, err := account.SignTransaction(txObj, privateKey, cryptoType == config.SM2)
		if err != nil {
			fmt.Println(err)
		}

		txID, err := connection.Core.SendRawTransaction(hexutil.Encode(res.Raw))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(txID)

		time.Sleep(5 * time.Second)

		receipt, err := connection.Core.GetTransactionReceipt(txID)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(receipt)
		if receipt != nil {
			fmt.Println(receipt.Status)
		}

		from, err := key(to, config.SECP256K1)
		if err != nil {
			fmt.Println(err)
		}

		number, _ := connection.Core.GetBlockNumber()
		balance, err := connection.Core.GetBalance(from.String(chainCode), hexutil.EncodeUint64(number.Uint64()))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(balance)
	}
}

func TestBalance(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(ip+":"+strconv.FormatUint(port, 10), 10, false))

	number, _ := connection.Core.GetBlockNumber()
	fmt.Println(number.String())

	from, err := key(privateKey, config.SECP256K1)
	if err != nil {
		fmt.Println(err)
	}

	balance, err := connection.Core.GetBalance(from.String(chainCode), hexutil.EncodeUint64(number.Uint64()))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(balance)

	for _, to := range tos {
		from, err := key(to, config.SECP256K1)
		if err != nil {
			fmt.Println(err)
		}

		balance, err := connection.Core.GetBalance(from.String(chainCode), hexutil.EncodeUint64(number.Uint64()))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(balance)
	}
}

func tx(connection *bif.Bif, from string, ct config.CryptoType, recipient utils.Address, amount *big.Int) *account.SignTxParams {
	addr, _ := key(from, ct)

	nonce, err := connection.Core.GetTransactionCount(addr.String(chainCode), block.PENDING)
	if err != nil {
		fmt.Println(err)
	}

	return &account.SignTxParams{
		Recipient: &recipient,
		Nonce:     nonce,
		GasPrice:  utils.Big0,
		GasLimit:  8000000000,
		Amount:    amount,
		Payload:   nil,
		ChainId:   chainId,
	}

}

func key(privateKey string, cryptoType config.CryptoType) (utils.Address, error) {
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = privateKey[2:]
	}
	privKey, err := crypto.HexToECDSA(privateKey, cryptoType)
	if err != nil {
		fmt.Println(err)
		return utils.Address{}, err
	}
	return crypto.PubkeyToAddress(privKey.PublicKey), nil
}
