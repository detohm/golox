package golox

type RuntimeError struct {
	Token   Token
	Message string
}

func NewRuntimeError(t Token, message string) RuntimeError {
	return RuntimeError{
		Token:   t,
		Message: message,
	}
}

func (e RuntimeError) Error() string {
	return e.Message
}
