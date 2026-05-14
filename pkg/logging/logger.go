package logging

type Logger interface {
	Trace() LogEvent
	Debug() LogEvent
	Info() LogEvent
	Warn() LogEvent
	Error() LogEvent
	Err(err error) LogEvent
}

type LogEvent interface {
	Msg(msg string)
	Str(key, value string) LogEvent
}
