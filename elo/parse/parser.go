package parse


import (
  "fmt"
  "github.com/glhrmfrts/elo-lang/elo/ast"
  "github.com/glhrmfrts/elo-lang/elo/token"
)

type parser struct {
  tok             token.Token
  literal         string
  ignoreNewlines  bool
  tokenizer       *tokenizer
}

type ParseError struct {
  guilty    token.Token
  line      int
  file      string
  message   string
}

func (err *ParseError) Error() string {
  return fmt.Sprintf("%s:%d: syntax error: %s", err.file, err.line, err.message)
}

//
// common productions
//

func (p *parser) error(msg string) error {
  t := p.tokenizer
  return &ParseError{guilty: p.tok, line: t.lineno, file: t.filename, message: msg}
}

func (p *parser) next() {
  p.tok, p.literal = p.tokenizer.nextToken()

  for p.ignoreNewlines && p.tok == token.NEWLINE {
    p.tok, p.literal = p.tokenizer.nextToken()
  }
}

func (p *parser) accept(toktype token.Token) bool {
  if p.tok == toktype {
    p.next()
    return true
  }
  return false
}

func (p *parser) idList() []*ast.Id {
  var list []*ast.Id

  for p.tok == token.ID {
    list = append(list, &ast.Id{Value: p.literal})

    p.next()
    if !p.accept(token.COMMA) {
      break
    }
  }

  return list
}

// check if an expression list contains only identifiers
func (p *parser) checkIdList(list []ast.Node) bool {
  for _, node := range list {
    if _, isId := node.(*ast.Id); !isId {
      return false
    }
  }

  return true
}

func (p *parser) exprList(inArray bool) ([]ast.Node, error) {
  var list []ast.Node
  for {
    // trailing comma check
    if inArray && p.tok == token.RBRACK {
      break
    }

    expr, err := p.expr()
    if err != nil {
      return nil, err
    }
    list = append(list, expr)
    if !p.accept(token.COMMA) {
      break
    }
  }

  return list, nil
}

func (p *parser) objectFieldList() ([]*ast.ObjectField, error) {
  var list []*ast.ObjectField
  for {
    // trailing comma check
    if p.tok == token.RBRACE {
      break
    }

    key, err := p.expr()
    if err != nil {
      return nil, err
    }
    if !p.accept(token.COLON) {
      list = append(list, &ast.ObjectField{Key: key})
    } else {
      value, err := p.expr()
      if err != nil {
        return nil, err
      }
      list = append(list, &ast.ObjectField{Key: key, Value: value})
    }

    if !p.accept(token.COMMA) {
      break
    }
  }

  return list, nil
}

//
// grammar rules
//

func (p *parser) array() (ast.Node, error) {
  p.next() // '['

  if p.accept(token.RBRACK) {
    // no elements
    return &ast.Array{}, nil
  }

  list, err := p.exprList(true)
  if err != nil {
    return nil, err
  }
  if !p.accept(token.RBRACK) {
    return nil, p.error(fmt.Sprintf("unexpected %s, expecting closing ']'", p.tok))
  }

  return &ast.Array{Values: list}, nil
}

func (p *parser) object() (ast.Node, error) {
  p.next() // '{'

  if p.accept(token.RBRACE) {
    // no elements
    return &ast.Object{}, nil
  }

  fields, err := p.objectFieldList()
  if err != nil {
    return nil, err
  }
  if !p.accept(token.RBRACE) {
    return nil, p.error(fmt.Sprintf("unexpected %s, expecting closing '}'", p.tok))
  }

  return &ast.Object{Fields: fields}, nil
}

func (p *parser) functionArgs() (ast.Node, error) {
  if !p.accept(token.LPAREN) {
    return nil, p.error(fmt.Sprintf("unexpected %s, expecting '('", p.tok))
  }
  
  var list []ast.Node
  if p.tok == token.RPAREN {
    // no arguments
    return list, nil
  }
  for p.tok == token.ID {
    arg := p.makeId()
    p.next()

    // '='
    if p.accept(token.EQ) {
      value, err := p.expr()
      if err != nil {
        return nil, err
      }
      arg = &ast.KwArg{Key: id.Value, Value: value}
    } else if p.accept(token.DOTDOTDOT) {
      arg = &ast.VarArg{Arg: arg}
    }

    list = append(list, arg)
    if !p.accept(token.COMMA) {
      break
    }
  }

  if !p.accept(token.RPAREN) {
    return nil, p.error(fmt.Sprintf("unexpected %s, expecting closing ')'", p.tok))
  }
  return list, nil
}

