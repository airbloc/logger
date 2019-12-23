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

func (s *subLogger) Recover(context Attrs) interface{} {
	mergedAttrs := s.defaultAttrs.Merge(context)
	return s.parent.Recover(mergedAttrs)
}

func (s *subLogger) WithAttrs(attrs Attrs) Logger {
	return &subLogger{
		parent:       s,
		defaultAttrs: attrs,
	}
}
