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

func mustMkdirp(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
}

func TestLocalRunCmdSync(t *testing.T) {
	output, err := localNode.RunCmdSync("", "echo abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalRunCmdSyncWD(t *testing.T) {
	mustMkdirp("test_folders/a")
	output, err := localNode.RunCmdSync("test_folders/a", "basename $(pwd)")
	assert.NoError(t, err)
	assert.Equal(t, "a\n", string(output))
}

func TestLocalRunCmd(t *testing.T) {
	err := localNode.RunCmd("", "echo", "abc")
	assert.NoError(t, err)
}

func TestLocalRunCmdSyncWithError(t *testing.T) {
	output, err := localNode.RunCmdSync("", "cat ./__not_exist__")
	// Output should not be empty even if error happened
	assert.Equal(t, string(output) != "", true)
	assert.Error(t, err)
}
