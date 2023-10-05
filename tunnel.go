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

// Calls node.RunUnsafe.
func (w *Tunnel) RunUnsafe(name string, arg ...string) error {
	cmdString := w.getCmdString(name, arg...)
	_, err := w.runCore(cmdString, func() ([]byte, error) {
		err := w.node.RunUnsafe(name, arg...)
		return nil, err
	})
	return err
}

// Calls node.RunSyncUnsafe.
func (w *Tunnel) RunSyncUnsafe(cmd string) ([]byte, error) {
	return w.runCore(cmd, func() ([]byte, error) {
		return w.node.RunSyncUnsafe(cmd)
	})
}

// Calls node.CDUnsafe.
func (w *Tunnel) CDUnsafe(dir string) error {
	return w.node.CDUnsafe(dir)
}

// Calls CDUnsafe and panics if there is an error.
func (w *Tunnel) CD(dir string) {
	err := w.CDUnsafe(dir)
	if err != nil {
		panic(err)
	}
}

// Calls RunUnsafe and panics if there is an error.
func (w *Tunnel) Run(name string, arg ...string) {
	err := w.RunUnsafe(name, arg...)
	if err != nil {
		panic(err)
	}
}

// Calls RunSyncUnsafe and panics if there is an error.
func (w *Tunnel) RunSync(cmd string) []byte {
	output, err := w.RunSyncUnsafe(cmd)
	if err != nil {
		panic(err)
	}
	return output
}

func (w *Tunnel) runCore(cmdText string, runCb func() ([]byte, error)) ([]byte, error) {
	w.logger.Log(LogLevelInfo, cmdText)
	var output []byte
	var err error
	output, err = runCb()
	if err != nil {
		if len(output) > 0 {
			w.logger.Log(LogLevelError, string(output))
		}
		w.logger.Log(LogLevelError, err.Error())
	} else {
		if len(output) > 0 {
			w.logger.Log(LogLevelVerbose, string(output))
		}
	}
	return output, err
}

func (w *Tunnel) getCmdString(name string, arg ...string) string {
	cmd := name
	for _, a := range arg {
		cmd += " " + a
	}
	return cmd
}
