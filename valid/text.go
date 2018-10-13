package valid

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailUserRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	emailHostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
)

var (
	ErrInvalidTextPattern = errors.New("invalid text pattern")
	ErrTextTooLong        = errors.New("text is too long")
	ErrTextTooShort       = errors.New("text is too short")
	ErrInvalidEmail       = errors.New("invalid email")
)

func Regexp(s string) Constraint {
	exp := regexp.MustCompile(s)
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}
		if !exp.MatchString(v) {
			return ErrInvalidTextPattern
		}
		return nil
	}
}

func MaxLength(n int) Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}

		if len(v) > n {
			return ErrTextTooLong
		}

		return nil
	}
}

func MinLength(n int) Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}

		if len(v) < n {
			return ErrTextTooShort
		}

		return nil
	}
}

func Email() Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}

		v = strings.TrimSpace(v)

		if len(v) < 6 || len(v) > 254 {
			return ErrInvalidEmail
		}

		at := strings.LastIndex(v, "@")
		if at <= 0 || at > len(v)-3 {
			return ErrInvalidEmail
		}

		user := v[:at]
		host := v[at+1:]

		if len(user) > 64 {
			return ErrInvalidEmail
		}

		if !emailUserRegexp.MatchString(user) || !emailHostRegexp.MatchString(host) {
			return ErrInvalidEmail
		}

		return nil
	}
}
