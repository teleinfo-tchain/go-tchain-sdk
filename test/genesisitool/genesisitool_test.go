package test

import (
	"github.com/bif/bif-sdk-go/genesistool"
	"testing"
)

func TestEncodeIstanbul(t *testing.T) {
	vanity := "0x00"
	validators := []string{"did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"}
	extraData, err := genesistool.EncodeIstanbul(vanity, validators)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	want := "0x0000000000000000000000000000000000000000000000000000000000000000f87fd9987366740ca125ea9c5e4a070dd2123351ecdcaac685d9dbf1b8620000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"

	if extraData != want {
		t.Logf("want extraData Istanbul encode is : %s , result is %s \n", want, extraData)
	}
}

func TestEncodeHotStuff(t *testing.T) {
	vanity := "0x00"
	validators := []string{"did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"}

	extraData, err := genesistool.EncodeHotStuff(vanity, validators)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	want := "0x0000000000000000000000000000000000000000000000000000000000000000f89a98000000000000000000000000000000000000000000000000d9987366740ca125ea9c5e4a070dd2123351ecdcaac685d9dbf1808080b8620000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

	if extraData != want {
		t.Logf("want extraData Istanbul encode is : %s , result is %s \n", want, extraData)
	}
}

func TestPeerIdGen(t *testing.T) {
	testNodePriStr := "bd0bc93d5bfd6f329f3b1fb61109fdb81290dfb9a78de491aa84276af4a713a2"
	peerId, err := genesistool.PeerIdGen(testNodePriStr)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	want := "16Uiu2HAmSgtVcHBHe79Ey3H3DHHxqbFBCFLL5UcEAUz8sBBxouui"

	if peerId != want {
		t.Logf("want peerId is : %s , result is %s \n", want, peerId)
	}

}

func TestRegEncode(t *testing.T) {
	regulators := []string{"did:bid:llj1:sfYVq8gWNHSFhwUtA5KcKCVMszR86Zgc"}
	regulatorsData := genesistool.RegEncode(regulators)
	want := "0x7366740ca125ea9c5e4a070dd2123351ecdcaac685d9dbf1"

	if regulatorsData != want {
		t.Logf("want regulators encode data is : %s , result is %s \n", want, regulatorsData)
	}
}
