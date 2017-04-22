// Package colon for parsing colon-separated strings like PATH
package colon

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	// "path/filepath"
	"strings"

	"github.com/mattn/go-shellwords"
)

var (
	// Separator is an identifier for delimiting the target string
	// Defaults to os.PathListSeparator but you can change anything you want
	Separator string = string(os.PathListSeparator)

	ErrInvalid = errors.New("invalid argument")
)

type Parser struct {
	Separator string
}

// Result is parsed result
type Result struct {
	// Index returns the number of the given character string
	// separated by Separator
	Index int

	Item    string
	Args    []string
	Command string

	// Errors stacks all errors that occurred during parsing
	Errors []error
}

type Results []Result

// NewParser creates Parser
func NewParser() *Parser {
	return &Parser{
		Separator,
	}
}

// Parser parses colon-separated string like PATH
func (p *Parser) Parse(str string) (*Results, error) {
	var rs Results
	if str == "" {
		return &rs, ErrInvalid
	}

	var (
		errStack []error
		isdir    bool
		first    string
		command  string
	)

	items := strings.Split(str, p.Separator)
	for index, item := range items {
		args, err := shellwords.Parse(item)
		if err != nil {
			errStack = append(errStack, err)
		}
		if len(args) == 0 {
			return &rs, ErrInvalid
		}
		first = args[0]
		isdir = isDir(first)
		if !isdir {
			// Discard err in order not to make an error
			// even if it is not found in your PATH
			command, _ = exec.LookPath(first)
		}
		// Error if it can not be found in your all PATHs
		// or in the current directory
		if command == "" && !isExist(first) {
			err = fmt.Errorf("%s: no such file or directory", first)
			errStack = append(errStack, err)
		}
		rs = append(rs, Result{
			Index:   index + 1,
			Item:    item,
			Args:    args,
			Command: command,
			Errors:  errStack,
		})
	}

	return &rs, nil
}

// Parse is public method exported for accessing from other package
func Parse(str string) (*Results, error) {
	return NewParser().Parse(str)
}

func isDir(name string) bool {
	if !strings.HasPrefix(name, "/") {
		return false
	}
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func isExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// Filter filters the parse result by condition
func (rs *Results) Filter(fn func(Result) bool) *Results {
	results := make(Results, 0)
	for _, result := range *rs {
		if fn(result) {
			results = append(results, result)
		}
	}
	return &results
}

// Executable returns objects whose first argument is in PATH
func (rs *Results) Executable() *Results {
	return rs.Filter(func(r Result) bool {
		return r.Command != ""
	})
}

func (rs *Results) First() (Result, error) {
	if len(*rs) == 0 {
		return Result{}, errors.New("no result")
	}
	return (*rs)[0], nil
}

func uniqueSlice(args Results) Results {
	rs := make(Results, 0, len(args))
	encountered := map[int]bool{}
	for _, arg := range args {
		if !encountered[arg.Index] {
			encountered[arg.Index] = true
			rs = append(rs, arg)
		}
	}
	return rs
}

func (rs *Results) Get(args ...interface{}) *Results {
	var results Results
	for _, arg := range args {
		switch arg.(type) {
		case string:
			results = append(results, *rs.Filter(func(r Result) bool {
				if len(r.Args) == 0 {
					return false
				}
				return r.Args[0] == arg.(string)
			})...)
		case int:
			results = append(results, *rs.Filter(func(r Result) bool {
				return r.Index == arg.(int)
			})...)
		}
	}
	// Remove it if there is the same
	results = uniqueSlice(results)
	return &results
}
