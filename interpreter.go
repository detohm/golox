package golox

import (
	"fmt"
	"math"
	"reflect"
)

type Interpreter struct {
	lox         *Lox
	environment *Environment
}

func NewInterpreter(lox *Lox) *Interpreter {
	return &Interpreter{
		lox:         lox,
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			// TODO - better type casting
			i.lox.RuntimeError(err.(RuntimeError))
			return
		}
	}
}

/* Intepreter implements on both ExprVisitor and StmtVisitor interfaces

type ExprVisitor interface {
  visitBinaryExpr(expr *Binary) (any, error)
  visitGroupingExpr(expr *Grouping) (any, error)
  visitLiteralExpr(expr *Literal) (any, error)
  visitUnaryExpr(expr *Unary) (any, error)
}

type StmtVisitor interface {
  visitExpressionStmt(stmt *Expression) (any, error)
  visitPrintStmt(stmt *Print) (any, error)
}

*/
func (i *Interpreter) visitExpressionStmt(stmt *Expression) (any, error) {
	_, err := i.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) visitIfStmt(stmt *If) (any, error) {
	cond, err := i.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(cond) {
		err = i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		err = i.execute(stmt.elseBranch)
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) visitPrintStmt(stmt *Print) (any, error) {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) visitVarStmt(stmt *Var) (any, error) {
	var value any = nil
	var err error = nil
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.define(stmt.name.lexeme, value)
	return nil, nil
}

func (i *Interpreter) visitWhileStmt(stmt *While) (any, error) {
	cond, err := i.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}

	for i.isTruthy(cond) {
		err := i.execute(stmt.body)
		if err != nil {
			return nil, err
		}
		cond, err = i.evaluate(stmt.condition)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) visitAssignExpr(expr *Assign) (any, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	err = i.environment.assign(expr.name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) visitLiteralExpr(expr *Literal) (any, error) {
	return expr.value, nil
}

func (i *Interpreter) visitLogicalExpr(expr *Logical) (any, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	if expr.operator.kind == TkOr {
		// short circuit for OR
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		// short circuit for AND
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.right)

}

func (i *Interpreter) visitGroupingExpr(expr *Grouping) (any, error) {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitUnaryExpr(expr *Unary) (any, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.kind {
	case TkBang:
		return !i.isTruthy(right), nil
	case TkMinus:
		err := i.checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -(right.(float64)), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) visitVariableExpr(expr *Variable) (any, error) {
	return i.environment.get(expr.name)
}

func (i *Interpreter) visitBinaryExpr(expr *Binary) (any, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.kind {

	case TkGreater:
		if err := i.checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case TkGreaterEqual:
		if err := i.checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case TkLess:
		if err := i.checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case TkLessEqual:
		if err := i.checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	// comparison case
	case TkBangEqual:
		return !i.isEqual(left, right), nil
	case TkEqualEqual:
		return i.isEqual(left, right), nil

	case TkMinus:
		if err := i.checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case TkPlus:

		// add
		if reflect.TypeOf(left).Kind() == reflect.Float64 &&
			reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) + right.(float64), nil
		}

		// concatinate
		if reflect.TypeOf(left).Kind() == reflect.String &&
			reflect.TypeOf(right).Kind() == reflect.String {
			return left.(string) + right.(string), nil
		}

		// otherwise, error
		return nil, NewRuntimeError(*expr.operator, "Operands must be two numbers or two strings.")

	case TkSlash:
		if err := i.checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case TkStar:
		if err := i.checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) checkNumberOperand(operator *Token, operand any) error {
	if reflect.TypeOf(operand).Kind() == reflect.Float64 {
		return nil
	}
	return NewRuntimeError(*operator, "Operand must be a number")
}

func (i *Interpreter) checkNumberOperands(operator *Token, left any, right any) error {
	if reflect.TypeOf(left).Kind() == reflect.Float64 &&
		reflect.TypeOf(right).Kind() == reflect.Float64 {
		return nil
	}
	return NewRuntimeError(*operator, "Operands must be numbers")
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) error {
	_, err := stmt.Accept(i)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) visitBlockStmt(stmt *Block) (any, error) {
	err := i.executeBlock(stmt.statements,
		NewEnvironmentWithEnclosing(i.environment))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := i.environment
	i.environment = environment
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	i.environment = previous
	return nil
}

func (i *Interpreter) isTruthy(ex any) bool {
	if ex == nil {
		return false
	}
	if reflect.TypeOf(ex).Kind() == reflect.Bool {
		return ex.(bool)
	}
	return true
}

func (i *Interpreter) isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	// TODO - ensure the golang comparison machanism
	return a == b
}

func (i *Interpreter) stringify(a any) string {
	if a == nil {
		return "nil"
	}
	if reflect.TypeOf(a).Kind() == reflect.Float64 {

		value := a.(float64)
		if math.Mod(value, 1.0) == 0 {
			return fmt.Sprintf("%.0f", value)
		}
		return fmt.Sprintf("%.2f", value)

	}

	return fmt.Sprintf("%v", a)
}
