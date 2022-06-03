// This file is generated from astgen.go
package golox

type Stmt interface {
  Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
  visitBlockStmt(stmt *Block) (any, error)
  visitExpressionStmt(stmt *Expression) (any, error)
  visitIfStmt(stmt *If) (any, error)
  visitPrintStmt(stmt *Print) (any, error)
  visitVarStmt(stmt *Var) (any, error)
  visitWhileStmt(stmt *While) (any, error)
}

type Block struct {
  statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
  return &Block{
    statements: statements,
  }
}

func (expr *Block) Accept(visitor StmtVisitor) (any, error) {
  return visitor.visitBlockStmt(expr)
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

type If struct {
  condition Expr
  thenBranch Stmt
  elseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) *If {
  return &If{
    condition: condition,
    thenBranch: thenBranch,
    elseBranch: elseBranch,
  }
}

func (expr *If) Accept(visitor StmtVisitor) (any, error) {
  return visitor.visitIfStmt(expr)
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

type While struct {
  condition Expr
  body Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
  return &While{
    condition: condition,
    body: body,
  }
}

func (expr *While) Accept(visitor StmtVisitor) (any, error) {
  return visitor.visitWhileStmt(expr)
}

