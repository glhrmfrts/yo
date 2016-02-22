// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

package went

import (
  "fmt"
  "math"
  "github.com/glhrmfrts/went/ast"
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
  blockContext int

  nameInfo struct {
    isConst bool
    value   Value // only set if isConst == true
    reg     int
    scope   scope
    block   *compilerBlock
  }

  loopInfo struct {
    breaks         []uint32
    continues      []uint32
    breakTarget    uint32
    continueTarget uint32
  }

  // lexical block structure for compiler
  compilerBlock struct {
    context  blockContext
    register int
    names    map[string]*nameInfo
    loop     *loopInfo
    proto    *FuncProto
    parent   *compilerBlock
  }

  compiler struct {
    lastLine int
    filename string
    mainFunc *FuncProto
    block    *compilerBlock
  }
)

// names lexical scopes
const (
  kScopeLocal scope = iota
  kScopeRef
  kScopeGlobal
)

// blocks context
const (
  kBlockContextFunc blockContext = iota
  kBlockContextLoop
  kBlockContextBranch
)

// How much registers an array can use at one time
// when it's created in literal form (see VisitArray)
const kArrayMaxRegisters = 10


func (err *CompileError) Error() string {
  return fmt.Sprintf("%s:%d: %s", err.File, err.Line, err.Message)
}

// compilerBlock

func newCompilerBlock(proto *FuncProto, context blockContext, parent *compilerBlock) *compilerBlock {
  return &compilerBlock{
    proto: proto,
    context: context,
    parent: parent,
    names: make(map[string]*nameInfo, 128),
  }
}

func (b *compilerBlock) nameInfo(name string) (*nameInfo, bool) {
  var closures int
  block := b
  for block != nil {
    info, ok := block.names[name]
    if ok {
      if closures > 0 && info.scope == kScopeLocal {
        info.scope = kScopeRef
      }
      return info, true
    }
    if block.context == kBlockContextFunc {
      closures++
    }
    block = block.parent
  }

  return nil, false
}

func (b *compilerBlock) addNameInfo(name string, info *nameInfo) {
  info.block = b
  b.names[name] = info
}

// compiler

func (c *compiler) error(line int, msg string) {
  panic(&CompileError{Line: line, File: c.filename, Message: msg})
}

