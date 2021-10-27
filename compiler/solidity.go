package compiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Output struct {
	Contracts map[string]*Artifact
	Sources   map[string]*Source
	Version   string
}

type Source struct {
	AST map[string]interface{}
}

type Artifact struct {
	Abi           string
	Bin           string
	BinRuntime    string `json:"bin-runtime"`
	SrcMap        string `json:"srcmap"`
	SrcMapRuntime string `json:"srcmap-runtime"`
}

// Solidity is the solidity compiler
type Solidity struct {
	path string
}

// NewSolidityCompiler instantiates a new solidity compiler
func NewSolidityCompiler(path string) *Solidity {
	return &Solidity{path}
}

// CompileCode compiles a solidity code
func (s *Solidity) CompileCode(code string) (*Output, error) {
	if code == "" {
		return nil, errors.New("code is empty")
	}
	output, err := s.compileImpl(code)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Compile implements the compiler interface
func (s *Solidity) Compile(files ...string) (*Output, error) {
	if len(files) == 0 {
		return nil, errors.New("no input files")
	}
	return s.compileImpl("", files...)
}

func (s *Solidity) compileImpl(code string, files ...string) (*Output, error) {
	args := []string{
		"--combined-json",
		"bin,bin-runtime,srcmap-runtime,abi,srcmap,ast",
	}
	if code != "" {
		args = append(args, "-")
	}
	if len(files) != 0 {
		args = append(args, files...)
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(s.path, args...)
	if code != "" {
		cmd.Stdin = strings.NewReader(code)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to compile: %s", string(stderr.Bytes()))
	}

	var output *Output
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, err
	}
	return output, nil
}
