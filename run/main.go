package main

import (
  "fmt"
  "os"
  "io/ioutil"
	"github.com/glhrmfrts/elo-lang/parse"
  "github.com/glhrmfrts/elo-lang/pretty"
  _ "github.com/glhrmfrts/elo-lang"
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

/*
  code, err := elo.Compile(root, filename)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  fmt.Println(pretty.Disasm(code))
*/
}