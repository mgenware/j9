package loggers

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mgenware/j9"
)

type ConsoleLogger struct {
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Log(level int, message string) {
	if level == j9.LogLevelVerbose {
		fmt.Println(message)
	} else {
		var console func(format string, a ...interface{})
		switch level {
		case j9.LogLevelError:
			console = color.Red
		case j9.LogLevelWarning:
			console = color.Yellow
		default:
			console = color.Cyan
		}

		console(message)
	}
}
