package main

import (
	"fmt"

	"github.com/b4b4r07/go-colon"
	"github.com/k0kubun/pp"
)

const COLON_STRING = "fzy:fzf --tac:/usr/local/bin/peco:not-found-cmd"

func main() {
	res, err := colon.Parse(COLON_STRING)
	if err != nil {
		panic(err)
	}
	pp.Println(res)
	fmt.Printf("%#v\n", res.First())
	fmt.Printf("%#v\n", res.Executable())
	// fmt.Printf("%#v\n", res.Executable().First())
	// fmt.Printf("%#v\n", res.Executable().Array())
}
