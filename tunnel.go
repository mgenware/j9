package j9

type Tunnel struct {
	node   Node
	logger Logger
}

func NewTunnel(node Node, logger Logger) *Tunnel {
	lg := logger
	if lg == nil {
		lg = &emptyLogger{}
	}
	return &Tunnel{node: node, logger: lg}
}

func (w *Tunnel) Node() Node {
	return w.node
}

func (w *Tunnel) Logger() Logger {
	return w.logger
}

func (w *Tunnel) Run(cmd string) []byte {
	output, err := w.run(false, cmd)
	if err != nil {
		panic(err)
	}
	return output
}

func (w *Tunnel) RunOrError(cmd string) ([]byte, error) {
	return w.run(true, cmd)
}

func (w *Tunnel) run(ignore bool, cmd string) ([]byte, error) {
	if ignore {
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
		if !ignore {
			panic(err)
		}
	} else {
		if len(output) > 0 {
			w.logger.Log(LogLevelVerbose, string(output))
		}
	}
	return output, err
}
