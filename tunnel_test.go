package j9

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newLocalTunnel() *Tunnel {
	return NewTunnel(NewLocalNode(), NewConsoleLogger())
}

func TestTunnelShellDirs(t *testing.T) {
	mustMkdirp("test_folders/a/b/c")

	tn := newLocalTunnel()
	assert.Equal(t, tn.Dir(), "")

	tn.CD("test_folders/a/b/c")
	assert.Equal(t, tn.Dir(), "test_folders/a/b/c")
	output := tn.Shell(&ShellOpt{Cmd: "basename $(pwd)"})
	assert.Equal(t, "c\n", string(output))

	tn.CD("..")
	assert.Equal(t, tn.Dir(), "test_folders/a/b")
	output = tn.Shell(&ShellOpt{Cmd: "basename $(pwd)"})
	assert.Equal(t, "b\n", string(output))

	tn.CD("/tmp")
	assert.Equal(t, tn.Dir(), "/tmp")
	output = tn.Shell(&ShellOpt{Cmd: "basename $(pwd)"})
	assert.Equal(t, "tmp\n", string(output))
}

func TestTunnelSpawnDirs(t *testing.T) {
	mustMkdirp("test_folders/a/b/c")

	tn := newLocalTunnel()
	assert.Equal(t, tn.Dir(), "")

	tn.CD("test_folders/a/b/c")
	assert.Equal(t, tn.Dir(), "test_folders/a/b/c")
	// WD is auto set by tunnel.Spawn
	_, d := spawnTunnelTestScript(tn, &SpawnOpt{Name: "../../../../test_spawn.sh"})
	assert.Equal(t, "c", d)

	tn.CD("..")
	assert.Equal(t, tn.Dir(), "test_folders/a/b")
	// WD is auto set by tunnel.Spawn
	_, d = spawnTunnelTestScript(tn, &SpawnOpt{Name: "../../../test_spawn.sh"})
	assert.Equal(t, "b", d)

	tn.CD("/tmp")
	assert.Equal(t, tn.Dir(), "/tmp")
	// WD is auto set by tunnel.Spawn
	// We are in /tmp, so locate the test_spawn.sh, we have to use absolute path.
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	shPath := filepath.Join(wd, "test_spawn.sh")
	_, d = spawnTunnelTestScript(tn, &SpawnOpt{Name: shPath})
	assert.Equal(t, "tmp", d)
}

func spawnTunnelTestScript(tn *Tunnel, opt *SpawnOpt) (string, string) {
	f, err := os.CreateTemp("", "j9_test_spawn")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	// Overwrite args to use the temp file.
	opt.Args = []string{f.Name()}

	tn.Spawn(opt)

	output, err := os.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	v := strings.TrimLeft(lines[0], "V=")
	d := strings.TrimLeft(lines[1], "D=")
	return v, d
}

func TestTunnelShellEnv(t *testing.T) {
	tn := newLocalTunnel()
	output := tn.Shell(&ShellOpt{Cmd: "echo $MY_ENV", Env: []string{"MY_ENV=abc"}})
	assert.Equal(t, "abc\n", string(output))
}

func TestTunnelSpawnEnv(t *testing.T) {
	tn := newLocalTunnel()
	env, wd := spawnTunnelTestScript(tn, &SpawnOpt{Name: "./test_spawn.sh", Env: []string{"VAR1=123"}})
	assert.Equal(t, "123", env)
	assert.Equal(t, "j9", wd)
}
