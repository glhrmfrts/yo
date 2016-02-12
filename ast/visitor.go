// Visitor interface

package ast

import (
)

type Visitor interface {
  VisitNil(node *Nil, data interface{})
  VisitBool(node *Bool, data interface{})
  VisitNumber(node *Number, data interface{})
  VisitId(node *Id, data interface{})
  VisitString(node *String, data interface{})
  VisitArray(node *Array, data interface{})
  VisitObjectField(node *ObjectField, data interface{})
  VisitObject(node *Object, data interface{})
  VisitFunction(node *Function, data interface{})
  VisitSelector(node *Selector, data interface{})
  VisitSubscript(node *Subscript, data interface{})
  VisitSlice(node *Slice, data interface{})
  VisitKwArg(node *KwArg, data interface{})
  VisitVarArg(node *VarArg, data interface{})
  VisitCallExpr(node *CallExpr, data interface{})
  VisitPostfixExpr(node *PostfixExpr, data interface{})
  VisitUnaryExpr(node *UnaryExpr, data interface{})
  VisitBinaryExpr(node *BinaryExpr, data interface{})
  VisitTernaryExpr(node *TernaryExpr, data interface{})
  VisitDeclaration(node *Declaration, data interface{})
  VisitAssignment(node *Assignment, data interface{})
  VisitBranchStmt(node *BranchStmt, data interface{})
  VisitReturnStmt(node *ReturnStmt, data interface{})
  VisitIfStmt(node *IfStmt, data interface{})
  VisitForIteratorStmt(node *ForIteratorStmt, data interface{})
  VisitForStmt(node *ForStmt, data interface{})
  VisitBlock(node *Block, data interface{})
}