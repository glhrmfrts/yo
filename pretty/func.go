// "Disassemble" a function prototype (bytecode chunk)
// into a human-readable string

package pretty

import (
  "fmt"
  "bytes"
  elo "github.com/glhrmfrts/elo-lang"
)

func doIndent(buf *bytes.Buffer, indent int) {
  for indent > 0 {
    buf.WriteString(" ")
    indent--
  }
}

func disasmImpl(f *elo.FuncProto, buf *bytes.Buffer, indent int) {
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
    if a >= elo.OpConstOffset {
      return fmt.Sprint(f.Consts[a-elo.OpConstOffset])
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
    opcode := elo.OpGetOpcode(instr)

    if lineChanged || i == 0 {
      buf.WriteString(fmt.Sprint(line.Line, "\t"))
    } else {
      buf.WriteString("\t")
    }

    buf.WriteString(fmt.Sprint(i))
    buf.WriteString(fmt.Sprint("\t", opcode, "\t"))

    switch opcode {
    case elo.OP_LOADNIL:
      buf.WriteString(fmt.Sprintf("\t!%d !%d", elo.OpGetA(instr), elo.OpGetB(instr)))
    case elo.OP_LOADCONST:
      buf.WriteString(fmt.Sprintf("!%d %s", elo.OpGetA(instr), f.Consts[elo.OpGetBx(instr)]))
    case elo.OP_NEG, elo.OP_NOT, elo.OP_CMPL:
      bx := elo.OpGetBx(instr)
      bstr := getRegOrConst(bx)
      buf.WriteString(fmt.Sprintf("\t!%d %s", elo.OpGetA(instr), bstr))
    case elo.OP_ADD, elo.OP_SUB, elo.OP_MUL, elo.OP_DIV, elo.OP_POW, elo.OP_SHL, elo.OP_SHR, 
        elo.OP_AND, elo.OP_OR, elo.OP_XOR, elo.OP_LE, elo.OP_LT, elo.OP_EQ, elo.OP_NE, elo.OP_GET, elo.OP_SET:
      a, b, c := elo.OpGetA(instr), elo.OpGetB(instr), elo.OpGetC(instr)
      bstr, cstr := getRegOrConst(b), getRegOrConst(c)
      buf.WriteString(fmt.Sprintf("\t!%d %s %s", a, bstr, cstr))
    case elo.OP_APPEND, elo.OP_RETURN:
      a, b := elo.OpGetA(instr), elo.OpGetB(instr)
      buf.WriteString(fmt.Sprintf("\t!%d #%d", a, b))
    case elo.OP_MOVE:
      a, b := elo.OpGetA(instr), elo.OpGetB(instr)
      buf.WriteString(fmt.Sprintf("\t!%d !%d", a, b))
    case elo.OP_LOADGLOBAL, elo.OP_SETGLOBAL:
      a, bx := elo.OpGetA(instr), elo.OpGetBx(instr)
      buf.WriteString(fmt.Sprintf("!%d %s", a, f.Consts[bx]))
    case elo.OP_LOADREF, elo.OP_SETREF:
      a, bx := elo.OpGetA(instr), elo.OpGetBx(instr)
      buf.WriteString(fmt.Sprintf("\t!%d %s", a, f.Consts[bx]))
    case elo.OP_CALL:
      a, b, c := elo.OpGetA(instr), elo.OpGetB(instr), elo.OpGetC(instr)
      buf.WriteString(fmt.Sprintf("\t!%d #%d #%d", a, b, c))
    case elo.OP_ARRAY, elo.OP_OBJECT:
      buf.WriteString(fmt.Sprintf("\t!%d", elo.OpGetA(instr)))
    case elo.OP_FUNC:
      buf.WriteString(fmt.Sprintf("\t!%d &%d", elo.OpGetA(instr), elo.OpGetBx(instr)))
    case elo.OP_JMP:
      sbx := elo.OpGetsBx(instr)
      buf.WriteString(fmt.Sprintf("\t->%d", i + 1 + sbx))
    case elo.OP_JMPTRUE, elo.OP_JMPFALSE:
      a, sbx := elo.OpGetA(instr), elo.OpGetsBx(instr)
      astr := getRegOrConst(a)
      buf.WriteString(fmt.Sprintf("%s ->%d", astr, i + 1 + sbx))
    }

    buf.WriteString("\n")
    doIndent(buf, indent)
  }

  buf.WriteString("}}}\n\n\n")
}

func Disasm(f *elo.FuncProto) string {
  var buf bytes.Buffer
  disasmImpl(f, &buf, 0)
  return buf.String()
}