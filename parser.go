package golox

import (
	"fmt"
)

/* Grammar Rules

program        -> declaration* EOF ;

declaration    -> funDecl
			   |varDecl
               | statement ;

funDecl        -> "fun" function ;
function       -> IDENTIFIER "(" parameters? ")" block ;

parameters	   -> IDENTIFIER ( "," IDENTIFIER )* ;

varDecl        -> "var" IDENTIFIER ( "=" expression )? ";" ;

statement      -> exprStmt
			   | forStmt
			   | ifStmt
               | printStmt
			   | whileStmt
			   | block ;

ifStmt         -> "if" "(" expression ")" statement
               ( "else" statement )? ;

block          -> "{" declaration* "}" ;

forStmt        -> "for" "(" (varDecl | exprStmt | ";")
			   expression? ";"
               expression? ")" statement ;

whileStmt	   -> "while" "(" expression ")" statement ;

exprStmt       -> expression ";" ;
printStmt      -> "print" expression ";" ;

expression     -> assignment ;
assignment     -> IDENTIFIER "=" assignment
               | logic_or ;

logic_or       -> logic_and ( "or" logic_and )* ;
logic_and      -> equality ( "and" equality )* ;

equality       -> comparison ( ( "!=" | "==" ) comparison )* ;
comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           -> factor ( ( "-" | "+" ) factor )* ;
factor         -> unary ( ( "/" | "*" ) unary )* ;
unary          -> ( "!" | "-" ) unary | call ;
call		   -> primary ( "(" arguments? ")" )* ;
               | primary ;
argument 	   -> expression ("," expression )* ;

primary        -> "true" | "false" | "nil"
			   | NUMBER | STRING |
               | "(" expression ")"
			   | IDENTIFIER ;
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

func (p *Parser) Parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			// TODO - better error handling
			fmt.Print(err)
			return statements
		}
		statements = append(statements, stmt)
	}
	return statements
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(TkFun) {
		return p.function("function")
	}
	if p.match(TkVar) {
		stmt, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, nil // cut off error propagation
		}
		return stmt, nil
	}

	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil, nil
	}
	return stmt, nil

}

func (p *Parser) function(kind string) (Stmt, error) {
	name, err := p.consume(TkIdentifier,
		fmt.Sprintf("Expect %s name.", kind),
	)
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TkLeftParen,
		fmt.Sprintf("Expect '(' after %s name", kind),
	)
	if err != nil {
		return nil, err
	}

	parameters := []*Token{}

	for {

		if len(parameters) >= 255 {
			p.error(p.peek(), "Can't have more than 255 parameters.")
		}

		t, err := p.consume(TkIdentifier, "Expect parameter name.")
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, t)

		if !p.match(TkComma) {
			break
		}
	}

	_, err = p.consume(TkRightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(TkLeftBrace,
		fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return NewFunction(name, parameters, body), nil

}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(TkIdentifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.match(TkEqual) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}

	}
	_, err = p.consume(TkSemicolon, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return NewVar(name, initializer), nil

}

func (p *Parser) statement() (Stmt, error) {
	if p.match(TkFor) {
		return p.forStatement()
	}
	if p.match(TkIf) {
		return p.ifStatement()
	}
	if p.match(TkPrint) {
		return p.printStatement()
	}
	if p.match(TkWhile) {
		return p.whileStatement()
	}
	if p.match(TkLeftBrace) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewBlock(statements), nil
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(TkLeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TkRightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt = nil
	if p.match(TkElse) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return NewIf(condition, thenBranch, elseBranch), nil

}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TkSemicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return NewPrint(expr), nil
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(TkLeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt = nil
	if p.match(TkSemicolon) {
		initializer = nil
	} else if p.match(TkVar) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr = nil
	if !p.check(TkSemicolon) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(TkSemicolon, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr = nil
	if !p.check(TkRightParen) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(TkRightParen, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		statements := []Stmt{
			body,
			NewExpression(increment),
		}
		body = NewBlock(statements)
	}

	if condition == nil {
		condition = NewLiteral(true)
	}
	body = NewWhile(condition, body)

	if initializer != nil {
		body = NewBlock([]Stmt{
			initializer,
			body,
		})
	}

	return body, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(TkLeftParen, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TkRightParen, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return NewWhile(condition, body), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(TkSemicolon, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}
	return NewExpression(expr), nil
}

func (p *Parser) block() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.check(TkRightBrace) && !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, dec)
	}
	_, err := p.consume(TkRightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(TkEqual) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		// only l-value is allowed
		v, ok := expr.(*Variable)
		if ok {
			name := v.name
			return NewAssign(name, value), nil
		}

		err = p.error(equals, "Invalid assignment target.")
		return nil, err
	}
	return expr, nil
}

// logic_or  ->  logic_and ( "or" logic_and )* ;
func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(TkOr) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}
	return expr, nil
}

// logic_and  ->  equality ( "and" equality )* ;
func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(TkAnd) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}
	return expr, nil
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
//        | call ;
func (p *Parser) unary() (Expr, error) {
	if p.match(TkBang, TkMinus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
	}
	return p.call()
}

// call		   -> primary ( "(" arguments? ")" )* ;
//             | primary ;
func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(TkLeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := []Expr{}
	if !p.check(TkRightParen) {

		for {

			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}

			ex, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, ex)

			if !p.match(TkComma) {
				break
			}
		}
	}
	paren, err := p.consume(TkRightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return NewCall(callee, paren, arguments), nil
}

// primary -> "true" | "false" | "nil"
// 			| NUMBER | STRING |
// 			| "(" expression ")"
// 			| IDENTIFIER ;
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

	if p.match(TkIdentifier) {
		return NewVariable(p.previous()), nil
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
