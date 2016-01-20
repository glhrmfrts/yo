// heavily inspired by Go's tokenizer :)

package parse

import (
  "unicode"
  "unicode/utf8"
  "os"
  "fmt"
)

type tokenizer struct {
  offset int
  readOffset int
  r rune
  src []byte
  filename string
  lineno int
}

const bom = 0xFEFF // byte-order mark, only allowed as the first character

func isLetter(ch rune) bool {
  return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-' || ch >= 0x80 && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
  return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}

func (t *tokenizer) error(msg string) {
  fmt.Printf("%s:%i: syntax error: %s\n", t.filename, t.lineno, msg)
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

  t.r = -1
  t.offset = len(t.src)
  return false
}

func (t *tokenizer) scanComment() bool {
  // initial '/' already consumed
  if t.r == '/' {
    for t.r != '\n' {
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
  return 16 // larger than any legal digit val
}

func (t *tokenizer) scanMantissa(base int) {
  for digitVal(t.r) < base {
    t.nextChar()
  }
}

func (t *tokenizer) scanNumber(seenDecimalPoint bool) string {
  // digitVal(t.r) < 10
  offs := t.offset

  if seenDecimalPoint {
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
    t.nextChar()
    t.scanMantissa(10)
  }

exponent:
  if t.r == 'e' || t.r == 'E' {
    t.nextChar()
    if t.r == '-' || t.r == '+' {
      t.nextChar()
    }
    t.scanMantissa(10)
  }

exit:
  return string(t.src[offs:t.offset])
}

func (t *tokenizer) skipWhitespace() {
  for t.r == ' ' || t.r == '\t' || t.r == '\n' {
    t.nextChar()
  }
}

func (t *tokenizer) maybe1(a token, c1 rune, t1 token) token {
  offset := t.readOffset
  defer func(t *tokenizer, offset int) {
    t.readOffset = offset
  }(t, offset)

  t.nextChar()
  if t.r == c1 {
    return t1
  }
  return a
}

func (t *tokenizer) maybe2(a token, c1 rune, t1 token, c2 rune, t2 token) token {
  offset := t.offset
  defer func(t *tokenizer, offset int) {
    t.offset = offset
  }(t, offset)

  t.nextChar()
  if t.r == c1 {
    return t1
  }
  if t.r == c2 {
    return t2
  }
  return a
}

func (t *tokenizer) nextToken() (token, string) {
  t.skipWhitespace()

  switch {
  case t.r != '-' && isLetter(t.r): // '-' is a letter but cannot start an identifier
    lit := t.scanIdentifier()
    return TOKEN_ID, lit
  case isDigit(t.r):
    lit := t.scanNumber(false)
    return TOKEN_NUMBER, lit
  default:
    // Always advance
    defer t.nextChar()

    if t.r == '/' {
      if t.scanComment() {
        return t.nextToken()
      }
      return TOKEN_DIV, string(t.r)
    }
    switch t.r {
    case '+': return TOKEN_PLUS, string(t.r)
    case '-': return TOKEN_MINUS, string(t.r)
    case '*': return TOKEN_MULT, string(t.r)
    case '<': return t.maybe2(TOKEN_LT, '=', TOKEN_LTEQ, '<', TOKEN_LTLT), string(t.r)
    case '>': return t.maybe2(TOKEN_GT, '=', TOKEN_GTEQ, '>', TOKEN_GTGT), string(t.r)
    case '=': return t.maybe1(TOKEN_EQ, '=', TOKEN_EQEQ), string(t.r)
    case ':': return TOKEN_COLON, string(t.r)
    case '[': return TOKEN_LBRACK, string(t.r)
    case ']': return TOKEN_RBRACK, string(t.r)
    case '{': return TOKEN_LBRACE, string(t.r)
    case '}': return TOKEN_RBRACE, string(t.r)
    }
  }

  if t.offset >= len(t.src) {
    return TOKEN_EOS, ""
  }

  return TOKEN_ILLEGAL, ""
}

func makeTokenizer(source, filename string) *tokenizer {
  tok := &tokenizer{
    src: []byte(source),
    filename: filename,
  }
  tok.nextChar()
  return tok
}