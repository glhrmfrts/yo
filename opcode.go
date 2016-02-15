// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package elo

// vm instruction details and implementation.
//
// The arguments comes first and the opcode comes last
// in the bits, e.g.:
//
// 9 | 9 | 8 | 6
// c | b | a | op
//  bx   | a | op
//
// I'm experimenting this way because we need to get
// the opcode more often than the arguments so it just
// avoids doing bit shifts and the opcode is available
// by just doing $instr & $opcodeMask

type Opcode uint

// NOTE: all register ranges are inclusive
const (
  OpLoadnil Opcode = iota  //  R(A) ... R(B) = nil
  OpLoadconst              //  R(A) = K(Bx)
  OpLoadglobal             //  R(A) = globals[K(Bx)]
  OpSetglobal              //  globals[K(Bx)] = R(A)
  OpLoadref                //  R(A) = refs[K(Bx)]
  OpSetref                 //  refs[K(Bx)] = R(A)

  OpNeg                    //  R(A) = -RK(Bx)
  OpNot                    //  R(A) = NOT RK(Bx)
  OpCmpl                   //  R(A) = ^RK(B)

  OpAdd                    //  R(A) = RK(B) + RK(C)
  OpSub                    //  R(A) = RK(B) - RK(C)
  OpMul                    //  R(A) = RK(B) * RK(C)
  OpDiv                    //  R(A) = RK(B) / RK(C)
  OpPow                    //  R(A) = pow(RK(B), RK(C))
  OpShl                    //  R(A) = RK(B) << RK(C)
  OpShr                    //  R(A) = RK(B) >> RK(C)
  OpAnd                    //  R(A) = RK(B) & RK(C)
  OpOr                     //  R(A) = RK(B) | RK(C)
  OpXor                    //  R(A) = RK(B) ^ RK(C)
  OpLt                     //  R(A) = RK(B) < RK(C)
  OpLe                     //  R(A) = RK(B) <= RK(C)
  OpEq                     //  R(A) = RK(B) == RK(C)
  OpNe                     //  R(A) = RK(B) != RK(C)

  OpMove                   //  R(A) = R(B)
  OpGet                    //  R(A) = R(B)[RK(C)]
  OpSet                    //  R(A)[RK(B)] = RK(C)
  OpAppend                 //  R(A) = append(R(A), R(A+1) ... R(A+B))

  OpCall                   //  R(A) ... R(A+B-1) = R(A)(R(A+B) ... R(A+B+C-1))
  OpCallmethod             //  same as OpCall, but first argument is the receiver
  OpArray                  //  R(A) = []
  OpObject                 //  R(A) = {}
  OpFunc                   //  R(A) = func() { proto = funcs[Bx] }

  OpJmp                    //  pc = pc + sBx
  OpJmptrue                //  pc = pc + sBx if RK(A) is not false or nil
  OpJmpfalse               //  pc = pc + sBx if RK(A) is false or nil
  OpReturn                 //  return R(A) ... R(A+B-1)
  OpForbegin               //  R(A), R(A+1) = objkeys(R(B)), len(objkeys(R(B))) if R(B) is an object
                            //  R(A), R(A+1) = R(B), len(R(B)) if R(B) is an array
                            //  error if not array

  OpForiter                //  R(A) = R(A+1)++ if R(B) is an array
                            //  R(A) = R(C)[R(A+1)++] if R(B) is an object (R(C) should be an array of keys of the object)
)

// instruction parameters
const (
  kOpcodeMask = 0x3f
  kArgAMask = 0xff
  kArgBCMask = 0x1ff
  kArgBxMask = (0x1ff << 9) | 0x1ff
  kArgsBxMask = kArgBxMask >> 1

  kOpcodeSize = 6
  kArgASize = 8
  kArgBCSize = 9

  kArgBOffset = kOpcodeSize + kArgASize
  kArgCOffset = kArgBOffset + kArgBCSize
)

// offset for RK
const OpConstOffset = 250

var (
  opStrings = map[Opcode]string{
    OpLoadnil: "loadnil",
    OpLoadconst: "loadconst",
    OpLoadglobal: "loadglobal",
    OpSetglobal: "setglobal",
    OpLoadref: "loadref",
    OpSetref: "setref",

    OpNeg: "neg",
    OpNot: "not",
    OpCmpl: "cmpl",

    OpAdd: "add",
    OpSub: "sub",
    OpMul: "mul",
    OpDiv: "div",
    OpPow: "pow",                 
    OpShl: "shl",         
    OpShr: "shr",        
    OpAnd: "and",      
    OpOr: "or",      
    OpXor: "xor",         
    OpLt: "lt",         
    OpLe: "le",           
    OpEq: "eq",           
    OpNe: "ne",          

    OpMove: "move",
    OpGet: "get",
    OpSet: "set",
    OpAppend: "append",

    OpCall: "call",
    OpCallmethod: "callmethod",
    OpArray: "array",
    OpObject: "object",
    OpFunc: "func",

    OpJmp: "jmp",
    OpJmptrue: "jmptrue",
    OpJmpfalse: "jmpfalse",
    OpReturn: "return",
    OpForbegin: "forbegin",
    OpForiter: "foriter",
  }
)

// Stringer interface
func (op Opcode) String() string {
  return opStrings[op]
}


// Instruction constructors.

func OpNew(op Opcode) uint32 {
  return uint32(op & kOpcodeMask)
}

func OpNewA(op Opcode, a int) uint32 {
  return uint32((a & kArgAMask) << kOpcodeSize | (int(op) & kOpcodeMask))
}

func OpNewAB(op Opcode, a, b int) uint32 {
  return uint32(((b & kArgBCMask) << kArgBOffset) | ((a & kArgAMask) << kOpcodeSize) | (int(op) & kOpcodeMask))
}

func OpNewABC(op Opcode, a, b, c int) uint32 {
  return uint32(((c & kArgBCMask) << kArgCOffset) | ((b & kArgBCMask) << kArgBOffset) | ((a & kArgAMask) << kOpcodeSize) | (int(op) & kOpcodeMask))
}

func OpNewABx(op Opcode, a, b int) uint32 {
  return uint32(((b & kArgBxMask) << kArgBOffset) | ((a & kArgAMask) << kOpcodeSize) | (int(op) & kOpcodeMask))
}

func OpNewAsBx(op Opcode, a, b int) uint32 {
  return OpNewABx(op, a, b + kArgsBxMask)
}

func OpGetOpcode(instr uint32) Opcode {
  return Opcode(instr & kOpcodeMask)
}

func OpGetA(instr uint32) uint {
  return uint((instr >> kOpcodeSize) & kArgAMask)
}

func OpGetB(instr uint32) uint {
  return uint((instr >> kArgBOffset) & kArgBCMask)
}

func OpGetC(instr uint32) uint {
  return uint((instr >> kArgCOffset) & kArgBCMask)
}

func OpGetBx(instr uint32) uint {
  return uint((instr >> kArgBOffset) & kArgBxMask)
}

func OpGetsBx(instr uint32) int {
  return int(OpGetBx(instr)) - kArgsBxMask
}