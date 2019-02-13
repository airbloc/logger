package logger

import (
	"fmt"
	"sync"
)

var (
	colors  sync.Map
	white   = "\033[37m"
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[90m"
	Red     = "\033[31m"
	Blue    = "\033[34m"
	Green   = "\033[32m"
	Cyan    = "\033[36m"
	Yellow  = "\033[33m"
	Magenta = "\033[35m"
)

func init() {
	colors.Store("index", 0)
	colors.Store("index:0", Blue)
	colors.Store("index:1", Green)
	colors.Store("index:2", Cyan)
	colors.Store("index:3", Yellow)
	colors.Store("index:4", Magenta)
	colors.Store("len", 5)
}

func nextColor() string {
	//colorIndex = colorIndex + 1
	//return colors[colorIndex%len(colors)]
	currentIndex, _ := colors.Load("index")
	len, _ := colors.Load("len")
	color, _ := colors.Load(fmt.Sprintf("index:%d", currentIndex.(int)%len.(int)))

	colors.Store("index", currentIndex.(int)+1)

	return color.(string)
}

func colorFor(key string) string {
	if color, ok := colors.Load(fmt.Sprintf("module:%s", key)); ok {
		return color.(string)
	}

	color := nextColor()
	colors.Store(fmt.Sprintf("module:%s", key), color)
	return color
}

func Colored(color string, msg string) string {
	return color + msg + reset
}
