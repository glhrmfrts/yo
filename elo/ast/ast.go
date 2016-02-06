// Abstract Syntax Tree

package ast

import (
  "github.com/glhrmfrts/elo-lang/elo/token"
)

type (
  Node interface {
    Accept(v Visitor)
  }

  //
  // expressions
  //

  Nil struct {  
  }

  Bool struct {
    Value bool
  }

  Number struct {
    Type  token.Token // int or float
    Value string
  }

  Id struct {
    Value string
  }

  String struct {
    Value string
  }

  Array struct {
    Values []Node
  }

  ObjectField struct {
    Key   Node
    Value Node
  }

  Object struct {
    Fields []*ObjectField
  }

  Function struct {
    Name Node
    Args []Node
    Body Node
  }

  Selector struct {
    Left  Node
    Value string
  }

  Subscript struct {
    Left  Node
    Right Node
  }

  Slice struct {
    Start Node
    End   Node
  }

  KwArg struct {
    Key   string
    Value Node
  }

  VarArg struct {
    Arg Node
  }

  CallExpr struct {
    Left  Node
    Args  []Node
  }

  UnaryExpr struct {
    Op    token.Token
    Right Node
  }

  BinaryExpr struct {
    Op    token.Token
    Left  Node
    Right Node
  }

  TernaryExpr struct {
    Cond Node
    Then Node
    Else Node
  }

  //
  // statements
  // 

  Declaration struct {
    IsConst bool
    Left    []*Id
    Right   []Node
  }

  Assignment struct {
    Op    token.Token
    Left  []Node
    Right []Node
  }

  BranchStmt struct {
    Type token.Token // BREAK, CONTINUE or FALLTHROUGH
  }

  ReturnStmt struct {
    Values []Node
  }

  IfStmt struct {
    Init *Assignment
    Cond Node
    Body Node
    Else Node
  }

  ForIteratorStmt struct {
    Iterator   *Id
    Collection Node
    Body       Node
  }

  ForStmt struct {
    Init *Assignment
    Cond Node
    Step Node
    Body Node
  }

  Block struct {
    Nodes []Node
  }
)


func (node *Nil) Accept(v Visitor) {
  v.VisitNil(node)
}

func (node *Bool) Accept(v Visitor) {
  v.VisitBool(node)
}

func (node *Number) Accept(v Visitor) {
  v.VisitNumber(node)
}

func (node *Id) Accept(v Visitor) {
  v.VisitId(node)
}

func (node *String) Accept(v Visitor) {
  v.VisitString(node)
}

func (node *Array) Accept(v Visitor) {
  v.VisitArray(node)
}

func (node *ObjectField) Accept(v Visitor) {
  v.VisitObjectField(node)
}

func (node *Object) Accept(v Visitor) {
  v.VisitObject(node)
}

func (node *Function) Accept(v Visitor) {
  v.VisitFunction(node)
}

func (node *Selector) Accept(v Visitor) {
  v.VisitSelector(node)
}

func (node *Subscript) Accept(v Visitor) {
  v.VisitSubscript(node)
}

func (node *Slice) Accept(v Visitor) {
  v.VisitSlice(node)
}

func (node *KwArg) Accept(v Visitor) {
  v.VisitKwArg(node)
}

func (node *VarArg) Accept(v Visitor) {
  v.VisitVarArg(node)
}

func (node *CallExpr) Accept(v Visitor) {
  v.VisitCallExpr(node)
}

func (node *UnaryExpr) Accept(v Visitor) {
  v.VisitUnaryExpr(node)
}

func (node *TernaryExpr) Accept(v Visitor) {
  v.VisitTernaryExpr(node)
}

func (node *BinaryExpr) Accept(v Visitor) {
  v.VisitBinaryExpr(node)
}

func (node *Declaration) Accept(v Visitor) {
  v.VisitDeclaration(node)
}

func (node *Assignment) Accept(v Visitor) {
  v.VisitAssignment(node)
}

func (node *BranchStmt) Accept(v Visitor) {
  v.VisitBranchStmt(node)
}

func (node *ReturnStmt) Accept(v Visitor) {
  v.VisitReturnStmt(node)
}

func (node *IfStmt) Accept(v Visitor) {
  v.VisitIfStmt(node)
}

func (node *ForIteratorStmt) Accept(v Visitor) {
  v.VisitForIteratorStmt(node)
}

func (node *ForStmt) Accept(v Visitor) {
  v.VisitForStmt(node)
}

func (node *Block) Accept(v Visitor) {
  v.VisitBlock(node)
}