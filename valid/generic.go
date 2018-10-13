package valid

import (
	"errors"
)

var (
	ErrNotEqual = errors.New("values are not equal")
	ErrNotNil   = errors.New("value is not nil")
	ErrNil      = errors.New("value is nil")
)

func Equal(value interface{}) Constraint {
	return func(v interface{}) error {
		if v != value {
			return ErrNotEqual
		}
		return nil
	}
}

func Nil() Constraint {
	return func(value interface{}) error {
		if value != nil {
			return ErrNotNil
		}
		return nil
	}
}

func NotNil() Constraint {
	return func(value interface{}) error {
		if value == nil {
			return ErrNil
		}
		return nil
	}
}
