// This file is generated from astgen.go
package golox

type Expr interface {
  Accept(visitor ExprVisitor) (any, error)
}

type ExprVisitor interface {
  visitAssignExpr(expr *Assign) (any, error)
  visitBinaryExpr(expr *Binary) (any, error)
  visitGroupingExpr(expr *Grouping) (any, error)
  visitLiteralExpr(expr *Literal) (any, error)
  visitLogicalExpr(expr *Logical) (any, error)
  visitUnaryExpr(expr *Unary) (any, error)
  visitVariableExpr(expr *Variable) (any, error)
}

type Assign struct {
  name *Token
  value Expr
}

func NewAssign(name *Token, value Expr) *Assign {
  return &Assign{
    name: name,
    value: value,
  }
}

func (expr *Assign) Accept(visitor ExprVisitor) (any, error) {
  return visitor.visitAssignExpr(expr)
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

func (expr *Binary) Accept(visitor ExprVisitor) (any, error) {
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

func (expr *Grouping) Accept(visitor ExprVisitor) (any, error) {
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

func (expr *Literal) Accept(visitor ExprVisitor) (any, error) {
  return visitor.visitLiteralExpr(expr)
}

type Logical struct {
  left Expr
  operator *Token
  right Expr
}

func NewLogical(left Expr, operator *Token, right Expr) *Logical {
  return &Logical{
    left: left,
    operator: operator,
    right: right,
  }
}

func (expr *Logical) Accept(visitor ExprVisitor) (any, error) {
  return visitor.visitLogicalExpr(expr)
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

func (expr *Unary) Accept(visitor ExprVisitor) (any, error) {
  return visitor.visitUnaryExpr(expr)
}

type Variable struct {
  name *Token
}

func NewVariable(name *Token) *Variable {
  return &Variable{
    name: name,
  }
}

func (expr *Variable) Accept(visitor ExprVisitor) (any, error) {
  return visitor.visitVariableExpr(expr)
}

