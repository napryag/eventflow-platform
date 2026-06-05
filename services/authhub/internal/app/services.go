package app

import (
	"context"

	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/pkg/logging"
	"github.com/napryag/eventflow-platform/pkg/postgres"
	"github.com/napryag/eventflow-platform/services/authhub/internal/application"
	"github.com/napryag/eventflow-platform/services/authhub/internal/infrastructure/bcrypt"
	"github.com/napryag/eventflow-platform/services/authhub/internal/infrastructure/postgres/user"
)

type Services struct {
	AuthUseCase *application.AuthUseCase
}

func BuildServices(ctx context.Context, cfg postgres.DatabaseConfig, logger logging.Logger) (*Services, error) {
	db, err := postgres.NewConnection(ctx, cfg, logger)
	if err != nil {
		return nil, errs.New("failed to create new connection to db").Wrap(err)
	}

	users := user.NewRepository(db)

	hasher := bcrypt.NewPasswordHasher()

	return &Services{AuthUseCase: application.NewAuthUseCase(users, hasher, logger)}, nil
}
