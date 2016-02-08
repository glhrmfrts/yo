// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>
// Copyright 2009 The Go Authors. All rights reserved.
//
// Some parts of this file were taken from Go's source code.
// For more details please take a look at the LICENSE file

package parse

import (
  "unicode"
  "unicode/utf8"
  "os"
  "fmt"
  "github.com/glhrmfrts/elo-lang/elo/ast"
)

type tokenizer struct {
  offset int
  readOffset int
  r rune
  src []byte
  filename string
  lineno int
}

const bom = 0xFEFF
const eof = -1

func isLetter(ch rune) bool {
  return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
  return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}

func (t *tokenizer) error(msg string) {
  fmt.Printf("%s:%d -> syntax error: %s\n", t.filename, t.lineno, msg)
  os.Exit(1)
}

func (t *tokenizer) nextChar() bool {
  if t.readOffset < len(t.src) {
    t.offset = t.readOffset
    ch := t.src[t.readOffset]
    
    r, w := rune(ch), 1
    switch {
    case r == 0:
      t.error("illegal character NUL")
    case r >= 0x80:
      // not ASCII
      r, w = utf8.DecodeRune(t.src[t.offset:])
      if r == utf8.RuneError && w == 1 {
        t.error("illegal UTF-8 encoding")
      } else if r == bom && t.offset > 0 {
        t.error("illegal byte order mark")
      }
    }

    if ch == '\n' {
      t.lineno++
    }

    t.r = r
    t.readOffset += w
    return true
  }

  t.r = eof
  t.offset = len(t.src)
  return false
}

func (t *tokenizer) scanComment() bool {
  // initial '/' already consumed
  if t.r == '/' {
    for t.r != eof && t.r != '\n' {
      t.nextChar()
    }

    return true
  }

  return false
}

func (t *tokenizer) scanIdentifier() string {
  offs := t.offset
  for isLetter(t.r) || isDigit(t.r) {
    t.nextChar()
  }
  return string(t.src[offs:t.offset])
}

func digitVal(ch rune) int {
  switch {
  case '0' <= ch && ch <= '9':
    return int(ch - '0')
  case 'a' <= ch && ch <= 'f':
    return int(ch - 'a' + 10)
  case 'A' <= ch && ch <= 'F':
    return int(ch - 'A' + 10)
  }
  return 16
}

func (t *tokenizer) scanMantissa(base int) {
  for digitVal(t.r) < base {
    t.nextChar()
  }
}

func (t *tokenizer) scanNumber(seenDecimalPoint bool) (ast.Token, string) {
  // digitVal(t.r) < 10
  offs := t.offset
  typ := ast.T_INT

  if seenDecimalPoint {
    typ = ast.T_FLOAT
    offs--
    t.scanMantissa(10)
    goto exponent
  }

  if t.r == '0' {
    // int or float
    offs := t.offset
    t.nextChar()
    if t.r == 'x' || t.r == 'X' {
      // hexadecimal int
      t.nextChar()
      t.scanMantissa(16)
      if t.offset-offs <= 2 {
        // only scanned "0x" or "0X"
        t.error("illegal hexadecimal number")
      }
    } else {
      // octal int or float
      seenDecimalDigit := false
      t.scanMantissa(8)
      if t.r == '8' || t.r == '9' {
        // illegal octal int or float
        seenDecimalDigit = true
        t.scanMantissa(10)
      }
      if t.r == '.' || t.r == 'e' || t.r == 'E' {
        goto fraction
      }
      // octal int
      if seenDecimalDigit {
        t.error("illegal octal number")
      }
    }
    goto exit
  }

  // decimal int or float
  t.scanMantissa(10)

fraction:
  if t.r == '.' {
    typ = ast.T_FLOAT
    t.nextChar()
    t.scanMantissa(10)
  }

exponent:
  if t.r == 'e' || t.r == 'E' {
    typ = ast.T_FLOAT
    t.nextChar()
    if t.r == '-' || t.r == '+' {
      t.nextChar()
    }
    t.scanMantissa(10)
  }

exit:
  return typ, string(t.src[offs:t.offset])
}

// scans a valid escape sequence and returns the evaluated value
func (t *tokenizer) scanEscape(quote rune) rune {

  var n int
  var base, max uint32
  var r rune

  switch t.r {
  case 'a': r = '\a'
  case 'b': r = '\b'
  case 'f': r = '\f'
  case 'n': r = '\n'
  case 'r': r = '\r'
  case 't': r = '\t'
  case 'v': r = '\v'
  case '\\': r = '\\'
  case quote: r = quote
  case '0', '1', '2', '3', '4', '5', '6', '7':
    n, base, max = 3, 8, 255
  case 'x':
    t.nextChar()
    n, base, max = 2, 16, 255
  case 'u':
    t.nextChar()
    n, base, max = 4, 16, unicode.MaxRune
  case 'U':
    t.nextChar()
    n, base, max = 8, 16, unicode.MaxRune
  default:
    msg := "unknown escape sequence"
    if t.r < 0 {
      msg = "escape sequence not terminated"
    }
    t.error(msg)
  }

  if r > 0 {
    return r
  }

  var x uint32
  for n > 0 {
    d := uint32(digitVal(t.r))
    if d >= base {
      msg := fmt.Sprintf("illegal character %#U in escape sequence", t.r)
      if t.r < 0 {
        msg = "escape sequence not terminated"
      }
      t.error(msg)
    }
    x = x*base + d
    t.nextChar()
    n--
    if n == 0 && base == 16 && max == 255 && t.r == '\\' {
      rd := t.readOffset
      t.nextChar()
      if t.r == 'x' {
        n = 2
        max = unicode.MaxRune
        t.nextChar()
      } else {
        t.readOffset = rd
      }
    }
  }

  if x > max || 0xD800 <= x && x < 0xE000 {
    t.error("escape sequence is invalid Unicode code point")
  }

  return rune(x)
}

