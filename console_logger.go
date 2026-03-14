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
	var console func(format string, a ...interface{})
	switch level {
	case LogLevelError:
		console = color.Red
	case LogLevelWarning:
		console = color.Yellow
	case LogLevelInfo:
		console = color.Cyan
	case LogLevelSuccess:
		console = color.Green
	}

	if console == nil {
		fmt.Println(message)
	} else {
		console(message)
	}
}
