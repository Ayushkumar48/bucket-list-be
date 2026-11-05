package utils

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

func ValidatePassword(p string) error {
	if utf8.RuneCountInString(p) < 12 {
		return errors.New("password must be at least 12 characters")
	}

	var hasLetter, hasDigit, hasSpecial bool
	for _, r := range p {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasLetter || !hasDigit || !hasSpecial {
		return errors.New("password must contain a letter, a number, and a special character")
	}
	return nil
}
