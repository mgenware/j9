package j9

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSHNode is a node that runs commands on a remote SSH server.
type SSHNode struct {
	config    *SSHConfig
	sshConfig *ssh.ClientConfig
	client    *ssh.Client

	Logger Logger
}

// Creates a new SSH node with the given configuration.
func NewSSHNode(config *SSHConfig) (*SSHNode, error) {
	if config.Host == "" {
		return nil, fmt.Errorf("config.Host cannot be empty")
	}
	if config.User == "" {
		return nil, fmt.Errorf("config.User cannot be empty")
	}
	if len(config.Auth) == 0 {
		return nil, fmt.Errorf("config.Auth cannot be empty")
	}
	if config.Port == 0 {
		config.Port = 22
	}
	hostKeyCallback := config.HostKeyCallback
	if hostKeyCallback == nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("resolve home directory for SSH known_hosts: %w", err)
		}
		hostKeyCallback, err = knownhosts.New(filepath.Join(home, ".ssh", "known_hosts"))
		if err != nil {
			return nil, fmt.Errorf("load SSH known_hosts: %w", err)
		}
	}
	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            config.Auth,
		HostKeyCallback: hostKeyCallback,
		Timeout:         30 * time.Second,
	}

	return &SSHNode{
		config:    config,
		sshConfig: sshConfig,
	}, nil
}

// Creates a new SSH node with the given configuration, panics if there is an error.
func MustCreateSSHNode(config *SSHConfig) *SSHNode {
	node, err := NewSSHNode(config)
	if err != nil {
		panic(err)
	}
	return node
}

func (node *SSHNode) Spawn(params *SpawnOpt) error {
	return errors.New("RunCmd is not supported in SSHNode")
}

func (node *SSHNode) Shell(params *ShellOpt) (string, error) {
	cmd, err := buildSSHCommand(params)
	if err != nil {
		return "", err
	}
	return node.runCore(func(session *ssh.Session) (string, error) {
		output, err := session.CombinedOutput(cmd)
		return string(output), err
	})
}

func buildSSHCommand(params *ShellOpt) (string, error) {
	var prefix strings.Builder
	if params.WorkingDir != "" {
		prefix.WriteString("cd ")
		prefix.WriteString(shellQuote(params.WorkingDir))
		prefix.WriteString(" && ")
	}
	if len(params.Env) > 0 {
		assignments := make([]string, 0, len(params.Env))
		for _, entry := range params.Env {
			name, value, ok := strings.Cut(entry, "=")
			if !ok || !validEnvName(name) {
				return "", fmt.Errorf("invalid environment variable assignment %q", entry)
			}
			assignments = append(assignments, name+"="+shellQuote(value))
		}
		prefix.WriteString("export ")
		prefix.WriteString(strings.Join(assignments, " "))
		prefix.WriteString(" && ")
	}
	prefix.WriteString(params.Cmd)
	return prefix.String(), nil
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func validEnvName(name string) bool {
	if name == "" || !(name[0] == '_' || name[0] >= 'A' && name[0] <= 'Z' || name[0] >= 'a' && name[0] <= 'z') {
		return false
	}
	for i := 1; i < len(name); i++ {
		char := name[i]
		if char != '_' && !(char >= 'A' && char <= 'Z') && !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}

func (node *SSHNode) runCore(sessionCb func(*ssh.Session) (string, error)) (string, error) {
	session, err := node.prepareSession()
	if err != nil {
		return "", err
	}
	return sessionCb(session)
}

func (node *SSHNode) prepareSession() (*ssh.Session, error) {
	if node.client != nil {
		node.client.Close()
	}

	cfg := node.config
	addr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	node.log(LogLevelInfo, fmt.Sprintf("SSH: Connecting to %v\n", addr))

	client, err := ssh.Dial("tcp", addr, node.sshConfig)
	if err != nil {
		return nil, err
	}
	node.client = client
	session, err := node.client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (node *SSHNode) log(logLevel int, message string) {
	if node.Logger != nil {
		node.Logger.Log(logLevel, message)
	}
}
