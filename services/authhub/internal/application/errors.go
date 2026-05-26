package application

import "github.com/napryag/eventflow-platform/pkg/errs"

var (
	ErrUserAlreadyExists = errs.New("user already exists")
	ErrUserNotFound      = errs.New("user not found")
)
