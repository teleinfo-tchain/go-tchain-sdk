package test

import (
	"fmt"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/utils"
	"math/big"
	"testing"
)

func TestSignTransactionWithKey(t *testing.T) {
	from := common.StringToAddress("did:bid:73f6a70d05af2141dd4ad995")
	to := common.StringToAddress("did:bid:73890cf407f6c883e9a42735")
	tx := &utils.Txdata{
		AccountNonce: 1,
		Price:        big.NewInt(10),
		GasLimit:     2000000,
		Sender:       &from,
		Recipient:    &to,
		Amount:       big.NewInt(500000),
		Payload:      nil,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
		T:            big.NewInt(0),
	}

	privKey := "eea9354b98fd51d7b962cb4c7e61d691e4d540951e1cf277dd72f2a37544c1da"
	signTx, err := utils.SignTransactionWithKey(tx, privKey, 0, 445)
	if err != nil {
		fmt.Println("signtransactionWithKey err:", err)
		return
	}
	fmt.Println("signTx:", signTx)
}
