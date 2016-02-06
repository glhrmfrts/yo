package vm

// Contains executable code by the VM and 
// static information generated at compilation time.
// All runtime functions has one of these
type FuncProto struct {
  Source    string
  NumConsts int32
  NumCode   int32
  Consts    []Value
  Code      []uint32
  Lines     []int32
}

func newFuncProto(source string) *FuncProto {
  return &FuncProto{
    Source: source,
  }
}

func (f *FuncProto) AddInstruction(instr uint32, line int) {
  f.Code = append(f.Code, instr)
  f.NumCode++

  // TODO: find an alternative to this
  f.Lines = append(f.Lines, int32(line))
}