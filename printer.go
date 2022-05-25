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
	return expr.Accept(p).(string)
}

func (p *AstPrinter) visitBinaryExpr(expr *Binary) any {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p *AstPrinter) visitGroupingExpr(expr *Grouping) any {
	return p.parenthesize("group", expr.expression)
}

func (p *AstPrinter) visitLiteralExpr(expr *Literal) any {
	if expr.value == nil {
		return "nil"
	}
	switch vt := expr.value.(type) {
	case bool:
		if vt {
			return "true"
		} else {
			return "false"
		}
	case string:
		return vt
	case int:
		return fmt.Sprintf("%d", vt)
	case float64:
		return fmt.Sprintf("%.2f", vt)
	}
	return "err" // TODO - better error handling
}

func (p *AstPrinter) visitUnaryExpr(expr *Unary) any {
	return p.parenthesize(expr.operator.lexeme, expr.right)
}

// parenthesize - private helper function
func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder
	sb.WriteString("(" + name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(expr.Accept(p).(string))
	}
	sb.WriteString(")")
	return sb.String()
}
