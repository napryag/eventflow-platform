package config

import (
	libsconfig "github.com/napryag/eventflow-platform/libs/platform/config"
	"github.com/napryag/eventflow-platform/libs/platform/errs"
)

type Config struct {
	LogLevel int
}

func Load() (*Config, error) {
	var cfg Config

	if err := libsconfig.LoadEnv("./deployments/env/authhub.env"); err != nil {
		return nil, errs.New("failed to load env").Wrap(err)
	}

	logLevelString := libsconfig.GetString("LOG_LEVEL")
	if logLevelString == "" {
		return nil, errs.New("missing LOG_LEVEL")
	}
	logLevel, err := libsconfig.MustInt(logLevelString, "LOG_LEVEL")
	if err != nil {
		return nil, errs.New("failed variable check").Wrap(err)
	}
	switch logLevel {
	case -1, 0, 1, 2, 3:
	default:
		return nil, errs.New("unsupported LOG_LEVEL")
	}
	cfg.LogLevel = logLevel

	return &cfg, nil
}
