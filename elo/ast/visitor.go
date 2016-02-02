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
  VisitUnaryExpr(node *UnaryExpr)
  VisitBinaryExpr(node *BinaryExpr)
  VisitDeclaration(node *Declaration)
}