// Copyright 2016 Guilherme Nemeth <guilherme.nemeth@gmail.com>

// token types and representations

package ast

type Token int

const (
	TokenEos Token = iota
	TokenNewline
	TokenNil
	TokenTrue
	TokenFalse
	TokenIf
	TokenElse
	TokenFor
	TokenWhen
	TokenFunc
	TokenConst
	TokenVar
	TokenBreak
	TokenContinue
	TokenFallthrough
	TokenTry
	TokenRecover
	TokenFinally
	TokenPanic
	TokenReturn
	TokenNot
	TokenIn
	TokenId
	TokenString
	TokenInt
	TokenFloat

	// assignment operators
	assignOpBegin
	TokenEq
	TokenColoneq
	TokenPluseq
	TokenMinuseq
	TokenTimeseq
	TokenDiveq
	TokenPipeeq
	TokenAmpeq
	TokenTildeeq
	assignOpEnd

	// operators in order of precedence from lower-to-higher
	binaryOpBegin
	TokenPipepipe

	TokenAmpamp

	TokenEqeq
	TokenBangeq

	TokenLt
	TokenLteq
	TokenGt
	TokenGteq

	TokenPlus
	TokenMinus
	TokenPipe
	TokenTilde

	TokenTimes
	TokenTimestimes
	TokenDiv
	TokenLtlt
	TokenGtgt
	TokenAmp
	TokenMod
	binaryOpEnd

	TokenPlusplus
	TokenMinusminus
	TokenMinusgt
	TokenColon
	TokenSemicolon
	TokenComma
	TokenDot
	TokenDotdotdot
	TokenBang
	TokenQuestion
	TokenLparen
	TokenRparen
	TokenLbrack
	TokenRbrack
	TokenLbrace
	TokenRbrace
	TokenIllegal
)

var (
	// map keywords to their respective types
	keywords = map[string]Token{
		"nil":         TokenNil,
		"true":        TokenTrue,
		"false":       TokenFalse,
		"if":          TokenIf,
		"else":        TokenElse,
		"for":         TokenFor,
		"when":        TokenWhen,
		"func":        TokenFunc,
		"const":       TokenConst,
		"var":         TokenVar,
		"break":       TokenBreak,
		"continue":    TokenContinue,
		"fallthrough": TokenFallthrough,
		"try":         TokenTry,
		"recover":     TokenRecover,
		"finally":     TokenFinally,
		"panic":       TokenPanic,
		"return":      TokenReturn,
		"not":         TokenNot,
		"in":          TokenIn,
	}

	// descriptive representation of tokens
	strings = map[Token]string{
		TokenEos:         "end of source",
		TokenNil:         "nil",
		TokenTrue:        "true",
		TokenFalse:       "false",
		TokenIf:          "if",
		TokenElse:        "else",
		TokenFor:         "for",
		TokenFunc:        "func",
		TokenWhen:        "when",
		TokenConst:       "const",
		TokenVar:         "var",
		TokenBreak:       "break",
		TokenContinue:    "continue",
		TokenFallthrough: "fallthrough",
		TokenTry:         "try",
		TokenRecover:     "recover",
		TokenFinally:     "finally",
		TokenPanic:       "panic",
		TokenReturn:      "return",
		TokenNot:         "not",
		TokenIn:          "in",
		TokenId:          "identifier",
		TokenString:      "string",
		TokenInt:         "int",
		TokenFloat:       "float",
		TokenPlus:        "+",
		TokenMinus:       "-",
		TokenTimes:       "*",
		TokenDiv:         "/",
		TokenAmpamp:      "&&",
		TokenPipepipe:    "||",
		TokenAmp:         "&",
		TokenPipe:        "|",
		TokenTilde:       "^",
		TokenMod:         "%",
		TokenLt:          "<",
		TokenLteq:        "<=",
		TokenLtlt:        "<<",
		TokenGt:          ">",
		TokenGteq:        ">=",
		TokenGtgt:        ">>",
		TokenEq:          "=",
		TokenBangeq:      "!=",
		TokenColoneq:     ":=",
		TokenPluseq:      "+=",
		TokenMinuseq:     "-=",
		TokenTimeseq:     "*=",
		TokenDiveq:       "/=",
		TokenAmpeq:       "&=",
		TokenPipeeq:      "|=",
		TokenTildeeq:     "^=",
		TokenEqeq:        "==",
		TokenPlusplus:    "++",
		TokenMinusminus:  "--",
		TokenMinusgt:     "->",
		TokenColon:       ":",
		TokenSemicolon:   ";",
		TokenComma:       ",",
		TokenDot:         ".",
		TokenDotdotdot:   "...",
		TokenBang:        "!",
		TokenLparen:      "(",
		TokenRparen:      ")",
		TokenLbrack:      "[",
		TokenRbrack:      "]",
		TokenLbrace:      "{",
		TokenRbrace:      "}",
		TokenIllegal:     "illegal",
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
	return (tok == TokenPlusplus || tok == TokenMinusminus)
}

func IsUnaryOp(tok Token) bool {
	return IsPostfixOp(tok) ||
		(tok == TokenNot || tok == TokenBang || tok == TokenMinus || tok == TokenPlus || tok == TokenTilde)
}

func IsBinaryOp(tok Token) bool {
	return tok >= binaryOpBegin && tok <= binaryOpEnd
}

func IsAssignOp(tok Token) bool {
	return tok >= assignOpBegin && tok <= assignOpEnd
}

func CompoundOp(tok Token) Token {
	switch tok {
	case TokenPluseq:
		return TokenPlus
	case TokenMinuseq:
		return TokenMinus
	case TokenTimeseq:
		return TokenTimes
	case TokenDiveq:
		return TokenDiv
	case TokenAmpeq:
		return TokenAmp
	case TokenPipeeq:
		return TokenPipe
	case TokenTildeeq:
		return TokenTilde
	}
	return Token(-1)
}

func Precedence(tok Token) int {
	return precedences[int(tok-binaryOpBegin-1)]
}

func RightAssociative(tok Token) bool {
	return tok == TokenEq
}

// method for Stringer interface
func (tok Token) String() string {
	return strings[tok]
}
