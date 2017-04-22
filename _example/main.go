package main

import (
	"github.com/b4b4r07/go-colon"
	"github.com/k0kubun/pp"
)

// const COLON_STRING = "fzy:fzf --tac:/usr/local/bin/peco:not-found-cmd:/bin:/unko"
const COLON_STRING = "/usr/local/bin/peco"

func main() {
	results, err := colon.Parse(COLON_STRING)
	if err != nil {
		panic(err)
	}
	// pp.Println(res)
	result, err := results.Executable().First()
	if err != nil {
		panic(err)
	}
	pp.Println(result)
}
