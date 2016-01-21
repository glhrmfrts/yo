package main

import (
  "fmt"
	"github.com/glhrmfrts/elo-lang/elo/parse"
  "github.com/glhrmfrts/elo-lang/elo/ast"
)

func main() {
	root := parse.Parse("player.pos.x = 3", "test")
  out := ast.Prettyprint(root)
  fmt.Println(out)
}