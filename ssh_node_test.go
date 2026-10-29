package j9

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestSSHNodeUsesConfiguredHostKeyCallback(t *testing.T) {
	wantErr := errors.New("host key rejected")
	hostKeyCallback := func(string, net.Addr, ssh.PublicKey) error {
		return wantErr
	}
	node, err := NewSSHNode(&SSHConfig{
		Host:            "example.com",
		User:            "user",
		Auth:            []ssh.AuthMethod{ssh.Password("password")},
		HostKeyCallback: hostKeyCallback,
	})
	require.NoError(t, err)
	require.NotNil(t, node.sshConfig.HostKeyCallback)
	assert.ErrorIs(t, node.sshConfig.HostKeyCallback("example.com:22", nil, nil), wantErr)
}

func TestSSHNodeLoadsKnownHostsByDefault(t *testing.T) {
	home := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(home, ".ssh"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(home, ".ssh", "known_hosts"), nil, 0600))
	t.Setenv("HOME", home)

	node, err := NewSSHNode(&SSHConfig{
		Host: "example.com",
		User: "user",
		Auth: []ssh.AuthMethod{ssh.Password("password")},
	})
	require.NoError(t, err)
	assert.NotNil(t, node.sshConfig.HostKeyCallback)
}

func TestSSHNodeRejectsMissingKnownHostsFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_, err := NewSSHNode(&SSHConfig{
		Host: "example.com",
		User: "user",
		Auth: []ssh.AuthMethod{ssh.Password("password")},
	})
	assert.ErrorContains(t, err, "load SSH known_hosts")
}

func TestBuildSSHCommandQuotesDirectoryAndEnvironment(t *testing.T) {
	workingDir := filepath.Join(t.TempDir(), "dir'quoted")
	require.NoError(t, os.Mkdir(workingDir, 0700))
	envValue := "literal'; printf HACKED; #"
	params := &ShellOpt{
		WorkingDir: workingDir,
		Env:        []string{"VALUE=" + envValue},
		Cmd:        `printf '%s|%s' "$VALUE" "$PWD"`,
	}

	cmd, err := buildSSHCommand(params)
	require.NoError(t, err)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	require.NoError(t, err, string(output))
	assert.Equal(t, envValue+"|"+workingDir, string(output))
	assert.Equal(t, `printf '%s|%s' "$VALUE" "$PWD"`, params.Cmd)
}

func TestBuildSSHCommandRejectsInvalidEnvironmentName(t *testing.T) {
	_, err := buildSSHCommand(&ShellOpt{
		Env: []string{"BAD;NAME=value"},
		Cmd: "true",
	})
	assert.ErrorContains(t, err, "invalid environment variable assignment")
}
