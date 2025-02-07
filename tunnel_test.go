package j9

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newLocalTunnel() *Tunnel {
	return NewTunnel(NewLocalNode(), NewConsoleLogger())
}

func TestTunnelCD(t *testing.T) {
	mustMkdirp("test_folders/a/b/c")

	tn := newLocalTunnel()
	assert.Equal(t, tn.LastDir(), "")

	tn.CD("test_folders/a/b/c")
	assert.Equal(t, tn.LastDir(), "test_folders/a/b/c")
	output := tn.RunSync("basename $(pwd)")
	assert.Equal(t, "c\n", string(output))

	tn.CD("..")
	assert.Equal(t, tn.LastDir(), "test_folders/a/b")
	output = tn.RunSync("basename $(pwd)")
	assert.Equal(t, "b\n", string(output))

	tn.CD("/tmp")
	assert.Equal(t, tn.LastDir(), "/tmp")
	output = tn.RunSync("basename $(pwd)")
	assert.Equal(t, "tmp\n", string(output))
}
