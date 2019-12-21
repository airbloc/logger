package logger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Attrs map[string]interface{}

type Log struct {
	Package     string    `json:"package"`
	Level       *LogLevel `json:"level"`
	Message     string    `json:"msg"`
	Attrs       *Attrs    `json:"attrs"`
	Time        int64     `json:"time"`
	Elapsed     int64     `json:"elapsed"`
	ElapsedNano int64     `json:"elapsed_nano"`

	// only shown in console (purged attributes)
	DisplayedAttrs *Attrs `json:"-"`
}

func (log *Log) End(msg string, args ...interface{}) {
	attrs := MergeAttrs(args)
	elapsed := Now() - log.Time
	formatted, purgedAttrs := Format(msg, *attrs)

	log.DisplayedAttrs = &purgedAttrs
	log.Attrs = attrs
	log.Elapsed = elapsed / 1000000
	log.ElapsedNano = elapsed
	log.Message = formatted

	runtime.Log(log)
}

func MergeAttrs(v []interface{}) *Attrs {
	if len(v) == 0 {
		return new(Attrs)
	}

	attrs, ok := v[len(v)-1].(Attrs)
	if !ok {
		// use empty one
		attrs = Attrs{}
	} else {
		// remove last argument
		v = v[:len(v)-1]
	}

	// convert list to map with string index (e.g. ["foo", "bar"] -> {"0": "foo", "1": "bar"}
	for i, value := range v {
		key := strconv.Itoa(i)
		attrs[key] = value
	}
	return &attrs
}

// Format formats string with Python style (PEP 3101, especially key-value formatting through brackets).
// returns unmatched attributes, for better printing.
func Format(format string, attrs Attrs) (formatted string, purged Attrs) {
	formatted = format
	purged = Attrs{}

	for key, value := range attrs {
		placeholder := "{" + key + "}"
		if strings.Contains(format, placeholder) && !strings.Contains(format, "{{"+key+"}}") {
			valueStr := fmt.Sprintf("%v", value)
			formatted = strings.Replace(formatted, placeholder, valueStr, -1)
		} else {
			purged[key] = value
		}
	}

	// assign positional arguments for empty brackets
	if _, ok := purged["0"]; ok {
		index := 0
		for strings.Contains(formatted, "{}") {
			key := strconv.Itoa(index)
			value := fmt.Sprintf("%v", purged[key])
			formatted = strings.Replace(formatted, "{}", value, 1)

			delete(purged, key)
			index++
		}
	}

	// escape by doubling brackets
	escape := strings.NewReplacer("{{", "{", "}}", "}")
	formatted = escape.Replace(formatted)
	return
}

// Now is a shortcut for returning the current time in Unix nanoseconds.
func Now() int64 {
	return time.Now().UnixNano()
}
