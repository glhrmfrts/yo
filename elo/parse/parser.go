package parse

import (
  //"fmt"
  "github.com/glhrmfrts/elo-lang/elo/ast"
)

type parser struct {
  tok token
  literal string

  tokenizer *tokenizer
}

func (p *parser) next() {
  p.tok, p.literal = p.tokenizer.nextToken()
}

func (p *parser) is(toktype token) bool {
  return p.tok == toktype
}

func (p *parser) accept(toktype token) bool {
  defer p.next()
  if p.is(toktype) {
    return true
  }
  return false
}

func (p *parser) primaryExpr() ast.Node {
  defer p.next()

  switch p.tok {
  case TOKEN_NUMBER:
    return &ast.Number{Value: p.literal}
  case TOKEN_ID:
    return &ast.Id{Value: p.literal}
  }

  return nil
}

func (p *parser) program() ast.Node {
  p.next()
  return p.primaryExpr()
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