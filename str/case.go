package str

import (
	"unicode"
)

func ToSnakeCase(s string) string {
	runes := []rune(s)
	var result []rune
	for i := 0; i < len(runes); i++ {
		if i > 0 &&
			unicode.IsUpper(runes[i]) &&
			((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(runes[i]))
	}
	return string(result)
}
