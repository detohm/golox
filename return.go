package golox

// ReturnValue - act as Runtime Exception for jacked up the control flow
// when the return is executed, the interpreter nees to jump all the wat out
// of the current context and complete the function call
type ReturnValue struct {
	value any
}

func NewReturnValue(value any) ReturnValue {
	return ReturnValue{value}
}

func (r ReturnValue) Error() string {
	return r.value.(string)
}
