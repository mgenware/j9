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

func TestLocalRunSync(t *testing.T) {
	output, err := localNode.RunSyncUnsafe("echo abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalRun(t *testing.T) {
	err := localNode.RunUnsafe("echo", "abc")
	assert.NoError(t, err)
}

func TestLocalRunSyncWithError(t *testing.T) {
	output, err := localNode.RunSyncUnsafe("cat ./__not_exist__")
	// Output should not be empty even if error happened
	assert.Equal(t, string(output) != "", true)
	assert.Error(t, err)
}

func TestLocalRunCD(t *testing.T) {
	localNode.CDUnsafe("/")
	output := runSync(localNode, "pwd")
	assert.Equal(t, "/\n", string(output))
	localNode.CDUnsafe("~")
	output = runSync(localNode, "pwd")
	assert.Equal(t, os.Getenv("HOME")+"\n", string(output))
}

func TestLocalLastDir(t *testing.T) {
	localNode.CDUnsafe("/")
	output := runSync(localNode, "pwd")
	assert.Equal(t, "/\n", string(output))
	// Double check if last dir is kept on subsequent commands
	output = runSync(localNode, "pwd")
	assert.Equal(t, "/\n", string(output))
}

func runSync(node *LocalNode, cmd string) string {
	output, err := node.RunSyncUnsafe(cmd)
	if err != nil {
		panic(err)
	}
	return string(output)
}
