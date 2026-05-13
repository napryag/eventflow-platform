package errs

import (
	"fmt"
	"runtime"
	"strings"
)

// Error represents an application error with optional
// wrapped cause, structured arguments, and stack trace.
type Error struct {
	msg   string
	cause error
	args  map[string]any
	stack stack
}

// New creates a new Error with the provided message
// and captures the current call stack.
func New(msg string) *Error {
	return &Error{
		msg:   msg,
		stack: callers(3),
	}
}

// Wrap attaches the given error as a cause.
//
// If err is nil, Wrap returns nil.
func (e *Error) Wrap(err error) *Error {
	if err == nil {
		return nil
	}

	if e == nil {
		return &Error{
			cause: err,
			stack: callers(3),
		}
	}

	e.cause = err
	if len(e.stack.pcs) == 0 {
		e.stack = callers(3)
	}
	return e
}

// Arg adds a structured argument to the error.
//
// Arguments are included into the formatted error message.
func (e *Error) Arg(key string, value any) *Error {
	if e == nil {
		return nil
	}
	if e.args == nil {
		e.args = make(map[string]any)
	}
	e.args[key] = value
	return e
}

// Error returns the formatted error message.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	msg := strings.TrimSpace(e.msg)

	switch {
	case msg != "" && e.cause != nil:
		msg = msg + ": " + e.cause.Error()
	case msg == "" && e.cause != nil:
		msg = e.cause.Error()
	case msg == "":
		msg = "error"
	}

	if len(e.args) > 0 {
		msg += " | " + formatArgs(e.args)
	}
	return msg
}

// String returns the string representation of the error.
func (e *Error) String() string {
	return e.Error()
}

// Unwrap returns the wrapped error cause.
func (e *Error) Unwrap() error {
	return e.cause
}

// Format implements fmt.Formatter.
//
// Supported verbs:
//   - %s prints the error message
//   - %q prints the quoted error message
//   - %+v prints the error message with stack trace
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			write(s, e.Error())
			if len(e.stack.pcs) > 0 {
				write(s, "\nstack:\n")
				write(s, e.stack.String())
			}
			return
		}
		fallthrough
	case 's':
		write(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func formatArgs(args map[string]any) string {
	var b strings.Builder
	first := true
	for k, v := range args {
		if !first {
			b.WriteString(", ")
		}
		first = false
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(fmt.Sprint(v))
	}
	return b.String()
}

func write(w interface{ Write([]byte) (int, error) }, s string) {
	_, _ = w.Write([]byte(s))
}

// stack contains program counters collected
// from the current call stack.
type stack struct {
	pcs []uintptr
}

// callers collects the current call stack and returns it
// as a stack structure.
//
// The skip parameter specifies how many initial stack frames
// should be skipped (for example, callers itself and wrapper functions).
//
// Function program counters are obtained using runtime.Callers.
// The maximum captured stack depth is limited by depth.
func callers(skip int) stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return stack{pcs: pcs[:n]}
}

// String returns formated string with function/file/line information
// for each stack frame.
//
// strings.Builder used to minimize memory copying.
func (s stack) String() string {
	var b strings.Builder
	frames := runtime.CallersFrames(s.pcs)

	for {
		f, more := frames.Next()
		b.WriteString(f.Function)
		b.WriteString("\n\t")
		b.WriteString(f.File)
		b.WriteString(fmt.Sprintf(":%d\n", f.Line))
		if !more {
			break
		}
	}
	return b.String()
}
