package colon

import (
	"errors"
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

type Objects []Object

type Result struct {
	String  string
	Objects Objects
}

type Object struct {
	Command string
	Options []string
	Error   error
	Path    string
}

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
		args, err := shellwords.Parse(item)
		objs = append(objs, Object{
			Command: args[0],
			Options: args[1:],
			Error:   err,
			Path:    "",
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

func (r *Result) First() Object {
	if len(r.Objects) > 0 {
		return r.Objects[0]
	}
	return Object{}
}

func (r *Result) Executable() Objects {
	var objs Objects
	for _, obj := range r.Objects {
		path, err := exec.LookPath(obj.Command)
		if err != nil {
			// TODO: add option
			// continue
		}
		obj.Path = path
		obj.Error = err
		objs = append(objs, obj)
	}
	return objs
}

// TODO: imple
// ArgsN
// Arg
// ArgN
func (objs Objects) Args(n int) []string {
	var ret []string
	for i, obj := range objs {
		if i == n {
			ret = []string{obj.Command}
			ret = append(ret, obj.Options...)
		}
	}
	return ret
}
