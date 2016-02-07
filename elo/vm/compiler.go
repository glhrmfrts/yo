package vm

import (
  "fmt"
  "strconv"
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
    rega int
    regb int
  }

  // lexical block structure for compiler
  compilerblock struct {
    registerId int
    names      map[string]int
    proto      *FuncProto
  }

  compiler struct {
    lastLine int
    filename string
    mainFunc *FuncProto
    block    *compilerblock
  }
)

func (err *CompileError) Error() string {
  return fmt.Sprintf("%s:%d: error: %s", err.File, err.Line, err.Message)
}


func newCompilerBlock(proto *FuncProto) *compilerblock {
  return &compilerblock{
    proto: proto,
  }
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
    c.lastLine = line
  }
}

func (c *compiler) emitAB(op Opcode, a, b int, line int) {
  c.emitInstruction(opNewAB(op, a, b), line)
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
  var reg, value int
  expr, ok := data.(*exprdata)
  if !ok {
    reg = c.genRegisterId()
  } else {
    reg = expr.rega
  }
  if node.Value {
    value = 1
  }
  c.emitAB(OP_LOADBOOL, reg, value, node.NodeInfo.Line)
}

func (c *compiler) VisitNumber(node *ast.Number, data interface{}) {
  var reg int
  var value Value
  expr, ok := data.(*exprdata)
  if !ok {
    reg = c.genRegisterId()
  } else {
    reg = expr.rega
  }
  if node.Type == ast.T_FLOAT {
    f, err := strconv.ParseFloat(node.Value, 64)
    if err != nil {
      panic(err)
    }
    value = Number(f)
  } else {
    i, err := strconv.Atoi(node.Value)
    if err != nil {
      panic(err)
    }
    value = Number(float64(i))
  }
  c.emitAB(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
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
  c.emitAB(OP_LOADCONST, reg, c.addConst(value), node.NodeInfo.Line)
}

func (c *compiler) VisitId(node *ast.Id, data interface{}) {

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
 
}

func (c *compiler) VisitBinaryExpr(node *ast.BinaryExpr, data interface{}) {

}

func (c *compiler) VisitTernaryExpr(node *ast.TernaryExpr, data interface{}) {

}

func (c *compiler) VisitDeclaration(node *ast.Declaration, data interface{}) {
 
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
  c.block = newCompilerBlock(c.mainFunc)
  
  switch node := root.(type) {
  case *ast.Block:
    for _, stmt := range node.Nodes {
      stmt.Accept(&c, nil)
    }
  default:
    node.Accept(&c, nil)
  }

  res = c.mainFunc
  return
}