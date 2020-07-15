package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/common"
	"github.com/bif/bif-sdk-go/providers"
	"github.com/bif/bif-sdk-go/test/resources"
	"testing"
)

func TestInitializationDDO(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	initializationDDOHash, err := doc.InitializationDDO(common.StringToAddress(coinBase), 2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(initializationDDOHash)
}

func TestSetBidName(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	setBidHash, err := doc.SetBidName(common.StringToAddress(coinBase), "testTele")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(setBidHash)
}

// 如何查找？？？？
func TestFindDDOByType(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	DDO, err := doc.FindDDOByType(common.StringToAddress(coinBase), 0, "testTele")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(DDO)
}

func TestAddPublicKey(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	AddPublicKeyHash, err := doc.AddPublicKey(common.StringToAddress(coinBase), "test", "1","123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(AddPublicKeyHash)
}

func TestDeletePublicKey(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	deletePublicKeyHash, err := doc.DeletePublicKey(common.StringToAddress(coinBase), "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(deletePublicKeyHash)
}

func TestAddProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	addProofHash, err := doc.AddProof(common.StringToAddress(coinBase), "123", "testProof")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(addProofHash)
}

func TestDeleteProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	deleteProofHash, err := doc.DeleteProof(common.StringToAddress(coinBase), "testProof")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(deleteProofHash)
}

func TestAddAttribute(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	addAttributeHash, err := doc.AddAttribute(common.StringToAddress(coinBase), "testAttr", "attr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(addAttributeHash)
}

func TestDeleteAttribute(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	deleteAttributeHash, err := doc.DeleteAttribute(common.StringToAddress(coinBase), "testAttr", "attr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(deleteAttributeHash)
}

func TestEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	enableHash, err := doc.Enable(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(enableHash)
}

func TestDisable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	disableHash, err := doc.Disable(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(disableHash)
}

func TestIsEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinbase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc, err := connection.System.NewDoc()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	isEnable, err := doc.IsEnable(common.StringToAddress(coinBase), 1, "test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}
