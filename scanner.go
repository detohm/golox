package golox

import (
	"strconv"
)

type Scanner struct {
	lox     *Lox
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"and":    TkAnd,
	"class":  TkClass,
	"else":   TkElse,
	"false":  TkFalse,
	"for":    TkFor,
	"fun":    TkFun,
	"if":     TkIf,
	"nil":    TkNil,
	"or":     TkOr,
	"print":  TkPrint,
	"return": TkReturn,
	"super":  TkSuper,
	"this":   TkThis,
	"true":   TkTrue,
	"var":    TkVar,
	"while":  TkWhile,
}

func NewScanner(lox *Lox, source string) *Scanner {
	return &Scanner{
		lox:     lox,
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	return s.tokens

}

func (s *Scanner) scanToken() {

	c := s.advance()

	switch c {
	case '(':
		s.addToken(TkLeftParen)
	case ')':
		s.addToken(TkRightParen)
	case '{':
		s.addToken(TkLeftBrace)
	case '}':
		s.addToken(TkRightBrace)
	case ',':
		s.addToken(TkComma)
	case '.':
		s.addToken(TkDot)
	case '-':
		s.addToken(TkMinus)
	case '+':
		s.addToken(TkPlus)
	case ';':
		s.addToken(TkSemicolon)
	case '*':
		s.addToken(TkStar)
	// one or two characters
	case '!':
		if s.match('=') {
			s.addToken(TkBangEqual)
		} else {
			s.addToken(TkBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(TkEqualEqual)
		} else {
			s.addToken(TkEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(TkLessEqual)
		} else {
			s.addToken(TkLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(TkGreaterEqual)
		} else {
			s.addToken(TkGreater)
		}
	// longer lexemes
	case '/':
		if s.match('/') {
			// keep consuming characters until the end of line
			// no addToken call
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TkSlash)
		}
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		// string literal
		s.string()

	default:
		if s.isDigit(c) {
			// number literal
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.lox.Error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = TkIdentifier
	}
	s.addToken(tokenType)
}

// string - consume string literal
func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.lox.Error(s.line, "Unterminated string.")
		return
	}

	// the closing "
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(TkString, value)
}

// number - consume number literal
func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.lox.Error(s.line, "Parse number literal error.")
		return
	}
	s.addTokenWithLiteral(TkNumber, num)
}

// match - is like conditional advance()
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

// peek - like advance but not consume the character
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

// peekNext - lookahead for the next characters of peek
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(kind TokenType) {
	s.addTokenWithLiteral(kind, nil)
}

func (s *Scanner) addTokenWithLiteral(kind TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{kind, text, literal, s.line})
}
