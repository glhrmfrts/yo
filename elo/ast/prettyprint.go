// Pretty-print the AST to a buffer

package ast

import (
  "fmt"
  "bytes"
  //"github.com/glhrmfrts/elo-lang/elo/token"
)

type Prettyprinter struct {
  indent int
  indentSize int
  buf bytes.Buffer
}

func (p *Prettyprinter) doIndent() {
  for i := 0; i < p.indent * p.indentSize; i++ {
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
  p.buf.WriteString(fmt.Sprintf("(%s %s)", node.Type, node.Value))
}

func (p *Prettyprinter) VisitId(node *Id) {
  p.buf.WriteString("(id " + node.Value + ")")
}

func (p *Prettyprinter) VisitString(node *String) {
  p.buf.WriteString("(string \""+ node.Value + "\")")
}

func (p *Prettyprinter) VisitSelector(node *Selector) {
  p.buf.WriteString("(selector\n")

  p.indent++
  p.doIndent()

  node.Left.Accept(p)

  p.buf.WriteString("\n")
  p.doIndent()

  p.indent--
  p.buf.WriteString("'" + node.Key + "')")
}

func (p *Prettyprinter) VisitSubscript(node *Subscript) {
  p.buf.WriteString("(subscript\n")

  p.indent++
  p.doIndent()

  node.Left.Accept(p)

  p.buf.WriteString("\n")
  p.doIndent()

  node.Right.Accept(p)

  p.indent--
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitSlice(node *Slice) {
  p.buf.WriteString("(slice\n")

  p.indent++
  p.doIndent()

  node.Start.Accept(p)

  p.buf.WriteString("\n")
  p.doIndent()

  node.End.Accept(p)

  p.indent--
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitUnaryExpr(node *UnaryExpr) {
  p.buf.WriteString(fmt.Sprintf("(unary %s\n", node.Op))
  
  p.indent++
  p.doIndent()

  node.Right.Accept(p)

  p.indent--
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitBinaryExpr(node *BinaryExpr) {
  p.buf.WriteString(fmt.Sprintf("(binary %s\n", node.Op))

  p.indent++
  p.doIndent()

  node.Left.Accept(p)

  p.buf.WriteString("\n")
  p.doIndent()

  node.Right.Accept(p)

  p.indent--
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitDeclaration(node *Declaration) {
  keyword := "var"
  if node.IsConst {
    keyword = "const"
  }

  p.buf.WriteString(fmt.Sprintf("(%s", keyword))
  p.indent++

  for _, id := range node.Left {
    p.buf.WriteString("\n")
    p.doIndent()
    id.Accept(p)
  }

  for _, node := range node.Right {
    p.buf.WriteString("\n")
    p.doIndent()
    node.Accept(p)
  }

  p.indent--
  p.buf.WriteString(")") 
}

func (p *Prettyprinter) VisitAssignment(node *Assignment) {
  p.buf.WriteString("(assignment")
  p.indent++

  for _, node := range node.Left {
    p.buf.WriteString("\n")
    p.doIndent()
    node.Accept(p)
  }

  p.buf.WriteString("\n")
  p.doIndent()
  p.buf.WriteString(node.Op.String())

  for _, node := range node.Right {
    p.buf.WriteString("\n")
    p.doIndent()
    node.Accept(p)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *Prettyprinter) VisitBlock(node *Block) {
  for _, n := range node.Nodes {
    p.doIndent()
    n.Accept(p)
    p.buf.WriteString("\n")
  }
}

func Prettyprint(root Node, indentSize int) string {
  v := Prettyprinter{indentSize: indentSize}
  root.Accept(&v)
  return v.buf.String()
}