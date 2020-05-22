package loggergin

type options struct {
	disableColor  bool
	disable404Log bool
}

type Option func(o *options)

func WithDisable404Logging() Option {
	return func(o *options) {
		o.disable404Log = true
	}
}

func WithDisableColors() Option {
	return func(o *options) {
		o.disableColor = true
	}
}

func buildOptions(opts []Option) *options {
	o := new(options)
	for _, optFn := range opts {
		optFn(o)
	}
	return o
}
