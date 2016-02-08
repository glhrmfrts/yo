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
  b.names[name] = info
}


func (c *compiler) error(line int, msg string) {
  panic(&CompileError{Line: line, File: c.filename, Message: msg})
}

func (c *compiler) emitInstruction(instr uint32, line int) {
  f := c.block.proto
  f.Code = append(f.Code, instr)
  f.NumCode++

  if line != c.lastLine {
    f.Lines = append(f.Lines, LineInfo{f.NumCode - 1, uint16(line)})
    f.NumLines++
    c.lastLine = line
  }
}

func (c *compiler) emitAB(op Opcode, a, b, line int) {
  c.emitInstruction(opNewAB(op, a, b), line)
}

func (cc *compiler) emitABC(op Opcode, a, b, c, line int) {
  cc.emitInstruction(opNewABC(op, a, b, c), line)
}

func (c *compiler) emitABx(op Opcode, a, b, line int) {
  c.emitInstruction(opNewABx(op, a, b), line)
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
  var value Value
  expr, ok := data.(*exprdata)
  if !ok {
    reg = c.genRegisterId()
  } else {
    reg = expr.rega
  }
  value = String(node.Value)
  c.emitABx(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitId(node *ast.Id, data interface{}) {
  var reg int
  var scope scope
  expr, exprok := data.(*exprdata)
  if !exprok {
    reg = c.genRegisterId()
  } else {
    reg = expr.rega
  }
  info, ok := c.block.nameInfo(node.Value)
  if ok && info.isConst {
    c.emitABx(OP_LOADCONST, reg, c.addConst(info.value), node.NodeInfo.Line)
    scope = -1
  } else if ok {
    scope = info.scope
  } else {
    // assume global if it can't be found
    scope = kScopeGlobal
  }
  switch scope {
  case kScopeLocal:
    c.emitAB(OP_MOVE, reg, info.reg, node.NodeInfo.Line)
  case kScopeUpval:
    //c.emitAB(OP_LOADUPVAL, reg, c.addUpval(node.Value), node.NodeInfo.Line)
    break
  case kScopeGlobal:
    c.emitABx(OP_LOADGLOBAL, reg, c.addConst(String(node.Value)), node.NodeInfo.Line)
  }
  if exprok && expr.propagate {
    expr.regb = reg
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
    // TODO: generate short-circuit code for && and ||
    
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
    }
    exprdata := exprdata{true, reg, 0}
    node.Left.Accept(c, &exprdata)
    left := exprdata.regb

    exprdata.rega += 1
    node.Right.Accept(c, &exprdata)
    right := exprdata.regb

    c.emitABC(op, reg, left, right, node.NodeInfo.Line)
    if exprok && expr.propagate {
      expr.regb = reg
    }
  }
}

func (c *compiler) VisitTernaryExpr(node *ast.TernaryExpr, data interface{}) {

}

func (c *compiler) VisitDeclaration(node *ast.Declaration, data interface{}) {
  if node.IsConst {
    valuesCount := len(node.Right)
    for i, id := range node.Left {
      if i >= valuesCount {
        c.error(node.NodeInfo.Line, "const without initializer")
      }
      value, ok := c.constFold(node.Right[i])
      if !ok {
        c.error(node.NodeInfo.Line, "const initializer is not a constant")
      }
      _, ok = c.block.names[id.Value]
      if ok {
        c.error(node.NodeInfo.Line, "cannot redeclare const")
      }
      c.block.addNameInfo(id.Value, &nameinfo{true, value, 0, kScopeLocal})
    }
  }
}

func (c *compiler) VisitAssignment(node *ast.Assignment, data interface{}) {
 
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