// This file is generated from astgen.go
package golox

type Expr interface {}
type Binary struct {
  left Expr
  operator Token
  right Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
  return &Binary{
    left: left,
    operator: operator,
    right: right,
  }
}

type Grouping struct {
  expression Expr
}

func NewGrouping(expression Expr) *Grouping {
  return &Grouping{
    expression: expression,
  }
}

type Literal struct {
  value any
}

func NewLiteral(value any) *Literal {
  return &Literal{
    value: value,
  }
}

type Unary struct {
  operator Token
  right Expr
}

func NewUnary(operator Token, right Expr) *Unary {
  return &Unary{
    operator: operator,
    right: right,
  }
}

