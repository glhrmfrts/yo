package vm

import (
  "fmt"
  "math"
  "github.com/glhrmfrts/elo-lang/elo/ast"
)

type (
  CompileError struct {
    Line    int
    File    string
    Message string
  }

  // holds registers for a expression
  exprdata struct {
    propagate bool
    rega      int // rega is default for write
    regb      int // regb is default for read
  }

  // lexical scope of a name
  scope int

  // lexical context of a block, (function, loop, branch...)
  blockcontext int

  // information of a name in the program
  nameinfo struct {
    isConst bool
    value   Value // only set if isConst == true
    reg     int
    scope   scope
    block   *compilerblock
  }

  // lexical block structure for compiler
  compilerblock struct {
    context    blockcontext
    registerId int
    names      map[string]*nameinfo
    proto      *FuncProto
    parent     *compilerblock
  }

  compiler struct {
    lastLine int
    filename string
    mainFunc *FuncProto
    block    *compilerblock
  }
)

// names lexical scopes
const (
  kScopeLocal scope = iota
  kScopeUpval
  kScopeGlobal
)

// blocks context
const (
  kContextFunc blockcontext = iota
  kContextLoop
  kContextBranch
)

func (err *CompileError) Error() string {
  return fmt.Sprintf("%s:%d: error: %s", err.File, err.Line, err.Message)
}


func newCompilerBlock(proto *FuncProto, context blockcontext) *compilerblock {
  return &compilerblock{
    proto: proto,
    context: context,
    names: make(map[string]*nameinfo, 128),
  }
}

func (b *compilerblock) nameInfo(name string) (*nameinfo, bool) {
  var closures int
  block := b
  for block != nil {
    info, ok := block.names[name]
    if ok {
      if closures > 0 {
        info.scope = kScopeUpval
      }
      return info, true
    }
    if block.context == kContextFunc {
      closures++
    }
    block = block.parent
  }

  return nil, false
}

func (b *compilerblock) addNameInfo(name string, info *nameinfo) {
  info.block = b
  b.names[name] = info
}


func (c *compiler) error(line int, msg string) {
  panic(&CompileError{Line: line, File: c.filename, Message: msg})
}

func (c *compiler) emitInstruction(instr uint32, line int) int {
  f := c.block.proto
  f.Code = append(f.Code, instr)
  f.NumCode++

  if line != c.lastLine {
    f.Lines = append(f.Lines, LineInfo{f.NumCode - 1, uint16(line)})
    f.NumLines++
    c.lastLine = line
  }
  return int(f.NumCode - 1)
}

func (c *compiler) modifyInstruction(index int, instr uint32) bool {
  f := c.block.proto
  if uint32(index) < f.NumCode {
    f.Code[index] = instr
    return true
  }
  return false
}

func (c *compiler) emitAB(op Opcode, a, b, line int) int {
  return c.emitInstruction(opNewAB(op, a, b), line)
}

func (cc *compiler) emitABC(op Opcode, a, b, c, line int) int {
  return cc.emitInstruction(opNewABC(op, a, b, c), line)
}

func (c *compiler) emitABx(op Opcode, a, b, line int) int {
  return c.emitInstruction(opNewABx(op, a, b), line)
}

func (c *compiler) modifyABx(index int, op Opcode, a, b int) bool {
  return c.modifyInstruction(index, opNewABx(op, a, b))
}

func (c *compiler) genRegisterId() int {
  id := c.block.registerId
  c.block.registerId++
  fmt.Printf("genRegisterId: %d\n", id)
  return id
}

// Add a constant to the current prototype's constant pool
// and return it's index
func (c *compiler) addConst(value Value) int {
  f := c.block.proto
  valueType := value.Type()
  for i, c := range f.Consts {
    if c.Type() == valueType && c == value {
      return i
    }
  }
  if f.NumConsts > funcMaxConsts - 1 {
    c.error(0, "too many constants") // should never happen
  }
  f.Consts = append(f.Consts, value)
  f.NumConsts++
  return int(f.NumConsts - 1)
}

