package mail

import (
	"net/mail"
	"strings"

	"github.com/napryag/eventflow-platform/pkg/errs"
)

func ValidateAndNormalizeEmail(email string) (string, error) {
	if email == "" {
		return "", errs.New("empty email")
	}

	if _, err := mail.ParseAddress(strings.TrimSpace(strings.ToLower(email))); err != nil {
		return "", errs.New("invalid email").Wrap(err)
	}

	return email, nil
}