func (c *compiler) emitInstruction(instr uint32, line int) int {
  f := c.block.proto
  f.Code = append(f.Code, instr)
  f.NumCode++

  if line != c.lastLine || f.NumLines == 0 {
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
  return c.emitInstruction(OpNewAB(op, a, b), line)
}

func (cc *compiler) emitABC(op Opcode, a, b, c, line int) int {
  return cc.emitInstruction(OpNewABC(op, a, b, c), line)
}

func (c *compiler) emitABx(op Opcode, a, b, line int) int {
  return c.emitInstruction(OpNewABx(op, a, b), line)
}

func (c *compiler) emitAsBx(op Opcode, a, b, line int) int {
  return c.emitInstruction(OpNewAsBx(op, a, b), line)
}

func (c *compiler) modifyABx(index int, op Opcode, a, b int) bool {
  return c.modifyInstruction(index, OpNewABx(op, a, b))
}

func (c *compiler) modifyAsBx(index int, op Opcode, a, b int) bool {
  return c.modifyInstruction(index, OpNewAsBx(op, a, b))
}

func (c *compiler) newLabel() uint32 {
  return c.block.proto.NumCode
}

func (c *compiler) labelOffset(label uint32) int {
  return int(c.block.proto.NumCode - label)
}

func (c *compiler) genRegister() int {
  id := c.block.register
  c.block.register++
  return id
}

func (c *compiler) declareLocalVar(name string, reg int) {
  if _, ok := c.block.names[name]; ok {
    c.error(c.lastLine, fmt.Sprintf("cannot redeclare '%s'", name))
  }
  c.block.addNameInfo(name, &nameInfo{false, nil, reg, kScopeLocal, c.block})
}

func (c *compiler) enterBlock(context blockContext) {
  assert(c.block != nil, "c.block enterBlock")
  block := newCompilerBlock(c.block.proto, context, c.block)
  block.register = c.block.register

  if context == kBlockContextLoop {
    block.loop = &loopInfo{}
  } else if c.block.loop != nil {
    block.loop = c.block.loop
  }
  c.block = block
}

func (c *compiler) leaveBlock() {
  block := c.block
  if block.context == kBlockContextLoop {
    loop := block.loop
    for _, index := range loop.breaks {
      c.modifyAsBx(int(index), OpJmp, 0, int(loop.breakTarget - index - 1))
    }
    for _, index := range loop.continues {
      c.modifyAsBx(int(index), OpJmp, 0, int(loop.continueTarget - index - 1))
    }
  }
  c.block = block.parent
}

func (c *compiler) insideLoop() bool {
  block := c.block
  for block != nil {
    if block.context == kBlockContextLoop {
      return true
    }
    if block.context == kBlockContextFunc {
      return false
    }
    block = block.parent
  }
  return false
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

// Try to "constant fold" an expression
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
  case *ast.CallExpr:
    id, ok := t.Left.(*ast.Id)
    if ok && len(t.Args) == 1 {
      rv, ok := c.constFold(t.Args[0])
      if !ok {
        return nil, false
      }
      if id.Value == "string" {
        return String(rv.String()), true
      } else if id.Value == "number" {
        switch rv.Type() {
        case ValueNumber:
          return rv, true
        case ValueString:
          n, err := parseNumber(rv.String())
          if err != nil {
            return nil, false
          }
          return Number(n), true
        }
      } else if id.Value == "bool" {
        return Bool(rv.ToBool()), true
      }
    }
  case *ast.UnaryExpr:
    if t.Op == ast.TokenMinus {
      val, ok := c.constFold(t.Right)
      if ok && val.Type() == ValueNumber {
        f64, _ := val.assertFloat64()
        return Number(-f64), true
      }
      return nil, false
    } else {
      // 'not' operator
      val, ok := c.constFold(t.Right)
      if ok && val.Type() == ValueBool {
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
      case ast.TokenPlus:
        ret = Number(lf64 + rf64)
      case ast.TokenMinus:
        ret = Number(lf64 - rf64)
      case ast.TokenTimes:
        ret = Number(lf64 * rf64)
      case ast.TokenDiv:
        ret = Number(lf64 / rf64)
      case ast.TokenTimestimes:
        ret = Number(math.Pow(lf64, rf64))
      case ast.TokenLt:
        ret = Bool(lf64 < rf64)
      case ast.TokenLteq:
        ret = Bool(lf64 <= rf64)
      case ast.TokenGt:
        ret = Bool(lf64 > rf64)
      case ast.TokenGteq:
        ret = Bool(lf64 >= rf64)
      case ast.TokenEqeq:
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
      case ast.TokenAmpamp:
        return Bool(lb && rb), true
      case ast.TokenPipepipe:
        return Bool(lb || rb), true
      }

    stringOps:
      ls, ok := left.assertString()
      rs, _ := right.assertString()
      if !ok {
        return nil, false
      }

      switch t.Op {
      case ast.TokenPlus:
        return String(ls + rs), true
      case ast.TokenLt:
        ret = Bool(ls < rs)
      case ast.TokenLteq:
        ret = Bool(ls <= rs)
      case ast.TokenGt:
        ret = Bool(ls > rs)
      case ast.TokenGteq:
        ret = Bool(ls >= rs)
      case ast.TokenEqeq:
        return Bool(ls == rs), true
      case ast.TokenBangeq:
        return Bool(ls != rs), true
      }
    }
  }
  return nil, false
}

