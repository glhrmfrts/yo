// token types and representations

package token


type Token int

const (
  EOS Token = iota
  NEWLINE
  NIL
  TRUE
  FALSE
  IF
  ELSE
  FOR
  FUNC
  CONST
  VAR
  BREAK
  CONTINUE
  FALLTHROUGH
  RETURN
  NOT
  IN
  ID
  STRING
  INT
  FLOAT

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

  EQGT
  COLON
  SEMICOLON
  COMMA
  DOT
  DOTDOTDOT
  BANG
  QUESTION
  LPAREN
  RPAREN
  LBRACK
  RBRACK
  LBRACE
  RBRACE
  ILLEGAL
)

var (
  // map keywords to their respective types
  keywords = map[string]Token{
    "nil": NIL,
    "true": TRUE,
    "false": FALSE,
    "if": IF,
    "else": ELSE,
    "for": FOR,
    "func": FUNC,
    "const": CONST,
    "var": VAR,
    "break": BREAK,
    "continue": CONTINUE,
    "fallthrough": FALLTHROUGH,
    "return": RETURN,
    "not": NOT,
    "in": IN,
  }

  // descriptive representation of tokens
  strings = map[Token]string{
  	EOS: "end of source",
    NIL: "nil",
    TRUE: "true",
    FALSE: "false",
    IF: "if",
    ELSE: "else",
    FOR: "for",
    FUNC: "func",
    CONST: "const",
    VAR: "var",
    BREAK: "break",
    CONTINUE: "continue",
    FALLTHROUGH: "fallthrough",
    RETURN: "return",
    NOT: "not",
    IN: "in",
    ID: "identifier",
    STRING: "string",
    INT: "int",
    FLOAT: "float",
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
    PLUSEQ: "+=",
    MINUSEQ: "-=",
    TIMESEQ: "*=",
    DIVEQ: "/=",
    AMPEQ: "&=",
    PIPEEQ: "|=",
    TILDEEQ: "^=",
    EQEQ: "==",
    EQGT: "=>",
    COLON: ":",
    SEMICOLON: ";",
    COMMA: ",",
    DOT: ".",
    DOTDOTDOT: "...",
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