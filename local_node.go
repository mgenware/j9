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

func (node *LocalNode) RunCmdSync(wd string, cmd string) ([]byte, error) {
	c := exec.Command("bash", "-c", cmd)
	c.Dir = wd
	output, err := c.CombinedOutput()
	if err != nil {
		return output, err
	}
	return output, nil
}

func (node *LocalNode) RunCmd(wd string, name string, arg ...string) error {
	c := exec.Command(name, arg...)
	c.Dir = wd
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
