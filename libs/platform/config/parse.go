package config

import (
	"strconv"

	"github.com/napryag/eventflow-platform/libs/platform/errs"
)

// Checks if ENV varibable with name "field" is integer.
//
//	func MustInt("1", "LOG_LEVEL")
func MustInt(value, field string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, errs.New(field + " must be numeric").Wrap(err)
	}

	return n, nil
}
