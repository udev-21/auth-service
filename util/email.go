package util

import (
	"net/mail"

	myErrors "udev21/auth/error"
)

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return myErrors.ErrInvalidEmail
	}
	return nil
}
