package query

import (
	"errors"
)

var (
	ErrInvalidCondition = errors.New("invalid condition")
	ErrInvalidDialect   = errors.New("invalid dialect")
)
