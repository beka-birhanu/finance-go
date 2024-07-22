package domain_errors

import "errors"

var (
	ErrUsernameConflict      = errors.New("username already taken")
	ErrUsernameTooShort      = errors.New("username is too short")
	ErrUsernameTooLong       = errors.New("username is too long")
	ErrWeakPassword          = errors.New("password is too weak")
	ErrUsernameInvalidFormat = errors.New("username has an invalid format")
)

