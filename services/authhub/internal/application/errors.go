package application

import "github.com/napryag/eventflow-platform/pkg/errs"

var (
	ErrPasswordTooShort          = errs.New("password too short")
	ErrInvalidCredentials        = errs.New("invalid email or password")
	ErrMismatchedHashAndPassword = errs.New("mismatched hash and password")
	ErrHashTooShort              = errs.New("hash too short")
	ErrPasswordTooLong           = errs.New("password too long")
	ErrEmptyPassword             = errs.New("empty password")
	ErrUserAlreadyExists         = errs.New("user already exists")
	ErrUserNotFound              = errs.New("user not found")
	ErrEmptyEmail                = errs.New("empty email")
	ErrInvalidEmail              = errs.New("invalid email")
)
