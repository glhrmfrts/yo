package main

import (
  "fmt"
  "os"
  "bytes"
  "io/ioutil"
	"github.com/glhrmfrts/elo-lang/elo/parse"
  "github.com/glhrmfrts/elo-lang/elo/ast"
  "github.com/glhrmfrts/elo-lang/elo/vm"
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

  code, err := vm.Compile(root, filename)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  buf := new(bytes.Buffer)
  vm.Disasm(code, buf)
  fmt.Println(buf.String())
}