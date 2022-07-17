package j9

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var localNode *LocalNode

func init() {
	localNode = NewLocalNode()
}

func TestLocalRunOrError(t *testing.T) {
	output, err := localNode.RunOrError("echo abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalRunOrErrorWithError(t *testing.T) {
	output, err := localNode.RunOrError("cat ./__not_exist__")
	// Output should not be empty even if error happened
	assert.Equal(t, string(output) != "", true)
	assert.Error(t, err)
}

func TestLocalRunCD(t *testing.T) {
	_ = run(localNode, "cd /")
	output := run(localNode, "pwd")
	assert.Equal(t, "/\n", string(output))
	_ = run(localNode, "cd ~")
	output = run(localNode, "pwd")
	assert.Equal(t, os.Getenv("HOME")+"\n", string(output))
}

func TestLocalLastDir(t *testing.T) {
	tmp := "/"
	_ = run(localNode, "cd "+tmp)
	output := run(localNode, "pwd")
	assert.Equal(t, tmp+"\n", string(output))
	// Double check if last dir is kept on subsequent commands
	output = run(localNode, "pwd")
	assert.Equal(t, tmp+"\n", string(output))
}

func run(node *LocalNode, cmd string) string {
	output, err := node.RunOrError(cmd)
	if err != nil {
		panic(err)
	}
	return string(output)
}
