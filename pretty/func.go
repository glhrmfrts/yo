// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>
//
// "Disassemble" a function prototype (bytecode chunk)
// into a human-readable string

package pretty

import (
	"bytes"
	"fmt"
	"github.com/glhrmfrts/yo"
)

func doIndent(buf *bytes.Buffer, indent int) {
	for indent > 0 {
		buf.WriteString(" ")
		indent--
	}
}

func disasmImpl(f *yo.Bytecode, buf *bytes.Buffer, indent int) {
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
		if a >= yo.OpConstOffset {
			return fmt.Sprint(f.Consts[a-yo.OpConstOffset])
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
		opcode := yo.OpGetOpcode(instr)
		pc := i + 1

		if lineChanged || i == 0 {
			buf.WriteString(fmt.Sprint(line.Line, "\t"))
		} else {
			buf.WriteString("\t")
		}

		buf.WriteString(fmt.Sprint(i))
		buf.WriteString(fmt.Sprint("\t", opcode, "\t"))

		switch opcode {
		case yo.OpLoadnil:
			buf.WriteString(fmt.Sprintf("\t!%d !%d", yo.OpGetA(instr), yo.OpGetB(instr)))
		case yo.OpLoadconst:
			buf.WriteString(fmt.Sprintf("!%d %s", yo.OpGetA(instr), f.Consts[yo.OpGetBx(instr)]))
		case yo.OpUnm, yo.OpNot, yo.OpCmpl:
			bx := yo.OpGetBx(instr)
			bstr := getRegOrConst(bx)
			buf.WriteString(fmt.Sprintf("\t!%d %s", yo.OpGetA(instr), bstr))
		case yo.OpAdd, yo.OpSub, yo.OpMul, yo.OpDiv, yo.OpPow, yo.OpShl, yo.OpShr,
			yo.OpAnd, yo.OpOr, yo.OpXor, yo.OpLe, yo.OpLt, yo.OpEq, yo.OpNe,
			yo.OpGetIndex, yo.OpSetIndex:
			a, b, c := yo.OpGetA(instr), yo.OpGetB(instr), yo.OpGetC(instr)
			bstr, cstr := getRegOrConst(b), getRegOrConst(c)
			buf.WriteString(fmt.Sprintf("\t!%d %s %s", a, bstr, cstr))
		case yo.OpAppend, yo.OpReturn:
			a, b := yo.OpGetA(instr), yo.OpGetB(instr)
			buf.WriteString(fmt.Sprintf("\t!%d #%d", a, b))
		case yo.OpMove:
			a, b := yo.OpGetA(instr), yo.OpGetB(instr)
			buf.WriteString(fmt.Sprintf("\t!%d !%d", a, b))
		case yo.OpLoadglobal, yo.OpSetglobal:
			a, bx := yo.OpGetA(instr), yo.OpGetBx(instr)
			buf.WriteString(fmt.Sprintf("!%d %s", a, f.Consts[bx]))
		case yo.OpLoadFree, yo.OpSetFree:
			a, bx := yo.OpGetA(instr), yo.OpGetBx(instr)
			buf.WriteString(fmt.Sprintf("\t!%d %s", a, f.Consts[bx]))
		case yo.OpCall, yo.OpCallmethod:
			a, b, c := yo.OpGetA(instr), yo.OpGetB(instr), yo.OpGetC(instr)
			if opcode == yo.OpCall {
				buf.WriteString(fmt.Sprintf("\t!%d #%d #%d", a, b, c))
			} else {
				buf.WriteString(fmt.Sprintf("!%d #%d #%d", a, b, c))
			}
		case yo.OpArray, yo.OpObject:
			buf.WriteString(fmt.Sprintf("\t!%d", yo.OpGetA(instr)))
		case yo.OpFunc:
			buf.WriteString(fmt.Sprintf("\t!%d &%d", yo.OpGetA(instr), yo.OpGetBx(instr)))
		case yo.OpJmp:
			sbx := yo.OpGetsBx(instr)
			buf.WriteString(fmt.Sprintf("\t->%d", pc+sbx))
		case yo.OpJmptrue, yo.OpJmpfalse:
			a, sbx := yo.OpGetA(instr), yo.OpGetsBx(instr)
			astr := getRegOrConst(a)
			buf.WriteString(fmt.Sprintf("%s ->%d", astr, pc+sbx))
		case yo.OpForbegin:
			a, b := yo.OpGetA(instr), yo.OpGetB(instr)
			buf.WriteString(fmt.Sprintf("!%d !%d", a, b))
		case yo.OpForiter:
			a, b, c := yo.OpGetA(instr), yo.OpGetB(instr), yo.OpGetC(instr)
			buf.WriteString(fmt.Sprintf("\t!%d !%d !%d", a, b, c))
		}

		buf.WriteString("\n")
		doIndent(buf, indent)
	}

	buf.WriteString("}}}\n\n\n")
}

func Disasm(f *yo.Bytecode) string {
	var buf bytes.Buffer
	disasmImpl(f, &buf, 0)
	return buf.String()
}
