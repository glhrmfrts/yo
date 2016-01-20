package parse

import (
  "fmt"
  "github.com/glhrmfrts/elo-lang/elo/ast"
)

type parser struct {
  tok token
  literal string

  tokenizer *tokenizer
}

func (p *parser) error(msg string) {
  p.tokenizer.error(msg)
}

func (p *parser) next() {
  p.tok, p.literal = p.tokenizer.nextToken()
}

func (p *parser) is(toktype token) bool {
  return p.tok == toktype
}

func (p *parser) accept(toktype token) bool {
  if p.is(toktype) {
    p.next()
    return true
  }
  return false
}

func (p *parser) maybeFuncDef(arg ast.Node) ast.Node {
  if p.is(TOKEN_LBRACE) {
    //return p.funcDef(arg)
  }
  return arg
}

func (p *parser) primaryExpr() ast.Node {
  switch p.tok {
  case TOKEN_LPAREN:
    p.next()
    expr := p.expr()
    if p.accept(TOKEN_COMMA) {
      //return p.funcDef(expr)
    } else if p.accept(TOKEN_RPAREN) {
      return p.maybeFuncDef(expr)
    }
    p.error(fmt.Sprintf("unexpected %s", p.literal))
  default:
    defer p.next()
    switch p.tok {
    case TOKEN_NUMBER:
      return &ast.Number{Value: p.literal}
    case TOKEN_ID:
      return &ast.Id{Value: p.literal}
    }
  }

  return nil
}

func (p *parser) expr() ast.Node {
  return p.primaryExpr()
} 

func (p *parser) program() ast.Node {
  p.next()
  return p.expr()
}

func makeParser(source, filename string) *parser {
  p := &parser{
    tokenizer: makeTokenizer(source, filename),
  }
  return p
}

func Parse(source, filename string) ast.Node {
  p := makeParser(source, filename)
  return p.program()
}