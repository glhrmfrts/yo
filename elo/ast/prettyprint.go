// Pretty-print the AST to a buffer

package ast

import (
  "fmt"
  "bytes"
  //"github.com/glhrmfrts/elo-lang/elo/token"
)

type prettyprinter struct {
  indent int
  indentSize int
  buf bytes.Buffer
}

func (p *prettyprinter) doIndent() {
  for i := 0; i < p.indent * p.indentSize; i++ {
    p.buf.WriteString(" ")
  }
}

func (p *prettyprinter) VisitNil(node *Nil, data interface{}) {
  p.buf.WriteString("(nil)")
}

func (p *prettyprinter) VisitBool(node *Bool, data interface{}) {
  var val string
  if node.Value {
    val = "true"
  } else {
    val = "false"
  }
  p.buf.WriteString("(" + val + ")")
}

func (p *prettyprinter) VisitNumber(node *Number, data interface{}) {
  p.buf.WriteString(fmt.Sprintf("(number %f)", node.Value))
}

func (p *prettyprinter) VisitId(node *Id, data interface{}) {
  p.buf.WriteString("(id " + node.Value + ")")
}

func (p *prettyprinter) VisitString(node *String, data interface{}) {
  p.buf.WriteString("(string \""+ node.Value + "\")")
}

