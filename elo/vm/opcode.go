package vm

// vm instruction details and implementation.
//
// The instruction design is different from the usual,
// the arguments comes first and the opcode comes last
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

const (
  OP_LOADNIL Opcode = iota  // set range (inclusive) R(A) .. R(B) = nil
  OP_LOADCONST              // set R(A) = K(Bx)
  OP_LOADGLOBAL             // set R(A) = globals[K(Bx)]

  // This opcodes might be temporary, if the language turn out to be OO
  OP_NEG                    // set R(A) = -RK(Bx)
  OP_NOT                    // set R(A) = NOT RK(Bx)
  OP_CMPL                   // set R(A) = ^RK(B)

  OP_ADD                    // set R(A) = RK(B) + RK(C)
  OP_SUB                    // set R(A) = RK(B) - RK(C)
  OP_MUL                    // set R(A) = RK(B) * RK(C)
  OP_DIV                    // set R(A) = RK(B) / RK(C)
  OP_POW                    // set R(A) = pow(RK(B), RK(C))
  OP_SHL                    // set R(A) = RK(B) << RK(C)
  OP_SHR                    // set R(A) = RK(B) >> RK(C)
  OP_AND                    // set R(A) = RK(B) & RK(C)
  OP_OR                     // set R(A) = RK(B) | RK(C)
  OP_XOR                    // set R(A) = RK(B) ^ RK(C)
  OP_LT                     // set R(A) = RK(B) < RK(C)
  OP_LE                     // set R(A) = RK(B) <= RK(C)
  OP_EQ                     // set R(A) = RK(B) == RK(C)
  OP_NE                     // set R(A) = RK(B) != RK(C)

  OP_MOVE                   // set R(A) = R(B)
  OP_GET                    // set R(A) = R(B)[RK(C)]
  OP_SET                    // set R(A)[RK(B)] = RK(C)

  OP_CALL                   // set R(A) ... R(A+B-1) = R(A)(R(A+B) ... R(A+B+C-1))

  OP_JMP                    // set pc = pc + Bx
  OP_JMPTRUE                // set pc = pc + Bx if R(A) is not false or nil
  OP_JMPFALSE               // set pc = pc + Bx if R(A) is false or nil
)

// instruction parameters
const (
  kOpcodeMask = 0x3f
  kArgAMask = 0xff
  kArgBCMask = 0x1ff
  kArgBxMask = (0x1ff << 9) | 0x1ff

  kOpcodeSize = 6
  kArgASize = 8
  kArgBCSize = 9

  kArgBOffset = kOpcodeSize + kArgASize
  kArgCOffset = kArgBOffset + kArgBCSize

  // offset for RK
  kConstOffset = 250
)

var (
  opStrings = map[Opcode]string{
    OP_LOADNIL: "LOADNIL",
    OP_LOADCONST: "LOADCONST",
    OP_LOADGLOBAL: "LOADGLOBAL",

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

    OP_CALL: "CALL",

    OP_JMP: "JMP",
    OP_JMPTRUE: "JMPTRUE",
    OP_JMPFALSE: "JMPFALSE",
  }
)

// Stringer interface
func (op Opcode) String() string {
  return opStrings[op]
}


// Instruction constructors.

func opNew(op Opcode) uint32 {
  return uint32(op & kOpcodeMask)
}

func opNewA(op Opcode, a int) uint32 {
  return uint32((a & kArgAMask) << kOpcodeSize | (int(op) & kOpcodeMask))
}

func opNewAB(op Opcode, a, b int) uint32 {
  return uint32(((b & kArgBCMask) << kArgBOffset) | ((a & kArgAMask) << kOpcodeSize) | (int(op) & kOpcodeMask))
}

func opNewABC(op Opcode, a, b, c int) uint32 {
  return uint32(((c & kArgBCMask) << kArgCOffset) | ((b & kArgBCMask) << kArgBOffset) | ((a & kArgAMask) << kOpcodeSize) | (int(op) & kOpcodeMask))
}

func opNewABx(op Opcode, a, b int) uint32 {
  return uint32(((b & kArgBxMask) << kArgBOffset) | ((a & kArgAMask) << kOpcodeSize) | (int(op) & kOpcodeMask))
}

func opGetOpcode(instr uint32) Opcode {
  return Opcode(instr & kOpcodeMask)
}

func opGetA(instr uint32) uint {
  return uint((instr >> kOpcodeSize) & kArgAMask)
}

func opGetB(instr uint32) uint {
  return uint((instr >> kArgBOffset) & kArgBCMask)
}

func opGetC(instr uint32) uint {
  return uint((instr >> kArgCOffset) & kArgBCMask)
}

func opGetBx(instr uint32) uint {
  return uint((instr >> kArgBOffset) & kArgBxMask)
}

func opGetsBx(instr uint32) int {
  return int((instr >> kArgBOffset) & kArgBxMask) 
}