// declare local variables
// assignments are done in sequence, since the registers are created as needed
func (c *compiler) declare(names []*ast.Id, values []ast.Node) {
  var isCall, isUnpack bool
  nameCount, valueCount := len(names), len(values)
  if valueCount > 0 {
    _, isCall = values[valueCount - 1].(*ast.CallExpr)
    _, isUnpack = values[valueCount - 1].(*ast.VarArg)
  }
  start := c.block.register
  end := start + nameCount - 1
  for i, id := range names {
    _, ok := c.block.names[id.Value]
    if ok {
      c.error(id.NodeInfo.Line, fmt.Sprintf("cannot redeclare '%s'", id.Value))
    }
    reg := c.genRegister()
    c.block.addNameInfo(id.Value, &nameInfo{false, nil, reg, kScopeLocal, c.block})

    exprdata := exprdata{false, reg, reg}
    if i == valueCount - 1 && (isCall || isUnpack) {
      // last expression receives all the remaining registers
      // in case it's a function call with multiple return values
      rem := i + 1
      for rem < nameCount {
        // reserve the registers
        id := names[rem]
        _, ok := c.block.names[id.Value]
        if ok {
          c.error(id.NodeInfo.Line, fmt.Sprintf("cannot redeclare '%s'", id.Value))
        }
        end = c.genRegister()
        c.block.addNameInfo(id.Value, &nameInfo{false, nil, end, kScopeLocal, c.block})
        rem++
      }
      exprdata.regb, start = end, end + 1
      values[i].Accept(c, &exprdata)
      break
    }
    if i < valueCount {
      values[i].Accept(c, &exprdata)
      start = reg + 1
    }
  }
  if end >= start {
    // variables without initializer are set to nil
    c.emitAB(OpLoadnil, start, end, names[0].NodeInfo.Line)
  }
}

func (c *compiler) assignmentHelper(left ast.Node, assignReg int, valueReg int) {
  switch v := left.(type) {
  case *ast.Id:
    var scope scope
    info, ok := c.block.nameInfo(v.Value)
    if !ok {
      scope = kScopeGlobal
    } else {
      scope = info.scope
    }
    switch scope {
    case kScopeLocal:
      c.emitAB(OpMove, info.reg, valueReg, v.NodeInfo.Line)
    case kScopeRef, kScopeGlobal:
      op := OpSetglobal
      if scope == kScopeRef {
        op = OpSetref
      }
      c.emitABx(op, valueReg, c.addConst(String(v.Value)), v.NodeInfo.Line)
    }
  case *ast.Subscript:
    arrData := exprdata{true, assignReg, assignReg}
    v.Left.Accept(c, &arrData)
    arrReg := arrData.regb

    subData := exprdata{true, assignReg, assignReg}
    v.Right.Accept(c, &subData)
    subReg := subData.regb
    c.emitABC(OpSet, arrReg, subReg, valueReg, v.NodeInfo.Line)
  case *ast.Selector:
    objData := exprdata{true, assignReg, assignReg}
    v.Left.Accept(c, &objData)
    objReg := objData.regb
    key := OpConstOffset + c.addConst(String(v.Value))

    c.emitABC(OpSet, objReg, key, valueReg, v.NodeInfo.Line)
  }
}

func (c *compiler) branchConditionHelper(cond, then, else_ ast.Node, reg int) {
  ternaryData := exprdata{true, reg + 1, reg + 1}
  cond.Accept(c, &ternaryData)
  condr := ternaryData.regb
  jmpInstr := c.emitAsBx(OpJmpfalse, condr, 0, c.lastLine)
  thenLabel := c.newLabel()

  ternaryData = exprdata{false, reg, reg}
  then.Accept(c, &ternaryData)
  c.modifyAsBx(jmpInstr, OpJmpfalse, condr, c.labelOffset(thenLabel))

  if else_ != nil {
    successInstr := c.emitAsBx(OpJmp, 0, 0, c.lastLine)
    
    elseLabel := c.newLabel()
    ternaryData = exprdata{false, reg, reg}
    else_.Accept(c, &ternaryData)

    c.modifyAsBx(successInstr, OpJmp, 0, c.labelOffset(elseLabel))
  }
}

