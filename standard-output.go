package logger

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	isterminal "github.com/azer/is-terminal"
)

func NewStandardOutput(file *os.File, levelSettings, filterSettings string) StandardWriter {
	var writer = StandardWriter{
		ColorsEnabled: isterminal.IsTerminal(int(file.Fd())),
		Target:        file,
	}

	if os.Getenv("LOG_LEVEL") != "" {
		levelSettings = os.Getenv("LOG_LEVEL")
	}
	if os.Getenv("LOG") != "" {
		filterSettings = os.Getenv("LOG")
	} else if filterSettings == "" {
		filterSettings = "*"
	}
	defaultOutputSettings := parseVerbosityLevel(levelSettings)
	writer.Verbosities = parsePackageSettings(filterSettings, defaultOutputSettings)

	return writer
}

type StandardWriter struct {
	ColorsEnabled bool
	Target        *os.File
	Verbosities   map[string]LogPriority
}

func (sw StandardWriter) Init() {}

func (sw StandardWriter) Write(log *Log) {
	if sw.IsEnabled(log.Package, log.Level) {
		fmt.Fprintln(sw.Target, sw.Format(log))
	}
}

func (sw *StandardWriter) IsEnabled(logger string, level *LogLevel) bool {
	verbosity := sw.LogVerbosityOfPackage(logger)
	return level.Priority >= verbosity
}

func (sw *StandardWriter) LogVerbosityOfPackage(p string) LogPriority {
	if settings, ok := sw.Verbosities[p]; ok {
		return settings
	}

	// If there is a "*" (Select all) setting, return that
	if settings, ok := sw.Verbosities["*"]; ok {
		return settings
	}
	return mutePriority
}

func (sw *StandardWriter) Format(log *Log) string {
	if sw.ColorsEnabled {
		return sw.PrettyFormat(log)
	} else {
		return sw.JSONFormat(log)
	}
}

func (sw *StandardWriter) JSONFormat(log *Log) string {
	str, err := json.Marshal(log)
	if err != nil {
		return fmt.Sprintf(`{ "logger-error": "%v" }`, err)
	}

	return string(str)
}

func (sw *StandardWriter) PrettyFormat(log *Log) string {
	msg := fmt.Sprintf(
		"%s │ %s%s: %s%s",
		log.Level.Symbol(),
		log.Package,
		sw.PrettyLabelExt(log),
		log.Message,
		sw.PrettyAttrs(log),
	)
	return fmt.Sprintf(
		"%s %s",
		Colored(dim, time.Now().Format("2006-01-02 15:04:05.000")),
		Colored(log.Level.Color, msg),
	)
}

func (sw *StandardWriter) PrettyAttrs(log *Log) string {
	if *log.DisplayedAttrs == nil {
		return ""
	}

	result := ""
	for key, val := range *log.DisplayedAttrs {
		if byteval, ok := val.([]byte); ok {
			val = hex.EncodeToString(byteval)
		}
		result = fmt.Sprintf("%s %s=%v", result, key, val)
	}

	if log.Level == Fatal {
		result = Colored(Red, result)
	}
	return result
}

func (sw *StandardWriter) PrettyLabel(log *Log) string {
	return fmt.Sprintf("%s%s │ %s%s:%s",
		log.Level.Color,
		log.Level.Symbol(),
		log.Package,
		sw.PrettyLabelExt(log),
		Reset)
}

func (sw *StandardWriter) PrettyLabelExt(log *Log) string {
	if log.Level == Timer {
		return fmt.Sprintf("(%v)", time.Duration(log.ElapsedNano))
	}
	return ""
}

// Accepts: foo,bar,qux@timer
//          *
//          *@error
//          *@error,database@timer
func parsePackageSettings(input string, defaultVerbosity LogPriority) map[string]LogPriority {
	all := map[string]LogPriority{}
	items := strings.Split(input, ",")

	for _, item := range items {
		name, verbosity := parsePackageName(item)
		if verbosity == -1 {
			verbosity = defaultVerbosity
		}
		all[name] = verbosity
	}
	return all
}

// Accepts: users
//          database@timer
//          server@error
func parsePackageName(input string) (string, LogPriority) {
	parsed := strings.Split(input, "@")
	name := strings.TrimSpace(parsed[0])

	if len(parsed) > 1 {
		return name, parseVerbosityLevel(parsed[1])
	}
	return name, -1
}

func parseVerbosityLevel(val string) LogPriority {
	val = strings.ToUpper(strings.TrimSpace(val))
	if lvl, ok := logLevelNameMap[val]; !ok {
		if val == "MUTE" {
			return Fatal.Priority
		}
		// "*" or unknown level: verbose
		return 0
	} else {
		return lvl.Priority
	}
}
