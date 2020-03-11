package logger

import (
	"bytes"
	"fmt"
	"github.com/maruel/panicparse/stack"
	"log"
	"os"
	goruntime "runtime"
	"strings"
)

const (
	maxPanicLines = 300
)

type PanicError struct {
	Reason string
	Stack  string

	GoroutineBuckets []*stack.Bucket
}

func (pe PanicError) Error() string {
	return pe.Reason
}

func (pe PanicError) Pretty() string {
	return fmt.Sprintf("%s\n\n%s", pe.Reason, pe.Stack)
}

// WrapRecover wraps panic with prettified stack.
// You must provide result of recover() as an argument. If nothing is recovered, it returns nil.
//
// Example:
//   defer func() {
//      if err := logger.WrapRecover(recover()); err != nil {
//          logger.Fatal("Failed with {}", err.Pretty())
//      }
//   }()
func WrapRecover(r interface{}) *PanicError {
	if r == nil {
		return nil
	}
	reason := fmt.Sprintf("panic: %s", r)

	st := make([]byte, 1024)
	for {
		n := goruntime.Stack(st, true)
		if n < len(st) {
			st = st[:n]
			break
		}
		st = make([]byte, 2*len(st))
	}
	c, err := stack.ParseDump(bytes.NewReader(st), os.Stdout, true)
	if err != nil {
		log.Printf("warning: unable to parse panic stacktrace: %v\n", err)
		return &PanicError{
			Reason: reason,
			Stack:  string(st),
		}
	}

	// Find out similar goroutine traces and group them into buckets.
	buckets := stack.Aggregate(c.Goroutines, stack.AnyValue)

	// Calculate alignment.
	srcLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(line.SrcLine()); l > srcLen {
				srcLen = l
			}
		}
	}

	prettyStack := ""
	var totalLines int
	for i, bucket := range buckets {
		panicIndex := -1
		for i, call := range bucket.Stack.Calls {
			if call.Func.Name() == "panic" {
				panicIndex = i
			}
		}
		if i == 0 {
			// remove stacks before main panic
			bucket.Stack.Calls = bucket.Stack.Calls[panicIndex+1:]
		} else if i != 0 {
			// dim text color
			prettyStack += "\n\x1b[2m"
		}

		// Print the goroutine header.
		extra := ""
		if s := bucket.SleepString(); s != "" {
			extra += "[" + s + "]"
		}
		if bucket.Locked {
			extra += "[locked]"
		}
		if c := bucket.CreatedByString(false); c != "" {
			extra += "[created by " + c + "]"
		}
		goroutineIDs := ""
		if len(bucket.IDs) < 3 {
			var ids []string
			for _, id := range bucket.IDs {
				ids = append(ids, fmt.Sprintf("#%d", id))
			}
			goroutineIDs = strings.Join(ids, ", ")
		} else {
			goroutineIDs = fmt.Sprintf("Group of %d goroutines", len(bucket.IDs))
		}

		prettyStack += fmt.Sprintf("%s: %s %s", goroutineIDs, bucket.State, extra)
		if totalLines >= maxPanicLines {
			prettyStack += " (...)\n"
			continue
		} else {
			prettyStack += "\n"
			totalLines += 1
		}

		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			prettyStack += fmt.Sprintf(
				"    %-*s  %s(%s)\n",
				srcLen, line.SrcLine(),
				line.Func.PkgDotName(), &line.Args)
			totalLines += 1
		}
		if bucket.Stack.Elided {
			prettyStack += "    (...)\n"
			totalLines += 1
		}

		if i != 0 {
			// reset text color
			prettyStack += "\x1b[0m"
			totalLines += 1
		}
	}
	return &PanicError{
		Reason:           reason,
		Stack:            prettyStack,
		GoroutineBuckets: buckets,
	}
}
