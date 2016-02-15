// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package main

import (
  "fmt"
  "os"
  "io/ioutil"
	"github.com/glhrmfrts/elo/parse"
  "github.com/glhrmfrts/elo/pretty"
  "github.com/glhrmfrts/elo"
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

  code, err := elo.Compile(root, filename)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  fmt.Println(pretty.Disasm(code))

  state := elo.State{}
  state.RunProto(code)
}