// Abstract Syntax Tree

package ast

import (
)

type (
  Node interface {
    Accept(v Visitor, data interface{})
  }

  NodeInfo struct {
    Line int
  }

  //
  // expressions
  //

  Nil struct {
    NodeInfo 
  }

  Bool struct {
    NodeInfo
    Value bool
  }

  Number struct {
    NodeInfo
    Value float64
  }

  Id struct {
    NodeInfo
    Value string
  }

  String struct {
    NodeInfo
    Value string
  }

  Array struct {
    NodeInfo
    Elements []Node
  }

  ObjectField struct {
    NodeInfo
    Key   string
    Value Node
  }

  Object struct {
    NodeInfo
    Fields []*ObjectField
  }

  Function struct {
    NodeInfo
    Name Node
    Args []Node
    Body Node
  }

  Selector struct {
    NodeInfo
    Left  Node
    Value string
  }

  Subscript struct {
    NodeInfo
    Left  Node
    Right Node
  }

  Slice struct {
    NodeInfo
    Start Node
    End   Node
  }

  KwArg struct {
    NodeInfo
    Key   string
    Value Node
  }

  VarArg struct {
    NodeInfo
    Arg Node
  }

  CallExpr struct {
    NodeInfo
    Left  Node
    Args  []Node
  }

  PostfixExpr struct {
    NodeInfo
    Op   Token
    Left Node
  }

  UnaryExpr struct {
    NodeInfo
    Op    Token
    Right Node
  }

  BinaryExpr struct {
    NodeInfo
    Op    Token
    Left  Node
    Right Node
  }

  TernaryExpr struct {
    NodeInfo
    Cond Node
    Then Node
    Else Node
  }

  //
  // statements
  // 

  Declaration struct {
    NodeInfo
    IsConst bool
    Left    []*Id
    Right   []Node
  }

  Assignment struct {
    NodeInfo
    Op    Token
    Left  []Node
    Right []Node
  }

  BranchStmt struct {
    NodeInfo
    Type Token // BREAK, CONTINUE or FALLTHROUGH
  }

  ReturnStmt struct {
    NodeInfo
    Values []Node
  }

  IfStmt struct {
    NodeInfo
    Init *Assignment
    Cond Node
    Body Node
    Else Node
  }

  ForIteratorStmt struct {
    NodeInfo
    Iterator   *Id
    Collection Node
    When       Node
    Body       Node
  }

  ForStmt struct {
    NodeInfo
    Init *Assignment
    Cond Node
    Step Node
    Body Node
  }

  Block struct {
    NodeInfo
    Nodes []Node
  }
)


func (node *Nil) Accept(v Visitor, data interface{}) {
  v.VisitNil(node, data)
}

func (node *Bool) Accept(v Visitor, data interface{}) {
  v.VisitBool(node, data)
}

func (node *Number) Accept(v Visitor, data interface{}) {
  v.VisitNumber(node, data)
}

func (node *Id) Accept(v Visitor, data interface{}) {
  v.VisitId(node, data)
}

func (node *String) Accept(v Visitor, data interface{}) {
  v.VisitString(node, data)
}

func (node *Array) Accept(v Visitor, data interface{}) {
  v.VisitArray(node, data)
}

func (node *ObjectField) Accept(v Visitor, data interface{}) {
  v.VisitObjectField(node, data)
}

func (node *Object) Accept(v Visitor, data interface{}) {
  v.VisitObject(node, data)
}

func (node *Function) Accept(v Visitor, data interface{}) {
  v.VisitFunction(node, data)
}

func (node *Selector) Accept(v Visitor, data interface{}) {
  v.VisitSelector(node, data)
}

func (node *Subscript) Accept(v Visitor, data interface{}) {
  v.VisitSubscript(node, data)
}

func (node *Slice) Accept(v Visitor, data interface{}) {
  v.VisitSlice(node, data)
}

func (node *KwArg) Accept(v Visitor, data interface{}) {
  v.VisitKwArg(node, data)
}

func (node *VarArg) Accept(v Visitor, data interface{}) {
  v.VisitVarArg(node, data)
}

func (node *CallExpr) Accept(v Visitor, data interface{}) {
  v.VisitCallExpr(node, data)
}

func (node *PostfixExpr) Accept(v Visitor, data interface{}) {
  v.VisitPostfixExpr(node, data)
}

func (node *UnaryExpr) Accept(v Visitor, data interface{}) {
  v.VisitUnaryExpr(node, data)
}

func (node *TernaryExpr) Accept(v Visitor, data interface{}) {
  v.VisitTernaryExpr(node, data)
}

func (node *BinaryExpr) Accept(v Visitor, data interface{}) {
  v.VisitBinaryExpr(node, data)
}

func (node *Declaration) Accept(v Visitor, data interface{}) {
  v.VisitDeclaration(node, data)
}

func (node *Assignment) Accept(v Visitor, data interface{}) {
  v.VisitAssignment(node, data)
}

func (node *BranchStmt) Accept(v Visitor, data interface{}) {
  v.VisitBranchStmt(node, data)
}

func (node *ReturnStmt) Accept(v Visitor, data interface{}) {
  v.VisitReturnStmt(node, data)
}

func (node *IfStmt) Accept(v Visitor, data interface{}) {
  v.VisitIfStmt(node, data)
}

func (node *ForIteratorStmt) Accept(v Visitor, data interface{}) {
  v.VisitForIteratorStmt(node, data)
}

func (node *ForStmt) Accept(v Visitor, data interface{}) {
  v.VisitForStmt(node, data)
}

func (node *Block) Accept(v Visitor, data interface{}) {
  v.VisitBlock(node, data)
}

// return true if the given node is a statement
func IsStmt(node Node) bool {
  switch node.(type) {
  case *Assignment, *IfStmt, *ForStmt, *ForIteratorStmt,
       *BranchStmt, *ReturnStmt, *Declaration:
    return true
  default:
    return false
  }
}