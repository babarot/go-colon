package main

import (
	"github.com/b4b4r07/go-colon"
	"github.com/k0kubun/pp"
)

const COLON_STRING = "fzy:fzf --tac:/usr/local/bin/peco:not-found-cmd:/bin:/unko"

func main() {
	res, err := colon.Parse(COLON_STRING)
	if err != nil {
		panic(err)
	}
	// pp.Println(res)
	// pp.Println(res.WithoutErrors())
	pp.Println(res.Executable())
	pp.Println(res.Directories())
}
