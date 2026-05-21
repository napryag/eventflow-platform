package postgres

import (
	"context"
	"time"

	"github.com/napryag/eventflow-platform/pkg/logging"
	"github.com/napryag/eventflow-platform/services/authhub/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(
	ctx context.Context,
	cfg config.DatabaseConfig,
	logger logging.Logger,
) (*gorm.DB, error) {
	dial := postgres.Open(cfg.DSN())

	db, err := gorm.Open(dial)
	if err != nil {
		logger.Err(err).Msg("failed to initialize db session")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Err(err).Msg("failed to get sql db")
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		logger.Err(err).Msg("failed to ping db")
		return nil, err
	}

	logger.Info().Msg("db initialized successfully")

	return db, nil
}
