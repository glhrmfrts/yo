// Visitor interface

package ast

import (
)

type Visitor interface {
  VisitNil(node *Nil)
  VisitBool(node *Bool)
  VisitNumber(node *Number)
  VisitId(node *Id)
  VisitString(node *String)
  VisitArray(node *Array)
  VisitObjectField(node *ObjectField)
  VisitObject(node *Object)
  VisitFunction(node *Function)
  VisitSelector(node *Selector)
  VisitSubscript(node *Subscript)
  VisitSlice(node *Slice)
  VisitKwArg(node *KwArg)
  VisitVarArg(node *VarArg)
  VisitCallExpr(node *CallExpr)
  VisitUnaryExpr(node *UnaryExpr)
  VisitBinaryExpr(node *BinaryExpr)
  VisitDeclaration(node *Declaration)
  VisitAssignment(node *Assignment)
  VisitBranchStmt(node *BranchStmt)
  VisitReturnStmt(node *ReturnStmt)
  VisitIfStmt(node *IfStmt)
  VisitForIteratorStmt(node *ForIteratorStmt)
  VisitForStmt(node *ForStmt)
  VisitBlock(node *Block)
}