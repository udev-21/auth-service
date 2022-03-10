package util

import "unicode"

func ValidatePasswordStrength(s string) bool {
	eightOrMore := len(s) >= 8
	number := false
	upper := false
	special := false
	lower := false

	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		default:
			return false
		}
	}

	return eightOrMore && lower && number && upper && special
}
