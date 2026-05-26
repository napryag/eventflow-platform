package user

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, email, password string, created, updated time.Time) (*User, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}

	email = strings.ToLower(email)

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
