// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package elo

import (
  "fmt"
)

type ValueType int

const (
  ValueNil ValueType = iota
  ValueBool
  ValueNumber
  ValueString
)

type (
  Value interface {
    // manual type retrieval seems to be faster than
    // go's interface type assertion (runtime.I2T2)
    Type() ValueType
    assertFloat64() (float64, bool)
    assertBool() (bool, bool)
    assertString() (string, bool)
  }

  Nil struct{}

  Bool bool

  Number float64

  String string

  Func struct {
    Proto *FuncProto
  }

  Array []Value

  Object struct {
    Parent *Object
    Fields map[string]Value
  }

  Channel chan Value
)

func (v Nil) Type() ValueType                { return ValueNil }
func (v Nil) assertFloat64() (float64, bool) { return 0, false }
func (v Nil) assertBool() (bool, bool)       { return false, false }
func (v Nil) assertString() (string, bool)   { return "", false }

func (v Nil) String() string {
  return "nil"
}

func (v Bool) Type() ValueType                { return ValueBool }
func (v Bool) assertFloat64() (float64, bool) { return 0, false }
func (v Bool) assertBool() (bool, bool)       { return bool(v), true }
func (v Bool) assertString() (string, bool)   { return "", false }

func (v Bool) String() string {
  if bool(v) {
    return "true"
  }
  return "false"
}

func (v Number) Type() ValueType                { return ValueNumber }
func (v Number) assertFloat64() (float64, bool) { return float64(v), true }
func (v Number) assertBool() (bool, bool)       { return false, false }
func (v Number) assertString() (string, bool)   { return "", false }

func (v Number) String() string {
  return fmt.Sprint(float64(v))
}

func (v String) Type() ValueType                { return ValueString }
func (v String) assertFloat64() (float64, bool) { return 0, false }
func (v String) assertBool() (bool, bool)       { return false, false }
func (v String) assertString() (string, bool)   { return string(v), true }

func (v String) String() string {
  return string(v)
}