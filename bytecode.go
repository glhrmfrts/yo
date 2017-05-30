// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package yo

type LineInfo struct {
	Instr uint32 // the instruction index
	Line  uint16
}

// Contains executable code by the VM and
// static information generated at compilation time.
// All runtime functions reference one of these
type Bytecode struct {
	Source    string
	NumConsts uint32
	NumCode   uint32
	NumLines  uint32
	NumFuncs  uint32
	Consts    []Value
	Code      []uint32
	Lines     []LineInfo
	Funcs     []*Bytecode
}

const (
	bytecodeMaxConsts = 0xffff
)

func newBytecode(source string) *Bytecode {
	return &Bytecode{
		Source: source,
	}
}
