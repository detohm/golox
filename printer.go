package golox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(expr Expr) string {
	str, err := expr.Accept(p)
	if err != nil {
		// quick and dirty return
		return "error" + err.Error()
	}
	return str.(string)
}

func (p *AstPrinter) visitLogicalExpr(expr *Logical) (any, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p *AstPrinter) visitBinaryExpr(expr *Binary) (any, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p *AstPrinter) visitGroupingExpr(expr *Grouping) (any, error) {
	return p.parenthesize("group", expr.expression)
}

func (p *AstPrinter) visitLiteralExpr(expr *Literal) (any, error) {
	if expr.value == nil {
		return "nil", nil
	}
	switch vt := expr.value.(type) {
	case bool:
		if vt {
			return "true", nil
		} else {
			return "false", nil
		}
	case string:
		return vt, nil
	case int:
		return fmt.Sprintf("%d", vt), nil
	case float64:
		return fmt.Sprintf("%.2f", vt), nil
	}
	// TODO - better error handling
	return "err", nil
}

func (p *AstPrinter) visitUnaryExpr(expr *Unary) (any, error) {
	return p.parenthesize(expr.operator.lexeme, expr.right)
}

func (p *AstPrinter) visitVariableExpr(expr *Variable) (any, error) {
	return "(var)", nil
}

func (p *AstPrinter) visitAssignExpr(expr *Assign) (any, error) {
	return fmt.Sprintf("(assign %s:%s)", expr.name.lexeme, expr.value), nil
}

func (p *AstPrinter) visitCallExpr(expr *Call) (any, error) {
	return fmt.Sprintf("(call %s)", expr.callee), nil
}

// parenthesize - private helper function
func (p *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	var sb strings.Builder
	sb.WriteString("(" + name)
	for _, expr := range exprs {
		sb.WriteString(" ")

		v, err := expr.Accept(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(v.(string))
	}
	sb.WriteString(")")
	return sb.String(), nil
}
