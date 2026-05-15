package config

import (
	"github.com/joho/godotenv"
	"github.com/napryag/eventflow-platform/pkg/config"
	"github.com/napryag/eventflow-platform/pkg/errs"
)

type HTTPConfig struct {
	Host string
	Port string
}
type Config struct {
	LogLevel int
	HTTP     HTTPConfig
}

func Load() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, errs.New("failed to load env").Wrap(err)
	}

	host, err := config.GetEnvString("HTTP_HOST")
	if err != nil {
		return nil, errs.New("failed to get host").Wrap(err)
	}

	cfg.HTTP.Host = host

	port, err := config.GetHTTPPort()
	if err != nil {
		return nil, errs.New("failed to get http port").Wrap(err)
	}

	cfg.HTTP.Port = port

	logLevel, err := config.GetLogLevel()
	if err != nil {
		return nil, errs.New("failed to get log level").Wrap(err)
	}

	cfg.LogLevel = logLevel

	return &cfg, nil
}
