package logger

type subLogger struct {
	parent       Logger
	defaultAttrs Attrs
}

func (s *subLogger) mergeWithDefaultAttrs(args []interface{}) []interface{} {
	attrs := Attrs{}
	if len(args) > 0 {
		lastIndex := len(args) - 1
		if att, ok := args[lastIndex].(Attrs); ok {
			attrs = att
			args = args[:lastIndex]
		}
	}
	mergedAttrs := s.defaultAttrs.Merge(attrs)
	return append(args, mergedAttrs)
}

func (s *subLogger) Log(level *LogLevel, message string, args []interface{}) {
	vv := s.mergeWithDefaultAttrs(args)
	s.parent.Log(level, message, vv)
}

func (s *subLogger) Verbose(msg string, v ...interface{}) {
	vv := s.mergeWithDefaultAttrs(v)
	s.parent.Verbose(msg, vv...)
}

func (s *subLogger) Debug(msg string, v ...interface{}) {
	vv := s.mergeWithDefaultAttrs(v)
	s.parent.Debug(msg, vv...)
}

func (s *subLogger) Info(msg string, v ...interface{}) {
	vv := s.mergeWithDefaultAttrs(v)
	s.parent.Info(msg, vv...)
}

func (s *subLogger) Warn(msg string, v ...interface{}) {
	vv := s.mergeWithDefaultAttrs(v)
	s.parent.Warn(msg, vv...)
}

func (s *subLogger) Error(msg string, v ...interface{}) {
	vv := s.mergeWithDefaultAttrs(v)
	s.parent.Error(msg, vv...)
}

func (s *subLogger) Timer() *Log {
	// TODO: implement
	return s.parent.Timer()
}

func (s *subLogger) Fatal(v ...interface{}) {
	vv := s.mergeWithDefaultAttrs(v)
	s.Fatal(vv...)
}

func (s *subLogger) Wtf(v ...interface{}) {
	vv := s.mergeWithDefaultAttrs(v)
	s.Wtf(vv...)
}

func (s *subLogger) Recover(optionalContext ...Attrs) *PanicError {
	if len(optionalContext) > 0 {
		mergedAttrs := s.defaultAttrs.Merge(optionalContext[0])
		return s.parent.Recover(mergedAttrs)
	}
	return s.parent.Recover()
}

func (s *subLogger) WithAttrs(attrs Attrs) Logger {
	return &subLogger{
		parent:       s,
		defaultAttrs: attrs,
	}
}
