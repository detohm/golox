package astgen

import (
	"fmt"
	"os"
	"strings"
)

func Generate(outputDir string) {

	defineAst(outputDir, "Expr", []string{
		"Assign : name *Token, value Expr",
		"Binary : left Expr, operator *Token, right Expr",
		"Grouping : expression Expr",
		"Literal : value any",
		"Logical : left Expr, operator *Token, right Expr",
		"Unary : operator *Token, right Expr",
		"Variable : name *Token",
	})

	defineAst(outputDir, "Stmt", []string{
		"Block : statements []Stmt",
		"Expression : expression Expr",
		"If : condition Expr, thenBranch Stmt, elseBranch Stmt",
		"Print : expression Expr",
		"Var : name *Token, initializer Expr",
		"While : condition Expr, body Stmt",
	})

}

func defineAst(outputDir string, baseName string, types []string) error {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.WriteString("// This file is generated from astgen.go\n")
	file.WriteString("package golox\n\n")

	// Expr interface
	file.WriteString(fmt.Sprintf("type %s interface {\n", baseName))
	file.WriteString(fmt.Sprintf("  Accept(visitor %sVisitor) (any, error)\n", baseName))
	file.WriteString("}\n\n")

	defineVisitor(file, baseName, types)

	// The AST type
	for _, v := range types {
		typeName := strings.TrimSpace(strings.Split(v, ":")[0])
		fields := strings.TrimSpace(strings.Split(v, ":")[1])
		defineType(file, baseName, typeName, fields)
	}

	file.Close()
	return nil
}

func defineVisitor(file *os.File, baseName string, types []string) {
	file.WriteString(fmt.Sprintf("type %sVisitor interface {\n", baseName))
	for _, v := range types {
		typeName := strings.TrimSpace(strings.Split(v, ":")[0])
		file.WriteString(fmt.Sprintf("  visit%s%s(%s *%s) (any, error)\n",
			typeName,
			baseName,
			strings.ToLower(baseName),
			typeName,
		))

	}
	file.WriteString("}\n\n")
}

func defineType(file *os.File, baseName string, typeName string, fieldList string) {

	fields := strings.Split(fieldList, ", ")

	// type
	file.WriteString(fmt.Sprintf("type %s struct {\n", typeName))
	for _, field := range fields {
		file.WriteString(fmt.Sprintf("  %s\n", field))
	}
	file.WriteString("}\n\n")

	// constructor
	file.WriteString(fmt.Sprintf("func New%s(%s) *%s {\n", typeName, fieldList, typeName))

	file.WriteString(fmt.Sprintf("  return &%s{\n", typeName))

	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		file.WriteString(fmt.Sprintf("    %s: %s,\n", name, name))
	}

	file.WriteString("  }\n")
	file.WriteString("}\n\n")

	// visitor pattern
	file.WriteString(fmt.Sprintf("func (expr *%s) Accept(visitor %sVisitor) (any, error) {\n", typeName, baseName))
	file.WriteString(fmt.Sprintf("  return visitor.visit%s%s(expr)\n",
		typeName,
		baseName,
	))

	file.WriteString("}\n\n")

}
