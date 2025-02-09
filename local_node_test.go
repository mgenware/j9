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

func TestLocalShellCmd(t *testing.T) {
	output, err := localNode.Shell(&ShellParams{Cmd: "echo abc"})
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalShellWD(t *testing.T) {
	mustMkdirp("test_folders/a")
	output, err := localNode.Shell(&ShellParams{Cmd: "basename $(pwd)", WorkingDir: "test_folders/a"})
	assert.NoError(t, err)
	assert.Equal(t, "a\n", string(output))
}

func TestLocalSpawnCmd(t *testing.T) {
	err := localNode.Spawn(&SpawnParams{Name: "echo", Args: []string{"abc"}})
	assert.NoError(t, err)
}

func TestLocalShellError(t *testing.T) {
	output, err := localNode.Shell(&ShellParams{Cmd: "cat ./__not_exist__"})
	// Output should not be empty even if error happened
	assert.Equal(t, string(output) != "", true)
	assert.Error(t, err)
}

func TestLocalShellEnv(t *testing.T) {
	output, err := localNode.Shell(&ShellParams{Cmd: "echo $MY_ENV", Env: []string{"MY_ENV=abc"}})
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalSpawnEnv(t *testing.T) {
	err := localNode.Spawn(&SpawnParams{Name: "sh", Args: []string{"-c", "echo $MY_ENV"}, Env: []string{"MY_ENV=abc"}})
	assert.NoError(t, err)
}
