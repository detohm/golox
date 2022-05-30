// This file is generated from astgen.go
package golox

type Stmt interface {
  Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
  visitExpressionStmt(stmt *Expression) (any, error)
  visitPrintStmt(stmt *Print) (any, error)
  visitVarStmt(stmt *Var) (any, error)
}

type Expression struct {
  expression Expr
}

func NewExpression(expression Expr) *Expression {
  return &Expression{
    expression: expression,
  }
}

func (expr *Expression) Accept(visitor StmtVisitor) (any, error) {
  return visitor.visitExpressionStmt(expr)
}

type Print struct {
  expression Expr
}

func NewPrint(expression Expr) *Print {
  return &Print{
    expression: expression,
  }
}

func (expr *Print) Accept(visitor StmtVisitor) (any, error) {
  return visitor.visitPrintStmt(expr)
}

type Var struct {
  name *Token
  initializer Expr
}

func NewVar(name *Token, initializer Expr) *Var {
  return &Var{
    name: name,
    initializer: initializer,
  }
}

func (expr *Var) Accept(visitor StmtVisitor) (any, error) {
  return visitor.visitVarStmt(expr)
}

