// Pretty-print the AST to a buffer

package ast

import (
  "fmt"
  "bytes"
  //"github.com/glhrmfrts/elo-lang/elo/token"
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

func (p *Prettyprinter) VisitUnaryExpr(node *UnaryExpr) {
  p.buf.WriteString(fmt.Sprintf("(unary %s ", node.Op))
  node.Right.Accept(p)
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitBinaryExpr(node *BinaryExpr) {
  p.buf.WriteString(fmt.Sprintf("(binary %s ", node.Op))
  node.Left.Accept(p)
  p.buf.WriteString(" ")
  node.Right.Accept(p)
  p.buf.WriteString(")") 
}

func (p *Prettyprinter) VisitDeclaration(node *Declaration) {
  keyword := "var"
  if node.IsConst {
    keyword = "const"
  }

  p.buf.WriteString(fmt.Sprintf("(decl %s\n", keyword))
  p.indent++

  for i, id := range node.Left {
    p.doIndent()
    p.buf.WriteString("(" + id.Value)

    if i < len(node.Right) {
      p.buf.WriteString(" = ")
      node.Right[i].Accept(p)
    }

    p.buf.WriteString(")\n")
  }

  p.indent--
  p.doIndent()
  p.buf.WriteString(")\n") 
}

func Prettyprint(root Node) string {
  v := Prettyprinter{}
  root.Accept(&v)
  return v.buf.String()
}