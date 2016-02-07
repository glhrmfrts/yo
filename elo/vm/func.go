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

  // information used only in compile-time
  lastLine  int
}

func newFuncProto(source string) *FuncProto {
  return &FuncProto{
    Source: source,
  }
}

func (f *FuncProto) addInstruction(instr uint32, line int) {
  f.Code = append(f.Code, instr)
  f.NumCode++

  if line != f.lastLine {
    f.Lines = append(f.Lines, LineInfo{f.NumCode - 1, uint16(line)})
    f.lastLine = line
  }
}