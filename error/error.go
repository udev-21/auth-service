package error

import (
	"errors"
)

type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(err error) *Error {
	return &Error{
		Message: err.Error(),
	}
}

var (
	ErrUnauthorized = NewError(errors.New("unauthorized"))
	ErrNotFound     = NewError(errors.New("not found"))
	ErrNotCreated   = NewError(errors.New("not created"))
	ErrNotUpdated   = NewError(errors.New("not updated"))
	ErrNotDeleted   = NewError(errors.New("not deleted"))
	ErrInvalidID    = NewError(errors.New("invalid id"))

	ErrObjectNil    = NewError(errors.New("object is nil"))
	ErrInvalidValue = NewError(errors.New("invalid value"))

	ErrSomethingWentWrong = NewError(errors.New("something went wrong"))
	ErrEmailAlreadyExists = NewError(errors.New("email already exists"))

	ErrAtLeastOneFieldRequired = NewError(errors.New("at least one field required"))

	ErrPasswordOrEmailIncorrect = NewError(errors.New("password or email incorrect"))
	ErrPasswordConfirmNotMatch  = NewError(errors.New("password confirmation do not match"))

	ErrInvalidInput       = NewError(errors.New("invalid input"))
	ErrInvalidPassword    = NewError(errors.New("invalid password"))
	ErrInvalidEmail       = NewError(errors.New("invalid email"))
	ErrInvalidContentType = NewError(errors.New("invalid content type"))

	ErrExpiredToken     = NewError(errors.New("expired token"))
	ErrInvalidToken     = NewError(errors.New("invalid token"))
	ErrAuthTokenExpired = NewError(errors.New("token has expired"))
)
