package errs

import "errors"

var (
<<<<<<< HEAD
	ErrExistingEmail    = errors.New("email already registered")
	ErrInternalServer   = errors.New("internal server error")
	ErrAccountNotActive = errors.New("email not activate yet")
	ErrAccountActivated = errors.New("account already activated")
	ErrTokenExpired     = errors.New("token is expired")
	ErrInvalidOTP       = errors.New("invalid otp")
=======
	ErrExistingEmail  = errors.New("email already registered")
	ErrInternalServer = errors.New("internal server error")
>>>>>>> b9ee6f3b7daa7e17199dec072791cf7dbe5d369b

	ErrInvalidEmail     = errors.New("email format is invalid")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)
