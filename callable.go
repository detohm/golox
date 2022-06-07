package golox

import "time"

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []any) (any, error)
}

type clock struct{}

func NewClock() *clock {
	return &clock{}
}

func (c *clock) arity() int {
	return 0
}

func (c *clock) call(i *Interpreter, arguments []any) (any, error) {
	return time.Now().Unix(), nil
}

func (c *clock) String() string {
	return "<native fn>"
}
