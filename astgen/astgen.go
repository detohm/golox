package astgen

import (
	"fmt"
	"os"
	"strings"
)

func Generate(outputDir string) {
	defineAst(outputDir, "Expr", []string{
		"Binary : left Expr, operator *Token, right Expr",
		"Grouping : expression Expr",
		"Literal : value any",
		"Unary : operator *Token, right Expr",
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
	file.WriteString("type Expr interface {\n")
	file.WriteString("  Accept(visitor Visitor) any\n")
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
	file.WriteString("type Visitor interface {\n")
	for _, v := range types {
		typeName := strings.TrimSpace(strings.Split(v, ":")[0])
		file.WriteString(fmt.Sprintf("  visit%s%s(%s *%s) any\n",
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
	file.WriteString(fmt.Sprintf("func (expr *%s) Accept(visitor Visitor) any {\n", typeName))
	file.WriteString(fmt.Sprintf("  return visitor.visit%s%s(expr)\n",
		typeName,
		baseName,
	))

	file.WriteString("}\n\n")

}
