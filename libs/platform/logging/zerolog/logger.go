package zerolog

import (
	"os"

	"github.com/rs/zerolog"
)

type zeroLogger struct {
	logger zerolog.Logger
}

type zeroEvent struct {
	event *zerolog.Event
}

func New(logLevel int) Logger {
	zerolog.TimeFieldFormat = "02-01-2006 15:04:05"

	return &zeroLogger{
		logger: zerolog.New(os.Stdout).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Logger(),
	}
}

func (zl *zeroLogger) Info() LogEvent {
	return &zeroEvent{
		event: zl.logger.Info(),
	}
}

func (zl *zeroLogger) Error() LogEvent {
	return &zeroEvent{
		event: zl.logger.Error(),
	}
}

func (ze *zeroEvent) Msg(msg string) {
	ze.event.Msg(msg)
}

func (ze *zeroEvent) Str(key, value string) LogEvent {
	ze.event.Str(key, value)
	return ze
}
