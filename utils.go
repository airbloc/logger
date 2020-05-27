package logger

import (
	"gopkg.in/yaml.v2"
	"strings"
)

// Inspect dumps given struct content to the Logger, with verbose message.
func Inspect(log Logger, prefixMsg string, structVal interface{}) {
	raw, err := yaml.Marshal(structVal)
	if err != nil {
		log.Warn("Unable to inspect struct (message was {})", err, prefixMsg)
	}
	logs := prefixMsg
	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		if strings.HasSuffix(line, `""`) {
			// skip empty value
			continue
		}
		logs += "\n" + strings.Repeat(" ", 25) + " │   " + line
	}
	log.Verbose(logs)
}
