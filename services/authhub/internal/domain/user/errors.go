package user

import "github.com/napryag/eventflow-platform/pkg/errs"

var (
	ErrEmptyEmail    = errs.New("empty email")
	ErrInvalidEmail  = errs.New("invalid email")
	ErrEmptyPassword = errs.New("empty password hash")
)