// try to "constant fold" an expression
func (c *compiler) constFold(node ast.Node) (Value, bool) {
  switch t := node.(type) {
  case *ast.Number:
    return Number(t.Value), true
  case *ast.Bool:
    return Bool(t.Value), true
  case *ast.String:
    return String(t.Value), true
  case *ast.Id:
    info, ok := c.block.nameInfo(t.Value)
    if ok && info.isConst {
      return info.value, true
    }
  case *ast.UnaryExpr:
    if t.Op == ast.T_MINUS {
      val, ok := c.constFold(t.Right)
      if ok && val.Type() == VALUE_NUMBER {
        f64, _ := val.assertFloat64()
        return Number(-f64), true
      }
      return nil, false
    } else {
      // 'not' operator
      val, ok := c.constFold(t.Right)
      if ok && val.Type() == VALUE_BOOL {
        bool_, _ := val.assertBool()
        return Bool(!bool_), true
      }
      return nil, false
    }
  case *ast.BinaryExpr:
    left, leftOk := c.constFold(t.Left)
    right, rightOk := c.constFold(t.Right)
    if leftOk && rightOk {
      var ret Value
      if left.Type() != right.Type() {
        return nil, false
      }
      lf64, ok := left.assertFloat64()
      rf64, _ := right.assertFloat64()
      if !ok {
        goto boolOps
      }

      // first check all arithmetic/relational operations
      switch t.Op {
      case ast.T_PLUS:
        ret = Number(lf64 + rf64)
      case ast.T_MINUS:
        ret = Number(lf64 - rf64)
      case ast.T_TIMES:
        ret = Number(lf64 * rf64)
      case ast.T_DIV:
        ret = Number(lf64 / rf64)
      case ast.T_TIMESTIMES:
        ret = Number(math.Pow(lf64, rf64))
      case ast.T_LT:
        ret = Bool(lf64 < rf64)
      case ast.T_LTEQ:
        ret = Bool(lf64 <= rf64)
      case ast.T_GT:
        ret = Bool(lf64 > rf64)
      case ast.T_GTEQ:
        ret = Bool(lf64 >= rf64)
      case ast.T_EQEQ:
        ret = Bool(lf64 == rf64)
      }
      if ret != nil {
        return ret, true
      }

    boolOps:
      // not arithmetic/relational, maybe logic?
      lb, ok := left.assertBool()
      rb, _ := right.assertBool()
      if !ok {
        goto stringOps
      }

      switch t.Op {
      case ast.T_AMPAMP:
        return Bool(lb && rb), true
      case ast.T_PIPEPIPE:
        return Bool(lb || rb), true
      }

    stringOps:
      ls, ok := left.assertString()
      rs, _ := right.assertString()
      if !ok {
        return nil, false
      }

      switch t.Op {
      case ast.T_EQEQ:
        return Bool(ls == rs), true
      case ast.T_BANGEQ:
        return Bool(ls != rs), true
      }
    }
  }
  return nil, false
}

func (c *compiler) VisitNil(node *ast.Nil, data interface{}) {
  var rega, regb int
  expr, ok := data.(*exprdata)
  if ok {
    rega, regb = expr.rega, expr.regb
    if rega > regb {
      regb = rega
    }
  } else {
    rega = c.genRegisterId()
    regb = rega
  }
  c.emitAB(OP_LOADNIL, rega, regb, node.NodeInfo.Line)
}

