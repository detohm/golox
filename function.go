package golox

import "fmt"

type loxFunction struct {
	declaration *Function
}

func NewLoxFunction(declaration *Function) *loxFunction {
	return &loxFunction{
		declaration: declaration,
	}
}

func (f *loxFunction) call(interpreter *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironmentWithEnclosing(interpreter.globals)
	for i := 0; i < len(f.declaration.params); i++ {
		environment.define(
			f.declaration.params[i].lexeme,
			arguments[i],
		)
	}
	err := interpreter.executeBlock(f.declaration.body, environment)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *loxFunction) arity() int {
	return len(f.declaration.params)
}

func (f *loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
}
