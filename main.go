package main

import (
	"github.com/glhrmfrts/elo-lang/elo"
)

func main() {
	ast := elo.Parse("a + 2")
	elo.Prettyprint(ast)
}