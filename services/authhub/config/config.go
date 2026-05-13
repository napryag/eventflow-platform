package config

import (
	"github.com/joho/godotenv"
	"github.com/napryag/eventflow-platform/pkg/config"
	"github.com/napryag/eventflow-platform/pkg/errs"
)

type Config struct {
	LogLevel int
}

func Load() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, errs.New("failed to load env").Wrap(err)
	}

	logLevel, err := config.GetLogLevel()
	if err != nil {
		return nil, errs.New("failed to get log level").Wrap(err)
	}

	cfg.LogLevel = logLevel

	return &cfg, nil
}
