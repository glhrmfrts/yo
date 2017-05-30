// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package went

import (
	"fmt"
)

type (
	// ValueType is an internal enumeration which identifies
	// the type of the values without the need for type assertion.
	ValueType int

	// Value interface
	Value interface {
		// manual type retrieval seems to be faster than
		// go's interface type assertion (runtime.I2T2)
		assertFloat64() (float64, bool)
		assertBool() (bool, bool)
		assertString() (string, bool)

		// Type returns the corresponding ValueType enumeration value
		// of this value.
		Type() ValueType

		// String returns the representation of this value as a string,
		// it is the default way to convert a value to a string.
		String() string

		// ToBool converts the value to a boolean value.
		// If it's is nil or false it returns false, otherwise returns true.
		ToBool() bool
	}

	// Nil represents the absence of a usable value.
	Nil struct{}

	// Bool can be either true or false.
	Bool bool

	// A Number is a double-precision floating point number.
	Number float64

	// A String is an immutable sequence of bytes.
	String string

	// GoFunc is a function defined in the host application, which
	// is callable from the script.
	GoFunc func(*FuncCall)

	// Func is a function defined in the script.
	Func struct {
		Proto  *FuncProto
		isGo   bool
		goFunc GoFunc
	}

	// Array is a collection of Values stored contiguously in memory,
	// it's index starts at 0.
	Array []Value

	// Object is a map that maps strings to Values, and may have a
	// parent Object which is used to look for keys that are not in it's own map.
	Object struct {
		Parent *Object
		Fields map[string]Value
	}

	// GoObject is an object that allows the host application to maintain
	// custom data throughout the script, the script user will see it as
	// a regular object.
	GoObject struct {
		Object
		Data interface{}
	}

	// Chan is an object that allows goroutines to
	// communicate/send Values to one another.
	Chan chan Value
)

const (
	ValueNil ValueType = iota
	ValueBool
	ValueNumber
	ValueString
	ValueFunc
	ValueArray
	ValueObject
	ValueChan
)

var (
	valueTypeNames = [8]string{"nil", "bool", "number", "string", "func", "array", "object", "chan"}
)

func (t ValueType) String() string {
	return valueTypeNames[t]
}

// Nil

func (v Nil) assertFloat64() (float64, bool) { return 0, false }
func (v Nil) assertBool() (bool, bool)       { return false, false }
func (v Nil) assertString() (string, bool)   { return "", false }

func (v Nil) Type() ValueType { return ValueNil }
func (v Nil) ToBool() bool    { return false }
func (v Nil) String() string {
	return "nil"
}

// Bool

func (v Bool) assertFloat64() (float64, bool) { return 0, false }
func (v Bool) assertBool() (bool, bool)       { return bool(v), true }
func (v Bool) assertString() (string, bool)   { return "", false }

func (v Bool) Type() ValueType { return ValueBool }
func (v Bool) ToBool() bool    { return bool(v) }
func (v Bool) String() string {
	if bool(v) {
		return "true"
	}
	return "false"
}

// Number

func (v Number) assertFloat64() (float64, bool) { return float64(v), true }
func (v Number) assertBool() (bool, bool)       { return false, false }
func (v Number) assertString() (string, bool)   { return "", false }

func (v Number) Type() ValueType { return ValueNumber }
func (v Number) ToBool() bool    { return true }
func (v Number) String() string {
	return fmt.Sprint(float64(v))
}

// String

func (v String) assertFloat64() (float64, bool) { return 0, false }
func (v String) assertBool() (bool, bool)       { return false, false }
func (v String) assertString() (string, bool)   { return string(v), true }

func (v String) Type() ValueType { return ValueString }
func (v String) ToBool() bool    { return true }
func (v String) String() string {
	return string(v)
}

// GoFunc 

func (v GoFunc) assertFloat64() (float64, bool) { return 0, false }
func (v GoFunc) assertBool() (bool, bool) { return false, false }
func (v GoFunc) assertString() (string, bool) { return "", false }

func (v GoFunc) Type() ValueType { return ValueFunc }
func (v GoFunc) ToBool() bool    { return true }
func (v GoFunc) String() string  { return "func" }

// Func

func (v Func) assertFloat64() (float64, bool) { return 0, false }
func (v Func) assertBool() (bool, bool) { return false, false }
func (v Func) assertString() (string, bool) { return "", false }

func (v Func) Type() ValueType { return ValueFunc }
func (v Func) ToBool() bool    { return true }
func (v Func) String() string  { return "func" }

// Array

func (v Array) assertFloat64() (float64, bool) { return 0, false }
func (v Array) assertBool() (bool, bool) { return false, false }
func (v Array) assertString() (string, bool) { return "", false }

func (v Array) Type() ValueType { return ValueArray }
func (v Array) ToBool() bool    { return true }
func (v Array) String() string  { return fmt.Sprintf("%v", v) }

// Object

func (v Object) assertFloat64() (float64, bool) { return 0, false }
func (v Object) assertBool() (bool, bool) { return false, false }
func (v Object) assertString() (string, bool) { return "", false }

func (v Object) Type() ValueType { return ValueObject }
func (v Object) ToBool() bool    { return true }
func (v Object) String() string  { return fmt.Sprintf("%v", v.Fields) }

func NewObject(parent *Object, fields map[string]Value) *Object {
	return &Object{
		Parent: parent,
		Fields: fields,
	}
}
