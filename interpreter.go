package golox

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	lox *Lox
}

func NewInterpreter(lox *Lox) *Interpreter {
	return &Interpreter{
		lox: lox,
	}
}

func (i *Interpreter) interpret(expression Expr) {
	value, err := i.evaluate(expression)
	if err != nil {
		// TODO - better type casting
		i.lox.RuntimeError(err.(RuntimeError))
	}
	fmt.Printf("%v\n", value)
}

/*
type Visitor interface {
  visitBinaryExpr(expr *Binary) any
  visitGroupingExpr(expr *Grouping) any
  visitLiteralExpr(expr *Literal) any
  visitUnaryExpr(expr *Unary) any
}
*/

func (i *Interpreter) visitLiteralExpr(expr *Literal) (any, error) {
	return expr.value, nil
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