func (c *compiler) VisitBool(node *ast.Bool, data interface{}) {
  var reg int
  value := Bool(node.Value)
  expr, ok := data.(*exprdata)
  if ok && expr.propagate {
    expr.regb = kConstOffset + c.addConst(value)
    return
  } else if ok {
    reg = expr.rega
  } else {
    reg = c.genRegisterId()
  }
  c.emitABx(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitNumber(node *ast.Number, data interface{}) {
  var reg int
  value := Number(node.Value)
  expr, ok := data.(*exprdata)
  if ok && expr.propagate {
    expr.regb = kConstOffset + c.addConst(value)
    return
  } else if ok {
    reg = expr.rega
  } else {
    reg = c.genRegisterId()
  }
  c.emitABx(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitString(node *ast.String, data interface{}) {
  var reg int
  value := String(node.Value)
  expr, ok := data.(*exprdata)
  if ok && expr.propagate {
    expr.regb = kConstOffset + c.addConst(value)
    return
  } else if ok {
    reg = expr.rega
  } else {
    reg = c.genRegisterId()
  }
  c.emitABx(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitId(node *ast.Id, data interface{}) {
  var reg int
  var scope scope = -1
  expr, exprok := data.(*exprdata)
  if !exprok {
    reg = c.genRegisterId()
  } else {
    reg = expr.rega
  }
  info, ok := c.block.nameInfo(node.Value)
  if ok && info.isConst {
    if exprok && expr.propagate {
      expr.regb = kConstOffset + c.addConst(info.value)
      return
    }
    c.emitABx(OP_LOADCONST, reg, c.addConst(info.value), node.NodeInfo.Line)
  } else if ok {
    scope = info.scope
  } else {
    // assume global if it can't be found
    scope = kScopeGlobal
  }
  switch scope {
  case kScopeLocal:
    if exprok && expr.propagate {
      expr.regb = info.reg
      return
    }
    c.emitAB(OP_MOVE, reg, info.reg, node.NodeInfo.Line)
  case kScopeUpval, kScopeGlobal:
    c.emitABx(OP_LOADGLOBAL, reg, c.addConst(String(node.Value)), node.NodeInfo.Line)
    if exprok && expr.propagate {
      expr.regb = reg
    }
  }
}

func (c *compiler) VisitArray(node *ast.Array, data interface{}) {

}

func (c *compiler) VisitObjectField(node *ast.ObjectField, data interface{}) {

}

func (c *compiler) VisitObject(node *ast.Object, data interface{}) {

}

func (c *compiler) VisitFunction(node *ast.Function, data interface{}) {
 
}

func (c *compiler) VisitSelector(node *ast.Selector, data interface{}) {
 
}

func (c *compiler) VisitSubscript(node *ast.Subscript, data interface{}) {

}

func (c *compiler) VisitSlice(node *ast.Slice, data interface{}) {

}

func (c *compiler) VisitKwArg(node *ast.KwArg, data interface{}) {
  
}

func (c *compiler) VisitVarArg(node *ast.VarArg, data interface{}) {

}

func (c *compiler) VisitCallExpr(node *ast.CallExpr, data interface{}) {

}

func (c *compiler) VisitUnaryExpr(node *ast.UnaryExpr, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegisterId()
  }
  value, ok := c.constFold(node)
  if ok {
    if exprok && expr.propagate {
      expr.regb = kConstOffset + c.addConst(value)
      return
    }
    c.emitABx(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
  } else {
    var op Opcode
    switch node.Op {
    case ast.T_MINUS:
      op = OP_NEGATE
    case ast.T_NOT, ast.T_BANG:
      op = OP_NOT
    }
    exprdata := exprdata{true, 0, 0}
    node.Right.Accept(c, &exprdata)
    c.emitABx(op, reg, exprdata.regb, node.NodeInfo.Line)
    if exprok && expr.propagate {
      expr.regb = reg
    }
  }
}

func (c *compiler) VisitBinaryExpr(node *ast.BinaryExpr, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegisterId()
  }
  value, ok := c.constFold(node)
  if ok {
    if exprok && expr.propagate {
      expr.regb = kConstOffset + c.addConst(value)
      return
    }
    c.emitABx(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
  } else {
    if isAnd, isOr := node.Op == ast.T_AMPAMP, node.Op == ast.T_PIPEPIPE; isAnd || isOr {
      var op Opcode
      if isAnd {
        op = OP_JMPFALSE
      } else {
        op = OP_JMPTRUE
      }
      exprdata := exprdata{true, reg, 0}
      node.Left.Accept(c, &exprdata)
      left := exprdata.regb

      jmpInstr := c.emitABx(op, left, 0, node.NodeInfo.Line)
      size := c.block.proto.NumCode

      exprdata.propagate = false
      node.Right.Accept(c, &exprdata)
      c.modifyABx(jmpInstr, op, left, int(c.block.proto.NumCode - size) + 1)
      return
    }
    
    var op Opcode
    switch node.Op {
    case ast.T_PLUS:
      op = OP_ADD
    case ast.T_MINUS:
      op = OP_SUB
    case ast.T_TIMES:
      op = OP_MUL
    case ast.T_DIV:
      op = OP_DIV
    case ast.T_LT, ast.T_GTEQ:
      op = OP_LT
    case ast.T_LTEQ, ast.T_GT:
      op = OP_LE
    case ast.T_EQ:
      op = OP_EQ
    case ast.T_BANGEQ:
      op = OP_NEQ
    }
    
    exprdata := exprdata{true, reg, 0}
    node.Left.Accept(c, &exprdata)
    left := exprdata.regb

    // temp register for right expression
    exprdata.rega += 1
    node.Right.Accept(c, &exprdata)
    right := exprdata.regb

    if node.Op == ast.T_GT || node.Op == ast.T_GTEQ {
      // invert operands
      c.emitABC(op, reg, right, left, node.NodeInfo.Line)  
    } else {
      c.emitABC(op, reg, left, right, node.NodeInfo.Line)
    }

    if exprok && expr.propagate {
      expr.regb = reg
    }
  }
}

func (c *compiler) VisitTernaryExpr(node *ast.TernaryExpr, data interface{}) {

}

// VisitDeclaration generates code for variable declaration.
// For consts declaration no code is generated, they are only kept
// in the current block's local symbol table.
func (c *compiler) VisitDeclaration(node *ast.Declaration, data interface{}) {
  namesCount, valuesCount := len(node.Left), len(node.Right)
  if node.IsConst {
    for i, id := range node.Left {
      _, ok := c.block.names[id.Value]
      if ok {
        c.error(node.NodeInfo.Line, fmt.Sprintf("cannot redeclare '%s'", id.Value))
      }
      if i >= valuesCount {
        c.error(node.NodeInfo.Line, fmt.Sprintf("const '%s' without initializer", id.Value))
      }
      value, ok := c.constFold(node.Right[i])
      if !ok {
        c.error(node.NodeInfo.Line, fmt.Sprintf("const '%s' initializer is not a constant", id.Value))
      }
      c.block.addNameInfo(id.Value, &nameinfo{true, value, 0, kScopeLocal, c.block})
    }
  } else {
    // declare local variables
    // TODO: multiple return values
    if valuesCount + namesCount > 2 {
      c.error(node.NodeInfo.Line, fmt.Sprintf("multiple variable declaration is not yet supported"))
    }
    start := c.block.registerId
    end := start - 1
    for i, id := range node.Left {
      _, ok := c.block.names[id.Value]
      if ok {
        c.error(node.NodeInfo.Line, fmt.Sprintf("cannot redeclare '%s'", id.Value))
      }
      reg := c.genRegisterId()
      if i >= valuesCount {
        end = reg
      } else {
        // load value into new register
        exprdata := exprdata{false, reg, 0}
        node.Right[i].Accept(c, &exprdata)
        start = reg + 1
      }
      c.block.addNameInfo(id.Value, &nameinfo{false, nil, reg, kScopeLocal, c.block})
    }
    if end >= start {
      // variables without initializer are set to nil
      c.emitAB(OP_LOADNIL, start, end, node.NodeInfo.Line)
    }
  }
}

func (c *compiler) VisitAssignment(node *ast.Assignment, data interface{}) {
  namesCount, valuesCount := len(node.Left), len(node.Right)
  if valuesCount + namesCount > 2 {
    c.error(node.NodeInfo.Line, fmt.Sprintf("multiple assignment is not yet supported"))
  }
  value := node.Right[0]
  if node.Op == ast.T_COLONEQ {
    // short variable declaration
    id := node.Left[0].(*ast.Id)
    _, ok := c.block.names[id.Value]
    if ok {
      c.error(node.NodeInfo.Line, fmt.Sprintf("cannot redeclare '%s'", id.Value))
    }
    reg := c.genRegisterId()
    exprdata := exprdata{false, reg, 0}
    value.Accept(c, &exprdata)
    c.block.addNameInfo(id.Value, &nameinfo{false, nil, reg, kScopeLocal, c.block})
    return
  }
  // regular assignment, if the left-side is an identifier
  // it has to be declared already
  id, ok := node.Left[0].(*ast.Id)
  if ok {
    var scope scope
    info, ok := c.block.names[id.Value]
    if !ok {
      scope = kScopeGlobal
    } else {
      scope = info.scope
    }
    switch scope {
    case kScopeLocal:
      exprdata := exprdata{false, info.reg, info.reg}
      value.Accept(c, &exprdata)
    case kScopeUpval, kScopeGlobal:
      // temp register for the value
      reg := c.block.registerId + 1
      exprdata := exprdata{true, reg, reg}
      value.Accept(c, &exprdata)
      // c.emitABx(OP_SETGLOBAL, exprdata.regb, c.addConst(String(id.Value)), node.NodeInfo.Line)
    }
    return
  }
  // TODO: object key and array index assignment
}

func (c *compiler) VisitBranchStmt(node *ast.BranchStmt, data interface{}) {

}

func (c *compiler) VisitReturnStmt(node *ast.ReturnStmt, data interface{}) {

}

func (c *compiler) VisitIfStmt(node *ast.IfStmt, data interface{}) {
 
}

func (c *compiler) VisitForIteratorStmt(node *ast.ForIteratorStmt, data interface{}) {

}

func (c *compiler) VisitForStmt(node *ast.ForStmt, data interface{}) {

}

func (c *compiler) VisitBlock(node *ast.Block, data interface{}) {
  for _, stmt := range node.Nodes {
    stmt.Accept(c, nil)

    if !ast.IsStmt(stmt) {
      c.block.registerId -= 1
    }
  }
}

// Compile receives the root node of the AST and generates code 
// for the "main" function from it.
// Any type of Node is accepted, either a block representing the program
// or a single expression.
func Compile(root ast.Node, filename string) (res *FuncProto, err error) {
  defer func() {
    if r := recover(); r != nil {
      if cerr, ok := r.(*CompileError); ok {
        err = cerr
      } else {
        panic(r)
      }
    }
  }()

  var c compiler
  c.filename = filename
  c.mainFunc = newFuncProto(filename)
  c.block = newCompilerBlock(c.mainFunc, kContextFunc)
  
  root.Accept(&c, nil)

  res = c.mainFunc
  return
}