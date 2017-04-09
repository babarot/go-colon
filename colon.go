// Package colon for parsing colon-separated strings like PATH
package colon

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// Object shows parsing result
// e.g. "/bin:/usr/bin:..."
// ......^^^^ Object
type Object struct {
	// Index returns the number of the given character string
	// separated by Separator
	Index int

	// Attr ...
	Attr Attribute

	// Errors stacks all errors that occurred during parsing
	Errors []error

	// Command returns the full path of the command
	// if its first argument is found in PATH
	Command string

	// IsDir returns true if the first argument is a directory
	IsDir bool
}

type Attribute struct {
	First     string
	Other     []string
	Args      []string
	Base, Dir string
}

type Result []Object

// NewParser creates Parser
func NewParser() *Parser {
	return &Parser{
		Separator,
	}
}

// Parser parses colon-separated string like PATH
func (p *Parser) Parse(str string) (*Result, error) {
	if str == "" {
		return &Result{}, ErrInvalid
	}
	var objs Result

	items := strings.Split(str, p.Separator)
	for index, item := range items {
		if item == "" {
			continue
		}
		var (
			errStack []error
			command  string
			first    string
		)
		args, err := shellwords.Parse(item)
		if len(args) > 0 {
			first = args[0]
		}
		if err != nil {
			errStack = append(errStack, err)
		}
		isDir := isDir(first)
		if !isDir {
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
		attr := Attribute{
			First: args[0],
			Other: args[1:],
			Args:  args,
			Base:  filepath.Base(args[0]),
			Dir: func(arg string) string {
				dir := filepath.Dir(arg)
				if dir == "." {
					return ""
				}
				return dir
			}(args[0]),
		}
		objs = append(objs, Object{
			Index:   index + 1,
			Attr:    attr,
			Errors:  errStack,
			Command: command,
			IsDir:   isDir,
		})
	}

	return &objs, nil
}

// Parse is public method exported for accessing from other package
func Parse(str string) (*Result, error) {
	return NewParser().Parse(str)
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

// Filter filters the parse result by condition
func (r *Result) Filter(fn func(Object) bool) *Result {
	ret := make(Result, 0)
	for _, obj := range *r {
		if fn(obj) {
			ret = append(ret, obj)
		}
	}
	return &ret
}

// Get returns one object containing the given string
func (r *Result) Get(str string) *Result {
	return r.Filter(func(o Object) bool {
		return strings.Contains(strings.Join(o.Attr.Args, " "), str)
	})
}

// WithoutErrors returns objects with no errors
func (r *Result) WithoutErrors() *Result {
	return r.Filter(func(o Object) bool {
		return len(o.Errors) == 0
	})
}

// Executable returns objects whose first argument is in PATH
func (r *Result) Executable() *Result {
	return r.Filter(func(o Object) bool {
		return o.Command != ""
	})
}

// Directories returns the objects that the first argument is a directory
func (r *Result) Directories() *Result {
	return r.Filter(func(o Object) bool {
		return o.IsDir
	})
}

// One returns the first Object. If not, returns empty
func (r *Result) One() Object {
	if len(*r) == 0 {
		return Object{}
	}
	return (*r)[0]
}
