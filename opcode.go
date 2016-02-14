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
  OP_LOADNIL Opcode = iota  //  R(A) ... R(B) = nil
  OP_LOADCONST              //  R(A) = K(Bx)
  OP_LOADGLOBAL             //  R(A) = globals[K(Bx)]
  OP_SETGLOBAL              //  globals[K(Bx)] = R(A)
  OP_LOADREF                //  R(A) = refs[K(Bx)]
  OP_SETREF                 //  refs[K(Bx)] = R(A)

  OP_NEG                    //  R(A) = -RK(Bx)
  OP_NOT                    //  R(A) = NOT RK(Bx)
  OP_CMPL                   //  R(A) = ^RK(B)

  OP_ADD                    //  R(A) = RK(B) + RK(C)
  OP_SUB                    //  R(A) = RK(B) - RK(C)
  OP_MUL                    //  R(A) = RK(B) * RK(C)
  OP_DIV                    //  R(A) = RK(B) / RK(C)
  OP_POW                    //  R(A) = pow(RK(B), RK(C))
  OP_SHL                    //  R(A) = RK(B) << RK(C)
  OP_SHR                    //  R(A) = RK(B) >> RK(C)
  OP_AND                    //  R(A) = RK(B) & RK(C)
  OP_OR                     //  R(A) = RK(B) | RK(C)
  OP_XOR                    //  R(A) = RK(B) ^ RK(C)
  OP_LT                     //  R(A) = RK(B) < RK(C)
  OP_LE                     //  R(A) = RK(B) <= RK(C)
  OP_EQ                     //  R(A) = RK(B) == RK(C)
  OP_NE                     //  R(A) = RK(B) != RK(C)

  OP_MOVE                   //  R(A) = R(B)
  OP_GET                    //  R(A) = R(B)[RK(C)]
  OP_SET                    //  R(A)[RK(B)] = RK(C)
  OP_APPEND                 //  R(A) = append(R(A), R(A+1) ... R(A+B))

  OP_CALL                   //  R(A) ... R(A+B-1) = R(A)(R(A+B) ... R(A+B+C-1))
  OP_ARRAY                  //  R(A) = []
  OP_OBJECT                 //  R(A) = {}
  OP_FUNC                   //  R(A) = func() { proto = funcs[Bx] }

  OP_JMP                    //  pc = pc + sBx
  OP_JMPTRUE                //  pc = pc + sBx if RK(A) is not false or nil
  OP_JMPFALSE               //  pc = pc + sBx if RK(A) is false or nil
  OP_RETURN                 //  return R(A) ... R(A+B-1)
  OP_FORBEGIN               // R(A), R(A+1) = objkeys(R(B)), len(objkeys(R(B))) if R(B) is an object
                            // R(A), R(A+1) = R(B), len(R(B)) if R(B) is an array
                            // error if not array

  OP_FORITER                // R(A) = R(A+1)++ if R(B) is an array
                            // R(A) = R(C)[R(A+1)++] if R(B) is an object (R(C) should be an array of keys of the object)
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
    OP_LOADNIL: "LOADNIL",
    OP_LOADCONST: "LOADCONST",
    OP_LOADGLOBAL: "LOADGLOBAL",
    OP_SETGLOBAL: "SETGLOBAL",
    OP_LOADREF: "LOADREF",
    OP_SETREF: "SETREF",

    OP_NEG: "NEG",
    OP_NOT: "NOT",
    OP_CMPL: "CMPL",

    OP_ADD: "ADD",
    OP_SUB: "SUB",
    OP_MUL: "MUL",
    OP_DIV: "DIV",
    OP_POW: "POW",                 
    OP_SHL: "SHL",         
    OP_SHR: "SHR",        
    OP_AND: "AND",      
    OP_OR: "OR",      
    OP_XOR: "XOR",         
    OP_LT: "LT",         
    OP_LE: "LE",           
    OP_EQ: "EQ",           
    OP_NE: "NE",          

    OP_MOVE: "MOVE",
    OP_GET: "GET",
    OP_SET: "SET",
    OP_APPEND: "APPEND",

    OP_CALL: "CALL",
    OP_ARRAY: "ARRAY",
    OP_OBJECT: "OBJECT",
    OP_FUNC: "FUNC",

    OP_JMP: "JMP",
    OP_JMPTRUE: "JMPTRUE",
    OP_JMPFALSE: "JMPFALSE",
    OP_RETURN: "RETURN",
    OP_FORBEGIN: "FORBEGIN",
    OP_FORITER: "FORITER",
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