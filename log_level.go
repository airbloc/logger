package logger

import (
	"encoding/json"
)

type LogPriority int

const (
	// mute only emits fatal log
	mutePriority = 99
)

type LogLevel struct {
	Name     string
	Color    string
	Priority LogPriority
}

func (lvl *LogLevel) String() string {
	return lvl.Name
}

func (lvl *LogLevel) Symbol() string {
	return string(lvl.Name[0])
}

func (lvl *LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(lvl.String())
}

var (
	Verbose = addLogLevel("Verbose", dim, 0)
	Debug   = addLogLevel("DEBUG", white, 1)
	Info    = addLogLevel("INFO", Reset, 2)
	Timer   = addLogLevel("TIMER", Green, 3)
	Warn    = addLogLevel("WARN", Yellow, 5)
	Error   = addLogLevel("ERROR", Red, 10)
	Fatal   = addLogLevel("FATAL", Red, 99)

	logLevelNameMap = map[string]*LogLevel{}
)

func addLogLevel(name, color string, priority LogPriority) *LogLevel {
	lvl := &LogLevel{
		Name:     name,
		Color:    color,
		Priority: priority,
	}
	logLevelNameMap[name] = lvl
	return lvl
}
