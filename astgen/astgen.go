package astgen

import (
	"fmt"
	"os"
	"strings"
)

func Generate(outputDir string) {
	defineAst(outputDir, "expr", []string{
		"Binary : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal : value any",
		"Unary : operator Token, right Expr",
	})
}

func defineAst(outputDir string, baseName string, types []string) error {
	path := outputDir + "/" + baseName + ".go"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.WriteString("// This file is generated from astgen.go\n")
	file.WriteString("package golox\n\n")

	// Expr interface
	file.WriteString("type Expr interface {}\n")

	for _, v := range types {
		typeName := strings.TrimSpace(strings.Split(v, ":")[0])
		fields := strings.TrimSpace(strings.Split(v, ":")[1])
		defineType(file, baseName, typeName, fields)
	}

	file.Close()
	return nil
}

func defineType(file *os.File, baseName string, typeName string, fieldList string) {

	fields := strings.Split(fieldList, ", ")

	file.WriteString(fmt.Sprintf("type %s struct {\n", typeName))
	for _, field := range fields {
		file.WriteString(fmt.Sprintf("  %s\n", field))
	}
	file.WriteString("}\n\n")

	file.WriteString(fmt.Sprintf("func New%s(%s) *%s {\n", typeName, fieldList, typeName))

	file.WriteString(fmt.Sprintf("  return &%s{\n", typeName))

	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		file.WriteString(fmt.Sprintf("    %s: %s,\n", name, name))
	}

	file.WriteString("  }\n")
	file.WriteString("}\n\n")
}
