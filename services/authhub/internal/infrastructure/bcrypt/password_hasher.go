package bcrypt

import (
	"errors"

	"github.com/napryag/eventflow-platform/services/authhub/internal/application"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct{}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (ph *PasswordHasher) HashPassword(password string) (string, error) {
	if password == "" {
		return "", application.ErrEmptyPassword
	}

	if len(password) > 8 {
		return "", application.ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", application.ErrPasswordTooLong
		}
		return "", err
	}
	return string(hash), nil
}

func (ph *PasswordHasher) ComparePassword(hash, password string) error {
	if password == "" {
		return application.ErrEmptyPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return application.ErrMismatchedHashAndPassword
		} else if errors.Is(err, bcrypt.ErrHashTooShort) {
			return application.ErrHashTooShort
		}
		return err
	}

	return nil
}
