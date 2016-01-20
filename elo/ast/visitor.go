// Visitor interface

package ast

import (
)

type Visitor interface {
  VisitNumber(node *Number)
  VisitId(node *Id)
}