func (p *parser) functionBody() (ast.Node, error) {
  if p.accept(token.TILDE) {
    // curried function
    if !p.accept(token.LPAREN) {
      return nil, p.error(fmt.Sprintf("unexpected %s, expecting '('", p.tok))
    }
    args, err := p.functionArgs()
    if err != nil {
      return nil, err
    }

    body, err := p.functionBody()
    if err != nil {
      return nil, err
    }

    fn := &ast.Function{Args: args, Body: body}
    return &ast.ReturnStmt{Values: []ast.Node{fn}}
  } else if p.accept(token.EQGT) {
    // short function
    list, err := p.exprList()
    if err != nil {
      return nil, err
    }
    return &ast.ReturnStmt{Values: list}
  } else if p.tok == token.LBRACE {
    // regular function body
    return p.block()
  }
}

func (p *parser) function() (ast.Node, error) {
  p.next() // 'func'

  var name ast.Node
  if p.accept(token.ID) {
    name := p.makeId()
    if p.accept(token.DOT) && p.tok == token.ID {
      name = p.makeSelector(name)
      p.next()
    }
  }
  args, err := p.functionArgs()
  if err != nil {
    return nil, err
  }

  body, err := p.functionBody()
  if err != nil {
    return nil, err
  }
  return &ast.Function{Name: name, Args: args, Body: body}
}

func (p *parser) primaryExpr() (ast.Node, error) {
  // these first productions before the second 'switch'
  // handle the ending token themselves, so 'defer p.next()'
  // needs to be after them
  switch p.tok {
  case token.LBRACK:
    return p.array()
  case token.LBRACE:
    return p.object()
  case token.LPAREN:
    p.next()
    expr, err := p.expr()

    if err != nil {
      return nil, err
    }

    if !p.accept(token.RPAREN) {
      return nil, p.error(fmt.Sprintf("unexpected %s, expecting closing ')'", p.tok))
    }

    return expr, nil
  default:
    defer p.next()
    switch p.tok {
    case token.INT, token.FLOAT:
      return &ast.Number{Type: p.tok, Value: p.literal}, nil
    case token.ID:
      return &ast.Id{Value: p.literal}, nil
    case token.STRING:
      return &ast.String{Value: p.literal}, nil
    case token.TRUE, token.FALSE:
      return &ast.Bool{Value: p.tok == token.TRUE}, nil
    case token.NIL:
      return &ast.Nil{}, nil
    }
  }

  return nil, p.error(fmt.Sprintf("unexpected %s", p.tok))
}

func (p *parser) selectorExpr(left ast.Node) (ast.Node, error) {
  if !(p.tok == token.ID) {
    return nil, p.error(fmt.Sprintf("unexpected %s, expecting identifier", p.tok))
  }

  defer p.next()
  return &ast.Selector{Left: left, Key: p.literal}, nil
}

func (p *parser) subscriptExpr(left ast.Node) (ast.Node, error) {
  expr, err := p.expr()
  if err != nil {
    return nil, err
  }

  sub := &ast.Subscript{Left: left, Right: expr}
  if p.accept(token.COLON) {
    expr2, err := p.expr()
    if err != nil {
      return nil, err
    }

    sub.Right = &ast.Slice{Start: expr, End: expr2}
  }

  if !p.accept(token.RBRACK) {
    return nil, p.error(fmt.Sprintf("unexpected %s, expecting closing ']'", p.tok))
  }

  return sub, nil
}

func (p *parser) selectorOrSubscriptExpr(left ast.Node) (ast.Node, error) {
  var err error

  if left == nil {
    left, err = p.primaryExpr()
    if err != nil {
      return nil, err
    }
  }

  for {
    if dot, lBrack := p.tok == token.DOT, p.tok == token.LBRACK; dot || lBrack {
      old := p.ignoreNewlines
      p.ignoreNewlines = false
      p.next()
      if p.tok == token.NEWLINE || p.tok == token.EOS {
        return nil, p.error("expression not terminated")
      }
      p.ignoreNewlines = old

      if dot {
        left, err = p.selectorExpr(left)
      } else {
        left, err = p.subscriptExpr(left)
      }

      if err != nil {
        return nil, err
      }
    } else {
      break
    }
  }

  return left, nil
}

