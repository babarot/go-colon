package colon

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

var (
	Separator string = ":"
	UseOption bool   = false
)

type Parser struct {
	Separator string
	UseOption bool
}

type Result struct {
	String  string
	Objects Objects
}

type Object struct {
	Args    []string
	First   string
	Errors  []error
	Command string
	IsDir   bool
}

type Objects []Object

func NewParser() *Parser {
	return &Parser{
		Separator,
		UseOption,
	}
}

func (p *Parser) Parse(str string) (*Result, error) {
	if str == "" {
		return &Result{}, errors.New("invalid strings")
	}
	var objs Objects

	items := strings.Split(str, p.Separator)
	for _, item := range items {
		if item == "" {
			continue
		}
		var (
			errStack []error
			command  string
		)
		args, err := shellwords.Parse(item)
		if err != nil {
			errStack = append(errStack, err)
		}
		isDir := isDir(args[0])
		if !isDir {
			// Discard err in order not to make an error
			// even if it is not found in your PATH
			command, _ = exec.LookPath(args[0])
		}
		// Error if it can not be found in your all PATHs
		// or in the current directory
		if command == "" && !isExist(args[0]) {
			err = fmt.Errorf("%s: no such file or directory", args[0])
			errStack = append(errStack, err)
		}
		objs = append(objs, Object{
			Args:    args,
			First:   args[0],
			Errors:  errStack,
			Command: command,
			IsDir:   isDir,
		})
	}

	return &Result{
		String:  str,
		Objects: objs,
	}, nil
}

func Parse(str string) (*Result, error) {
	return NewParser().Parse(str)
}

func (r *Result) Filter(fn func(Object) bool) Objects {
	ret := make(Objects, 0)
	for _, obj := range r.Objects {
		if fn(obj) {
			ret = append(ret, obj)
		}
	}
	return ret
}

func (r *Result) Get(str string) Objects {
	return r.Filter(func(o Object) bool {
		return strings.Contains(strings.Join(o.Args, " "), str)
	})
}

func (r *Result) WithoutErrors() Objects {
	return r.Filter(func(o Object) bool {
		return len(o.Errors) == 0
	})
}

func (r *Result) Executable() Objects {
	return r.Filter(func(o Object) bool {
		return o.Command != ""
	})
}

func isDir(name string) bool {
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
