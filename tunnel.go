package j9

// Tunnel is a wrapper for a node that provides a logger.
type Tunnel struct {
	node   Node
	logger Logger
}

// NewTunnel creates a new tunnel with the given node and logger.
func NewTunnel(node Node, logger Logger) *Tunnel {
	lg := logger
	if lg == nil {
		lg = &emptyLogger{}
	}
	return &Tunnel{node: node, logger: lg}
}

// Node returns the node.
func (w *Tunnel) Node() Node {
	return w.node
}

// Logger returns the logger.
func (w *Tunnel) Logger() Logger {
	return w.logger
}

// Runs the given command on the node and panics if there is an error.
func (w *Tunnel) Run(cmd string) []byte {
	output, err := w.run(false, cmd)
	if err != nil {
		panic(err)
	}
	return output
}

// Runs the given command on the node and returns the output and error.
func (w *Tunnel) RunOrError(cmd string) ([]byte, error) {
	return w.run(true, cmd)
}

func (w *Tunnel) run(ignoreError bool, cmd string) ([]byte, error) {
	if ignoreError {
		w.logger.Log(LogLevelInfo, "🚙 "+cmd)
	} else {
		w.logger.Log(LogLevelInfo, "🚗 "+cmd)
	}
	output, err := w.node.RunOrError(cmd)
	if err != nil {
		if len(output) > 0 {
			w.logger.Log(LogLevelError, string(output))
		}
		w.logger.Log(LogLevelError, err.Error())
		if !ignoreError {
			panic(err)
		}
	} else {
		if len(output) > 0 {
			w.logger.Log(LogLevelVerbose, string(output))
		}
	}
	return output, err
}
