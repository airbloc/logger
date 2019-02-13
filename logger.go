package logger

import (
	"fmt"
	"runtime/debug"
)

// New returns a logger bound to the given name.
func New(name string) *Logger {
	return &Logger{
		Name: name,
	}
}

// Logger is the unit of the logger package, a smart, pretty-printing gate between
// the program and the output stream.
type Logger struct {
	// Name by which the logger is identified when enabling or disabling it, and by envvar.
	Name string
}

func (logger *Logger) Log(level, message string, args []interface{}) {
	attrs := MergeAttrs(args)
	formatted, purgedAttrs := Format(message, *attrs)

	runtime.Log(&Log{
		Package: logger.Name,
		Level:   level,
		Message: formatted,
		Time:    Now(),
		Attrs:   attrs,

		DisplayedAttrs: &purgedAttrs,
	})
}

// Info prints log information to the screen that is informational in nature.
func (l *Logger) Info(msg string, v ...interface{}) {
	l.Log("INFO", msg, v)
}

// Error logs an debug message.
func (l *Logger) Debug(msg string, v ...interface{}) {
	l.Log("DEBUG", msg, v)
}

// Error logs an error message.
// If error has been given as a first argument, the error will be logged also.
func (l *Logger) Error(msg string, v ...interface{}) {
	if len(v) > 0 {
		if err, hasErr := v[0].(error); hasErr {
			msg = fmt.Sprintf("%s: %v", msg, err)
			v = v[1:]
		}
	}
	l.Log("ERROR", msg, v)
}

// Wtf logs error detailed, and reports error to error transport.
// Error message (format) is optional, so you can call the method just like `Wtf(error)`
// TODO: add transports for errors reported with WTF level
func (l *Logger) Wtf(v ...interface{}) {
	msg := ""
	if m, ok := v[0].(string); ok {
		msg = m
		v = v[1:]
	}
	if len(v) > 0 {
		if err, hasErr := v[0].(error); hasErr {
			if msg != "" {
				msg = msg + ": "
			}
			msg = msg + fmt.Sprintf("%+v", err)
			v = v[1:]
		}
	}
	l.Log("FATAL", Colored(Red, msg), v)
}

func (l *Logger) Recover(context Attrs) interface{} {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			r = err.Error()
		}
		l.Wtf("panic: {}\n{}", r, Colored(dim, string(debug.Stack())), context)
		return r
	}
	return nil
}

// Timer returns a timer sub-logger.
func (l *Logger) Timer() *Log {
	return &Log{
		Package: l.Name,
		Level:   "TIMER",
		Time:    Now(),
	}
}
