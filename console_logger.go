package j9

import (
	"fmt"

	"github.com/fatih/color"
)

type ConsoleLogger struct {
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Log(level int, message string) {
	if level == LogLevelVerbose {
		fmt.Println(message)
	} else {
		var console func(format string, a ...interface{})
		switch level {
		case LogLevelError:
			console = color.Red
		case LogLevelWarning:
			console = color.Yellow
		default:
			console = color.Cyan
		}

		console(message)
	}
}
