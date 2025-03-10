package j9

import (
	"encoding/json"
	"path/filepath"
)

// Tunnel is a wrapper for a node that provides a logger.
type Tunnel struct {
	node   Node
	logger Logger

	dir string
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

func (w *Tunnel) Dir() string {
	return w.dir
}

func (w *Tunnel) SpawnRaw(params *SpawnOpt) error {
	logString := params.Name
	for _, arg := range params.Args {
		argJson, _ := json.Marshal(arg)
		logString += " " + string(argJson)
	}

	_, err := w.logAndCall(logString, func() (string, error) {
		// Update working dir if needed.
		if params.WorkingDir == "" {
			params.WorkingDir = w.dir
		}
		err := w.node.Spawn(params)
		return "", err
	})
	return err
}

func (w *Tunnel) ShellRaw(params *ShellOpt) (string, error) {
	return w.logAndCall(params.Cmd, func() (string, error) {
		// Update working dir if needed.
		if params.WorkingDir == "" {
			params.WorkingDir = w.dir
		}
		return w.node.Shell(params)
	})
}

func (w *Tunnel) CD(dir string) {
	w.logger.Log(LogLevelInfo, "cd "+dir)
	if filepath.IsAbs(dir) {
		w.dir = dir
	} else {
		w.dir = filepath.Join(w.dir, dir)
	}
}

func (w *Tunnel) Spawn(params *SpawnOpt) {
	err := w.SpawnRaw(params)
	if err != nil {
		panic(err)
	}
}

func (w *Tunnel) Shell(params *ShellOpt) string {
	output, err := w.ShellRaw(params)
	if err != nil {
		panic(err)
	}
	return output
}

func (w *Tunnel) logAndCall(cmdLog string, runCb func() (string, error)) (string, error) {
	w.logger.Log(LogLevelInfo, cmdLog)
	var output string
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
