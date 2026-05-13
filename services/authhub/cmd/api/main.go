package main

import (
	"os"

	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/pkg/logging/zerolog"
	"github.com/napryag/eventflow-platform/services/authhub/config"
	http "github.com/napryag/eventflow-platform/services/authhub/internal/interfaces"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		err = errs.New("failed to load config").Wrap(err)
		log.Err(err).Send()
		os.Exit(1)
	}

	logger := zerolog.New(cfg.LogLevel)

	logger.Info().Str("service", "authhub").Msg("api initialized")

	router := http.SetupRouter()
	router.GET("/health", http.HealthHandler)

	if err := router.Run(cfg.HTTP.Host + ":" + cfg.HTTP.Port); err != nil {
		logger.Error().Msg("failed to start server")
	}
}
