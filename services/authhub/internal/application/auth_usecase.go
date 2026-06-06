package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/pkg/logging"
	"github.com/napryag/eventflow-platform/pkg/mail"
	"github.com/napryag/eventflow-platform/services/authhub/internal/domain/user"
)

type AuthUseCase struct {
	users  user.Repository
	hasher PasswordHasher
	logger logging.Logger
}

func NewAuthUseCase(
	users user.Repository,
	hasher PasswordHasher,
	logger logging.Logger,
) *AuthUseCase {
	return &AuthUseCase{
		users:  users,
		hasher: hasher,
		logger: logger,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, password string) (*user.User, error) {
	normalizedEmail, err := mail.ValidateAndNormalizeEmail(email)
	if err != nil {
		uc.logger.Err(err).Str("email", email).Msg("failed to validate email")

		return nil, err
	}

	hashedPassword, err := uc.hasher.HashPassword(password)
	if err != nil {
		uc.logger.Err(err).Msg("failed to hash password")

		return nil, err
	}

	user, err := user.NewUser(uuid.New(), normalizedEmail, hashedPassword, time.Now(), time.Now())
	if err != nil {
		uc.logger.Err(err).Msg("failed to create user domain object")

		return nil, err
	}

	if err := uc.users.Create(ctx, *user); err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			uc.logger.Warn().Str("email", normalizedEmail).Msg("user already exists")

			return nil, ErrUserAlreadyExists
		}
		uc.logger.Err(err).Msg("failed to create user in database")

		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, email, password string) (*user.User, error) {
	normalizedEmail, err := mail.ValidateAndNormalizeEmail(email)
	if err != nil {
		uc.logger.Err(err).Str("email", email).Msg("failed to validate email")

		return nil, err
	}

	user, err := uc.users.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			uc.logger.Warn().Str("email", normalizedEmail).Msg("user not found")

			return nil, ErrInvalidCredentials
		}
		uc.logger.Err(err).Str("email", normalizedEmail).Msg("failed to get user by email")

		return nil, err
	}

	if err := uc.hasher.ComparePassword(user.PasswordHash, password); err != nil {
		uc.logger.Err(err).Msg("failed to compare password")

		return nil, ErrInvalidCredentials
	}

	return user, nil
}
