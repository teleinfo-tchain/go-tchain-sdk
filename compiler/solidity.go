package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type solcOutput struct {
	Contracts map[string]*solcContract
	Version   string
}

type solcContract struct {
	BinRuntime string `json:"bin-runtime"`
	Bin        string
	Abi        string
}

// Solidity is the solidity compiler
type Solidity struct {
	path string
}

// NewSolidityCompiler instantiates a new solidity compiler
func NewSolidityCompiler(path string) Compiler {
	return &Solidity{path}
}

// CompileCode compiles a solidity code
func (s *Solidity) CompileCode(demoPath string) (map[string]*Artifact, error) {
	if demoPath == "" {
		return nil, fmt.Errorf("demoPath is empty and files not exist")
	}
	artifacts, err := s.compileImpl(demoPath)
	if err != nil {
		return nil, err
	}
	return artifacts, nil
}

// Compile implements the compiler interface
func (s *Solidity) Compile(demoPath string) (map[string]*Artifact, error) {
	if len(demoPath) == 0 {
		return nil, fmt.Errorf("no input files")
	}
	return s.compileImpl(demoPath)
}

func (s *Solidity) compileImpl(demoPath string) (map[string]*Artifact, error) {
	args := []string{
		"--combined-json",
		"bin,bin-runtime,abi",
	}

	args = append(args, demoPath)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("D:\\Go\\src\\github.com\\bif\\bif-sdk-go\\compiler\\tmp\\solc.exe", args...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to compile: %s", string(stderr.Bytes()))
	}

	var output *solcOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, err
	}

	artifacts := map[string]*Artifact{}
	for name, i := range output.Contracts {
		artifacts[name] = &Artifact{
			Bin:        i.Bin,
			BinRuntime: i.BinRuntime,
			Abi:        i.Abi,
		}
	}
	return artifacts, nil
}

// DownloadSolidity downloads the solidity compiler
func DownloadSolidity(version string, dst string, renameDst bool) error {
	url := "https://github.com/ethereum/solidity/releases/download/v" + version + "/solc-windows.exe"

	// check if the dst is correct
	exists := false
	fi, err := os.Stat(dst)
	if err == nil {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			exists = true
		case mode.IsRegular():
			return fmt.Errorf("dst is a file")
		}
	} else {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat dst '%s': %v", dst, err)
		}
	}

	// create the destiny path if does not exists
	if !exists {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return fmt.Errorf("cannot create dst path: %v", err)
		}
	}

	// rename binary
	name := "solidity"
	if renameDst {
		name += "-" + version
	}

	// tmp folder to download the binary
	tmpDir, err := ioutil.TempDir("/tmp", "solc-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, name)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// make binary executable
	if err := os.Chmod(path, 0755); err != nil {
		return err
	}

	// move file to dst
	if err := os.Rename(path, filepath.Join(dst, name)); err != nil {
		return err
	}
	return nil
}
