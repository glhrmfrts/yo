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
  VisitKeyword(node *Keyword)
  VisitAtomKeyword(node *AtomKeyword)
  VisitCallArgs(node *CallArgs)
  VisitCall(node *Call)
}