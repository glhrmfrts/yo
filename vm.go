package yo

import (
	"math"
	"github.com/glhrmfrts/yo/parse"
)

const (
	MaxRegisters  = 249
	CallStackSize = 255
)

type opHandler func(*VM, *callFrame, uint32) int

var opTable [kOpCount]opHandler

type callFrame struct {
	pc         int
	line       int
	canRecover bool
	fn         *Func
	r          [MaxRegisters]Value
}

type callFrameStack struct {
	sp    int
	stack [CallStackSize]callFrame
}

func (stack *callFrameStack) New() *callFrame {
	stack.sp += 1
	return &stack.stack[stack.sp-1]
}

func (stack *callFrameStack) Last() *callFrame {
	if stack.sp == 0 {
		return nil
	}
	return &stack.stack[stack.sp-1]
}

type FuncCall struct {
	Args          []Value
	ExpectResults uint
	NumArgs       uint
	NumResults    uint

	results []Value
}

func (c *FuncCall) PushReturnValue(v Value) {
	c.results = append(c.results, v)
	c.NumResults++
}

type VM struct {
	Globals      map[string]Value

	currentFrame *callFrame
	calls        callFrameStack
	error        error
}

func (vm *VM) Define(name string, v Value) {
	vm.Globals[name] = v
}

func (vm *VM) RunString(source []byte, filename string) error {
	nodes, err := parse.ParseFile(source, filename)
	if err != nil {
		return err
	}

	code, err := Compile(nodes, filename)
	if err != nil {
		return err
	}

	return vm.RunBytecode(code)
}

func (vm *VM) RunBytecode(b *Bytecode) error {
	vm.currentFrame = vm.calls.New()
	vm.currentFrame.fn = &Func{b}

	return mainLoop(vm)
}

func NewVM() *VM {
	vm := &VM{
		Globals: make(map[string]Value, 128),
	}

	defineBuiltins(vm)

	return vm
}

func init() {
	opTable = [kOpCount]opHandler{
		func(vm *VM, cf *callFrame, instr uint32) int { // OpLoadNil
			a, b := OpGetA(instr), OpGetB(instr)
			for a <= b {
				cf.r[a] = Nil{}
				a++
			}
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpLoadConst
			a, bx := OpGetA(instr), OpGetBx(instr)
			cf.r[a] = cf.fn.Bytecode.Consts[bx]
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpLoadGlobal
			a, bx := OpGetA(instr), OpGetBx(instr)
			str := cf.fn.Bytecode.Consts[bx].String()
			if g, ok := vm.Globals[str]; ok {
				cf.r[a] = g
			} else {
				//vm.setError("undefined global %s", str)
				return 1
			}
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpSetGlobal
			// TODO: remove this instruction
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpLoadRef
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpSetRef
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpUnm
			a, bx := OpGetA(instr), OpGetBx(instr)
			var bv Value
			if bx >= OpConstOffset {
				bv = cf.fn.Bytecode.Consts[bx-OpConstOffset]
			} else {
				bv = cf.r[bx]
			}
			f, ok := bv.assertFloat64()
			if !ok {
				//vm.setError("cannot perform unary minus on %s", bv.Type())
				return 1
			}
			cf.r[a] = Number(-f)
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpNot
			a, bx := OpGetA(instr), OpGetBx(instr)
			var bv Value
			if bx >= OpConstOffset {
				bv = cf.fn.Bytecode.Consts[bx-OpConstOffset]
			} else {
				bv = cf.r[bx]
			}
			cf.r[a] = Bool(!bv.ToBool())
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpCmpl
			a, bx := OpGetA(instr), OpGetBx(instr)
			var bv Value
			if bx >= OpConstOffset {
				bv = cf.fn.Bytecode.Consts[bx-OpConstOffset]
			} else {
				bv = cf.r[bx]
			}
			f, ok := bv.assertFloat64()
			if !ok || isInt(f) {
				//vm.setError("cannot perform complement on %s", bv.Type())
				return 1
			}
			cf.r[a] = Number(float64(^int(f)))
			return 0
		},
		opArith, // OpAdd
		opArith, // OpSub
		opArith, // OpMul
		opArith, // OpDiv
		opArith, // OpPow
		opArith, // OpShl
		opArith, // OpShr
		opArith, // OpAnd
		opArith, // OpOr
		opArith, // OpXor
		opCmp,   // OpLt
		opCmp,   // OpLe
		opCmp,   // OpEq
		opCmp,   // opNe
		func(vm *VM, cf *callFrame, instr uint32) int { // OpMove
			a, b := OpGetA(instr), OpGetB(instr)
			cf.r[a] = cf.r[b]
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpGetIndex
			a, b, c := OpGetA(instr), OpGetB(instr), OpGetC(instr)
			v := cf.r[b]
			if v.Type() == ValueArray {
				var index Value
				if c >= OpConstOffset {
					index = cf.fn.Bytecode.Consts[c - OpConstOffset]
				} else {
					index = cf.r[c]
				}

				arr := []Value(v.(Array))
				n, ok := index.assertFloat64()
				if !ok {
					// panic
					return 1
				}

				cf.r[a] = arr[int(n)]
			}

			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpSetIndex
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpAppend
			a, b := OpGetA(instr), OpGetB(instr)
			from := a + 1
			to := from + b
			arr := &cf.r[a]
			*arr = append((*arr).(Array), cf.r[from:to]...)
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpCall
			a, b, c := OpGetA(instr), OpGetB(instr), OpGetC(instr)
			fn := cf.r[a]
			if fn.Type() == ValueGoFunc {
				callGoFunc(vm, cf, fn.(GoFunc), a, b, c)
			}
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpCallMethod
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpArray
			cf.r[OpGetA(instr)] = Array([]Value{})
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpObject
			cf.r[OpGetA(instr)] = NewObject(nil, make(map[string]Value))
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpFunc
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpJmp
			cf.pc += OpGetsBx(instr)
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpJmpTrue
			a := OpGetA(instr)
			var val Value
			if a >= OpConstOffset {
				val = cf.fn.Bytecode.Consts[a-OpConstOffset]
			} else {
				val = cf.r[a]
			}
			if val.ToBool() {
				cf.pc += OpGetsBx(instr)
			}
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpJmpFalse
			a := OpGetA(instr)
			var val Value
			if a >= OpConstOffset {
				val = cf.fn.Bytecode.Consts[a-OpConstOffset]
			} else {
				val = cf.r[a]
			}
			if !val.ToBool() {
				cf.pc += OpGetsBx(instr)
			}
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpReturn
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpForBegin
			return 0
		},
		func(vm *VM, cf *callFrame, instr uint32) int { // OpForIter
			return 0
		},
	}
}

