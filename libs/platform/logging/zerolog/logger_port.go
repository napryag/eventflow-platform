package zerolog

type Logger interface {
	Info() LogEvent
	Error() LogEvent
}

type LogEvent interface {
	Msg(msg string)
	Str(key, value string) LogEvent
}
