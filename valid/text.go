package valid

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	emailUserRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	emailHostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	slugRegexp      = regexp.MustCompile("^[a-z0-9]+(?:-[a-z0-9]+)*$")
)

var (
	ErrInvalidTextPattern = errors.New("invalid text pattern")
	ErrTextTooLong        = errors.New("text is too long")
	ErrTextTooShort       = errors.New("text is too short")
	ErrTextMissingUpper   = errors.New("missing uppercase character")
	ErrTextMissingLower   = errors.New("missing lowercase character")
	ErrTextMissingSpecial = errors.New("missing special character")
	ErrTextMissingNumber  = errors.New("missing number")

	ErrInvalidEmail    = errors.New("invalid email")
	ErrTextInvalidSlug = errors.New("invalid slug")
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

func MinMaxLength(min, max int) Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}

		if len(v) < min {
			return ErrTextTooShort
		}

		if len(v) > max {
			return ErrTextTooLong
		}

		return nil
	}
}

func HasUpperCase() Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}
		for _, r := range []rune(v) {
			if unicode.IsUpper(r) {
				return nil
			}
		}
		return ErrTextMissingUpper
	}
}

func HasLowerCase() Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}
		for _, r := range []rune(v) {
			if unicode.IsLower(r) {
				return nil
			}
		}
		return ErrTextMissingLower
	}
}

func HasSpecial() Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}
		for _, r := range []rune(v) {
			if unicode.IsPunct(r) || unicode.IsSymbol(r) {
				return nil
			}
		}
		return ErrTextMissingSpecial
	}
}

func HasNumber() Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}
		for _, r := range []rune(v) {
			if unicode.IsNumber(r) {
				return nil
			}
		}
		return ErrTextMissingNumber
	}
}

func Email() Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}

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

func Slug() Constraint {
	return func(value interface{}) error {
		v, ok := value.(string)
		if !ok {
			return ErrIncompatibleTypes
		}

		if !slugRegexp.MatchString(v) {
			return ErrTextInvalidSlug
		}

		return nil
	}
}
