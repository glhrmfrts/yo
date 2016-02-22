// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>
//
// "Disassemble" a function prototype (bytecode chunk)
// into a human-readable string

package pretty

import (
	"bytes"
	"fmt"
	"github.com/glhrmfrts/went"
)

func doIndent(buf *bytes.Buffer, indent int) {
	for indent > 0 {
		buf.WriteString(" ")
		indent--
	}
}

func disasmImpl(f *went.FuncProto, buf *bytes.Buffer, indent int) {
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
		disasmImpl(f, buf, indent+2)
	}

	doIndent(buf, indent)
	buf.WriteString("line\t#\topcode\t\targs\n")

	getRegOrConst := func(a uint) string {
		if a >= went.OpConstOffset {
			return fmt.Sprint(f.Consts[a-went.OpConstOffset])
		} else {
			return fmt.Sprintf("!%d", a)
		}
	}

	doIndent(buf, indent)
	var currentLine uint32
	for i, instr := range f.Code {
		lineChanged := false
		if currentLine+1 < f.NumLines && (i >= int(f.Lines[currentLine+1].Instr)) {
			currentLine += 1
			lineChanged = true
		}

		line := f.Lines[currentLine]
		opcode := went.OpGetOpcode(instr)
		pc := i + 1

		if lineChanged || i == 0 {
			buf.WriteString(fmt.Sprint(line.Line, "\t"))
		} else {
			buf.WriteString("\t")
		}

		buf.WriteString(fmt.Sprint(i))
		buf.WriteString(fmt.Sprint("\t", opcode, "\t"))

		switch opcode {
		case went.OpLoadnil:
			buf.WriteString(fmt.Sprintf("\t!%d !%d", went.OpGetA(instr), went.OpGetB(instr)))
		case went.OpLoadconst:
			buf.WriteString(fmt.Sprintf("!%d %s", went.OpGetA(instr), f.Consts[went.OpGetBx(instr)]))
		case went.OpUnm, went.OpNot, went.OpCmpl:
			bx := went.OpGetBx(instr)
			bstr := getRegOrConst(bx)
			buf.WriteString(fmt.Sprintf("\t!%d %s", went.OpGetA(instr), bstr))
		case went.OpAdd, went.OpSub, went.OpMul, went.OpDiv, went.OpPow, went.OpShl, went.OpShr,
			went.OpAnd, went.OpOr, went.OpXor, went.OpLe, went.OpLt, went.OpEq, went.OpNe, went.OpGet, went.OpSet:
			a, b, c := went.OpGetA(instr), went.OpGetB(instr), went.OpGetC(instr)
			bstr, cstr := getRegOrConst(b), getRegOrConst(c)
			buf.WriteString(fmt.Sprintf("\t!%d %s %s", a, bstr, cstr))
		case went.OpAppend, went.OpReturn:
			a, b := went.OpGetA(instr), went.OpGetB(instr)
			buf.WriteString(fmt.Sprintf("\t!%d #%d", a, b))
		case went.OpMove:
			a, b := went.OpGetA(instr), went.OpGetB(instr)
			buf.WriteString(fmt.Sprintf("\t!%d !%d", a, b))
		case went.OpLoadglobal, went.OpSetglobal:
			a, bx := went.OpGetA(instr), went.OpGetBx(instr)
			buf.WriteString(fmt.Sprintf("!%d %s", a, f.Consts[bx]))
		case went.OpLoadref, went.OpSetref:
			a, bx := went.OpGetA(instr), went.OpGetBx(instr)
			buf.WriteString(fmt.Sprintf("\t!%d %s", a, f.Consts[bx]))
		case went.OpCall, went.OpCallmethod:
			a, b, c := went.OpGetA(instr), went.OpGetB(instr), went.OpGetC(instr)
			if opcode == went.OpCall {
				buf.WriteString(fmt.Sprintf("\t!%d #%d #%d", a, b, c))
			} else {
				buf.WriteString(fmt.Sprintf("!%d #%d #%d", a, b, c))
			}
		case went.OpArray, went.OpObject:
			buf.WriteString(fmt.Sprintf("\t!%d", went.OpGetA(instr)))
		case went.OpFunc:
			buf.WriteString(fmt.Sprintf("\t!%d &%d", went.OpGetA(instr), went.OpGetBx(instr)))
		case went.OpJmp:
			sbx := went.OpGetsBx(instr)
			buf.WriteString(fmt.Sprintf("\t->%d", pc+sbx))
		case went.OpJmptrue, went.OpJmpfalse:
			a, sbx := went.OpGetA(instr), went.OpGetsBx(instr)
			astr := getRegOrConst(a)
			buf.WriteString(fmt.Sprintf("%s ->%d", astr, pc+sbx))
		case went.OpForbegin:
			a, b := went.OpGetA(instr), went.OpGetB(instr)
			buf.WriteString(fmt.Sprintf("!%d !%d", a, b))
		case went.OpForiter:
			a, b, c := went.OpGetA(instr), went.OpGetB(instr), went.OpGetC(instr)
			buf.WriteString(fmt.Sprintf("\t!%d !%d !%d", a, b, c))
		}

		buf.WriteString("\n")
		doIndent(buf, indent)
	}

	buf.WriteString("}}}\n\n\n")
}

func Disasm(f *went.FuncProto) string {
	var buf bytes.Buffer
	disasmImpl(f, &buf, 0)
	return buf.String()
}