func (c *compiler) functionReturnGuard() {
  last := c.block.proto.Code[c.block.proto.NumCode-1]
  if OpGetOpcode(last) != OpReturn {
    c.emitAB(OpReturn, 0, 0, c.lastLine)
  }
}

//
// visitor interface
//

func (c *compiler) VisitNil(node *ast.Nil, data interface{}) {
  var rega, regb int
  expr, ok := data.(*exprdata)
  if ok {
    rega, regb = expr.rega, expr.regb
    if rega > regb {
      regb = rega
    }
  } else {
    rega = c.genRegister()
    regb = rega
  }
  c.emitAB(OpLoadnil, rega, regb, node.NodeInfo.Line)
}

func (c *compiler) VisitBool(node *ast.Bool, data interface{}) {
  var reg int
  value := Bool(node.Value)
  expr, ok := data.(*exprdata)
  if ok && expr.propagate {
    expr.regb = OpConstOffset + c.addConst(value)
    return
  } else if ok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  c.emitABx(OpLoadconst, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitNumber(node *ast.Number, data interface{}) {
  var reg int
  value := Number(node.Value)
  expr, ok := data.(*exprdata)
  if ok && expr.propagate {
    expr.regb = OpConstOffset + c.addConst(value)
    return
  } else if ok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  c.emitABx(OpLoadconst, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitString(node *ast.String, data interface{}) {
  var reg int
  value := String(node.Value)
  expr, ok := data.(*exprdata)
  if ok && expr.propagate {
    expr.regb = OpConstOffset + c.addConst(value)
    return
  } else if ok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  c.emitABx(OpLoadconst, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitId(node *ast.Id, data interface{}) {
  var reg int
  var scope scope = -1
  expr, exprok := data.(*exprdata)
  if !exprok {
    reg = c.genRegister()
  } else {
    reg = expr.rega
  }
  info, ok := c.block.nameInfo(node.Value)
  if ok && info.isConst {
    if exprok && expr.propagate {
      expr.regb = OpConstOffset + c.addConst(info.value)
      return
    }
    c.emitABx(OpLoadconst, reg, c.addConst(info.value), node.NodeInfo.Line)
  } else if ok {
    scope = info.scope
  } else {
    // assume global if it can't be found in the lexical scope
    scope = kScopeGlobal
  }
  switch scope {
  case kScopeLocal:
    if exprok && expr.propagate {
      expr.regb = info.reg
      return
    }
    c.emitAB(OpMove, reg, info.reg, node.NodeInfo.Line)
  case kScopeRef, kScopeGlobal:
    op := OpLoadglobal
    if scope == kScopeRef {
      op = OpLoadref
    }
    c.emitABx(op, reg, c.addConst(String(node.Value)), node.NodeInfo.Line)
    if exprok && expr.propagate {
      expr.regb = reg
    }
  }
}

func (c *compiler) VisitArray(node *ast.Array, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  length := len(node.Elements)
  c.emitAB(OpArray, reg, 0, node.NodeInfo.Line)

  times := length / kArrayMaxRegisters + 1
  for t := 0; t < times; t++ {
    start, end := t * kArrayMaxRegisters, (t+1) * kArrayMaxRegisters
    end = int(math.Min(float64(end - start), float64(length - start)))
    if end == 0 {
      break
    }
    for i := 0; i < end; i++ {
      el := node.Elements[start + i]
      exprdata := exprdata{false, reg + i + 1, reg + i + 1}
      el.Accept(c, &exprdata)
    }
    c.emitAB(OpAppend, reg, end, node.NodeInfo.Line)
  }
  if exprok && expr.propagate {
    expr.regb = reg
  }
}

func (c *compiler) VisitObjectField(node *ast.ObjectField, data interface{}) {
  expr, exprok := data.(*exprdata)
  assert(exprok, "ObjectField exprok")
  objreg := expr.rega
  key := OpConstOffset + c.addConst(String(node.Key))

  valueData := exprdata{true, objreg + 1, objreg + 1}
  node.Value.Accept(c, &valueData)
  value := valueData.regb

  c.emitABC(OpSet, objreg, key, value, node.NodeInfo.Line)
}

func (c *compiler) VisitObject(node *ast.Object, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  c.emitAB(OpObject, reg, 0, node.NodeInfo.Line)
  for _, field := range node.Fields {
    fieldData := exprdata{false, reg, reg}
    field.Accept(c, &fieldData)
  }
  if exprok && expr.propagate {
    expr.regb = reg
  }
}

func (c *compiler) VisitFunction(node *ast.Function, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  parent := c.block.proto
  proto := newFuncProto(parent.Source)

  block := newCompilerBlock(proto, kBlockContextFunc, c.block)
  c.block = block

  index := int(parent.NumFuncs)
  parent.Funcs = append(parent.Funcs, proto)
  parent.NumFuncs++

  // insert 'this' into scope
  c.declareLocalVar("this", c.genRegister())

  // insert arguments into scope
  for _, n := range node.Args {
    switch arg := n.(type) {
    case *ast.Id:
      reg := c.genRegister()
      c.block.addNameInfo(arg.Value, &nameInfo{false, nil, reg, kScopeLocal, c.block})
    }
  }

  node.Body.Accept(c, nil)
  c.functionReturnGuard()

  c.block = c.block.parent
  c.emitABx(OpFunc, reg, index, node.NodeInfo.Line)

  if node.Name != nil {
    switch name := node.Name.(type) {
    case *ast.Id:
      c.declareLocalVar(name.Value, reg)
    default:  
      c.assignmentHelper(name, reg + 1, reg)
    }
  }
  if exprok && expr.propagate {
    expr.regb = reg
  }
}

func (c *compiler) VisitSelector(node *ast.Selector, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  objData := exprdata{true, reg + 1, reg + 1}
  node.Left.Accept(c, &objData)
  objReg := objData.regb

  key := OpConstOffset + c.addConst(String(node.Value))
  c.emitABC(OpGet, reg, objReg, key, node.NodeInfo.Line)
  if exprok && expr.propagate {
    expr.regb = objReg
  }
}

func (c *compiler) VisitSubscript(node *ast.Subscript, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  arrData := exprdata{true, reg + 1, reg + 1}
  node.Left.Accept(c, &arrData)
  arrReg := arrData.regb

  _, ok := node.Right.(*ast.Slice)
  if ok {
    // TODO: generate code for slice
    return
  }

  indexData := exprdata{true, reg + 1, reg + 1}
  node.Right.Accept(c, &indexData)
  indexReg := indexData.regb
  c.emitABC(OpGet, reg, arrReg, indexReg, node.NodeInfo.Line)

  if exprok && expr.propagate {
    expr.regb = reg
  }
}

func (c *compiler) VisitSlice(node *ast.Slice, data interface{}) {

}

func (c *compiler) VisitKwArg(node *ast.KwArg, data interface{}) {
  
}

func (c *compiler) VisitVarArg(node *ast.VarArg, data interface{}) {

}

func (c *compiler) VisitCallExpr(node *ast.CallExpr, data interface{}) {
  var startReg, endReg, resultCount int
  expr, exprok := data.(*exprdata)
  if exprok {
    startReg, endReg = expr.rega, expr.regb
    resultCount = endReg - startReg + 1
  } else {
    startReg = c.genRegister()
    endReg = startReg
    resultCount = 1
  }

  // check if it's a type conversion (string, number, bool)
  v, ok := c.constFold(node)
  if ok {
    c.emitABx(OpLoadconst, startReg, c.addConst(v), node.NodeInfo.Line)
    return
  }

  argCount := len(node.Args)
  var op Opcode
  switch node.Left.(type) {
  case *ast.Selector:
    op = OpCallmethod
    callerData := exprdata{true, startReg, startReg}
    node.Left.Accept(c, &callerData)
    objReg := callerData.regb

    // insert object as first argument
    endReg += 1
    argCount += 1
    c.emitAB(OpMove, endReg, objReg, node.NodeInfo.Line)
  default:
    op = OpCall
    callerData := exprdata{false, startReg, startReg}
    node.Left.Accept(c, &callerData)
  }

  for i, arg := range node.Args {
    reg := endReg + i + 1
    argData := exprdata{false, reg, reg}
    arg.Accept(c, &argData)
  }

  c.emitABC(op, startReg, resultCount, argCount, node.NodeInfo.Line)
}

func (c *compiler) VisitPostfixExpr(node *ast.PostfixExpr, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  var op Opcode
  switch node.Op {
  case ast.TokenPlusplus:
    op = OpAdd
  case ast.TokenMinusminus:
    op = OpSub
  }
  leftdata := exprdata{true, reg, reg}
  node.Left.Accept(c, &leftdata)
  left := leftdata.regb
  one := OpConstOffset + c.addConst(Number(1))

  // don't bother moving if we're not in an expression
  if exprok {
    c.emitAB(OpMove, reg, left, node.NodeInfo.Line)
  }
  c.emitABC(op, left, left, one, node.NodeInfo.Line)
}

func (c *compiler) VisitUnaryExpr(node *ast.UnaryExpr, data interface{}) {
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  value, ok := c.constFold(node)
  if ok {
    if exprok && expr.propagate {
      expr.regb = OpConstOffset + c.addConst(value)
      return
    }
    c.emitABx(OpLoadconst, reg, c.addConst(value), node.NodeInfo.Line)
  } else if ast.IsPostfixOp(node.Op) {
    op := OpAdd
    if node.Op == ast.TokenMinusminus {
      op = OpSub
    }
    exprdata := exprdata{true, reg, reg}
    node.Right.Accept(c, &exprdata)
    one := OpConstOffset + c.addConst(Number(1))
    c.emitABC(op, exprdata.regb, exprdata.regb, one, node.NodeInfo.Line)

    // don't bother moving if we're not in an expression
    if exprok {
      c.emitAB(OpMove, reg, exprdata.regb, node.NodeInfo.Line)
    }
  } else {
    var op Opcode
    switch node.Op {
    case ast.TokenMinus:
      op = OpUnm
    case ast.TokenNot, ast.TokenBang:
      op = OpNot
    case ast.TokenTilde:
      op = OpCmpl
    }
    exprdata := exprdata{true, reg, reg}
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
    reg = c.genRegister()
  }
  value, ok := c.constFold(node)
  if ok {
    if exprok && expr.propagate {
      expr.regb = OpConstOffset + c.addConst(value)
      return
    }
    c.emitABx(OpLoadconst, reg, c.addConst(value), node.NodeInfo.Line)
  } else {
    if isAnd, isOr := node.Op == ast.TokenAmpamp, node.Op == ast.TokenPipepipe; isAnd || isOr {
      var op Opcode
      if isAnd {
        op = OpJmpfalse
      } else {
        op = OpJmptrue
      }
      exprdata := exprdata{true, reg, reg}
      node.Left.Accept(c, &exprdata)
      left := exprdata.regb

      jmpInstr := c.emitAsBx(op, left, 0, node.NodeInfo.Line)
      size := c.block.proto.NumCode

      exprdata.propagate = false
      node.Right.Accept(c, &exprdata)
      c.modifyAsBx(jmpInstr, op, left, int(c.block.proto.NumCode - size))
      return
    }
    
    var op Opcode
    switch node.Op {
    case ast.TokenPlus:
      op = OpAdd
    case ast.TokenMinus:
      op = OpSub
    case ast.TokenTimes:
      op = OpMul
    case ast.TokenDiv:
      op = OpDiv
    case ast.TokenTimestimes:
      op = OpPow
    case ast.TokenLtlt:
      op = OpShl
    case ast.TokenGtgt:
      op = OpShr
    case ast.TokenAmp:
      op = OpAnd
    case ast.TokenPipe:
      op = OpOr
    case ast.TokenTilde:
      op = OpXor
    case ast.TokenLt, ast.TokenGteq:
      op = OpLt
    case ast.TokenLteq, ast.TokenGt:
      op = OpLe
    case ast.TokenEq:
      op = OpEq
    case ast.TokenBangeq:
      op = OpNe
    }

    exprdata := exprdata{true, reg, 0}
    node.Left.Accept(c, &exprdata)
    left := exprdata.regb

    // temp register for right expression
    exprdata.rega += 1
    node.Right.Accept(c, &exprdata)
    right := exprdata.regb

    if node.Op == ast.TokenGt || node.Op == ast.TokenGteq {
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
  var reg int
  expr, exprok := data.(*exprdata)
  if exprok {
    reg = expr.rega
  } else {
    reg = c.genRegister()
  }
  c.branchConditionHelper(node.Cond, node.Then, node.Else, reg)
}

func (c *compiler) VisitDeclaration(node *ast.Declaration, data interface{}) {
  valueCount := len(node.Right)
  if node.IsConst {
    for i, id := range node.Left {
      _, ok := c.block.names[id.Value]
      if ok {
        c.error(node.NodeInfo.Line, fmt.Sprintf("cannot redeclare '%s'", id.Value))
      }
      if i >= valueCount {
        c.error(node.NodeInfo.Line, fmt.Sprintf("const '%s' without initializer", id.Value))
      }
      value, ok := c.constFold(node.Right[i])
      if !ok {
        c.error(node.NodeInfo.Line, fmt.Sprintf("const '%s' initializer is not a constant", id.Value))
      }
      c.block.addNameInfo(id.Value, &nameInfo{true, value, 0, kScopeLocal, c.block})
    }
    return
  }
  c.declare(node.Left, node.Right)
}

func (c *compiler) VisitAssignment(node *ast.Assignment, data interface{}) {
  if node.Op == ast.TokenColoneq {
    // short variable declaration
    var names []*ast.Id
    for _, id := range node.Left {
      names = append(names, id.(*ast.Id))
    }
    c.declare(names, node.Right)
    return
  }
  // regular assignment, if the left-side is an identifier
  // then it has to be declared already
  varCount, valueCount := len(node.Left), len(node.Right)
  _, isCall := node.Right[valueCount - 1].(*ast.CallExpr)
  _, isUnpack := node.Right[valueCount - 1].(*ast.VarArg)
  start := c.block.register
  current := start
  end := start + varCount - 1

  // evaluate all expressions first with temp registers
  for i, _ := range node.Left {
    reg := start + i
    exprdata := exprdata{false, reg, reg}
    if i == valueCount - 1 && (isCall || isUnpack) {
      exprdata.regb, current = end, end
      node.Right[i].Accept(c, &exprdata)
      break
    }
    if i < valueCount {
      node.Right[i].Accept(c, &exprdata)
      current = reg + 1
    }
  }
  // assign the results to the variables
  for i, variable := range node.Left {
    valueReg := start + i

    // don't touch variables without a corresponding value
    if valueReg >= current {
      break
    }
    c.assignmentHelper(variable, current + 1, valueReg)
  }
}

func (c *compiler) VisitBranchStmt(node *ast.BranchStmt, data interface{}) {
  if !c.insideLoop() {
    c.error(node.NodeInfo.Line, fmt.Sprintf("%s outside loop", node.Type))
  }
  instr := c.emitAsBx(OpJmp, 0, 0, node.NodeInfo.Line)
  switch node.Type {
  case ast.TokenContinue:
    c.block.loop.continues = append(c.block.loop.continues, uint32(instr))
  case ast.TokenBreak:
    c.block.loop.breaks = append(c.block.loop.breaks, uint32(instr))
  }
}

func (c *compiler) VisitReturnStmt(node *ast.ReturnStmt, data interface{}) {
  start := c.block.register
  for _, v := range node.Values {
    reg := c.genRegister()
    data := exprdata{false, reg, reg}
    v.Accept(c, &data)
  }
  c.emitAB(OpReturn, start, len(node.Values), node.NodeInfo.Line)
}

func (c *compiler) VisitIfStmt(node *ast.IfStmt, data interface{}) {
  _, ok := data.(*exprdata)
  if !ok {
    c.enterBlock(kBlockContextBranch)
    defer c.leaveBlock()
  }
  if node.Init != nil {
    node.Init.Accept(c, nil)
  }
  c.branchConditionHelper(node.Cond, node.Body, node.Else, c.block.register)
}

func (c *compiler) VisitForIteratorStmt(node *ast.ForIteratorStmt, data interface{}) {
  c.enterBlock(kBlockContextLoop)
  defer c.leaveBlock()
  
  arrReg := c.genRegister()
  lenReg := c.genRegister()
  keyReg := c.genRegister()
  idxReg := c.genRegister()
  valReg := c.genRegister()
  colReg := c.genRegister()

  collectionData := exprdata{false, colReg, colReg}
  node.Collection.Accept(c, &collectionData)
  c.emitAB(OpForbegin, arrReg, colReg, node.NodeInfo.Line)
  c.emitABx(OpLoadconst, idxReg, c.addConst(Number(0)), c.lastLine)

  if node.Value == nil {
    c.declareLocalVar(node.Key.Value, valReg)
  } else {
    c.declareLocalVar(node.Key.Value, keyReg)
    c.declareLocalVar(node.Value.Value, valReg)
  }

  testLabel := c.newLabel()
  testReg := c.block.register
  c.emitABC(OpLt, testReg, idxReg, lenReg, c.lastLine)
  jmpInstr := c.emitAsBx(OpJmpfalse, testReg, 0, c.lastLine)

  c.emitABC(OpForiter, keyReg, colReg, arrReg, node.NodeInfo.Line)
  c.emitABC(OpGet, valReg, colReg, keyReg, c.lastLine)

  node.Body.Accept(c, nil)
  c.block.loop.continueTarget = c.newLabel()

  c.emitAsBx(OpJmp, 0, -c.labelOffset(testLabel) - 1, c.lastLine)
  c.block.loop.breakTarget = c.newLabel()

  c.modifyAsBx(jmpInstr, OpJmpfalse, testReg, c.labelOffset(uint32(jmpInstr) + 1))
}

func (c *compiler) VisitForStmt(node *ast.ForStmt, data interface{}) {
  c.enterBlock(kBlockContextLoop)
  defer c.leaveBlock()

  hasCond := node.Cond != nil
  if node.Init != nil {
    node.Init.Accept(c, nil)
  }
 
  startLabel := c.newLabel()

  var cond, jmpInstr int
  var jmpLabel uint32
  if hasCond {
    reg := c.block.register
    condData := exprdata{true, reg, reg}
    node.Cond.Accept(c, &condData)

    cond = condData.regb
    jmpInstr = c.emitAsBx(OpJmpfalse, cond, 0, c.lastLine)
    jmpLabel = c.newLabel()
  }

  node.Body.Accept(c, nil)
  c.block.loop.continueTarget = c.newLabel()

  if node.Step != nil {
    node.Step.Accept(c, nil)
    c.block.register -= 1 // discard register consumed by Step
  } else {
    c.block.loop.continueTarget = startLabel // saves one jump
  }

  c.emitAsBx(OpJmp, 0, -c.labelOffset(startLabel) - 1, c.lastLine)

  if hasCond {
    c.modifyAsBx(jmpInstr, OpJmpfalse, cond, c.labelOffset(jmpLabel))
  }
  c.block.loop.breakTarget = c.newLabel()
}

func (c *compiler) VisitBlock(node *ast.Block, data interface{}) {
  for _, stmt := range node.Nodes {
    stmt.Accept(c, nil)

    if !ast.IsStmt(stmt) {
      c.block.register -= 1
    }
  }
}

// Compile receives the root node of the AST and generates code
// for the "main" function from it.
// Any type of Node is accepted, either a block representing the program
// or a single expression.
//
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
  c.block = newCompilerBlock(c.mainFunc, kBlockContextFunc, nil)
  
  root.Accept(&c, nil)
  c.functionReturnGuard()

  res = c.mainFunc
  return
}