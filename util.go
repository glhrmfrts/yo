// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package elo

import (
  "fmt"
  "os"
)

func assert(cond bool, msg string) {
  if !cond {
    fmt.Printf("assertion failed: %s\n", msg)
    os.Exit(1)
  }
}