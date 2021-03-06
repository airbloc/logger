package logger

import (
	"os"
)

var (
	runtime *Runtime
)

func init() {
	runtime = &Runtime{
		Writers: []OutputWriter{
			NewStandardOutput(os.Stderr, "", "*"),
		},
	}
}

type OutputWriter interface {
	Init()
	Write(log *Log)
}

type Runtime struct {
	Writers []OutputWriter
}

func (runtime *Runtime) Log(log *Log) {
	if len(runtime.Writers) == 0 {
		return
	}

	// Avoid getting into a loop if there is just one writer
	if len(runtime.Writers) == 1 {
		runtime.Writers[0].Write(log)
	} else {
		for _, w := range runtime.Writers {
			w.Write(log)
		}
	}
}

// Add a new writer
func Hook(writer OutputWriter) {
	writer.Init()
	runtime.Writers = append(runtime.Writers, writer)
}

// Legacy method
func SetLogger(writer OutputWriter) {
	runtime.Writers[0] = writer
}
