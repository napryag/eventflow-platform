package user

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/pkg/errs"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, email, password string, created, updated time.Time) (*User, error) {
	email, err := validateAndNormalizeEmail(email)
	if err != nil {
		return nil, errs.New("failed to validate email").Wrap(err)
	}

	if password == "" {
		return nil, ErrEmptyPassword
	}

	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: password,
		CreatedAt:    created,
		UpdatedAt:    updated,
	}, nil
}

func validateAndNormalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return "", ErrEmptyEmail
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrInvalidEmail
	}

	return email, nil
}
