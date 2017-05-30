// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package main

import (
	"fmt"
	"github.com/glhrmfrts/yo"
	"github.com/glhrmfrts/yo/parse"
	"github.com/glhrmfrts/yo/pretty"
	"io/ioutil"
	"os"
)

func main() {
	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	root, err := parse.ParseFile(source, filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//fmt.Println(pretty.SyntaxTree(root, 2))

	code, err := yo.Compile(root, filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(pretty.Disasm(code))

	vm := yo.NewVM()
	vm.RunBytecode(code)
}
