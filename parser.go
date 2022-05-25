package golox

/* Grammar Rules
expression     -> equality ;
equality       -> comparison ( ( "!=" | "==" ) comparison )* ;
comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           -> factor ( ( "-" | "+" ) factor )* ;
factor         -> unary ( ( "/" | "*" ) unary )* ;
unary          -> ( "!" | "-" ) unary
               | primary ;
primary        -> NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
*/
type Parser struct {
	lox     *Lox
	tokens  []Token
	current int
}

type ParseError struct{}

func (e ParseError) Error() string { return "" }

func NewParser(lox *Lox, tokens []Token) *Parser {
	return &Parser{
		lox:     lox,
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() Expr {
	expr, err := p.expression()
	if err != nil {
		return nil
	}
	return expr
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

// equality  ->  comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(TkBangEqual, TkEqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		// left-associative nested tree
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// comparison  ->  term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(TkGreater, TkGreaterEqual, TkLess, TkLessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// term  ->  factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(TkMinus, TkPlus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// factor  ->  unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(TkSlash, TkStar) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}
	return expr, nil
}

// unary  ->  ( "!" | "-" ) unary
//        | primary ;
func (p *Parser) unary() (Expr, error) {
	if p.match(TkBang, TkMinus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
	}
	return p.primary()
}

// primary  ->  NUMBER | STRING | "true" | "false" | "nil"
//          | "(" expression ")" ;
func (p *Parser) primary() (Expr, error) {

	if p.match(TkFalse) {
		return NewLiteral(false), nil
	}
	if p.match(TkTrue) {
		return NewLiteral(true), nil
	}
	if p.match(TkNil) {
		return NewLiteral(nil), nil
	}
	if p.match(TkNumber, TkString) {
		return NewLiteral(p.previous().literal), nil
	}
	if p.match(TkLeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(TkRightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewGrouping(expr), nil
	}

	return nil, p.error(p.peek(), "Expect expression.")

}

func (p *Parser) match(kinds ...TokenType) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(kind TokenType, message string) (*Token, error) {
	if p.check(kind) {
		return p.advance(), nil
	}
	return nil, p.error(p.peek(), message)
}

func (p *Parser) check(kind TokenType) bool {

	if p.isAtEnd() {
		return false
	}
	return p.peek().kind == kind
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().kind == TkEof
}

func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) error(token *Token, message string) ParseError {
	p.lox.ErrorWithToken(*token, message)
	return ParseError{}
}

func (p *Parser) synchronize() {
	p.advance()
	// find the beginning of the next statement
	for !p.isAtEnd() {

		if p.previous().kind == TkSemicolon {
			return
		}

		switch p.peek().kind {
		case
			TkClass,
			TkFun,
			TkVar,
			TkFor,
			TkIf,
			TkWhile,
			TkPrint,
			TkReturn:
			return
		}

		p.advance()
	}
}
