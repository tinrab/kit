package kit

import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Parent  error
	Code    uint64
	Message string
	Args    []interface{}
	File    string
	Line    int
}

func NewCodeError(code uint64, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewMessageError(message string) *Error {
	return &Error{
		Message: message,
	}
}

func (e *Error) With(args ...interface{}) *Error {
	e.Args = args
	return e
}

func (e *Error) Throw(args ...interface{}) *Error {
	e.Args = args
	_, file, line, _ := runtime.Caller(1)
	e.Line = line

	fe := strings.LastIndex(file, "/") + 1
	if fe < len(file) {
		e.File = file[fe:]
	} else {
		e.File = file
	}

	return e
}

func (e *Error) Wrap(err error) *Error {
	e.Parent = err
	return e
}

func (e *Error) Walk(step func(err error)) {
	err := e
	for err != nil {
		step(err)
		err, _ = e.Parent.(*Error)
	}
	if e.Parent != nil {
		step(e.Parent)
	}
}

func (e *Error) Error() string {
	if e.File == "" {
		if len(e.Args) == 0 {
			return e.Message
		}
		return fmt.Sprintf(e.Message, e.Args...)
	}
	if len(e.Args) == 0 {
		return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
	}
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, fmt.Sprintf(e.Message, e.Args...))
}
