package errs

import "errors"

var (
	ErrExistingEmail      = errors.New("email already registered")
	ErrInternalServer     = errors.New("internal server error")
	ErrAccountNotActive   = errors.New("email not activate yet")
	ErrAccountActivated   = errors.New("account already activated")
	ErrTokenExpired       = errors.New("token is expired")
	ErrInvalidOTP         = errors.New("invalid otp")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotFound      = errors.New("email not found")

	ErrInvalidEmail     = errors.New("email format is invalid")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)
