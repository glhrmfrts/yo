// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>
//
// Main interpreter loop and core operations

package elo

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
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpSetGlobal
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpLoadRef
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpSetRef
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpNeg
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpNot
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpCmpl
      return 0
    },
    opArith, // OpAdd
    opArith, // OpSub
    opArith, // OpMul
    opArith, // OpDiv
    opArith, // OpPow
    func(state *State, cf *callFrame, instr uint32) int { // OpShl
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpShr
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpAnd
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpOr
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpXor
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpLt
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpLe
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpEq
      return 0
    },
    func(state *State, cf *callFrame, instr uint32) int { // OpNe
      return 0
    },
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
  if b > OpConstOffset {
    vb = cf.fn.Proto.Consts[b-OpConstOffset]
  } else {
    vb = cf.r[b]
  }
  if c > OpConstOffset {
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
  default:
    return 0
  }
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