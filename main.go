package main

import (
  "fmt"
	"github.com/glhrmfrts/elo-lang/elo/parse"
  "github.com/glhrmfrts/elo-lang/elo/ast"
)

func main() {
	root := parse.Parse("hello", "test")
  out := ast.Prettyprint(root)
  fmt.Println(out)
}