package j9

import (
	"os"
	"strings"
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
	output, err := localNode.Shell(&ShellOpt{Cmd: "echo abc"})
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalShellWD(t *testing.T) {
	mustMkdirp("test_folders/a")
	output, err := localNode.Shell(&ShellOpt{Cmd: "basename $(pwd)", WorkingDir: "test_folders/a"})
	assert.NoError(t, err)
	assert.Equal(t, "a\n", string(output))
}

func TestLocalSpawnWD(t *testing.T) {
	mustMkdirp("test_folders/b")
	_, wd := spawnNodeTestScript(localNode, &SpawnOpt{Name: "../../test_spawn.sh", WorkingDir: "test_folders/b"})
	assert.Equal(t, "b", wd)
}

func TestLocalSpawnCmd(t *testing.T) {
	err := localNode.Spawn(&SpawnOpt{Name: "echo", Args: []string{"abc"}})
	assert.NoError(t, err)
}

func TestLocalShellError(t *testing.T) {
	output, err := localNode.Shell(&ShellOpt{Cmd: "cat ./__not_exist__"})
	// Output should not be empty even if error happened
	assert.Equal(t, string(output) != "", true)
	assert.Error(t, err)
}

func TestLocalShellEnv(t *testing.T) {
	output, err := localNode.Shell(&ShellOpt{Cmd: "echo $MY_ENV", Env: []string{"MY_ENV=abc"}})
	assert.NoError(t, err)
	assert.Equal(t, "abc\n", string(output))
}

func TestLocalSpawnEnv(t *testing.T) {
	env, wd := spawnNodeTestScript(localNode, &SpawnOpt{Name: "./test_spawn.sh", Env: []string{"VAR1=123"}})
	assert.Equal(t, "123", env)
	assert.Equal(t, "j9", wd)
}

func spawnNodeTestScript(node *LocalNode, opt *SpawnOpt) (string, string) {
	f, err := os.CreateTemp("", "j9_test_spawn")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	// Overwrite args to use the temp file.
	opt.Args = []string{f.Name()}

	err = node.Spawn(opt)
	if err != nil {
		panic(err)
	}

	output, err := os.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	v := strings.TrimLeft(lines[0], "V=")
	d := strings.TrimLeft(lines[1], "D=")
	return v, d
}
