package golox

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    make(map[string]any),
	}
}

func NewEnvironmentWithEnclosing(enc *Environment) *Environment {
	en := NewEnvironment()
	en.enclosing = enc
	return en
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) (any, error) {
	value, ok := e.values[name.lexeme]
	if ok {
		return value, nil
	}

	// find a var in outter scope
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, NewRuntimeError(*name,
		fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}

func (e *Environment) assign(name *Token, value any) error {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		err := e.enclosing.assign(name, value)
		return err
	}

	return NewRuntimeError(*name,
		fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}
