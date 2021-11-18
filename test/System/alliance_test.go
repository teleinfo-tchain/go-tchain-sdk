package System

import (
	"github.com/bif/bif-sdk-go"
	"github.com/bif/bif-sdk-go/dto"
	"github.com/bif/bif-sdk-go/test/resources"
	"path"
	"testing"
	"time"
)

func TestRegisterDirector(t *testing.T) {
	file := path.Join(bif.GetCurrentAbPath(), "test", "resources", "keystore", resources.TestAddressRegulatoryFile)
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	registerDirector := new(dto.AllianceInfo)
	registerDirector.Id = resources.RegisterAllianceTwo
	registerDirector.PublicKey = resources.RegisterAllianceTwoPubKey
	registerDirector.CompanyName = "teleInfo"
	registerDirector.CompanyCode = "110112"

	registerDirectorHash, err := ali.RegisterDirector(sigPara, registerDirector)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(registerDirectorHash, err)

	time.Sleep(8*time.Second)

	log, err := con.System.SystemLogDecode(registerDirectorHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestUpgradeDirector(t *testing.T) {
	file := path.Join(bif.GetCurrentAbPath(), "test", "resources", "keystore", resources.TestAddressRegulatoryFile)
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	transactionHash, err := ali.UpgradeDirector(sigPara, resources.RegisterAllianceTwo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)

	time.Sleep(8*time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestRevoke(t *testing.T) {
	file := path.Join(bif.GetCurrentAbPath(), "test", "resources", "keystore", resources.TestAddressRegulatoryFile)
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	revokeReason := "不合法"

	transactionHash, err := ali.Revoke(sigPara, resources.RegisterAllianceOne, revokeReason)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)

	time.Sleep(8*time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestSetWeights(t *testing.T) {
	file := path.Join(bif.GetCurrentAbPath(), "test", "resources", "keystore", resources.TestAddressRegulatoryFile)
	con, sigPara, err := connectWithSig(resources.TestAddressRegulatory, file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	transactionHash, err := ali.SetWeights(sigPara, 2, 3, 4)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(transactionHash, err)

	time.Sleep(8*time.Second)

	log, err := con.System.SystemLogDecode(transactionHash)

	if err != nil {
		t.Errorf("err log : %v ", err)
		t.FailNow()
	}

	if !log.Status {
		t.Errorf("err, method is %s , err is %s ", log.Method, log.Result)
	}
}

func TestAllDirectors(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	directors, err := ali.AllDirectors()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range directors{
		t.Logf("directors is %+v \n", v)
	}
}

func TestAllVices(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	vices, err := ali.AllVices()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range vices{
		t.Logf("vices is %+v \n", v)
	}
}

func TestAllAlliance(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	alliances, err := ali.AllAlliances()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range alliances{
		t.Logf("alliances is %+v \n", v)
	}
}

func TestGetAlliance(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	alliance, err := ali.GetAlliance(resources.RegisterAllianceOne)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("aliance is %+v \n", alliance)
}

func TestGetWeights(t *testing.T) {
	con, err := connectBif()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ali := con.System.NewAlliance()

	weights, err := ali.GetWeights()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("weights is %v \n", weights)
}
