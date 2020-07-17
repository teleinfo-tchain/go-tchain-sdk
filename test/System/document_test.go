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
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	initializationDDOHash, err := doc.InitializationDDO(common.StringToAddress(coinBase), 2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(initializationDDOHash)
}

func TestSetBidName(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	setBidHash, err := doc.SetBidName(common.StringToAddress(coinBase), "testTele")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(setBidHash)
}

func TestGetDocument(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	document, err := doc.GetDocument(0, "testTele")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(document)
}

func TestAddPublicKey(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	AddPublicKeyHash, err := doc.AddPublicKey(common.StringToAddress(coinBase), "test", "1", "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(AddPublicKeyHash)
}

func TestDeletePublicKey(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	deletePublicKeyHash, err := doc.DeletePublicKey(common.StringToAddress(coinBase), "123")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(deletePublicKeyHash)
}

func TestAddProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	addProofHash, err := doc.AddProof(common.StringToAddress(coinBase), "123", "testProof")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(addProofHash)
}

func TestDeleteProof(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	deleteProofHash, err := doc.DeleteProof(common.StringToAddress(coinBase), "testProof")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(deleteProofHash)
}

func TestAddAttribute(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	addAttributeHash, err := doc.AddAttribute(common.StringToAddress(coinBase), "testAttr", "attr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(addAttributeHash)
}

func TestDeleteAttribute(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	deleteAttributeHash, err := doc.DeleteAttribute(common.StringToAddress(coinBase), "testAttr", "attr")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(deleteAttributeHash)
}

func TestEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	enableHash, err := doc.Enable(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(enableHash)
}

func TestDisable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	coinBase, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	disableHash, err := doc.Disable(common.StringToAddress(coinBase))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(disableHash)
}

func TestIsEnable(t *testing.T) {
	var connection = bif.NewBif(providers.NewHTTPProvider(resources.IP+":"+resources.Port, 10, false))
	_, err := connection.Core.GetCoinBase()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	doc := connection.System.NewDoc()

	isEnable, err := doc.IsEnable(1, "test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(isEnable)
}
