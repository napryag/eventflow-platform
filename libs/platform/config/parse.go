package config

import (
	"os"
	"strconv"

	"github.com/napryag/eventflow-platform/libs/platform/errs"
)

// ParseLogLevel gets env variable with "key" name,
// converting to int and checks if log level supported by logger.
func ParseLogLevel(key string) (int, error) {
	logLevelString, ok := os.LookupEnv(key)
	if !ok {
		return 0, errs.New("missing LOG_LEVEL")
	}
	if logLevelString == "" {
		return 0, errs.New("LOG_LEVEL is not set")
	}

	logLevel, err := strconv.Atoi(logLevelString)
	if err != nil {
		return 0, errs.New("LOG_LEVEL must be numeric").Wrap(err)
	}

	switch logLevel {
	case -1, 0, 1, 2, 3:
	default:
		return 0, errs.New("unsupported LOG_LEVEL")
	}

	return logLevel, nil
}
