// Pretty-print the AST to a buffer

package ast

import (
  "bytes"
)

type Prettyprinter struct {
  buf bytes.Buffer
}

func (p *Prettyprinter) VisitNumber(node *Number) {
  p.buf.WriteString("[number " + node.Value + "]\n")
}

func (p *Prettyprinter) VisitId(node *Id) {
  p.buf.WriteString("[id " + node.Value + "]\n")
}

func Prettyprint(root Node) string {
  v := Prettyprinter{}
  root.Accept(&v)
  return v.buf.String()
}