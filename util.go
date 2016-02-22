// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package went

import (
  "fmt"
  "os"
  "math"
)

func assert(cond bool, msg string) {
  if !cond {
    fmt.Printf("assertion failed: %s\n", msg)
    os.Exit(1)
  }
}

func isInt(n float64) bool {
	return math.Trunc(n) == n
}