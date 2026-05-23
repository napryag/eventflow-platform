package postgres

import (
	"context"
	"time"

	"github.com/napryag/eventflow-platform/pkg/config"
	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/pkg/logging"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPingTimeout = 5 * time.Second

func NewConnection(
	ctx context.Context,
	cfg config.DatabaseConfig,
	logger logging.Logger,
) (*gorm.DB, error) {
	if err := cfg.Validate(); err != nil {
		err = errs.New("failed to validate config").Wrap(err)
		logger.Err(err).Send()

		return nil, err
	}

	db, err := gorm.Open(gormpostgres.Open(cfg.DSN()), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		err = errs.New("failed to initialized db session based on dialector").Wrap(err)
		logger.Err(err).Send()

		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		err = errs.New("failed to get postgres sql db").Wrap(err)
		logger.Err(err).Send()

		return nil, err
	}

	config.ConfigurePool(sqlDB, cfg)

	timeout := cfg.PingTimeout
	if timeout == 0 {
		timeout = defaultPingTimeout
	}

	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		err = errs.New("failed to ping postgres").
			Arg("host", cfg.Host).
			Arg("port", cfg.Port).
			Arg("database", cfg.Name).
			Wrap(err)
		logger.Err(err).Send()
		return nil, err
	}

	logger.Info().
		Str("host", cfg.Host).
		Str("port", cfg.Port).
		Str("database", cfg.Name).
		Msg("postgres connection initialized")

	return db, nil
}
