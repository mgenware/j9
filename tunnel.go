package j9

import (
	"encoding/json"
	"path/filepath"
)

// Tunnel is a wrapper for a node that provides a logger.
type Tunnel struct {
	node   Node
	logger Logger

	lastDir string
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

func (w *Tunnel) LastDir() string {
	return w.lastDir
}

// Runs the given command and arguments. Returns an error if the command fails.
func (w *Tunnel) RunRaw(name string, arg ...string) error {
	cmdString := w.getCmdLog(name, arg...)
	_, err := w.logAndCall(cmdString, func() ([]byte, error) {
		err := w.node.RunCmd(w.lastDir, name, arg...)
		return nil, err
	})
	return err
}

// Runs and given command string and returns the output.
// Returns an error if the command fails.
func (w *Tunnel) RunSyncRaw(cmd string) ([]byte, error) {
	return w.logAndCall(cmd, func() ([]byte, error) {
		return w.node.RunCmdSync(w.lastDir, cmd)
	})
}

func (w *Tunnel) CDRaw(dir string) error {
	w.logger.Log(LogLevelInfo, "cd "+dir)
	if filepath.IsAbs(dir) {
		w.lastDir = dir
	} else {
		w.lastDir = filepath.Join(w.lastDir, dir)
	}
	// Execute the command.
	_, err := w.node.RunCmdSync(w.lastDir, "pwd")
	if err != nil {
		return err
	}
	return nil
}

func (w *Tunnel) CD(dir string) {
	err := w.CDRaw(dir)
	if err != nil {
		panic(err)
	}
}

// Runs the given command and arguments. Panics if there is an error.
func (w *Tunnel) Run(name string, arg ...string) {
	err := w.RunRaw(name, arg...)
	if err != nil {
		panic(err)
	}
}

// Runs and given command string and returns the output.
// Panics if there is an error.
func (w *Tunnel) RunSync(cmd string) []byte {
	output, err := w.RunSyncRaw(cmd)
	if err != nil {
		panic(err)
	}
	return output
}

func (w *Tunnel) logAndCall(cmdLog string, runCb func() ([]byte, error)) ([]byte, error) {
	w.logger.Log(LogLevelInfo, cmdLog)
	var output []byte
	var err error
	output, err = runCb()
	if err != nil {
		if len(output) > 0 {
			w.logger.Log(LogLevelError, string(output))
		}
		w.logger.Log(LogLevelError, err.Error())
		return output, err
	}
	if len(output) > 0 {
		w.logger.Log(LogLevelVerbose, string(output))
	}
	return output, nil
}

func (w *Tunnel) getCmdLog(name string, args ...string) string {
	cmd := name
	for _, arg := range args {
		argJson, _ := json.Marshal(arg)
		cmd += " " + string(argJson)
	}
	return cmd
}
