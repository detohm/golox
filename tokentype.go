package golox

type TokenType int

// Prefix with Tk...
const (
	// Single-character token
	TkLeftParen TokenType = iota
	TkRightParen
	TkLeftBrace
	TkRightBrace
	TkComma
	TkDot
	TkMinus
	TkPlus
	TkSemicolon
	TkSlash
	TkStar

	// One or two character token
	TkBang
	TkBangEqual
	TkEqual
	TkEqualEqual
	TkGreater
	TkGreaterEqual
	TkLess
	TkLessEqual

	// Literal
	TkIdentifier
	TkString
	TkNumber

	// Keywords
	TkAnd
	TkClass
	TkElse
	TkFalse
	TkFun
	TkFor
	TkIf
	TkNil
	TkOr
	TkPrint
	TkReturn
	TkSuper
	TkThis
	TkTrue
	TkVar
	TkWhile

	TkEof
)

// TODO - implement string function for readable token type
// func (t TokenType) String() string {
// 	switch t {
// 	case TkLeftParen:
// 		return "LeftParen"
// 	}
// 	return "not recognised"
// }
