package main

import (
	"github.com/b4b4r07/go-colon"
	"github.com/k0kubun/pp"
)

const COLON_STRING = "fzy:fzf --tac:/usr/local/bin/peco:not-found-cmd"

func main() {
	res, err := colon.Parse(COLON_STRING)
	if err != nil {
		panic(err)
	}
	pp.Println(res.Executable())
	pp.Println(res.Executable().Args(1))
}
