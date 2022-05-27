package golox

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool
	interpreter     *Interpreter
}

func NewLox() *Lox {
	lox := &Lox{
		hadError:        false,
		hadRuntimeError: false,
	}
	lox.interpreter = NewInterpreter(lox)
	return lox
}

func (l *Lox) Main(args []string) {
	if len(args) > 2 {
		fmt.Println("Usage: golox [script]")
		return
	}
	if len(args) == 2 {
		if err := l.runFile(args[0]); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := l.runPrompt(); err != nil {
			fmt.Println(err)
		}
	}
}

func (l *Lox) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	l.run(string(bytes))
	if l.hadError {
		os.Exit(65)
	}
	if l.hadRuntimeError {
		os.Exit(70)
	}
	return nil
}

func (l *Lox) runPrompt() error {
	for {

		fmt.Print("> ")

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		l.run(line)
		l.hadError = false
	}
}

func (l *Lox) run(source string) {

	scanner := NewScanner(l, source)
	tokens := scanner.scanTokens()

	parser := NewParser(l, tokens)
	expression := parser.Parse()

	if l.hadError {
		return
	}

	l.interpreter.interpret(expression)
	// fmt.Println(NewAstPrinter().Print(expression))

}

func (l *Lox) Error(line int, message string) {
	l.Report(line, "", message)
}

func (l *Lox) RuntimeError(err RuntimeError) {
	fmt.Printf("%s\n[line %d]\n", err.Message, err.Token.line)
	l.hadRuntimeError = true
}

func (l *Lox) ErrorWithToken(token Token, message string) {
	if token.kind == TkEof {
		l.Report(token.line, " at end", message)
	} else {
		l.Report(token.line, fmt.Sprintf(" at '%s'", token.lexeme), message)
	}
}

func (l *Lox) Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}
