package main

import (
	"context"
	"os"

	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/pkg/logging/zerolog"
	"github.com/napryag/eventflow-platform/services/authhub/config"
	"github.com/napryag/eventflow-platform/services/authhub/internal/app"
	"github.com/napryag/eventflow-platform/services/authhub/internal/interfaces/http"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		err = errs.New("failed to load config").Wrap(err)
		log.Err(err).Send()
		os.Exit(1)
	}

	logger := zerolog.New(cfg.LogLevel)

	router := http.New()

	svc, err := app.BuildServices(ctx, cfg.Database, logger)
	if err != nil {
		logger.Err(err).Msg("failed to build application services")
		os.Exit(1)
	}

	_ = svc // will be passed later

	logger.Info().Str("service", "authhub").Msg("api initialized")

	if err := router.Run(cfg.HTTP); err != nil {
		logger.Err(err).Msg("failed to start server")
		os.Exit(1)
	}
}
