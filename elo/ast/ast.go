// Abstract Syntax Tree

package ast

import (
)

type Node interface {
  Accept(v Visitor)
}

type Number struct {
  Value string
}

type Id struct {
  Value string
}

func (node *Number) Accept(v Visitor) {
  v.VisitNumber(node)
}

func (node *Id) Accept(v Visitor) {
  v.VisitId(node)
}