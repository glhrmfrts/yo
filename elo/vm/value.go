package vm

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


func (v Nil) Type() ValueType {
  return VALUE_NIL
}

func (v Bool) Type() ValueType {
  return VALUE_STRING
}

func (v Number) Type() ValueType {
  return VALUE_NUMBER
}

func (v String) Type() ValueType {
  return VALUE_STRING
}