package vm

import (
  "fmt"
  "bytes"
)

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

func Disasm(f *FuncProto, buf *bytes.Buffer) {
  buf.WriteString(fmt.Sprintf("function at %s\n", f.Source))
  buf.WriteString(fmt.Sprintf("constants: %d\n", f.NumConsts))
  for _, c := range f.Consts {
    buf.WriteString(fmt.Sprint("\t", c))
  }
  buf.WriteString("\n\n")
  buf.WriteString("line\t#\topcode\t\targs\n")

  var currentLine uint32
  for i, instr := range f.Code {
    if currentLine + 1 < f.NumLines && (i >= int(f.Lines[currentLine + 1].Instr)) {
      currentLine += 1
    }

    line := f.Lines[currentLine]
    opcode := opGetOpcode(instr)

    buf.WriteString(fmt.Sprint(line.Line, "\t", i))
    buf.WriteString(fmt.Sprint("\t", opcode, "\t"))

    switch opcode {
    case OP_LOADNIL:
      buf.WriteString(fmt.Sprintf("!%d !%d", opGetA(instr), opGetB(instr)))
    case OP_LOADCONST:
      buf.WriteString(fmt.Sprintf("!%d %s", opGetA(instr), f.Consts[opGetBx(instr)]))
    case OP_NEGATE, OP_NOT:
      bx := opGetBx(instr)
      var bstr string
      if bx >= kConstOffset {
        bstr = fmt.Sprint(f.Consts[bx-kConstOffset])
      } else {
        bstr = fmt.Sprintf("!%d", bx)
      }
      buf.WriteString(fmt.Sprintf("\t!%d %s", opGetA(instr), bstr))
    case OP_ADD, OP_SUB, OP_MUL, OP_DIV:
      a, b, c := opGetA(instr), opGetB(instr), opGetC(instr)
      var bstr, cstr string
      if b >= kConstOffset {
        bstr = fmt.Sprint(f.Consts[b-kConstOffset])
      } else {
        bstr = fmt.Sprintf("!%d", b)
      }
      if c >= kConstOffset {
        cstr = fmt.Sprint(f.Consts[c-kConstOffset])
      } else {
        cstr = fmt.Sprintf("!%d", c)
      }
      buf.WriteString(fmt.Sprintf("\t!%d %s %s", a, bstr, cstr))
    case OP_MOVE:
      a, b := opGetA(instr), opGetB(instr)
      buf.WriteString(fmt.Sprintf("\t!%d !%d", a, b))
    case OP_LOADGLOBAL:
      a, bx := opGetA(instr), opGetBx(instr)
      buf.WriteString(fmt.Sprintf("!%d %s", a, f.Consts[bx]))
    case OP_JMP:
      sbx := opGetsBx(instr)
      buf.WriteString(fmt.Sprintf("->%d", i + sbx))
    case OP_JMPTRUE, OP_JMPFALSE:
      a, sbx := opGetA(instr), opGetsBx(instr)
      buf.WriteString(fmt.Sprintf("!%d ->%d", a, i + sbx))
    }

    buf.WriteString("\n")
  }
}