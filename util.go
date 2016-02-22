// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package went

import (
  "fmt"
  "os"
  "math"
  "strings"
  "strconv"
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

// Parse a number from a string
func parseNumber(number string) (float64, error) {
	var value float64
	number = strings.Trim(number, " \t\n")
	if v, err := strconv.ParseInt(number, 0, 64); err != nil {
		if v2, err2 := strconv.ParseFloat(number, 64); err2 != nil {
			return 0, err2
		} else {
			value = float64(v2)
		}
	} else {
		value = float64(v)
	}
	return value, nil
}