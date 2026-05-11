package config

import (
	"github.com/joho/godotenv"
	libsconfig "github.com/napryag/eventflow-platform/libs/platform/config"
	"github.com/napryag/eventflow-platform/libs/platform/errs"
)

type Config struct {
	LogLevel int
}

func Load() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, errs.New("failed to load env").Wrap(err)
	}

	logLevel, err := libsconfig.ParseLogLevel("LOG_LEVEL")
	if err != nil {
		return nil, errs.New("failed to parse LOG_LEVEL").Wrap(err)
	}

	cfg.LogLevel = logLevel

	return &cfg, nil
}
