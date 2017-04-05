package colon

import (
	"errors"
	"os/exec"
	"strings"
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
	String string
	Parsed []string
}

/*
type Result struct {
	String string
	Parsed Parsed
}

type Parsed struct {
	Command string
	Option []string
	Error error
}
*/

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

	var (
		items  []string = strings.Split(str, p.Separator)
		parsed []string
	)

	parsed = items
	if !p.UseOption {
		// initialize if this cond is false (default)
		parsed = []string{}
		for _, item := range items {
			parsed = append(parsed, strings.Split(item, " ")[0])
		}
	}

	return &Result{
		String: str,
		Parsed: parsed,
	}, nil
}

func Parse(str string) (*Result, error) {
	return NewParser().Parse(str)
}

func (r *Result) First() string {
	if len(r.Parsed) > 0 {
		return r.Parsed[0]
	}
	return ""
}

func (r *Result) Executable() Array {
	commands := []string{}
	for _, item := range r.Parsed {
		if item == "" {
			continue
		}
		// TODO
		// case: UseOption is true
		_, err := exec.LookPath(item)
		if err != nil {
			// not found in PATH
			continue
		}
		commands = append(commands, item)
	}
	// TODO (shoud be returned as []string?)
	return Array(commands)
}

type Array []string

func (a Array) Get(n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

func (a Array) First() string {
	return a.Get(0)
}

func (a Array) Array() []string {
	return []string(a)
}
