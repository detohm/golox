package golox

import "fmt"

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) (any, error) {
	value, ok := e.values[name.lexeme]
	if ok {
		return value, nil
	}

	return nil, NewRuntimeError(*name,
		fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}
