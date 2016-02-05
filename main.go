package main

import (
  "fmt"
  "os"
  "io/ioutil"
	"github.com/glhrmfrts/elo-lang/elo/parse"
  "github.com/glhrmfrts/elo-lang/elo/ast"
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

  out := ast.Prettyprint(root, 2)
  fmt.Println(out)
}