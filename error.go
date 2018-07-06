package util

import (
	"fmt"
	"runtime"
)

// Error is main error type in syncro package.
type Error struct {
	Message string
	Args    []interface{}
	File    string
	Line    int
}

// NewError creates new error with a domain and message.
func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}

// Throw sets the file name, line number and any arguments for the error.
func (e *Error) Throw(args ...interface{}) *Error {
	e.Args = args
	_, e.File, e.Line, _ = runtime.Caller(1)
	return e
}

func (e *Error) Error() string {
	if e.File == "" {
		if len(e.Args) == 0 {
			return e.Message
		}
		return fmt.Sprintf(e.Message, e.Args)
	}
	if len(e.Args) == 0 {
		return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
	}
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, fmt.Sprintf(e.Message, e.Args))
}