func (t *tokenizer) scanString(quote rune) string {
  var result string
  for {
    ch := t.r
    if ch < 0 {
      t.error("string literal not terminated")
    }
    t.nextChar()
    if ch == quote {
      break
    }
    if ch == '\\' {
      ch = t.scanEscape(quote)
    }
    result += string(ch)
  }
  return result
}

func (t *tokenizer) skipWhitespace() {
  for t.r == ' ' || t.r == '\t' || t.r == '\r' {
    t.nextChar()
  }
}

// functions that look 1 or 2 characters ahead,
// and return the given token types based on that

func (t *tokenizer) maybe1(a ast.Token, c1 rune, t1 ast.Token) ast.Token {
  offset := t.readOffset

  t.nextChar()
  if t.r == c1 {
    return t1
  }

  t.readOffset = offset
  return a
}

func (t *tokenizer) maybe2(a ast.Token, c1 rune, t1 ast.Token, c2 rune, t2 ast.Token) ast.Token {
  offset := t.readOffset

  t.nextChar()
  if t.r == c1 {
    return t1
  }
  if t.r == c2 {
    return t2
  }

  t.readOffset = offset
  return a
}

// does the actual scanning and return the type of the token
// and a literal string representing it
func (t *tokenizer) nextToken() (ast.Token, string) {
  t.skipWhitespace()

  switch ch := t.r; {
  case isLetter(t.r):
    lit := t.scanIdentifier()
    kwtype, ok := ast.Keyword(lit)
    if ok {
      return kwtype, lit
    }
    return ast.T_ID, lit
  case isDigit(t.r):
    return t.scanNumber(false)
  case t.r == '\'' || t.r == '"':
    t.nextChar()
    return ast.T_STRING, t.scanString(ch)
  default:
    if t.r == '/' {
      t.nextChar()
      if t.scanComment() {
        return t.nextToken()
      }

      if t.r == '=' {
        t.nextChar()
        return ast.T_DIVEQ, "/="
      }
      return ast.T_DIV, "/"
    }

    tok := ast.Token(-1)
    offs := t.offset

    switch t.r {
    case '\n': tok = ast.T_NEWLINE
    case '+': tok = t.maybe1(ast.T_PLUS, '=', ast.T_PLUSEQ)
    case '-': tok = t.maybe1(ast.T_MINUS, '=', ast.T_MINUSEQ)
    case '*': tok = t.maybe2(ast.T_TIMES, '=', ast.T_TIMESEQ, '*', ast.T_TIMESTIMES)
    case '&': tok = t.maybe2(ast.T_AMP, '=', ast.T_AMPEQ, '&', ast.T_AMPAMP)
    case '|': tok = t.maybe2(ast.T_PIPE, '=', ast.T_PIPEEQ, '|', ast.T_PIPEPIPE)
    case '^': tok = t.maybe1(ast.T_TILDE, '=', ast.T_TILDEEQ)
    case '<': tok = t.maybe2(ast.T_LT, '=', ast.T_LTEQ, '<', ast.T_LTLT)
    case '>': tok = t.maybe2(ast.T_GT, '=', ast.T_GTEQ, '>', ast.T_GTGT)
    case '=': tok = t.maybe2(ast.T_EQ, '=', ast.T_EQEQ, '>', ast.T_EQGT)
    case ':': tok = t.maybe1(ast.T_COLON, '=', ast.T_COLONEQ)
    case ';': tok = ast.T_SEMICOLON
    case ',': tok = ast.T_COMMA
    case '!': tok = t.maybe1(ast.T_BANG, '=', ast.T_BANGEQ)
    case '?': tok = ast.T_QUESTION
    case '(': tok = ast.T_LPAREN
    case ')': tok = ast.T_RPAREN
    case '[': tok = ast.T_LBRACK
    case ']': tok = ast.T_RBRACK
    case '{': tok = ast.T_LBRACE
    case '}': tok = ast.T_RBRACE
    case '.':
      t.nextChar()
      if isDigit(t.r) {
        return t.scanNumber(true)
      } else if t.r == '.' {
        t.nextChar()
        if t.r == '.' {
          t.nextChar()
          return ast.T_DOTDOTDOT, "..."
        }
      } else {
        return ast.T_DOT, "."
      }
    }

    if tok != eof {
      t.nextChar()
      return tok, string(t.src[offs:t.offset])
    }
  }

  if t.offset >= len(t.src) {
    return ast.T_EOS, "end"
  }

  fmt.Print(string(t.r))
  return ast.T_ILLEGAL, ""
}

func (t *tokenizer) init(source []byte, filename string) {
  t.src = source
  t.filename = filename
  t.lineno = 1

  // fetch the first char
  t.nextChar()
}