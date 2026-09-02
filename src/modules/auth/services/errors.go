package services

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already registered")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrInvalidOTP         = errors.New("invalid or expired code")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrRegistrationClosed = errors.New("registration is disabled")
	ErrWeakPassword       = errors.New("password does not meet requirements")
	ErrAlreadyVerified    = errors.New("email already verified")
	ErrInvalidPicture     = errors.New("invalid picture")
)
