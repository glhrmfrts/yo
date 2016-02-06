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
  tokenizer       tokenizer
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

func (p *parser) errorExpected(expected string) error {
  return p.error(fmt.Sprintf("unexpected %s, expected %s", p.tok, expected))
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

func (p *parser) makeId() *ast.Id {
  return &ast.Id{Value: p.literal}
}

func (p *parser) makeSelector(left ast.Node) *ast.Selector {
  return &ast.Selector{Left: left, Value: p.literal}
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

// check if an expression list contains only assignable values
func (p *parser) checkLhsList(list []ast.Node) bool {
  for _, node := range list {
    switch node.(type) {
    case *ast.Id, *ast.Selector, *ast.Subscript:
      continue
    default:
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
    return nil, p.errorExpected("closing ']'")
  }

  return &ast.Array{Values: list}, nil
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
    return nil, p.errorExpected("closing '}'")
  }

  return &ast.Object{Fields: fields}, nil
}

func (p *parser) functionArgs() ([]ast.Node, error) {
  if !p.accept(token.LPAREN) {
    return nil, p.errorExpected("'('")
  }
  
  var list []ast.Node
  if p.accept(token.RPAREN) {
    // no arguments
    return list, nil
  }

  var vararg, kwarg bool
  for p.tok == token.ID {
    if vararg {
      return nil, p.error("argument after variadic argument")
    }

    var arg ast.Node
    id := p.makeId()
    p.next()

    // '='
    if p.accept(token.EQ) {
      value, err := p.expr()
      if err != nil {
        return nil, err
      }
      arg = &ast.KwArg{Key: id.Value, Value: value}
      kwarg = true
    } else if p.accept(token.DOTDOTDOT) {
      arg = &ast.VarArg{Arg: id}
      vararg = true
    } else {
      if vararg {
        return nil, p.error("positional argument after variadic argument")
      }
      if kwarg {
        return nil, p.error("positional argument after keyword argument")
      }
      arg = id
    }

    list = append(list, arg)
    if !p.accept(token.COMMA) {
      break
    }
  }

  if !p.accept(token.RPAREN) {
    return nil, p.errorExpected("closing ')'")
  }
  return list, nil
}

func (p *parser) functionBody() (ast.Node, error) {
  if p.accept(token.TILDE) {
    // '^' curried function
    args, err := p.functionArgs()
    if err != nil {
      return nil, err
    }

    body, err := p.functionBody()
    if err != nil {
      return nil, err
    }

    fn := &ast.Function{Args: args, Body: body}
    return &ast.Block{
      Nodes: []ast.Node{ &ast.ReturnStmt{Values: []ast.Node{fn}} },
    }, nil
  } else if p.accept(token.EQGT) {
    // '=>' short function
    list, err := p.exprList(false)
    if err != nil {
      return nil, err
    }
    return &ast.Block{
      Nodes: []ast.Node{ &ast.ReturnStmt{Values: list} },
    }, nil
  } else if p.tok == token.LBRACE {
    // '{' regular function body
    return p.block()
  }

  return nil, p.errorExpected("'^', '=>' or '{'")
}

func (p *parser) function() (ast.Node, error) {
  p.next() // 'func'

  var name ast.Node
  if p.tok == token.ID {
    name = ast.Node(p.makeId())
    p.next()

    if p.accept(token.DOT) && p.tok == token.ID {
      name = ast.Node(p.makeSelector(name))
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
  return &ast.Function{Name: name, Args: args, Body: body}, nil
}

func (p *parser) primaryExpr() (ast.Node, error) {
  // these first productions before the second 'switch'
  // handle the ending token themselves, so 'defer p.next()'
  // needs to be after them
  switch p.tok {
  case token.FUNC:
    return p.function()
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
      return nil, p.errorExpected("closing ')'")
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
    return nil, p.errorExpected("identifier")
  }

  defer p.next()
  return p.makeSelector(left), nil
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
    return nil, p.errorExpected("closing ']'")
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
      return nil, p.errorExpected("closing ')'")
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

func (p *parser) ternaryExpr(left ast.Node) (ast.Node, error) {
  p.next() // '?'

  whenTrue, err := p.expr()
  if err != nil {
    return nil, err
  }
  if !p.accept(token.COLON) {
    return nil, p.errorExpected("':'")
  }
  whenFalse, err := p.expr()
  if err != nil {
    return nil, err
  }
  return &ast.TernaryExpr{Cond: left, Then: whenTrue, Else: whenFalse}, nil
}

func (p *parser) expr() (ast.Node, error) {
  left, err := p.unaryExpr()
  if err != nil {
    return nil, err
  }
  left, err = p.binaryExpr(left, 0)
  if err != nil {
    return nil,err
  }

  // avoid unecessary calls to ternaryExpr
  if p.tok == token.QUESTION {
    return p.ternaryExpr(left)
  }

  return left, nil
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
    if isIdList := p.checkIdList(left); !isIdList {
      return nil, p.error("non-identifier at left side of ':='")
    }
  } else {
    // validate left side of assignment
    if isLhsList := p.checkLhsList(left); !isLhsList {
      return nil, p.error("non-assignable at left side of '='")
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
  switch tok := p.tok; tok {
  case token.CONST, token.VAR:
    return p.declaration()
  case token.BREAK, token.CONTINUE, token.FALLTHROUGH:
    p.next()
    return &ast.BranchStmt{Type: tok}, nil
  case token.RETURN:
    p.next()
    values, err := p.exprList(false)
    if err != nil {
      return nil, err
    }
    return &ast.ReturnStmt{Values: values}, nil
  case token.IF:
    return p.ifStmt()
  case token.FOR:
    return p.forStmt()
  default:
    return p.assignment()
  }
}

func (p *parser) ifStmt() (ast.Node, error) {
  p.next() // 'if'

  var init *ast.Assignment
  var else_ ast.Node
  cond, err := p.assignment()
  if err != nil {
    return nil, err
  }

  init, ok := cond.(*ast.Assignment)
  if ok {
    if !p.accept(token.SEMICOLON) {
      return nil, p.errorExpected("';'")
    }
    cond, err = p.expr()
    if err != nil {
      return nil, err
    }
  }

  body, err := p.block()
  if err != nil {
    return nil, err
  }

  if p.accept(token.ELSE) {
    if p.tok == token.LBRACE {
      else_, err = p.block()
    } else if p.tok == token.IF {
      else_, err = p.ifStmt()
    } else {
      err = p.errorExpected("if or '{'")
    }

    if err != nil {
      return nil, err
    }
  }

  return &ast.IfStmt{Init: init, Cond: cond, Body: body, Else: else_}, nil
}

func (p *parser) forIteratorStmt(id *ast.Id) (ast.Node, error) {
  p.next() // 'in'

  coll, err := p.expr()
  if err != nil {
    return nil, err
  }

  body, err := p.block()
  return &ast.ForIteratorStmt{Iterator: id, Collection: coll, Body: body}, nil
}

func (p *parser) forStmt() (ast.Node, error) {
  p.next() // 'for'

  var init *ast.Assignment
  var cond ast.Node
  var step ast.Node
  var err error
  var ok bool
  if p.tok == token.LBRACE {
    goto parseBody
  }

  cond, err = p.assignment()
  if err != nil {
    return nil, err
  }

  init, ok = cond.(*ast.Assignment)
  if ok {
    if !p.accept(token.SEMICOLON) {
      return nil, p.errorExpected("';'")
    }
    if p.tok == token.LBRACE {
      goto parseBody
    }
    cond, err = p.expr()
    if err != nil {
      return nil, err
    }
  } else if id, ok := cond.(*ast.Id); ok && p.tok == token.IN {
    return p.forIteratorStmt(id)
  }

  if p.accept(token.SEMICOLON) && p.tok != token.LBRACE {
    step, err = p.assignment()
    if err != nil {
      return nil, err
    }
  }

parseBody:
  body, err := p.block()
  return &ast.ForStmt{Init: init, Cond: cond, Step: step, Body: body}, nil
}

func (p *parser) block() (ast.Node, error) {
  if !p.accept(token.LBRACE) {
    return nil, p.errorExpected("'{'")
  }

  var nodes []ast.Node
  for !(p.tok == token.RBRACE || p.tok == token.EOS) {
    stmt, err := p.stmt()
    if err != nil {
      return nil, err
    }

    nodes = append(nodes, stmt)
  }

  if !p.accept(token.RBRACE) {
    return nil, p.errorExpected("closing '}'")
  }
  return &ast.Block{Nodes: nodes}, nil
}

func (p *parser) program() (ast.Node, error) {
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

func (p *parser) init(source []byte, filename string) {
  p.ignoreNewlines = true
  p.tokenizer.init(source, filename)

  // fetch the first token
  p.next()
}

func ParseExpr(source []byte) (ast.Node, error) {
  var p parser
  p.init(source, "")
  return p.expr()
}

func ParseFile(source []byte, filename string) (ast.Node, error) {
  var p parser
  p.init(source, filename)
  return p.program()
}