func opArith(vm *VM, cf *callFrame, instr uint32) int {
	var vb, vc Value
	a, b, c := OpGetA(instr), OpGetB(instr), OpGetC(instr)
	if b >= OpConstOffset {
		vb = cf.fn.Bytecode.Consts[b-OpConstOffset]
	} else {
		vb = cf.r[b]
	}
	if c >= OpConstOffset {
		vc = cf.fn.Bytecode.Consts[c-OpConstOffset]
	} else {
		vc = cf.r[c]
	}
	fb, okb := vb.assertFloat64()
	fc, okc := vc.assertFloat64()
	if okb && okc {
		cf.r[a] = Number(numberArith(OpGetOpcode(instr), fb, fc))
	}
	return 0
}

func numberArith(op Opcode, a, b float64) float64 {
	switch op {
	case OpAdd:
		return a + b
	case OpSub:
		return a - b
	case OpMul:
		return a * b
	case OpDiv:
		return a / b
	case OpPow:
		return math.Pow(a, b)
	case OpShl:
		return float64(uint(a) << uint(b))
	case OpShr:
		return float64(uint(a) >> uint(b))
	case OpAnd:
		return float64(uint(a) & uint(b))
	case OpOr:
		return float64(uint(a) | uint(b))
	case OpXor:
		return float64(uint(a) ^ uint(b))
	default:
		return 0
	}
}

func opCmp(vm *VM, cf *callFrame, instr uint32) int {
	var vb, vc Value
	a, b, c := OpGetA(instr), OpGetB(instr), OpGetC(instr)
	if b >= OpConstOffset {
		vb = cf.fn.Bytecode.Consts[b-OpConstOffset]
	} else {
		vb = cf.r[b]
	}
	if c >= OpConstOffset {
		vc = cf.fn.Bytecode.Consts[c-OpConstOffset]
	} else {
		vc = cf.r[c]
	}

	if (vb.Type() != ValueNil && vc.Type() != ValueNil) && vb.Type() != vc.Type() {
		cf.r[a] = Bool(false)
		return 0
	}

	op := OpGetOpcode(instr)
	var res bool
	switch vb.Type() {
	case ValueNil:
		if op == OpEq || op == OpNe {
			res = vc.Type() == ValueNil
			if op == OpNe {
				res = !res
			}
		} else {
			// throw error
			return 1
		}
	case ValueBool:
		bb, _ := vb.assertBool()
		bc, _ := vc.assertBool()
		if eq, ne := op == OpEq, op == OpNe; eq || ne {
			res = bb == bc
			if ne {
				res = !res
			}
		} else {
			// throw error
			return 0
		}
	case ValueString:
		sb, sc := vb.String(), vc.String()
		switch op {
		case OpLt:
			res = sb < sc
		case OpLe:
			res = sb <= sc
		case OpEq:
			res = sb == sc
		case OpNe:
			res = sb != sc
		}
	case ValueNumber:
		numb, _ := vb.assertFloat64()
		numc, _ := vb.assertFloat64()
		switch op {
		case OpLt:
			res = numb < numc
		case OpLe:
			res = numb <= numc
		case OpEq:
			res = numb == numc
		case OpNe:
			res = numb != numc
		}
	}
	cf.r[a] = Bool(res)
	return 0
}

func callGoFunc(vm *VM, cf *callFrame, fn GoFunc, a, b, c uint) {
	ab := a + b
	ac := ab + c - 1
	call := FuncCall{
		Args: make([]Value, c),
		ExpectResults: ab - a,
		NumArgs: c,
	}

	for i, r := 0, ab; r <= ac; i, r = i+1, r+1 {
		call.Args[i] = cf.r[r]
	}
	fn(&call)

	nr := call.NumResults
	if nr != call.ExpectResults {
		nr = call.ExpectResults
	}

	for i := uint(0); i < nr; i++ {
		if int(i) >= len(call.results) {
			cf.r[a + 1] = Nil{}
		} else {
			cf.r[a + i] = call.results[i]
		}
	}
}

func mainLoop(vm *VM) error {
	var currentLine uint32
	cf := vm.currentFrame
	proto := cf.fn.Bytecode

	for cf.pc < int(proto.NumCode) {
		if currentLine+1 < proto.NumLines && (cf.pc >= int(proto.Lines[currentLine].Instr)) {
			currentLine += 1
		}

		instr := proto.Code[cf.pc]
		cf.pc++
		cf.line = int(proto.Lines[currentLine].Line)
		if opTable[int(instr&kOpcodeMask)](vm, cf, instr) == 1 {
			return vm.error
		}

		if vm.currentFrame != cf {
			currentLine = 0
		}
		cf = vm.currentFrame
		proto = cf.fn.Bytecode
	}

	return nil
}
