// token types and representations

package ast


type Token int

const (
  T_EOS Token = iota
  T_NEWLINE
  T_NIL
  T_TRUE
  T_FALSE
  T_IF
  T_ELSE
  T_FOR
  T_WHEN
  T_FUNC
  T_CONST
  T_VAR
  T_BREAK
  T_CONTINUE
  T_FALLTHROUGH
  T_RETURN
  T_NOT
  T_IN
  T_ID
  T_STRING
  T_INT
  T_FLOAT

  // assignment operators
  assignOpBegin
  T_EQ
  T_COLONEQ
  T_PLUSEQ
  T_MINUSEQ
  T_TIMESEQ
  T_DIVEQ
  T_PIPEEQ
  T_AMPEQ
  T_TILDEEQ
  assignOpEnd

  // operators in order of precedence from lower-to-higher
  binaryOpBegin
  T_PIPEPIPE

  T_AMPAMP

  T_EQEQ
  T_BANGEQ

  T_LT
  T_LTEQ
  T_GT
  T_GTEQ

  T_PLUS
  T_MINUS
  T_PIPE
  T_TILDE

  T_TIMES
  T_TIMESTIMES
  T_DIV
  T_LTLT
  T_GTGT
  T_AMP
  T_MOD
  binaryOpEnd

  T_PLUSPLUS
  T_MINUSMINUS
  T_EQGT
  T_COLON
  T_SEMICOLON
  T_COMMA
  T_DOT
  T_DOTDOTDOT
  T_BANG
  T_QUESTION
  T_LPAREN
  T_RPAREN
  T_LBRACK
  T_RBRACK
  T_LBRACE
  T_RBRACE
  T_ILLEGAL
)

var (
  // map keywords to their respective types
  keywords = map[string]Token{
    "nil": T_NIL,
    "true": T_TRUE,
    "false": T_FALSE,
    "if": T_IF,
    "else": T_ELSE,
    "for": T_FOR,
    "when": T_WHEN,
    "func": T_FUNC,
    "const": T_CONST,
    "var": T_VAR,
    "break": T_BREAK,
    "continue": T_CONTINUE,
    "fallthrough": T_FALLTHROUGH,
    "return": T_RETURN,
    "not": T_NOT,
    "in": T_IN,
  }

  // descriptive representation of tokens
  strings = map[Token]string{
    T_EOS: "end of source",
    T_NIL: "nil",
    T_TRUE: "true",
    T_FALSE: "false",
    T_IF: "if",
    T_ELSE: "else",
    T_FOR: "for",
    T_FUNC: "func",
    T_WHEN: "when",
    T_CONST: "const",
    T_VAR: "var",
    T_BREAK: "break",
    T_CONTINUE: "continue",
    T_FALLTHROUGH: "fallthrough",
    T_RETURN: "return",
    T_NOT: "not",
    T_IN: "in",
    T_ID: "identifier",
    T_STRING: "string",
    T_INT: "int",
    T_FLOAT: "float",
    T_PLUS: "+",
    T_MINUS: "-",
    T_TIMES: "*",
    T_DIV: "/",
    T_AMPAMP: "&&",
    T_PIPEPIPE: "||",
    T_AMP: "&",
    T_PIPE: "|",
    T_TILDE: "^",
    T_MOD: "%",
    T_LT: "<",
    T_LTEQ: "<=",
    T_LTLT: "<<",
    T_GT: ">",
    T_GTEQ: ">=",
    T_GTGT: ">>",
    T_EQ: "=",
    T_BANGEQ: "!=",
    T_COLONEQ: ":=",
    T_PLUSEQ: "+=",
    T_MINUSEQ: "-=",
    T_TIMESEQ: "*=",
    T_DIVEQ: "/=",
    T_AMPEQ: "&=",
    T_PIPEEQ: "|=",
    T_TILDEEQ: "^=",
    T_EQEQ: "==",
    T_PLUSPLUS: "++",
    T_MINUSMINUS: "--",
    T_EQGT: "=>",
    T_COLON: ":",
    T_SEMICOLON: ";",
    T_COMMA: ",",
    T_DOT: ".",
    T_DOTDOTDOT: "...",
    T_BANG: "!",
    T_LPAREN: "(",
    T_RPAREN: ")",
    T_LBRACK: "[",
    T_RBRACK: "]",
    T_LBRACE: "{",
    T_RBRACE: "}",
    T_ILLEGAL: "illegal",
  }

  // operators precedence
  precedences = []int{
    10,
    20,
    30, 30,
    40, 40, 40, 40,
    50, 50, 50, 50,
    60, 60, 60, 60, 60, 60, 60,
  }
)

func Keyword(lit string) (Token, bool) {
  t, ok := keywords[lit]
  return t, ok
}

func IsPostfixOp(tok Token) bool {
  return (tok == T_PLUSPLUS || tok == T_MINUSMINUS)
}

func IsUnaryOp(tok Token) bool {
  return IsPostfixOp(tok) || 
    (tok == T_NOT || tok == T_BANG || tok == T_MINUS || tok == T_PLUS || tok == T_TILDE)
}

func IsBinaryOp(tok Token) bool {
  return tok >= binaryOpBegin && tok <= binaryOpEnd
}

func IsAssignOp(tok Token) bool {
  return tok >= assignOpBegin && tok <= assignOpEnd
}

func Precedence(tok Token) int {
  return precedences[int(tok - binaryOpBegin - 1)]
}

func RightAssociative(tok Token) bool {
  return tok == T_EQ
}

// method for Stringer interface
func (tok Token) String() string {
  return strings[tok]
}