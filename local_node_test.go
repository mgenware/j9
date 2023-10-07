package j9

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var localNode *LocalNode

func init() {
	localNode = NewLocalNode()
}

func TestLocalRunSync(t *testing.T) {
	output, err := localNode.RunCmdSync("", "echo abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalRun(t *testing.T) {
	err := localNode.RunCmd("", "echo", "abc")
	assert.NoError(t, err)
}

func TestLocalRunSyncWithError(t *testing.T) {
	output, err := localNode.RunCmdSync("", "cat ./__not_exist__")
	// Output should not be empty even if error happened
	assert.Equal(t, string(output) != "", true)
	assert.Error(t, err)
}
