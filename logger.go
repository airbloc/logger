package logger

import (
	"fmt"
	"os"
	"runtime/debug"
)

// Logger is the unit of the logger package, a smart, pretty-printing gate between
// the program and the output stream.
type Logger interface {
	Debug(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Error(msg string, v ...interface{})

	Timer() *Log

	Fatal(v ...interface{})
	Wtf(v ...interface{})

	// WithAttrs returns a sub-logger with given attributes attached as a default.
	WithAttrs(attrs Attrs) Logger
}

// New returns a logger bound to the given name.
func New(name string) Logger {
	return &logger{
		Name: name,
	}
}

type logger struct {
	// Name by which the logger is identified when enabling or disabling it, and by envvar.
	Name string
}

func (l *logger) Log(level *LogLevel, message string, args []interface{}) {
	attrs := MergeAttrs(args)
	formatted, purgedAttrs := Format(message, *attrs)

	runtime.Log(&Log{
		Package: l.Name,
		Level:   level,
		Message: formatted,
		Time:    Now(),
		Attrs:   attrs,

		DisplayedAttrs: &purgedAttrs,
	})
}

// Error logs an debug message.
func (l *logger) Debug(msg string, v ...interface{}) {
	l.Log(Debug, msg, v)
}

// Info prints log information to the screen that is informational in nature.
func (l *logger) Info(msg string, v ...interface{}) {
	l.Log(Info, msg, v)
}

// Error logs an warning message.
func (l *logger) Warn(msg string, v ...interface{}) {
	l.Log(Warn, msg, v)
}

// Error logs an error message.
// If error has been given as a first argument, the error will be logged also.
func (l *logger) Error(msg string, v ...interface{}) {
	if len(v) > 0 {
		if err, hasErr := v[0].(error); hasErr {
			msg = fmt.Sprintf("%s: %v", msg, err)
			v = v[1:]
		}
	}
	l.Log(Error, msg, v)
}

// Wtf logs error detailed, and reports error to error transport.
// Error message (format) is optional, so you can call the method just like `Wtf(error)`
// TODO: add transports for errors reported with WTF level
func (l *logger) Wtf(v ...interface{}) {
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
	l.Log(Fatal, Colored(Red, msg), v)
}

// Fatal behaves same as Wtf, but it exits process with code 1
func (l *logger) Fatal(v ...interface{}) {
	l.Wtf(v...)
	os.Exit(1)
}

func (l *logger) Recover(context Attrs) interface{} {
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
func (l *logger) Timer() *Log {
	return &Log{
		Package: l.Name,
		Level:   Timer,
		Time:    Now(),
	}
}

func (l *logger) WithAttrs(attr Attrs) Logger {
	return &subLogger{
		parent:       l,
		defaultAttrs: attr,
	}
}
