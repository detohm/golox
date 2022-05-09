package golox

import (
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

func NewLox() *Lox {
	return &Lox{
		hadError: false,
	}
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
	return nil
}

func (l *Lox) runPrompt() error {
	for {
		var line string
		fmt.Print("> ")
		if _, err := fmt.Scan(&line); err != nil {
			return err
		}
		l.run(line)
		l.hadError = false
	}
}

func (l *Lox) run(source string) {
	scanner := NewScanner(l, source)
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Printf("%s\n", token)
	}
}

func (l *Lox) Error(line int, message string) {
	l.Report(line, "", message)
}

func (l *Lox) Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s", line, where, message)
	l.hadError = true
}
