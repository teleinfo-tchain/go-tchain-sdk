package compiler

import (
	"fmt"
	"github.com/bif/bif-sdk-go"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	solcDir  = path.Join(bif.GetCurrentAbPath(), "compiler", "tmp")
	solFile  = path.Join(bif.GetCurrentAbPath(), "compiler", "contract")
)


func TestSolidityInline(t *testing.T) {
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows"{
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	}else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solc := NewSolidityCompiler(solcDir)

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
			for i := range output.Contracts {
				result[strings.TrimPrefix(i, "<stdin>:")] = struct{}{}
			}

			// only one source file
			assert.Len(t, output.Sources, 1)

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
	sysType := runtime.GOOS
	version := "v0.5.5"
	if sysType == "windows"{
		solcDir = path.Join(solcDir, version, sysType, "solc.exe")
	}else {
		solcDir = path.Join(solcDir, version, sysType, "solc-static-linux")
	}

	solc := NewSolidityCompiler(solcDir)

	files := []string{
		path.Join(solFile, "ballot.sol"),
		path.Join(solFile, "simple_auction.sol"),
	}
	output, err := solc.Compile(files...)
	if err != nil {
		t.Fatal(err)
	}
	if len(output.Contracts) != 2 {
		t.Fatal("two expected")
	}

	for k, v := range output.Contracts {
		fmt.Println(k)
		fmt.Println(v)
	}


}
