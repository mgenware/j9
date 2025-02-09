package j9

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
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
	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            config.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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

func (node *SSHNode) Spawn(params *SpawnParams) error {
	return errors.New("RunCmd is not supported in SSHNode")
}

func (node *SSHNode) Shell(params *ShellParams) (string, error) {
	return node.runCore(func(session *ssh.Session) (string, error) {
		if params.WorkingDir != "" {
			params.Cmd = "cd '" + params.WorkingDir + "' && " + params.Cmd
		}
		output, err := session.CombinedOutput(params.Cmd)
		return string(output), err
	})
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
