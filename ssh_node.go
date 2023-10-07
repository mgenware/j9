package j9

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHNode is a node that runs commands on a remote SSH server.
type SSHNode struct {
	config    *SSHConfig
	sshConfig *ssh.ClientConfig
	client    *ssh.Client
	lastDir   string

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

func (node *SSHNode) RunUnsafe(name string, arg ...string) error {
	return errors.New("Run is not supported in SSHNode")
}

func (node *SSHNode) RunSyncUnsafe(cmd string) ([]byte, error) {
	return node.runCore(func(session *ssh.Session) ([]byte, error) {
		return session.CombinedOutput(cmd)
	})
}

func (node *SSHNode) runCore(sessionCb func(*ssh.Session) ([]byte, error)) ([]byte, error) {
	session, err := node.prepareSession()
	if err != nil {
		return nil, err
	}
	output, err := sessionCb(session)
	if err != nil {
		if len(output) > 0 {
			node.log(LogLevelWarning, string(output))
		}
		node.log(LogLevelWarning, err.Error())
		node.log(LogLevelWarning, "Reconnecting...")

		// set the previous client to nil
		if node.client != nil {
			err = node.client.Close()
			if err != nil {
				node.log(LogLevelWarning, "Error closing previous client: "+err.Error())
			}
			node.client = nil
		}
		session, err = node.prepareSession()
		if err != nil {
			return nil, err
		}

		output, err = sessionCb(session)
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

func (node *SSHNode) CDUnsafe(dir string) error {
	node.lastDir = filepath.Join(node.lastDir, dir)
	if dir != "" {
		// TODO: better handling of escaping path
		cmd := "cd '" + dir + "'"
		_, err := node.RunSyncUnsafe(cmd)
		return err
	}
	return nil
}

func (node *SSHNode) prepareSession() (*ssh.Session, error) {
	if node.client == nil {
		cfg := node.config
		addr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
		node.log(LogLevelInfo, fmt.Sprintf("SSH: Connecting to %v\n", addr))

		client, err := ssh.Dial("tcp", addr, node.sshConfig)
		if err != nil {
			return nil, err
		}
		node.client = client
	}

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
