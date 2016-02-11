package elo

import (
  "fmt"
)

type ValueType int

const (
  VALUE_NIL ValueType = iota
  VALUE_BOOL
  VALUE_NUMBER
  VALUE_STRING
)

type (
  Value interface {
    // manual type retrieval seems to be faster than
    // go's interface type assertion
    Type() ValueType
    assertFloat64() (float64, bool)
    assertBool() (bool, bool)
    assertString() (string, bool)
  }

  // Nil is just an empty struct
  Nil struct{}

  // Bool maps directly to bool
  Bool bool

  // Number is a double-precision floating point number
  Number float64

  // String also maps directly to string
  String string
)


func (v Nil) Type() ValueType                { return VALUE_NIL }
func (v Nil) assertFloat64() (float64, bool) { return 0, false }
func (v Nil) assertBool() (bool, bool)       { return false, false }
func (v Nil) assertString() (string, bool)   { return "", false }

func (v Nil) String() string {
  return "nil"
}

func (v Bool) Type() ValueType                { return VALUE_BOOL }
func (v Bool) assertFloat64() (float64, bool) { return 0, false }
func (v Bool) assertBool() (bool, bool)       { return bool(v), true }
func (v Bool) assertString() (string, bool)   { return "", false }

func (v Bool) String() string {
  if bool(v) {
    return "true"
  }
  return "false"
}

func (v Number) Type() ValueType                { return VALUE_NUMBER }
func (v Number) assertFloat64() (float64, bool) { return float64(v), true }
func (v Number) assertBool() (bool, bool)       { return false, false }
func (v Number) assertString() (string, bool)   { return "", false }

func (v Number) String() string {
  return fmt.Sprint(float64(v))
}

func (v String) Type() ValueType                { return VALUE_STRING }
func (v String) assertFloat64() (float64, bool) { return 0, false }
func (v String) assertBool() (bool, bool)       { return false, false }
func (v String) assertString() (string, bool)   { return string(v), true }

func (v String) String() string {
  return string(v)
}