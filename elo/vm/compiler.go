package vm

import (
  "fmt"
  "github.com/glhrmfrts/elo-lang/elo/ast"
)

type CompileError struct {
  Line    int
  File    string
  Message string
}

type symboltable map[string]int


// lexical block structure for compiler
type compilerblock struct {
  registerId int
  names      symboltable
  proto      *FuncProto
}

type compiler struct {
  line     int
  filename string
  mainFunc *FuncProto
  block    *compilerblock
}


func (err *CompileError) Error() string {
  return fmt.Sprintf("%s:%d: error: %s", err.File, err.Line, err.Message)
}


func newCompilerBlock(proto *FuncProto) *compilerblock {
  return &compilerblock{
    proto: proto,
  }
}

func (b *compilerblock) genRegisterId() int {
  id := b.registerId
  b.registerId++
  return id
}


func (c *compiler) error(msg string) {
  panic(&CompileError{Line: c.line, File: c.filename, Message: msg})
}

func (c *compiler) emitAB(op Opcode, a, b int, line int) {
  c.block.proto.AddInstruction(opNewAB(op, a, b), line)
}

// visitor interface functions

func (c *compiler) VisitNil(node *ast.Nil) {
  // load nil to register
  reg := c.block.genRegisterId()
  c.emitAB(OP_LOADNIL, reg, 0, node.NodeInfo.Line)
}

func (c *compiler) VisitBool(node *ast.Bool) {

}

func (c *compiler) VisitNumber(node *ast.Number) {

}

func (c *compiler) VisitId(node *ast.Id) {

}

func (c *compiler) VisitString(node *ast.String) {

}

func (c *compiler) VisitArray(node *ast.Array) {

}

func (c *compiler) VisitObjectField(node *ast.ObjectField) {

}

func (c *compiler) VisitObject(node *ast.Object) {

}

func (c *compiler) VisitFunction(node *ast.Function) {
 
}

func (c *compiler) VisitSelector(node *ast.Selector) {
 
}

func (c *compiler) VisitSubscript(node *ast.Subscript) {

}

func (c *compiler) VisitSlice(node *ast.Slice) {

}

func (c *compiler) VisitKwArg(node *ast.KwArg) {
  
}

func (c *compiler) VisitVarArg(node *ast.VarArg) {

}

func (c *compiler) VisitCallExpr(node *ast.CallExpr) {

}

func (c *compiler) VisitUnaryExpr(node *ast.UnaryExpr) {
 
}

func (c *compiler) VisitBinaryExpr(node *ast.BinaryExpr) {

}

func (c *compiler) VisitTernaryExpr(node *ast.TernaryExpr) {

}

func (c *compiler) VisitDeclaration(node *ast.Declaration) {
 
}

func (c *compiler) VisitAssignment(node *ast.Assignment) {
 
}

func (c *compiler) VisitBranchStmt(node *ast.BranchStmt) {

}

func (c *compiler) VisitReturnStmt(node *ast.ReturnStmt) {

}

func (c *compiler) VisitIfStmt(node *ast.IfStmt) {
 
}

func (c *compiler) VisitForIteratorStmt(node *ast.ForIteratorStmt) {

}

func (c *compiler) VisitForStmt(node *ast.ForStmt) {

}

func (c *compiler) VisitBlock(node *ast.Block) {
  
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
      stmt.Accept(&c)
    }
  default:
    node.Accept(&c)
  }

  res = c.mainFunc
  return
}