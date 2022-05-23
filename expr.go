// This file is generated from astgen.go
package golox

type Expr interface {
  Accept(visitor Visitor) any
}

type Visitor interface {
  visitBinaryExpr(expr *Binary) any
  visitGroupingExpr(expr *Grouping) any
  visitLiteralExpr(expr *Literal) any
  visitUnaryExpr(expr *Unary) any
}

type Binary struct {
  left Expr
  operator *Token
  right Expr
}

func NewBinary(left Expr, operator *Token, right Expr) *Binary {
  return &Binary{
    left: left,
    operator: operator,
    right: right,
  }
}

func (expr *Binary) Accept(visitor Visitor) any {
  return visitor.visitBinaryExpr(expr)
}

type Grouping struct {
  expression Expr
}

func NewGrouping(expression Expr) *Grouping {
  return &Grouping{
    expression: expression,
  }
}

func (expr *Grouping) Accept(visitor Visitor) any {
  return visitor.visitGroupingExpr(expr)
}

type Literal struct {
  value any
}

func NewLiteral(value any) *Literal {
  return &Literal{
    value: value,
  }
}

func (expr *Literal) Accept(visitor Visitor) any {
  return visitor.visitLiteralExpr(expr)
}

type Unary struct {
  operator *Token
  right Expr
}

func NewUnary(operator *Token, right Expr) *Unary {
  return &Unary{
    operator: operator,
    right: right,
  }
}

func (expr *Unary) Accept(visitor Visitor) any {
  return visitor.visitUnaryExpr(expr)
}

