package mail

import (
	"net/mail"
	"strings"

	"github.com/napryag/eventflow-platform/pkg/errs"
)

func ValidateAndNormalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return "", errs.New("empty email")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return "", errs.New("invalid email")
	}

	return email, nil
}
