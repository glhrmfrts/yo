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
  ID
  STRING
  NUMBER

  // operators in order of precedence from lower-to-higher
  binaryOpBegin
  EQ
  COLONEQ

  LT
  LTEQ
  LTLT
  GT
  GTEQ
  GTGT

  EQEQ

  PLUS
  MINUS

  MULT
  DIV
  binaryOpEnd

  COLON
  COMMA
  DOT
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
  ID: "identifier",
  STRING: "string",
  NUMBER: "number",
  PLUS: "+",
  MINUS: "-",
  MULT: "*",
  DIV: "/",
  LT: "<",
  LTEQ: "<=",
  LTLT: "<<",
  GT: ">",
  GTEQ: ">=",
  GTGT: ">>",
  EQ: "=",
  COLONEQ: ":=",
  EQEQ: "==",
  COLON: ":",
  COMMA: ",",
  DOT: ".",
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
	0, 0,
	10, 10, 10, 10, 10, 10,
	20,
	30, 30, 
	40, 40,
}

func Keyword(lit string) (Token, bool) {
	t, ok := keywords[lit]
	return t, ok
}

func IsBinaryOp(tok Token) bool {
	return tok >= binaryOpBegin && tok <= binaryOpEnd
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