package config

import (
	"os"
	"strconv"

	"github.com/napryag/eventflow-platform/pkg/errs"
)

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
