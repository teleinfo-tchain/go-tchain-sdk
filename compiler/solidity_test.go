package compiler

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var (
	solcDir  = "D:\\Go\\src\\github.com\\umbracle\\go-web3\\compiler\\tmp"
	solcPath = solcDir + "\\solc.exe"
)


func TestSolidityInline(t *testing.T) {
	solc := NewSolidityCompiler(solcPath).(*Solidity)

	cases := []struct {
		code      string
		contracts []string
	}{
		{
			`
		pragma solidity >0.0.0;
		contract foo{}
			`,
			[]string{
				"foo",
			},
		},
		{
			`
		pragma solidity >0.0.0;
		contract foo{}
		contract bar{}
			`,
			[]string{
				"bar",
				"foo",
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			output, err := solc.CompileCode(c.code)
			if err != nil {
				t.Fatal(err)
			}

			result := map[string]struct{}{}
			for i := range output {
				result[strings.TrimPrefix(i, "<stdin>:")] = struct{}{}
			}

			expected := map[string]struct{}{}
			for _, i := range c.contracts {
				expected[i] = struct{}{}
			}

			if !reflect.DeepEqual(result, expected) {
				t.Fatal("bad")
			}
		})
	}
}

func TestSolidity(t *testing.T) {
	solc := NewSolidityCompiler(solcPath)
	demoPath := "D:\\Go\\src\\github.com\\bif\\bif-sdk-go\\test\\resources\\simple-token.sol"
	output, err := solc.Compile(demoPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(output) != 2 {
		t.Fatal("two expected")
	}
}

func existsSolidity(t *testing.T, path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		t.Fatal(err)
	}

	cmd := exec.Command(path, "--version")
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		t.Fatalf("solidity version failed: %s", string(stderr.Bytes()))
	}
	if len(stdout.Bytes()) == 0 {
		t.Fatal("empty output")
	}
	return true
}

func TestDownloadSolidityCompiler(t *testing.T) {
	dst1, err := ioutil.TempDir("/tmp", "go-web3-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dst1)

	if err := DownloadSolidity("0.5.5", dst1, true); err != nil {
		t.Fatal(err)
	}
	if existsSolidity(t, filepath.Join(dst1, "solidity")) {
		t.Fatal("it should not exist")
	}
	if !existsSolidity(t, filepath.Join(dst1, "solidity-0.5.5")) {
		t.Fatal("it should exist")
	}

	dst2, err := ioutil.TempDir("/tmp", "go-web3-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dst2)

	if err := DownloadSolidity("0.5.5", dst2, false); err != nil {
		t.Fatal(err)
	}
	if !existsSolidity(t, filepath.Join(dst2, "solidity")) {
		t.Fatal("it should exist")
	}
	if existsSolidity(t, filepath.Join(dst2, "solidity-0.5.5")) {
		t.Fatal("it should not exist")
	}
}
