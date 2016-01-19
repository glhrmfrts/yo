// Copyright 2016 Guilherme Nemeth
//
// The scanning, parsing and compilation is done all in one
// pass since it's a small and simple language. That means no
// AST or any semantic analysis.
//
// (I guess?) the parser is not context-free, since the user
// can define infix and unary operators.
//

package elo

import (
	"unicode"
	"unicode/utf8"
)

type token int

type parser struct {
	tok Token
	literal string

	offset int
	src []byte
	r rune

	keywords map[string]Token
}

const (
	TOKEN_EOS = iota
	TOKEN_NIL
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_ID
	TOKEN_STRING
	TOKEN_NUMBER
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULT
	TOKEN_DIV
	TOKEN_LT
	TOKEN_LTEQ
	TOKEN_LTLT
	TOKEN_GT
	TOKEN_GTEQ
	TOKEN_GTGT
	TOKEN_EQ
	TOKEN_EQEQ
	TOKEN_COLON
	TOKEN_LBRACK
	TOKEN_RBRACK
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_ILLEGAL
)

const bom = 0xFEFF // byte-order mark, only allowed as the first character

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-' || ch >= 0x80 && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}

func (p *parser) nextChar() bool {
	if p.offset < len(p.src) {
		ch = p.src[p.offset + 1]
		
		r, w := rune(ch), 1
		switch {
		case r == 0:
			p.error("illegal character NUL")
		case r >= 0x80:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				p.error(p.offset, "illegal UTF-8 encoding")
			} else if r == bom && p.offset > 0 {
				p.error(p.offset, "illegal byte order mark")
			}
		}

		if ch == '\n' {
			p.fileInfo.line++
			p.fileInfo.col = 0
		} else {
			p.fileInfo.col++
		}

		p.r = r
		p.offset += w
		return true
	}

	return false
}

func (p *parser) scanComment() bool {
	// initial '/' already consumed
	if p.r == '/' {
		for p.r != '\n' {
			p.nextChar()
		}

		return true
	}

	return false
}

func (p *parser) scanIdentifier() string {
	offs := sp.offset
	for isLetter(p.r) || isDigit(p.r) {
		p.next()
	}
	return string(p.src[offs:p.offset])
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

func (p *parser) scanMantissa(base int) {
	for digitVal(p.r) < base {
		p.nextChar()
	}
}

func (p *parser) scanNumber() string {
	// digitVal(p.r) < 10
	offs := p.offset

	if seenDecimalPoint {
		offs--
		p.scanMantissa(10)
		goto exponent
	}

	if p.r == '0' {
		// int or float
		offs := p.offset
		p.nextChar()
		if p.r == 'x' || p.r == 'X' {
			// hexadecimal int
			p.nextChar()
			p.scanMantissa(16)
			if p.offset-offs <= 2 {
				// only scanned "0x" or "0X"
				p.error("illegal hexadecimal number")
			}
		} else {
			// octal int or float
			seenDecimalDigit := false
			p.scanMantissa(8)
			if p.r == '8' || p.r == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				p.scanMantissa(10)
			}
			if p.r == '.' || p.r == 'e' || p.r == 'E' {
				goto fraction
			}
			// octal int
			if seenDecimalDigit {
				p.error("illegal octal number")
			}
		}
		goto exit
	}

	// decimal int or float
	p.scanMantissa(10)

fraction:
	if p.r == '.' {
		p.nextChar()
		p.scanMantissa(10)
	}

exponent:
	if p.r == 'e' || p.r == 'E' {
		p.nextChar()
		if p.r == '-' || p.r == '+' {
			p.nextChar()
		}
		p.scanMantissa(10)
	}

exit:
	return string(p.src[offs:p.offset])
}

func (p *parser) skipWhitespace() {
	for p.r == ' ' || p.r == '\t' {
		p.nextChar()
	}
}

func (p *parser) nextToken() token {
	p.skipWhitespace()

	if p.offset >= len(p.src) {
		return TOKEN_EOS
	}

	// Always advance
	defer p.nextChar()

	switch {
	case isLetter(p.r):
		lit := p.scanIdentifier()
		if len(lit) > 1 {
			// a possible keyword
			keyword, ok := p.keywords[lit]
			if ok {
				return keyword
			}
		}
		p.literal = lit
		return TOKEN_ID
	case isDigit(p.r):
		lit := p.scanNumber()
		p.literal = lit
		return TOKEN_NUMBER
	default:
		if p.r == '/' {
			if p.scanComment() {
				return p.nextToken()
			}
			return TOKEN_DIV
		}
		switch p.r {
		case '+': return TOKEN_PLUS
		case '-': return TOKEN_MINUS
		case '*': return TOKEN_MULT
		case '<': return p.maybe2(TOKEN_LT, '=', TOKEN_LTEQ, '<', TOKEN_LTLT)
		case '>': return p.maybe2(TOKEN_GT, '=', TOKEN_GTEQ, '>', TOKEN_GTGT)
		case '=': return p.maybe1(TOKEN_EQ, '=', TOKEN_EQEQ)
		case ':': return TOKEN_COLON
		case '[': return TOKEN_LBRACK
		case ']': return TOKEN_RBRACK
		case '{': return TOKEN_LBRACE
		case '}': return TOKEN_RBRACE
		}
	}

	return TOKEN_ILLEGAL
}