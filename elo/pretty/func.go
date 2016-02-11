// "Disassemble" a function prototype (bytecode chunk)
// into a human-readable string

package pretty

import (
  "fmt"
  "bytes"
  "github.com/glhrmfrts/elo-lang/elo/vm"
)

func doIndent(buf *bytes.Buffer, indent int) {
  for indent > 0 {
    buf.WriteString(" ")
    indent--
  }
}

func disasmImpl(f *vm.FuncProto, buf *bytes.Buffer, indent int) {
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
    if a >= vm.OpConstOffset {
      return fmt.Sprint(f.Consts[a-vm.OpConstOffset])
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
    opcode := vm.OpGetOpcode(instr)

    if lineChanged || i == 0 {
      buf.WriteString(fmt.Sprint(line.Line, "\t"))
    } else {
      buf.WriteString("\t")
    }

    buf.WriteString(fmt.Sprint(i))
    buf.WriteString(fmt.Sprint("\t", opcode, "\t"))

    switch opcode {
    case vm.OP_LOADNIL:
      buf.WriteString(fmt.Sprintf("\t!%d !%d", vm.OpGetA(instr), vm.OpGetB(instr)))
    case vm.OP_LOADCONST:
      buf.WriteString(fmt.Sprintf("!%d %s", vm.OpGetA(instr), f.Consts[vm.OpGetBx(instr)]))
    case vm.OP_NEG, vm.OP_NOT, vm.OP_CMPL:
      bx := vm.OpGetBx(instr)
      bstr := getRegOrConst(bx)
      buf.WriteString(fmt.Sprintf("\t!%d %s", vm.OpGetA(instr), bstr))
    case vm.OP_ADD, vm.OP_SUB, vm.OP_MUL, vm.OP_DIV, vm.OP_POW, vm.OP_SHL, vm.OP_SHR, 
        vm.OP_AND, vm.OP_OR, vm.OP_XOR, vm.OP_LE, vm.OP_LT, vm.OP_EQ, vm.OP_NE, vm.OP_GET, vm.OP_SET:
      a, b, c := vm.OpGetA(instr), vm.OpGetB(instr), vm.OpGetC(instr)
      bstr, cstr := getRegOrConst(b), getRegOrConst(c)
      buf.WriteString(fmt.Sprintf("\t!%d %s %s", a, bstr, cstr))
    case vm.OP_APPEND:
      a, bx := vm.OpGetA(instr), vm.OpGetBx(instr)
      buf.WriteString(fmt.Sprintf("\t!%d #%d", a, bx))
    case vm.OP_MOVE:
      a, b := vm.OpGetA(instr), vm.OpGetB(instr)
      buf.WriteString(fmt.Sprintf("\t!%d !%d", a, b))
    case vm.OP_LOADGLOBAL, vm.OP_SETGLOBAL:
      a, bx := vm.OpGetA(instr), vm.OpGetBx(instr)
      buf.WriteString(fmt.Sprintf("!%d %s", a, f.Consts[bx]))
    case vm.OP_LOADREF, vm.OP_SETREF:
      a, bx := vm.OpGetA(instr), vm.OpGetBx(instr)
      buf.WriteString(fmt.Sprintf("\t!%d %s", a, f.Consts[bx]))
    case vm.OP_CALL:
      a, b, c := vm.OpGetA(instr), vm.OpGetB(instr), vm.OpGetC(instr)
      buf.WriteString(fmt.Sprintf("\t!%d #%d #%d", a, b, c))
    case vm.OP_ARRAY, vm.OP_OBJECT:
      buf.WriteString(fmt.Sprintf("\t!%d", vm.OpGetA(instr)))
    case vm.OP_FUNC:
      buf.WriteString(fmt.Sprintf("\t!%d &%d", vm.OpGetA(instr), vm.OpGetBx(instr)))
    case vm.OP_JMP:
      sbx := vm.OpGetsBx(instr)
      buf.WriteString(fmt.Sprintf("->%d", i + sbx))
    case vm.OP_JMPTRUE, vm.OP_JMPFALSE:
      a, sbx := vm.OpGetA(instr), vm.OpGetsBx(instr)
      buf.WriteString(fmt.Sprintf("!%d ->%d", a, i + sbx))
    }

    buf.WriteString("\n")
    doIndent(buf, indent)
  }

  buf.WriteString("}}}\n\n\n")
}

func Disasm(f *vm.FuncProto) string {
  var buf bytes.Buffer
  disasmImpl(f, &buf, 0)
  return buf.String()
}