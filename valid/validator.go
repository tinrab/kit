package valid

import (
	"errors"
)

var (
	ErrIncompatibleTypes = errors.New("incompatible types")
)

type Constraint func(value interface{}) error

type Validator struct {
	constraints []Constraint
}

func New(constraints ...Constraint) *Validator {
	return &Validator{
		constraints: constraints,
	}
}

func (v Validator) Validate(value interface{}) error {
	for _, c := range v.constraints {
		if err := c(value); err != nil {
			return err
		}
	}
	return nil
}
