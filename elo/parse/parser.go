package parse

import (
  "fmt"
  "github.com/glhrmfrts/elo-lang/elo/ast"
  "github.com/glhrmfrts/elo-lang/elo/token"
)

type parser struct {
  tok token.Token
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

func (p *parser) is(toktype token.Token) bool {
  return p.tok == toktype
}

func (p *parser) accept(toktype token.Token) bool {
  if p.is(toktype) {
    p.next()
    return true
  }
  return false
}

func (p *parser) primaryExpr() (ast.Node, error) {
  defer p.next()
  switch p.tok {
  case token.LPAREN:
    p.next()
    expr, err := p.expr()
    if err != nil {
      return nil, err
    }
    if !p.is(token.RPAREN) {
      return nil, p.error(fmt.Sprintf("unexpected %s", p.literal))
    }
    return expr, nil
  case token.NUMBER:
    return &ast.Number{Value: p.literal}, nil
  case token.ID:
    return &ast.Id{Value: p.literal}, nil
  case token.STRING:
    return &ast.String{Value: p.literal}, nil
  case token.TRUE, token.FALSE:
    return &ast.Bool{Value: p.tok == token.TRUE}, nil
  case token.NIL:
    return &ast.Nil{}, nil
  }

  return nil, p.error(fmt.Sprintf("unexpected %s", p.literal))
}

func (p *parser) unaryExpr() (ast.Node, error) {
  if token.IsUnaryOp(p.tok) {
    op := p.tok
    p.next()

    var right ast.Node
    var err error
    if op == token.NOT {
      right, err = p.expr()
    } else {
      right, err = p.primaryExpr()
    }

    if err != nil {
      return nil, err
    }

    return &ast.UnaryExpr{Op: op, Right: right}, nil
  }

  return p.primaryExpr()
}

func (p *parser) binaryExpr(left ast.Node, minPrecedence int) (ast.Node, error) {
  for token.IsBinaryOp(p.tok) && token.Precedence(p.tok) >= minPrecedence {
    op := p.tok
    opPrecedence := token.Precedence(op)

    // consume operator
    p.next()
    if p.is(token.EOS) {
      return nil, p.error("expression not terminated")
    }

    right, err := p.unaryExpr()
    if err != nil {
      return nil, err
    }

    for token.IsBinaryOp(p.tok) && token.Precedence(p.tok) > opPrecedence || 
        (token.RightAssociative(p.tok) && token.Precedence(p.tok) >= opPrecedence) {

      right, err = p.binaryExpr(right, token.Precedence(p.tok))
      if err != nil {
        return nil, err
      }
    }

    left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
  }

  return left, nil
}

func (p *parser) expr() (ast.Node, error) {
  left, err := p.unaryExpr()
  if err != nil {
    return nil, err
  }

  return p.binaryExpr(left, 0)
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