package j9

const (
	LogLevelError   = iota
	LogLevelWarning = iota
	LogLevelInfo    = iota
	LogLevelVerbose = iota
)

type Logger interface {
	Log(level int, message string)
}

type emptyLogger struct {
}

func (logger *emptyLogger) Log(level int, message string) {}
