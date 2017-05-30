// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package main

import (
	"fmt"
	"github.com/glhrmfrts/went"
	"github.com/glhrmfrts/went/parse"
	//"github.com/glhrmfrts/went/pretty"
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

	code, err := went.Compile(root, filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//fmt.Println(pretty.Disasm(code))

	state := went.NewState()
<<<<<<< HEAD
	went.DefineBaseLib(state)

	state.RunProto(code)
=======
	val, err := state.LoadProto(code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(val)
>>>>>>> 52602272bea9b6aba2391188bc0feff5e05be216
}
