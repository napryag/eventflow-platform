package config

import (
	"strconv"

	"github.com/joho/godotenv"
	"github.com/napryag/eventflow-platform/pkg/config"
	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/pkg/postgres"
	"github.com/rs/zerolog/log"
)

type HTTPConfig struct {
	Host string
	Port string
}
type Config struct {
	LogLevel int
	HTTP     HTTPConfig
	Database postgres.DatabaseConfig
}

func Load() (*Config, error) {
	var cfg Config
	var err error

	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("missing env file")
	}

	cfg.LogLevel, err = GetLogLevel()
	if err != nil {
		return nil, errs.New("failed to get log level").Wrap(err)
	}

	cfg.HTTP, err = loadHTTPConfig()
	if err != nil {
		return nil, errs.New("failed to load http config").Wrap(err)
	}

	cfg.Database, err = loadDatabaseConfig()
	if err != nil {
		return nil, errs.New("failed to load database config").Wrap(err)
	}

	return &cfg, nil
}

func GetLogLevel() (int, error) {
	value, err := config.GetEnvInt("LOG_LEVEL")
	if err != nil {
		return 0, err
	}

	switch value {
	case -1, 0, 1, 2, 3:
		return value, nil
	default:
		return 0, errs.New("unsupported LOG_LEVEL").Arg("value", value)
	}
}

func GetHTTPPort() (string, error) {
	valueString, err := config.GetEnvString("HTTP_PORT")
	if err != nil {
		return "", errs.New("failed to get port").Wrap(err)
	}

	value, err := strconv.Atoi(valueString)
	if err != nil {
		return "", errs.New("value must be numeric").Wrap(err)
	}

	if value < 1 || value > 65535 {
		return "", errs.New("invalid value")
	}

	return valueString, nil
}

func loadHTTPConfig() (HTTPConfig, error) {
	var cfg HTTPConfig
	var err error

	cfg.Host, err = config.GetEnvString("HTTP_HOST")
	if err != nil {
		return cfg, errs.New("failed to get host").Wrap(err)
	}

	cfg.Port, err = GetHTTPPort()
	if err != nil {
		return cfg, errs.New("failed to get http port").Wrap(err)
	}

	return cfg, nil
}

func loadDatabaseConfig() (postgres.DatabaseConfig, error) {
	var cfg postgres.DatabaseConfig
	var err error

	cfg.Host, err = config.GetEnvString("POSTGRES_HOST")
	if err != nil {
		return cfg, errs.New("failed to get postgres host").Wrap(err)
	}

	cfg.Port, err = getPostgresPort()
	if err != nil {
		return cfg, errs.New("failed to get postgres port").Wrap(err)
	}

	cfg.Name, err = config.GetEnvString("POSTGRES_DB")
	if err != nil {
		return cfg, errs.New("failed to get postgres db name").Wrap(err)
	}

	cfg.User, err = config.GetEnvString("POSTGRES_USER")
	if err != nil {
		return cfg, errs.New("failed to get postgres user").Wrap(err)
	}

	cfg.Password, err = config.GetEnvString("POSTGRES_PASSWORD")
	if err != nil {
		return cfg, errs.New("failed to get postgres password").Wrap(err)
	}

	cfg.SSLMode, err = getPostgresSSLMode()
	if err != nil {
		return cfg, errs.New("failed to get postgres ssl mode").Wrap(err)
	}

	return cfg, nil
}

func getPostgresPort() (string, error) {
	port, err := config.GetEnvInt("POSTGRES_PORT")
	if err != nil {
		return "", errs.New("invalid POSTGRES_PORT").Wrap(err)
	}

	if port < 1 || port > 65535 {
		return "", errs.New("POSTGRES_PORT out of range").Arg("value", port)
	}

	return strconv.Itoa(port), nil
}

func getPostgresSSLMode() (string, error) {
	mode, err := config.GetEnvString("POSTGRES_SSLMODE")
	if err != nil {
		return "", errs.New("failed to get POSTGRES_SSLMODE").Wrap(err)
	}

	switch mode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
		return mode, nil
	default:
		return "", errs.New("unsupported POSTGRES_SSLMODE").Arg("value", mode)
	}
}
