// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>
//
// "Disassemble" a function prototype (bytecode chunk)
// into a human-readable string

package pretty

import (
  "fmt"
  "bytes"
  elo "github.com/glhrmfrts/elo"
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
    pc := i + 1

    if lineChanged || i == 0 {
      buf.WriteString(fmt.Sprint(line.Line, "\t"))
    } else {
      buf.WriteString("\t")
    }

    buf.WriteString(fmt.Sprint(i))
    buf.WriteString(fmt.Sprint("\t", opcode, "\t"))

    switch opcode {
    case elo.OpLoadnil:
      buf.WriteString(fmt.Sprintf("\t!%d !%d", elo.OpGetA(instr), elo.OpGetB(instr)))
    case elo.OpLoadconst:
      buf.WriteString(fmt.Sprintf("!%d %s", elo.OpGetA(instr), f.Consts[elo.OpGetBx(instr)]))
    case elo.OpNeg, elo.OpNot, elo.OpCmpl:
      bx := elo.OpGetBx(instr)
      bstr := getRegOrConst(bx)
      buf.WriteString(fmt.Sprintf("\t!%d %s", elo.OpGetA(instr), bstr))
    case elo.OpAdd, elo.OpSub, elo.OpMul, elo.OpDiv, elo.OpPow, elo.OpShl, elo.OpShr, 
        elo.OpAnd, elo.OpOr, elo.OpXor, elo.OpLe, elo.OpLt, elo.OpEq, elo.OpNe, elo.OpGet, elo.OpSet:
      a, b, c := elo.OpGetA(instr), elo.OpGetB(instr), elo.OpGetC(instr)
      bstr, cstr := getRegOrConst(b), getRegOrConst(c)
      buf.WriteString(fmt.Sprintf("\t!%d %s %s", a, bstr, cstr))
    case elo.OpAppend, elo.OpReturn:
      a, b := elo.OpGetA(instr), elo.OpGetB(instr)
      buf.WriteString(fmt.Sprintf("\t!%d #%d", a, b))
    case elo.OpMove:
      a, b := elo.OpGetA(instr), elo.OpGetB(instr)
      buf.WriteString(fmt.Sprintf("\t!%d !%d", a, b))
    case elo.OpLoadglobal, elo.OpSetglobal:
      a, bx := elo.OpGetA(instr), elo.OpGetBx(instr)
      buf.WriteString(fmt.Sprintf("!%d %s", a, f.Consts[bx]))
    case elo.OpLoadref, elo.OpSetref:
      a, bx := elo.OpGetA(instr), elo.OpGetBx(instr)
      buf.WriteString(fmt.Sprintf("\t!%d %s", a, f.Consts[bx]))
    case elo.OpCall, elo.OpCallmethod:
      a, b, c := elo.OpGetA(instr), elo.OpGetB(instr), elo.OpGetC(instr)
      if opcode == elo.OpCall {
        buf.WriteString(fmt.Sprintf("\t!%d #%d #%d", a, b, c))
      } else {
        buf.WriteString(fmt.Sprintf("!%d #%d #%d", a, b, c))
      }
    case elo.OpArray, elo.OpObject:
      buf.WriteString(fmt.Sprintf("\t!%d", elo.OpGetA(instr)))
    case elo.OpFunc:
      buf.WriteString(fmt.Sprintf("\t!%d &%d", elo.OpGetA(instr), elo.OpGetBx(instr)))
    case elo.OpJmp:
      sbx := elo.OpGetsBx(instr)
      buf.WriteString(fmt.Sprintf("\t->%d", pc + sbx))
    case elo.OpJmptrue, elo.OpJmpfalse:
      a, sbx := elo.OpGetA(instr), elo.OpGetsBx(instr)
      astr := getRegOrConst(a)
      buf.WriteString(fmt.Sprintf("%s ->%d", astr, pc + sbx))
    case elo.OpForbegin:
      a, b := elo.OpGetA(instr), elo.OpGetB(instr)
      buf.WriteString(fmt.Sprintf("!%d !%d", a, b))
    case elo.OpForiter:
      a, b, c := elo.OpGetA(instr), elo.OpGetB(instr), elo.OpGetC(instr)
      buf.WriteString(fmt.Sprintf("\t!%d !%d !%d", a, b, c))
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