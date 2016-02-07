package vm

type LineInfo struct {
  Instr uint32 // the instruction index
  Line  uint16
}

// Contains executable code by the VM and 
// static information generated at compilation time.
// All runtime functions has one of these
type FuncProto struct {
  Source    string
  NumConsts uint32
  NumCode   uint32
  NumLines  uint32
  Consts    []Value
  Code      []uint32
  Lines     []LineInfo
}

const (
  funcMaxConsts = 0xffff
)

func newFuncProto(source string) *FuncProto {
  return &FuncProto{
    Source: source,
  }
}