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

func (p *parser) error(msg string) error {
  t := p.tokenizer
  return fmt.Errorf("%s:%d -> syntax error: %s", t.filename, t.lineno, msg)
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

func (p *parser) primaryExpr() (ast.Node, error) {
  defer p.next()
  switch p.tok {
  case TOKEN_LPAREN:
    p.next()
    expr, err := p.expr()
    if err != nil {
      return nil, err
    }
    if !p.is(TOKEN_RPAREN) {
      return nil, p.error(fmt.Sprintf("unexpected %s", p.literal))
    }
    return expr, nil
  case TOKEN_NUMBER:
    return &ast.Number{Value: p.literal}, nil
  case TOKEN_ID:
    return &ast.Id{Value: p.literal}, nil
  case TOKEN_STRING:
    return &ast.String{Value: p.literal}, nil
  case TOKEN_TRUE, TOKEN_FALSE:
    return &ast.Bool{Value: p.tok == TOKEN_TRUE}, nil
  case TOKEN_NIL:
    return &ast.Nil{}, nil
  }

  return nil, p.error(fmt.Sprintf("unexpected %s", p.literal))
}

func (p *parser) expr() (ast.Node, error) {
  return p.primaryExpr()
}

func (p *parser) program() (ast.Node, error) {
  p.next()
  return p.expr()
}

func makeParser(source []byte, filename string) *parser {
  p := &parser{
    tokenizer: makeTokenizer(source, filename),
  }
  return p
}

func Parse(source []byte, filename string) (ast.Node, error) {
  p := makeParser(source, filename)
  return p.program()
}