package loggergrpc

import "github.com/airbloc/logger"

type options struct {
	requestLogLevel *logger.LogLevel
	errorLogLevel   *logger.LogLevel
}

type Option func(o *options)

// WithRequestLogLevel sets level on logging successful request logs.
func WithRequestLogLevel(level *logger.LogLevel) Option {
	return func(o *options) {
		o.requestLogLevel = level
	}
}

// WithErrorLogLevel sets level on logging successful request logs.
func WithErrorLogLevel(level *logger.LogLevel) Option {
	return func(o *options) {
		o.errorLogLevel = level
	}
}

func createOptions(opts []Option) *options {
	opt := &options{
		requestLogLevel: logger.Verbose,
		errorLogLevel: logger.Error,
	}
	for _, fn := range opts {
		fn(opt)
	}
	return opt
}
