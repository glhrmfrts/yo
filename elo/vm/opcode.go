package vm

type Opcode int

const (
  OP_LOADNIL Opcode = iota  // set range (inclusive) R(A) .. R(B) to nil
  OP_LOADCONST              // set R(A) to K(Bx)

  OP_NEGATE                 // set R(A) to -RK(Bx)
  OP_NOT                    // set R(A) to NOT RK(Bx)
)

const (
  // offset for RK
  kConstOffset = 250
)

var (
  opStrings = map[Opcode]string{
    OP_LOADNIL: "OP_LOADNIL",
    OP_LOADCONST: "OP_LOADCONST",

    OP_NEGATE: "OP_NEGATE",
    OP_NOT: "OP_NOT",
  }
)

// Stringer interface
func (op Opcode) String() string {
  return opStrings[op]
}

// Instruction constructors.
// the instruction design is different from the usual,
// the arguments comes first and the opcode comes last
// in the bits, e.g.:
//
// 8 | 8 | 8 | 8
// c | b | a | op
//
// I'm experimenting this way because we need to get
// the opcode more often than the arguments so it just
// avoids doing bit shifts and the opcode is available
// by just doing $instr & 0xff
//

func opNew(op Opcode) uint32 {
  return uint32(op & 0xff)
}

func opNewA(op Opcode, a int) uint32 {
  return uint32((a & 0xff) << 8 | (int(op) & 0xff))
}

func opNewAB(op Opcode, a, b int) uint32 {
  return uint32(((b & 0xff) << 16) | ((a & 0xff) << 8) | (int(op) & 0xff))
}

func opNewABC(op Opcode, a, b, c int) uint32 {
  return uint32(((c & 0xff) << 24) | ((b & 0xff) << 16) | ((a & 0xff) << 8) | (int(op) & 0xff))
}

func opNewABx(op Opcode, a, b int) uint32 {
  return uint32(((b & 0xffff) << 16) | ((a & 0xff) << 8) | (int(op) & 0xff))
}

func opGetOpcode(instr uint32) Opcode {
  return Opcode(instr & 0xff)
}

func opGetA(instr uint32) int {
  return int((instr >> 8) & 0xff)
}

func opGetB(instr uint32) int {
  return int((instr >> 16) & 0xff)
}

func opGetC(instr uint32) int {
  return int((instr >> 24) & 0xff)
}

func opGetBx(instr uint32) int {
  return int((instr >> 16) & 0xffff)
}