package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/pkg/mail"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, email, password string, created, updated time.Time) (*User, error) {
	email, err := mail.ValidateAndNormalizeEmail(email)
	if err != nil {
		return nil, errs.New("failed to validate email").Wrap(err)
	}

	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: password,
		CreatedAt:    created,
		UpdatedAt:    updated,
	}, nil
}
