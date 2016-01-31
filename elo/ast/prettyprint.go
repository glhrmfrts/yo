// Pretty-print the AST to a buffer

package ast

import (
  "bytes"
)

type Prettyprinter struct {
  indent int
  buf bytes.Buffer
}

func (p *Prettyprinter) doIndent() {
  for i := 0; i < p.indent; i++ {
    p.buf.WriteString(" ")
  }
}

func (p *Prettyprinter) VisitNil(node *Nil) {
  p.buf.WriteString("(nil)")
}

func (p *Prettyprinter) VisitBool(node *Bool) {
  var val string
  if node.Value {
    val = "true"
  } else {
    val = "false"
  }
  p.buf.WriteString("(" + val + ")")
}

func (p *Prettyprinter) VisitNumber(node *Number) {
  p.buf.WriteString("(number " + node.Value + ")")
}

func (p *Prettyprinter) VisitId(node *Id) {
  p.buf.WriteString("(id " + node.Value + ")")
}

func (p *Prettyprinter) VisitString(node *String) {
  p.buf.WriteString("(string \""+ node.Value + "\")")
}

func (p *Prettyprinter) VisitKeyword(node *Keyword) {
  p.buf.WriteString("(kw ")
  node.Left.Accept(p)
  p.buf.WriteString(" = ")
  node.Right.Accept(p)
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitAtomKeyword(node *AtomKeyword) {
  p.buf.WriteString("(atom-kw ")
  node.Left.Accept(p)
  p.buf.WriteString(" = ")
  node.Right.Accept(p)
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitCallArgs(node *CallArgs) {
  p.buf.WriteString("(args\n")
  p.indent++
  for _, v := range node.Pos {
    p.doIndent()
    v.Accept(p)
    p.buf.WriteString("\n")
  }
  for _, v := range node.Keywords {
    p.doIndent()
    v.Accept(p)
    p.buf.WriteString("\n")
  }
  for _, v := range node.AtomKeywords {
    p.doIndent()
    v.Accept(p)
    p.buf.WriteString("\n")
  }
  p.indent--
  p.doIndent()
  p.buf.WriteString(")\n")
}

func (p *Prettyprinter) VisitCall(node *Call) {
  p.buf.WriteString("(call\n")
  p.indent++
  p.doIndent()
  node.Left.Accept(p)
  p.buf.WriteString("\n")
  p.doIndent()
  node.Args.Accept(p)
  p.indent--
  p.doIndent()
  p.buf.WriteString(")\n")
}

func Prettyprint(root Node) string {
  v := Prettyprinter{}
  root.Accept(&v)
  return v.buf.String()
}