package errs

import "errors"

var (
	ErrExistingEmail    = errors.New("email already registered")
	ErrInternalServer   = errors.New("internal server error")
	ErrAccountNotActive = errors.New("email not activate yet")

	ErrInvalidEmail     = errors.New("email format is invalid")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)
