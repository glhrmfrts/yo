// Visitor interface

package ast

import (
)

type Visitor interface {
  VisitNumber(node *Number)
  VisitId(node *Id)
  VisitAtom(node *Atom)
  VisitKeyword(node *Keyword)
  VisitAtomKeyword(node *AtomKeyword)
  VisitCallArgs(node *CallArgs)
  VisitCall(node *Call)
}