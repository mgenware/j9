package j9

import (
	"os"
	"os/exec"
)

// LocalNode is used for running commands locally.
type LocalNode struct {
	lastDir string
}

func NewLocalNode() *LocalNode {
	return &LocalNode{}
}

func (node *LocalNode) RunSyncUnsafe(cmd string) ([]byte, error) {
	c := exec.Command("bash", "-c", cmd)
	c.Dir = node.lastDir
	output, err := c.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}

func (node *LocalNode) RunUnsafe(name string, arg ...string) error {
	c := exec.Command(name, arg...)
	c.Dir = node.lastDir
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (node *LocalNode) CDUnsafe(dir string) error {
	node.lastDir = dir
	return nil
}
