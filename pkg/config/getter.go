package config

import (
	"os"
	"strconv"

	"github.com/napryag/eventflow-platform/pkg/errs"
)

func GetLogLevel() (int, error) {
	value, err := GetEnvInt("LOG_LEVEL")
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

func GetEnvInt(key string) (int, error) {
	valueString, err := GetEnvString(key)
	if err != nil {
		return 0, errs.New("failed to get string").Wrap(err)
	}

	value, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, errs.New("value must be numeric").Arg("key", key).Wrap(err)
	}

	return value, nil
}

func GetEnvString(key string) (string, error) {
	valueString, ok := os.LookupEnv(key)
	if !ok || valueString == "" {
		return "", errs.New("key is not set").Arg("key", key)
	}
	return valueString, nil
}

func GetHTTPPort() (string, error) {
	valueString, err := GetEnvString("HTTP_PORT")
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
