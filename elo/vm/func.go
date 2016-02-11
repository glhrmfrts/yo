package vm

// TODO: move disasm to it's own package

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
// All runtime functions reference one of these
type FuncProto struct {
  Source    string
  NumConsts uint32
  NumCode   uint32
  NumLines  uint32
  NumFuncs  uint32
  Consts    []Value
  Code      []uint32
  Lines     []LineInfo
  Funcs     []*FuncProto
}

const (
  funcMaxConsts = 0xffff
)

func newFuncProto(source string) *FuncProto {
  return &FuncProto{
    Source: source,
  }
}

func doIndent(buf *bytes.Buffer, indent int) {
  for indent > 0 {
    buf.WriteString(" ")
    indent--
  }
}

func disasmImpl(f *FuncProto, buf *bytes.Buffer, indent int) {
  buf.WriteString("\n\n")
  doIndent(buf, indent)
  buf.WriteString(fmt.Sprintf("function at %s {{{\n", f.Source))
  doIndent(buf, indent)
  buf.WriteString(fmt.Sprintf("constants: %d\n", f.NumConsts))
  doIndent(buf, indent)
  for _, c := range f.Consts {
    buf.WriteString(fmt.Sprint("\t", c))
  }
  buf.WriteString("\n\n")

  doIndent(buf, indent)
  buf.WriteString(fmt.Sprintf("funcs: %d\n", f.NumFuncs))
  for _, f := range f.Funcs {
    disasmImpl(f, buf, indent + 2)
  }

  doIndent(buf, indent)
  buf.WriteString("line\t#\topcode\t\targs\n")

  getRegOrConst := func(a uint) string {
    if a >= kConstOffset {
      return fmt.Sprint(f.Consts[a-kConstOffset])
    } else {
      return fmt.Sprintf("!%d", a)
    }
  }

  doIndent(buf, indent)
  var currentLine uint32
  for i, instr := range f.Code {
    lineChanged := false
    if currentLine + 1 < f.NumLines && (i >= int(f.Lines[currentLine + 1].Instr)) {
      currentLine += 1
      lineChanged = true
    }

    line := f.Lines[currentLine]
    opcode := opGetOpcode(instr)

    if lineChanged || i == 0 {
      buf.WriteString(fmt.Sprint(line.Line, "\t"))
    } else {
      buf.WriteString("\t")
    }

    buf.WriteString(fmt.Sprint(i))
    buf.WriteString(fmt.Sprint("\t", opcode, "\t"))

    switch opcode {
    case OP_LOADNIL:
      buf.WriteString(fmt.Sprintf("\t!%d !%d", opGetA(instr), opGetB(instr)))
    case OP_LOADCONST:
      buf.WriteString(fmt.Sprintf("!%d %s", opGetA(instr), f.Consts[opGetBx(instr)]))
    case OP_NEG, OP_NOT, OP_CMPL:
      bx := opGetBx(instr)
      bstr := getRegOrConst(bx)
      buf.WriteString(fmt.Sprintf("\t!%d %s", opGetA(instr), bstr))
    case OP_ADD, OP_SUB, OP_MUL, OP_DIV, OP_POW, OP_SHL, OP_SHR, 
        OP_AND, OP_OR, OP_XOR, OP_LE, OP_LT, OP_EQ, OP_NE, OP_GET, OP_SET:
      a, b, c := opGetA(instr), opGetB(instr), opGetC(instr)
      bstr, cstr := getRegOrConst(b), getRegOrConst(c)
      buf.WriteString(fmt.Sprintf("\t!%d %s %s", a, bstr, cstr))
    case OP_APPEND:
      a, bx := opGetA(instr), opGetBx(instr)
      buf.WriteString(fmt.Sprintf("\t!%d #%d", a, bx))
    case OP_MOVE:
      a, b := opGetA(instr), opGetB(instr)
      buf.WriteString(fmt.Sprintf("\t!%d !%d", a, b))
    case OP_LOADGLOBAL, OP_SETGLOBAL:
      a, bx := opGetA(instr), opGetBx(instr)
      buf.WriteString(fmt.Sprintf("!%d %s", a, f.Consts[bx]))
    case OP_LOADREF, OP_SETREF:
      a, bx := opGetA(instr), opGetBx(instr)
      buf.WriteString(fmt.Sprintf("\t!%d %s", a, f.Consts[bx]))
    case OP_CALL:
      a, b, c := opGetA(instr), opGetB(instr), opGetC(instr)
      buf.WriteString(fmt.Sprintf("\t!%d #%d #%d", a, b, c))
    case OP_ARRAY, OP_OBJECT:
      buf.WriteString(fmt.Sprintf("\t!%d", opGetA(instr)))
    case OP_FUNC:
      buf.WriteString(fmt.Sprintf("\t!%d &%d", opGetA(instr), opGetBx(instr)))
    case OP_JMP:
      sbx := opGetsBx(instr)
      buf.WriteString(fmt.Sprintf("->%d", i + sbx))
    case OP_JMPTRUE, OP_JMPFALSE:
      a, sbx := opGetA(instr), opGetsBx(instr)
      buf.WriteString(fmt.Sprintf("!%d ->%d", a, i + sbx))
    }

    buf.WriteString("\n")
    doIndent(buf, indent)
  }

  buf.WriteString("}}}\n\n\n")
}

func Disasm(f *FuncProto, buf *bytes.Buffer) {
  disasmImpl(f, buf, 0)
}