func (p *prettyprinter) VisitArray(node *Array, data interface{}) {
  p.buf.WriteString("(array")
  p.indent++

  for _, n := range node.Values {
    p.buf.WriteString("\n")
    p.doIndent()
    n.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitObjectField(node *ObjectField, data interface{}) {
  p.buf.WriteString("(field\n")
  p.indent++
  p.doIndent()

  if node.Key != nil {
    node.Key.Accept(p, nil)
  }

  p.buf.WriteString("\n")
  p.doIndent()

  if node.Value != nil {
    node.Value.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitObject(node *Object, data interface{}) {
  p.buf.WriteString("(object")
  p.indent++

  for _, f := range node.Fields {
    p.buf.WriteString("\n")
    p.doIndent()
    f.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitFunction(node *Function, data interface{}) {
  p.buf.WriteString("(func ")
  if node.Name != nil {
    node.Name.Accept(p, nil)
  }
  p.buf.WriteString("\n")
  p.indent++

  for _, a := range node.Args {
    p.doIndent()
    a.Accept(p, nil)
    p.buf.WriteString("\n")
  }

  p.doIndent()
  p.buf.WriteString("->\n")

  p.doIndent()
  node.Body.Accept(p, nil)

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitSelector(node *Selector, data interface{}) {
  p.buf.WriteString("(selector\n")

  p.indent++
  p.doIndent()

  node.Left.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()

  p.indent--
  p.buf.WriteString("'" + node.Value + "')")
}

func (p *prettyprinter) VisitSubscript(node *Subscript, data interface{}) {
  p.buf.WriteString("(subscript\n")

  p.indent++
  p.doIndent()

  node.Left.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()

  node.Right.Accept(p, nil)

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitSlice(node *Slice, data interface{}) {
  p.buf.WriteString("(slice\n")

  p.indent++
  p.doIndent()

  node.Start.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()

  node.End.Accept(p, nil)

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitKwArg(node *KwArg, data interface{}) {
  p.buf.WriteString("(kwarg\n")

  p.indent++
  p.doIndent()

  p.buf.WriteString("'" + node.Key + "'\n")

  p.doIndent()
  node.Value.Accept(p, nil)

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitVarArg(node *VarArg, data interface{}) {
  p.buf.WriteString("(vararg ")
  node.Arg.Accept(p, nil)
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitCallExpr(node *CallExpr, data interface{}) {
  p.buf.WriteString("(call\n")
  p.indent++
  p.doIndent()

  node.Left.Accept(p, nil)

  for _, arg := range node.Args {
    p.buf.WriteString("\n")
    p.doIndent()
    arg.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitUnaryExpr(node *UnaryExpr, data interface{}) {
  p.buf.WriteString(fmt.Sprintf("(unary %s\n", node.Op))
  
  p.indent++
  p.doIndent()

  node.Right.Accept(p, nil)

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitBinaryExpr(node *BinaryExpr, data interface{}) {
  p.buf.WriteString(fmt.Sprintf("(binary %s\n", node.Op))

  p.indent++
  p.doIndent()

  node.Left.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()

  node.Right.Accept(p, nil)

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitTernaryExpr(node *TernaryExpr, data interface{}) {
  p.buf.WriteString("(ternary\n")
  p.indent++
  p.doIndent()

  node.Cond.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()
  p.buf.WriteString("?\n")
  p.doIndent()

  node.Then.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()
  p.buf.WriteString(":\n")
  p.doIndent()

  node.Else.Accept(p, nil)

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitDeclaration(node *Declaration, data interface{}) {
  keyword := "var"
  if node.IsConst {
    keyword = "const"
  }

  p.buf.WriteString(fmt.Sprintf("(%s", keyword))
  p.indent++

  for _, id := range node.Left {
    p.buf.WriteString("\n")
    p.doIndent()
    id.Accept(p, nil)
  }

  p.buf.WriteString("\n")
  p.doIndent()
  p.buf.WriteString("=")

  for _, node := range node.Right {
    p.buf.WriteString("\n")
    p.doIndent()
    node.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")") 
}

func (p *prettyprinter) VisitAssignment(node *Assignment, data interface{}) {
  p.buf.WriteString("(assignment")
  p.indent++

  for _, node := range node.Left {
    p.buf.WriteString("\n")
    p.doIndent()
    node.Accept(p, nil)
  }

  p.buf.WriteString("\n")
  p.doIndent()
  p.buf.WriteString(node.Op.String())

  for _, node := range node.Right {
    p.buf.WriteString("\n")
    p.doIndent()
    node.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitBranchStmt(node *BranchStmt, data interface{}) {
  p.buf.WriteString(fmt.Sprintf("(%s)", node.Type))
}

func (p *prettyprinter) VisitReturnStmt(node *ReturnStmt, data interface{}) {
  p.buf.WriteString("(return")
  p.indent++

  for _, v := range node.Values {
    p.buf.WriteString("\n")
    p.doIndent()
    v.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitIfStmt(node *IfStmt, data interface{}) {
  p.buf.WriteString("(if\n")
  p.indent++
  
  if node.Init != nil {
    p.doIndent()
    node.Init.Accept(p, nil)
    p.buf.WriteString("\n")
  }

  p.doIndent()
  node.Cond.Accept(p, nil)
  p.buf.WriteString("\n")

  p.doIndent()
  node.Body.Accept(p, nil)
  p.buf.WriteString("\n")
  p.doIndent()

  if node.Else != nil {
    node.Else.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitForIteratorStmt(node *ForIteratorStmt, data interface{}) {
  p.buf.WriteString("(for iterator\n")
  p.indent++
  p.doIndent()
  node.Iterator.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()
  node.Collection.Accept(p, nil)

  p.buf.WriteString("\n")
  p.doIndent()
  node.Body.Accept(p, nil)

  p.indent--
  p.buf.WriteString("\n")
}

func (p *prettyprinter) VisitForStmt(node *ForStmt, data interface{}) {
  p.buf.WriteString("(for\n")
  p.indent++

  if node.Init != nil {
    p.doIndent()
    node.Init.Accept(p, nil)
    p.buf.WriteString("\n")
  }

  if node.Cond != nil {
    p.doIndent()
    node.Cond.Accept(p, nil)
    p.buf.WriteString("\n")
  }

  p.doIndent()
  if node.Step != nil {
    node.Step.Accept(p, nil)
    p.buf.WriteString("\n")
    p.doIndent()
  }

  node.Body.Accept(p, nil)
  p.indent--
  p.buf.WriteString(")")
}

func (p *prettyprinter) VisitBlock(node *Block, data interface{}) {
  p.buf.WriteString("(block")
  p.indent++

  for _, n := range node.Nodes {
    p.buf.WriteString("\n")
    p.doIndent()
    n.Accept(p, nil)
  }

  p.indent--
  p.buf.WriteString(")")
}

func Prettyprint(root Node, indentSize int) string {
  v := prettyprinter{indentSize: indentSize}
  root.Accept(&v, nil)
  return v.buf.String()
}