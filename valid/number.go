package valid

import (
	"errors"
	"reflect"
)

var (
	ErrNumberOutOfRange = errors.New("number is out of range")
)

func Range(min interface{}, max interface{}) Constraint {
	return func(value interface{}) error {
		valid := true

		if reflect.TypeOf(min) != reflect.TypeOf(max) || reflect.TypeOf(min) != reflect.TypeOf(value) {
			return ErrIncompatibleTypes
		}

		switch v := value.(type) {
		case int:
			valid = v >= min.(int) && v <= max.(int)
		case int8:
			valid = v >= min.(int8) && v <= max.(int8)
		case int16:
			valid = v >= min.(int16) && v <= max.(int16)
		case int32:
			valid = v >= min.(int32) && v <= max.(int32)
		case int64:
			valid = v >= min.(int64) && v <= max.(int64)
		case uint8:
			valid = v >= min.(uint8) && v <= max.(uint8)
		case uint16:
			valid = v >= min.(uint16) && v <= max.(uint16)
		case uint32:
			valid = v >= min.(uint32) && v <= max.(uint32)
		case uint64:
			valid = v >= min.(uint64) && v <= max.(uint64)
		case float32:
			valid = v >= min.(float32) && v <= max.(float32)
		case float64:
			valid = v >= min.(float64) && v <= max.(float64)
		}

		if !valid {
			return ErrNumberOutOfRange
		}
		return nil
	}
}
