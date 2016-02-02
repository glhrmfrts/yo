// token types and representations

package token


type Token int

const (
  EOS = iota
  NIL
  TRUE
  FALSE
  IF
  ELSE
  FOR
  FUNC
  CONST
  VAR
  NOT
  ID
  STRING
  NUMBER

  // assignment operators
  assignOpBegin
  EQ
  COLONEQ
  PLUSEQ
  MINUSEQ
  TIMESEQ
  DIVEQ
  PIPEEQ
  AMPEQ
  TILDEEQ
  assignOpEnd

  // operators in order of precedence from lower-to-higher
  binaryOpBegin
  PIPEPIPE

  AMPAMP

  EQEQ
  BANGEQ

  LT
  LTEQ
  GT
  GTEQ

  PLUS
  MINUS
  PIPE
  TILDE

  TIMES
  TIMESTIMES
  DIV
  LTLT
  GTGT
  AMP
  MOD
  binaryOpEnd

  COLON
  COMMA
  DOT
  BANG
  LPAREN
  RPAREN
  LBRACK
  RBRACK
  LBRACE
  RBRACE
  ILLEGAL
)

// map keywords to their respective types
var keywords = map[string]Token{
  "nil": NIL,
  "true": TRUE,
  "false": FALSE,
  "if": IF,
  "else": ELSE,
  "for": FOR,
  "func": FUNC,
  "const": CONST,
  "var": VAR,
  "not": NOT,
}

// descriptive representation of tokens
var strings = map[Token]string{
	EOS: "end of source",
  NIL: "nil",
  TRUE: "true",
  FALSE: "false",
  IF: "if",
  ELSE: "else",
  FOR: "for",
  FUNC: "func",
  NOT: "not",
  ID: "identifier",
  STRING: "string",
  NUMBER: "number",
  PLUS: "+",
  MINUS: "-",
  TIMES: "*",
  DIV: "/",
  AMPAMP: "&&",
  PIPEPIPE: "||",
  AMP: "&",
  PIPE: "|",
  TILDE: "^",
  MOD: "%",
  LT: "<",
  LTEQ: "<=",
  LTLT: "<<",
  GT: ">",
  GTEQ: ">=",
  GTGT: ">>",
  EQ: "=",
  BANGEQ: "!=",
  COLONEQ: ":=",
  EQEQ: "==",
  COLON: ":",
  COMMA: ",",
  DOT: ".",
  BANG: "!",
  LPAREN: "(",
  RPAREN: ")",
  LBRACK: "[",
  RBRACK: "]",
  LBRACE: "{",
  RBRACE: "}",
  ILLEGAL: "illegal",
}

// operators precedence
var precedences = []int{
	10,
	20,
	30, 30,
	40, 40, 40, 40,
	50, 50, 50, 50,
	60, 60, 60, 60, 60, 60, 60,
}

func Keyword(lit string) (Token, bool) {
	t, ok := keywords[lit]
	return t, ok
}

func IsUnaryOp(tok Token) bool {
	return (tok == NOT || tok == BANG || tok == MINUS || tok == PLUS)
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
	return tok == EQ
}

// method for Stringer interface
func (tok Token) String() string {
	return strings[tok]
}