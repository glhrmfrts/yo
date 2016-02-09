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
  OP_LOADNIL Opcode = iota  // set range (inclusive) R(A) .. R(B) to nil
  OP_LOADCONST              // set R(A) to K(Bx)

  // This opcodes might be temporary, if the language turn out to be OO
  OP_NEGATE                 // set R(A) to -RK(Bx)
  OP_NOT                    // set R(A) to NOT RK(Bx)

  OP_ADD                    // set R(A) to RK(B) + RK(C)
  OP_SUB                    // set R(A) to RK(B) - RK(C)
  OP_MUL                    // set R(A) to RK(B) * RK(C)
  OP_DIV                    // set R(A) to RK(B) / RK(C)
  OP_POW                    // set R(A) to pow(RK(B), RK(C))
  OP_SHL                    // set R(A) to RK(B) << RK(C)
  OP_SHR                    // set R(A) to RK(B) >> RK(C)
  OP_AND                    // set R(A) to RK(B) & RK(C)
  OP_OR                     // set R(A) to RK(B) | RK(C)
  OP_XOR                    // set R(A) to RK(B) ^ RK(C)
  OP_CMPL                   // set R(A) to ^RK(B)
  OP_LT                     // set R(A) to RK(B) < RK(C)
  OP_LE                     // set R(A) to RK(B) <= RK(C)
  OP_EQ                     // set R(A) to RK(B) == RK(C)
  OP_NEQ                    // set R(A) to RK(B) != RK(C)

  OP_MOVE                   // set R(A) to R(B)
  OP_LOADGLOBAL             // set R(A) to globals[K(Bx)]

  OP_JMP                    // set pc to pc + Bx
  OP_JMPTRUE                // set pc to pc + Bx if R(A) is not false or nil
  OP_JMPFALSE               // set pc to pc + Bx if R(A) is false or nil
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
    OP_LOADNIL: "OP_LOADNIL",
    OP_LOADCONST: "OP_LOADCONST",

    OP_NEGATE: "OP_NEGATE",
    OP_NOT: "OP_NOT",

    OP_ADD: "OP_ADD",
    OP_SUB: "OP_SUB",
    OP_MUL: "OP_MUL",
    OP_DIV: "OP_DIV",

    OP_MOVE: "OP_MOVE",
    OP_LOADGLOBAL: "OP_LOADGLOBAL",

    OP_JMP: "OP_JMP",
    OP_JMPTRUE: "OP_JMPTRUE",
    OP_JMPFALSE: "OP_JMPFALSE",
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