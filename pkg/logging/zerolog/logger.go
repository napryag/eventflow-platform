package zerolog

import (
	"os"

	"github.com/napryag/eventflow-platform/pkg/logging"
	"github.com/rs/zerolog"
)

type zeroLogger struct {
	logger zerolog.Logger
}

type zeroEvent struct {
	event *zerolog.Event
}

func New(logLevel int) logging.Logger {
	zerolog.TimeFieldFormat = "02-01-2006 15:04:05"

	return &zeroLogger{
		logger: zerolog.New(os.Stdout).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Logger(),
	}
}

func (zl *zeroLogger) Trace() logging.LogEvent {
	return &zeroEvent{
		event: zl.logger.Trace(),
	}
}

func (zl *zeroLogger) Debug() logging.LogEvent {
	return &zeroEvent{
		event: zl.logger.Debug(),
	}
}

func (zl *zeroLogger) Info() logging.LogEvent {
	return &zeroEvent{
		event: zl.logger.Info(),
	}
}

func (zl *zeroLogger) Warn() logging.LogEvent {
	return &zeroEvent{
		event: zl.logger.Warn(),
	}
}

func (zl *zeroLogger) Error() logging.LogEvent {
	return &zeroEvent{
		event: zl.logger.Error(),
	}
}

func (ze *zeroEvent) Msg(msg string) {
	ze.event.Msg(msg)
}

func (ze *zeroEvent) Str(key, value string) logging.LogEvent {
	ze.event.Str(key, value)
	return ze
}
