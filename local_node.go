package j9

import (
	"os"
	"os/exec"
)

// LocalNode is used for running commands locally.
type LocalNode struct {
}

func NewLocalNode() *LocalNode {
	return &LocalNode{}
}

func (node *LocalNode) Spawn(params *SpawnParams) error {
	c := exec.Command(params.Name, params.Args...)
	if params.WorkingDir != "" {
		c.Dir = params.WorkingDir
	}
	c.Env = params.Env
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (node *LocalNode) Shell(params *ShellParams) (string, error) {
	c := exec.Command("sh", "-c", params.Cmd)
	if params.WorkingDir != "" {
		c.Dir = params.WorkingDir
	}
	c.Env = params.Env
	output, err := c.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
