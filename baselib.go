package went

import (
  "fmt"
)

func DefineBaseLib(s *State) {
  s.DefineGlobal("println", GoFunc(basePrintln))
}

func basePrintln(call *FuncCall) {
  for i := uint(0); i < call.NumArgs; i++ {
    fmt.Printf("%v ", call.Args[i])
  }

  fmt.Println()
}