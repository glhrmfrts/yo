// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>
//
// Main interpreter loop and core operations

package went

import (
  "fmt"
  "math"
)

type opHandler func(*State, *callFrame, uint32) int

var opTable [kOpCount]opHandler

func init() {
  opTable = [kOpCount]opHandler{
    func(state *State, cf *callFrame, instr uint32) int { // OpLoadNil
      a, b := OpGetA(instr), OpGetB(instr)
      for a <= b {
        cf.r[a] = Nil{}
        a++
      }
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpLoadConst
      a, bx := OpGetA(instr), OpGetBx(instr)
      cf.r[a] = cf.fn.Proto.Consts[bx]
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpLoadGlobal
      a, bx := OpGetA(instr), OpGetBx(instr)
      str := cf.fn.Proto.Consts[bx].String()
      if g, ok := state.Globals[str]; ok {
        cf.r[a] = g
      } else {
        // throw error
        return 1
      }
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpSetGlobal
      a, bx := OpGetA(instr), OpGetBx(instr)
      str := cf.fn.Proto.Consts[bx].String()
      if _, ok := state.Globals[str]; ok {
        state.Globals[str] = cf.r[a]
      } else {
        // throw error
        return 1
      }
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpLoadRef
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpSetRef
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpUnm
      a, bx := OpGetA(instr), OpGetBx(instr)
      var bv Value
      if bx >= OpConstOffset {
        bv = cf.fn.Proto.Consts[bx-OpConstOffset]
      } else {
        bv = cf.r[bx]
      }
      f, ok := bv.assertFloat64()
      if !ok {
        // throw error
        return 1
      }
      cf.r[a] = Number(-f)
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpNot
      a, bx := OpGetA(instr), OpGetBx(instr)
      var bv Value
      if bx >= OpConstOffset {
        bv = cf.fn.Proto.Consts[bx-OpConstOffset]
      } else {
        bv = cf.r[bx]
      }
      cf.r[a] = Bool(!bv.ToBool())
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpCmpl
      a, bx := OpGetA(instr), OpGetBx(instr)
      var bv Value
      if bx >= OpConstOffset {
        bv = cf.fn.Proto.Consts[bx-OpConstOffset]
      } else {
        bv = cf.r[bx]
      }
      f, ok := bv.assertFloat64()
      if !ok || isInt(f) {
        // throw error
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
    func(state *State, cf *callFrame, instr uint32) int { // OpMove
      a, b := OpGetA(instr), OpGetB(instr)
      cf.r[a] = cf.r[b]
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpGet
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpSet
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpAppend
      return 0  
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpCall
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpCallMethod
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpArray
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpObject
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpFunc
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpJmp
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpJmpTrue
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpJmpFalse
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpReturn
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpForBegin
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpForIter
      return 0
    },
  }
}

func opArith(state *State, cf *callFrame, instr uint32) int {
  var vb, vc Value
  a, b, c := OpGetA(instr), OpGetB(instr), OpGetC(instr)
  if b >= OpConstOffset {
    vb = cf.fn.Proto.Consts[b-OpConstOffset]
  } else {
    vb = cf.r[b]
  }
  if c >= OpConstOffset {
    vc = cf.fn.Proto.Consts[c-OpConstOffset]
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

func opCmp(state *State, cf *callFrame, instr uint32) int {
  var vb, vc Value
  a, b, c := OpGetA(instr), OpGetB(instr), OpGetC(instr)
  if b >= OpConstOffset {
    vb = cf.fn.Proto.Consts[b-OpConstOffset]
  } else {
    vb = cf.r[b]
  }
  if c >= OpConstOffset {
    vc = cf.fn.Proto.Consts[c-OpConstOffset]
  } else {
    vc = cf.r[c]
  }
  if (vb.Type() != ValueNil && vc.Type() != ValueNil) && vb.Type() != vc.Type() {
    // throw error
    return 1
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

func execute(state *State) {
  var currentLine uint32
  cf := state.currentFrame
  proto := cf.fn.Proto

  for cf.pc < int(proto.NumCode) {
    if currentLine + 1 < proto.NumLines && (cf.pc >= int(proto.Lines[currentLine].Instr)) {
      currentLine += 1
    }

    instr := proto.Code[cf.pc]
    cf.pc++
    cf.line = int(proto.Lines[currentLine].Line)
    if opTable[int(instr & kOpcodeMask)](state, cf, instr) == 1 {
      break
    }

    if state.currentFrame != cf {
      currentLine = 0
    }
    cf = state.currentFrame
    proto = cf.fn.Proto
  }

  fmt.Println(cf.r)
  fmt.Println(cf.line)
}