func (p *parser) callArgs() ([]ast.Node, error) {
  var list []ast.Node
  if p.tok == token.RPAREN {
    // no arguments
    return list, nil
  }

  for {
    arg, err := p.expr()
    if err != nil {
      return nil, err
    }

    // '='
    if p.accept(token.EQ) {
      value, err := p.expr()
      if err != nil {
        return nil, err
      }

      if id, isId := arg.(*ast.Id); isId {
        arg = &ast.KwArg{Key: id.Value, Value: value}
      } else {
        return nil, p.error("non-identifier in left side of keyword argument")
      }
    } else if p.accept(token.DOTDOTDOT) {
      arg = &ast.VarArg{Arg: arg}
    }

    list = append(list, arg)
    if !p.accept(token.COMMA) {
      break
    }
  }

  return list, nil
}

func (p *parser) callExpr() (ast.Node, error) {
  left, err := p.selectorOrSubscriptExpr(nil)
  if err != nil {
    return nil, err
  }

  var args []ast.Node
  for p.accept(token.LPAREN) {
    args, err = p.callArgs()
    if err != nil {
      return nil, err
    }

    if !p.accept(token.RPAREN) {
      return nil, p.error(fmt.Sprintf("unexpected %s, expected closing ')'", p.tok))
    }
    left = &ast.CallExpr{Left: left, Args: args}
  }

  return p.selectorOrSubscriptExpr(left)
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
      right, err = p.callExpr()
    }

    if err != nil {
      return nil, err
    }

    return &ast.UnaryExpr{Op: op, Right: right}, nil
  }

  return p.callExpr()
}

// parse a binary expression using the legendary wikipedia's algorithm :)
func (p *parser) binaryExpr(left ast.Node, minPrecedence int) (ast.Node, error) {
  for token.IsBinaryOp(p.tok) && token.Precedence(p.tok) >= minPrecedence {
    op := p.tok
    opPrecedence := token.Precedence(op)

    // consume operator
    old := p.ignoreNewlines
    p.ignoreNewlines = false
    p.next()
    if p.tok == token.NEWLINE || p.tok == token.EOS {
      return nil, p.error("expression not terminated")
    }
    p.ignoreNewlines = old

    right, err := p.unaryExpr()
    if err != nil {
      return nil, err
    }

    for (token.IsBinaryOp(p.tok) && token.Precedence(p.tok) > opPrecedence) || 
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

func (p *parser) declaration() (ast.Node, error) {
  isConst := p.tok == token.CONST
  p.next()

  left := p.idList()

  // '='
  if (!p.accept(token.EQ)) {
    // a declaration without any values
    return &ast.Declaration{IsConst: isConst, Left: left}, nil
  }

  right, err := p.exprList(false)
  if err != nil {
    return nil, err
  }

  return &ast.Declaration{IsConst: isConst, Left: left, Right: right}, nil
}

func (p *parser) assignment() (ast.Node, error) {
  left, err := p.exprList(false)
  if err != nil {
    return nil, err
  }

  if !token.IsAssignOp(p.tok) {
    if len(left) > 1 {
      return nil, p.error("illegal expression")
    }

    return left[0], nil
  }

  // ':='
  if p.tok == token.COLONEQ {
    // a short variable declaration
    isIdList := p.checkIdList(left)

    if !isIdList {
      return nil, p.error("non-identifier at left side of ':='")
    }
  }

  op := p.tok
  p.next()

  right, err := p.exprList(false)
  if err != nil {
    return nil, err
  }

  return &ast.Assignment{Op: op, Left: left, Right: right}, nil
}

func (p *parser) stmt() (ast.Node, error) {
  defer p.accept(token.SEMICOLON)
  switch p.tok {
  case token.CONST, token.VAR:
    return p.declaration()
  default:
    return p.assignment()
  }
}

func (p *parser) program() (ast.Node, error) {
  p.next()

  var nodes []ast.Node
  for !(p.tok == token.EOS) {
    stmt, err := p.stmt()
    if err != nil {
      return nil, err
    }

    nodes = append(nodes, stmt)
  }

  return &ast.Block{Nodes: nodes}, nil
}

//
// initialization of parser
//

func makeParser(source []byte, filename string) *parser {
  p := &parser{
    ignoreNewlines: true,
    tokenizer: makeTokenizer(source, filename),
  }
  return p
}

func Parse(source []byte, filename string) (ast.Node, error) {
  p := makeParser(source, filename)
  return p.program()
}