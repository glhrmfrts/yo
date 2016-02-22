// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package main

import (
  "fmt"
  "os"
  "io/ioutil"
	"github.com/glhrmfrts/went/parse"
  "github.com/glhrmfrts/went/pretty"
  "github.com/glhrmfrts/went"
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

  fmt.Println(pretty.SyntaxTree(root, 2))

  code, err := went.Compile(root, filename)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  fmt.Println(pretty.Disasm(code))

  state := went.NewState()
  state.RunProto(code)
}